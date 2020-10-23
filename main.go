package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var mnzchain = NewMnzChain()

// MineHandler handles http request get ip:port/mine
func MineHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("mine handler")

	type Block struct {
		Index        int           `json:"index"`
		Message      string        `json:"message"`
		Transactions []transaction `json:"transactions"`
		Proof        int           `json:"proof"`
		PreviousHash string        `json:"previousHash"`
	}

	result := Block{}

	var code int

	code = http.StatusBadRequest
	result.Message = "unsupport"

	req.ParseForm()

	if req.Method == "GET" {

		lastBlock := mnzchain.lastBlock()
		lastproof := lastBlock.proof
		previousHash := lastBlock.previousHash

		log.Printf("last proof:%d\n", lastproof)
		log.Printf("previous hash:%s\n", previousHash)

		proof := mnzchain.proofOfWork(lastproof)

		var transactionReward transaction
		transactionReward.Sender = "0"
		transactionReward.Recipient = "a random address"
		transactionReward.Amount = 1

		block := mnzchain.NewBlock(proof, previousHash)

		code = http.StatusOK

		result.Message = "New Block Forged"
		result.Index = block.index
		result.Proof = block.proof
		result.Transactions = block.transactions
		result.PreviousHash = block.previousHash
	}

	bytes, _ := json.Marshal(result)
	w.WriteHeader(code)
	fmt.Fprintf(w, string(bytes))

}

// NewTransactionHandler handles new transaction post request ip:port/transactions/new
func NewTransactionHandler(w http.ResponseWriter, req *http.Request) {

	log.Println("new transactions handler")

	type Transaction struct {
		Sender    string `json:"sender"`
		Recipient string `json:"recipient"`
		Amount    int    `json:"amount"`
	}

	type ResponseJSON struct {
		Message string `json:"message"`
	}

	transaction := Transaction{}
	result := ResponseJSON{}

	var code int

	code = http.StatusBadRequest
	result.Message = "unsupport"

	req.ParseForm()

	b, _ := ioutil.ReadAll(req.Body)

	log.Printf("body:%s\n", b)

	if req.Method == "POST" {

		err := json.Unmarshal([]byte(b), &transaction)

		if err != nil {
			log.Println(err.Error())

			code = http.StatusInternalServerError
			result.Message = "json unparse failed"

		} else {

			log.Printf("%+v\n", transaction)

			index := mnzchain.NewTransaction(transaction.Sender, transaction.Recipient, transaction.Amount)

			code = http.StatusCreated
			result.Message = fmt.Sprintf("New nodes have been added to Block %d", index)

		}

	}

	bytes, _ := json.Marshal(result)
	w.WriteHeader(code)
	fmt.Fprintf(w, string(bytes))

}

// ChainHandler  get ip:port/chain  walks the chain and output it via http
func ChainHandler(w http.ResponseWriter, req *http.Request) {

	//TODO output in json to be interoperable
	log.Println("chain transaction handler")

	w.WriteHeader(http.StatusOK)

	req.ParseForm()

	if req.Method == "GET" {
		fmt.Fprintf(w, "Current Transactions: %v\n", mnzchain.currentTransactions)
		fmt.Fprintf(w, "Chain:\n")
		for _, blk := range mnzchain.chain {
			fmt.Fprintf(w, "Block: %s\n", blk.description())
		}
		fmt.Fprintf(w, "Nodes: %v\n", mnzchain.nodes)
	}

}

// RegisterNodesHandler register new network nodes :  post ip:port/nodes/register
func RegisterNodesHandler(w http.ResponseWriter, req *http.Request) {

	log.Println("nodes register handler")

	type Nodes struct {
		Nodes []string
	}

	type ResponseJSONBean struct {
		Message string   `json:"message"`
		Data    []string `json:"total_nodes"`
	}

	nodeGroup := Nodes{}
	result := ResponseJSONBean{}

	var code int

	code = http.StatusBadRequest
	result.Message = "unsupport"

	req.ParseForm()

	b, _ := ioutil.ReadAll(req.Body)
	log.Printf("body:%s\n", b)

	if req.Method == "POST" {

		err := json.Unmarshal([]byte(b), &nodeGroup)

		if err != nil {
			log.Println(err.Error())

			code = http.StatusInternalServerError
			result.Message = "json unparse failed"

		} else {

			result.Data = make([]string, 0)

			log.Printf("%+v\n", nodeGroup)

			for _, node := range nodeGroup.Nodes {

				log.Printf("node:%s\n", node)
				mnzchain.RegisterNode(node)
				result.Data = append(result.Data, node)
			}

			code = http.StatusCreated
			result.Message = "New nodes have been added"

		}

	}

	bytes, _ := json.Marshal(result)
	w.WriteHeader(code)
	fmt.Fprintf(w, string(bytes))

}

// ConsensusHandler TODO get ip:port/nodes/resolve
func ConsensusHandler(w http.ResponseWriter, req *http.Request) {

	log.Println("nodes resolve register")
	w.Write([]byte("sorry, not yet implemented!"))
	//TODO

}

func main() {

	port := flag.String("port", "8001", "use -port <port number>")
	flag.Parse()

	log.Printf("port is:%s\n", *port)

	http.HandleFunc("/mine", MineHandler)
	http.HandleFunc("/transactions/new", NewTransactionHandler)
	http.HandleFunc("/nodes/register", RegisterNodesHandler)
	http.HandleFunc("/nodes/resolve", ConsensusHandler)
	http.HandleFunc("/chain", ChainHandler)

	http.ListenAndServe(fmt.Sprintf(":%s", *port), nil)
}
