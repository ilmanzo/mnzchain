package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type transaction struct {
	Sender    string
	Recipient string
	Amount    int
}

type block struct {
	index        int
	timestamp    int64
	transactions []transaction
	proof        int
	previousHash string
}

func (b *block) description() string {
	return fmt.Sprintf("index: %d\n timestamp: %d\n transactions: %v\n proof: %d\n previous hash: %s", b.index, b.timestamp, b.transactions, b.proof, b.previousHash)
}

type MnzChain struct {
	currentTransactions []transaction
	chain               []block
	nodes               []string
}

// NewMnzChain creates a new chain from scratch
func NewMnzChain() *MnzChain {
	bc := &MnzChain{}

	bc.currentTransactions = make([]transaction, 0)
	bc.chain = make([]block, 0)
	bc.nodes = make([]string, 0)

	bc.NewBlock(1, "1")

	return bc
}

func hash(data block) string {

	log.Println("hash.....")
	timestamp := []byte(strconv.FormatInt(data.timestamp, 10))
	previousHash := []byte(data.previousHash)
	dataStr := bytes.Join([][]byte{previousHash, timestamp}, []byte{})
	hash := sha256.Sum256(dataStr)
	hashStr := string(hash[:])
	log.Printf("hash string:%s \n", hashStr)
	return hashStr
}

// in this implementation a valid proof must start with 00
func validProof(lastproof int, proof int) bool {

	// log.Printf("valid proof:%d\n", proof)

	bytelastproof := []byte(strconv.FormatInt(int64(lastproof), 10))
	byteproof := []byte(strconv.FormatInt(int64(proof), 10))
	data := bytes.Join([][]byte{bytelastproof, byteproof}, []byte{})
	guesshash := sha256.Sum256(data)
	log.Println("guess hash:", guesshash)
	return bytes.Equal(guesshash[:2], []byte("00"))
}

// NewBlock creates a new block and appends it to the chain
func (bc *MnzChain) NewBlock(proof int, previousHash string) *block {

	log.Println("new block")
	bl := &block{}
	bl.index = len(bc.chain) + 1
	bl.timestamp = time.Now().Unix()
	bl.transactions = bc.currentTransactions
	bl.proof = proof
	bl.previousHash = previousHash
	transactionLen := len(bc.currentTransactions)
	log.Printf("transaction Len:%d \n", transactionLen)
	log.Printf("index:%d \n", bl.index)
	bc.currentTransactions = bc.currentTransactions[transactionLen:] //clear
	bc.chain = append(bc.chain, *bl)
	return bl
}

func (bc *MnzChain) proofOfWork(lastProof int) int {

	log.Println("proof of work")
	proof := 0
	for !(validProof(lastProof, proof)) {
		proof++
	}
	return proof
}

func (bc *MnzChain) NewTransaction(sender string, recipient string, amount int) int {

	log.Println("new transaction")
	var trans transaction
	trans.Sender = sender
	trans.Recipient = recipient
	trans.Amount = amount
	bc.currentTransactions = append(bc.currentTransactions, trans)
	block := bc.lastBlock()
	log.Printf("index:%d\n", block.index)
	return block.index + 1
}

func (bc *MnzChain) lastBlock() block {

	log.Println("get last block")
	height := len(bc.chain)
	block := bc.chain[height-1]
	return block
}

func (bc *MnzChain) RegisterNode(address string) {

	log.Println("register node")
	u, err := url.Parse(address)
	if err != nil {
		panic(err)
	}
	bc.nodes = append(bc.nodes, u.Host)

}

func (bc *MnzChain) ValidChain(chain []block) bool {

	log.Println("valid chain")
	lastBlock := chain[0]
	currentIndex := 1
	for currentIndex < len(chain) {
		currentBlock := chain[currentIndex]
		log.Println("current block:", currentBlock)
		log.Println("last block:", lastBlock)
		log.Println("---------------------")
		if currentBlock.previousHash != hash(lastBlock) {
			log.Println("error, not a valid chain")
			return false
		}
		if !(validProof(lastBlock.proof, currentBlock.proof)) {
			return false
		}
		lastBlock = currentBlock
		currentIndex++
	}
	return true

}

// ResolveConflicts contacts other nodes to choose the best valid chain
func (bc *MnzChain) ResolveConflicts() bool {

	log.Println("resolving conflicts...")
	neighbours := bc.nodes
	var newChain []block
	// max_length := len(bc.chain)
	for _, node := range neighbours {
		log.Println("node:", node)
		url := fmt.Sprintf("http://%s/chain", node)
		log.Println("url:", url)
		resp, err := http.Get(url)
		if err != nil {
			log.Fatalf("Error %v: Cannot contact %s", err.Error(), url)
		}
		defer resp.Body.Close()
		log.Println("resp:", resp)
		if resp.StatusCode == 200 {
			//TODO
		}

	}

	if newChain != nil {
		bc.chain = newChain
		return true
	}

	return false

}
