package main

import (
	"log"
	"testing"
)

func TestValidNodes(t *testing.T) {

	mnzChain := NewMnzChain()
	mnzChain.RegisterNode("http://127.0.0.1:8002")
	neighbours := mnzChain.nodes
	found := false
	for _, node := range neighbours {
		log.Println("node:", node)
		if node == "127.0.0.1:8002" {
			found = true
		}
	}

	if !found {
		t.Error("not found added node")
	}

}

func TestCreateBlock(t *testing.T) {

	mnzChain := NewMnzChain()
	mnzChain.NewBlock(123, "sa sa prova")

	lastblock := mnzChain.lastBlock()
	if len(mnzChain.chain) != 2 {
		t.Error("error length at genesis block")
	}
	if lastblock.index != 2 {
		t.Error("genesis block index must be 2")
	}
	if lastblock.proof != 123 {
		t.Error("wrong proof")
	}
	if lastblock.previousHash != "sa sa prova" {
		t.Error("wrong previous hash")
	}

}

func TestCreateTransaction(t *testing.T) {

	mnzChain := NewMnzChain()
	mnzChain.NewTransaction("aa", "bb", 10)
	height := len(mnzChain.currentTransactions)
	transaction := mnzChain.currentTransactions[height-1]
	if transaction.Sender != "aa" {
		t.Error("wrong sender")
	}
	if transaction.Recipient != "bb" {
		t.Error("wrong recipient")
	}
	if transaction.Amount != 10 {
		t.Error("wrong amount")
	}
}
