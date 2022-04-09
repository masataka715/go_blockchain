package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

const (
	MINING_DIFFICULTY = 3                // nouceを求める際に、先頭3つが000の値を探す
	MINING_SENDER     = "THE BLOCKCHAIN" // マイニングする人(報酬を受け取る人)から見た、送信者（node側）のブロックチェーンアドレス
	MINING_REWARD     = 1.0              // マイニングに成功した場合の報酬
)

// Block
type Block struct {
	nonce        int
	previousHash [sha256.Size]byte
	timestamp    int64
	transactions []*Transaction
}

func NewBlock(nouce int, previousHash [sha256.Size]byte, transactions []*Transaction) *Block {
	b := new(Block) // newでpoint型を明示できる
	b.timestamp = time.Now().UnixNano()
	b.nonce = nouce
	b.previousHash = previousHash
	b.transactions = transactions
	return b
}

func (b *Block) Print() {
	fmt.Printf("timestamp        %d\n", b.timestamp)
	fmt.Printf("nonce            %d\n", b.nonce)
	fmt.Printf("previous_hash    %x\n", b.previousHash)
	for _, t := range b.transactions {
		t.Print()
	}
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
		Transactions []*Transaction    `json:"transactions"`
	}{
		Timestamp:    b.timestamp,
		Nonce:        b.nonce,
		PreviousHash: b.previousHash,
		Transactions: b.transactions,
	})
}

// Blockchain
type Blockchain struct {
	transactionPool   []*Transaction
	chain             []*Block
	blockchainAddress string
}

func NewBlockchain(blockchainAddress string) *Blockchain {
	b := &Block{}
	bc := new(Blockchain)
	bc.blockchainAddress = blockchainAddress
	bc.CreateBlock(0, b.Hash()) // 1個目のブロック
	return bc
}

func (bc *Blockchain) CreateBlock(nonce int, previousHash [sha256.Size]byte) *Block {
	// BlockchainのtransactionPoolから、Blockのtransactionsに渡す
	b := NewBlock(nonce, previousHash, bc.transactionPool)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*Transaction{} // 渡した後のtransactionPoolは空にする
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
	fmt.Printf("%s\n\n", strings.Repeat("*", 25))
}

func (bc *Blockchain) AddTransaction(sender string, recipient string, value float32) {
	t := NewTransaction(sender, recipient, value)
	bc.transactionPool = append(bc.transactionPool, t)
}

func (bc *Blockchain) CopyTransactionPool() []*Transaction {
	transactions := make([]*Transaction, 0)
	for _, t := range bc.transactionPool {
		transactions = append(transactions,
			NewTransaction(t.senderBlockchainAddress,
				t.recipientBlockchainAddress,
				t.value),
		)
	}
	return transactions
}

func (bc *Blockchain) ValidProof(nouce int, previousHash [sha256.Size]byte, transactions []*Transaction, difficulty int) bool {
	zeros := strings.Repeat("0", difficulty)
	guessBlock := Block{
		nonce:        nouce,
		previousHash: previousHash,
		timestamp:    0,
		transactions: transactions,
	}
	gueesHashStr := fmt.Sprintf("%x", guessBlock.Hash())
	return gueesHashStr[:difficulty] == zeros
}

// ProofOfWork
// nouceを求める演算処理
func (bc *Blockchain) ProofOfWork() int {
	transactions := bc.CopyTransactionPool()
	previousHash := bc.LastBlock().Hash()
	nonce := 0
	for !bc.ValidProof(nonce, previousHash, transactions, MINING_DIFFICULTY) {
		nonce += 1
	}
	return nonce
}

func (bc *Blockchain) Mining() bool {
	bc.AddTransaction(MINING_SENDER, bc.blockchainAddress, MINING_REWARD)
	nonce := bc.ProofOfWork()
	previousHash := bc.LastBlock().Hash()
	bc.CreateBlock(nonce, previousHash)
	log.Println("action=mining, status=success")
	return true
}

// CalculateTotalAmount
// 今持っているコインの合計を求める
func (bc *Blockchain) CalculateTotalAmount(blockchainAddress string) float32 {
	var totalAmount float32 = 0.0
	for _, b := range bc.chain {
		for _, t := range b.transactions {
			value := t.value
			if blockchainAddress == t.recipientBlockchainAddress {
				totalAmount += value
			}

			if blockchainAddress == t.senderBlockchainAddress {
				totalAmount -= value
			}
		}
	}
	return totalAmount
}

type Transaction struct {
	senderBlockchainAddress    string
	recipientBlockchainAddress string
	value                      float32 // 送金する額
}

func NewTransaction(sender string, recipient string, value float32) *Transaction {
	return &Transaction{sender, recipient, value}
}

func (t *Transaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	fmt.Printf(" sender_blockchain_address  %s\n", t.senderBlockchainAddress)
	fmt.Printf(" recipient_blockchain_address  %s\n", t.recipientBlockchainAddress)
	fmt.Printf(" value  %.1f\n", t.value)
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"sender_blockchain_address"`
		Recipient string  `json:"recipient_blockchain_address"`
		Value     float32 `json:"value"`
	}{
		Sender:    t.senderBlockchainAddress,
		Recipient: t.recipientBlockchainAddress,
		Value:     t.value,
	})
}

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	myBlockchainAddress := "my_blockchain_address" // マイニングする人のアドレス

	blockChain := NewBlockchain(myBlockchainAddress)
	blockChain.Print()

	// AさんがBさんに1.0コインを送金
	blockChain.AddTransaction("A", "B", 1.0)
	blockChain.Mining() // 2個目のブロック
	blockChain.Print()

	// CさんがDさんに2.0コインを送金
	blockChain.AddTransaction("C", "D", 2.0)
	// XさんがYさんに3.0コインを送金
	blockChain.AddTransaction("X", "Y", 3.0)
	blockChain.Mining() // 3個目のブロック
	blockChain.Print()

	fmt.Printf("my %.1f\n", blockChain.CalculateTotalAmount(myBlockchainAddress))
	fmt.Printf("C %.1f\n", blockChain.CalculateTotalAmount("C"))
	fmt.Printf("D %.1f\n", blockChain.CalculateTotalAmount("D"))
}
