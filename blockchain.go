package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

// Block
type Block struct {
	nonce        int
	previousHash [sha256.Size]byte
	timestamp    int64
	transactions []string
}

func NewBlock(nouce int, previousHash [sha256.Size]byte) *Block {
	b := new(Block) // newでpoint型を明示できる
	b.timestamp = time.Now().UnixNano()
	b.nonce = nouce
	b.previousHash = previousHash
	return b
}

func (b *Block) Print() {
	fmt.Printf("timestamp        %d\n", b.timestamp)
	fmt.Printf("nonce            %d\n", b.nonce)
	fmt.Printf("previous_hash    %x\n", b.previousHash)
	fmt.Printf("transactions     %s\n", b.transactions)
}

func (b *Block) Hash() [sha256.Size]byte {
	m, _ := json.Marshal(b)
	return sha256.Sum256([]byte(m))
}

// MarshalJSON
// Blockのfieldがprivateなので、json化時にpublicにすることが必要
func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp    int64             `json:"timestamp"`
		Nonce        int               `json:"nonce"`
		PreviousHash [sha256.Size]byte `json:"previous_hash"`
		Transactions []string          `json:"transactions"`
	}{
		Timestamp:    b.timestamp,
		Nonce:        b.nonce,
		PreviousHash: b.previousHash,
		Transactions: b.transactions,
	})
}

// Blockchain
type Blockchain struct {
	transactionPool []string
	chain           []*Block
}

func NewBlockchain() *Blockchain {
	b := &Block{}
	bc := new(Blockchain)
	bc.CreateBlock(0, b.Hash()) // 1個目のブロック
	return bc
}

func (bc *Blockchain) CreateBlock(nonce int, previousHash [sha256.Size]byte) *Block {
	b := NewBlock(nonce, previousHash)
	bc.chain = append(bc.chain, b)
	return b
}

func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *Blockchain) Print() {
	for i, block := range bc.chain {
		borderString := strings.Repeat("=", 25)
		fmt.Printf("%s Chain %d %s \n", borderString, i, borderString)
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 25))
}

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	blockChain := NewBlockchain()
	blockChain.Print()

	previousHash := blockChain.LastBlock().Hash()
	blockChain.CreateBlock(5, previousHash) // 2個目のブロック
	blockChain.Print()

	previousHash = blockChain.LastBlock().Hash()
	blockChain.CreateBlock(2, previousHash) // 3個目のブロック
	blockChain.Print()
}
