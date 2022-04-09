package main

import (
	"fmt"
	"go_blockchain/block"
	"go_blockchain/wallet"
	"log"
)

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	walletM := wallet.NewWallet() // Mさん（マイニングする人）
	walletA := wallet.NewWallet()
	walletB := wallet.NewWallet()

	// Blockchain
	blockchain := block.NewBlockchain(walletM.BlockchainAddress())

	// AさんがBさんに1.0コインを送る
	t := wallet.NewTransaction(walletA.PrivateKey(), walletA.PublicKey(), walletA.BlockchainAddress(), walletB.BlockchainAddress(), 1.0)
	isAdded := blockchain.AddTransaction(walletA.BlockchainAddress(), walletB.BlockchainAddress(), 1.0,
		walletA.PublicKey(), t.GenerateSignature())
	fmt.Println("Added? ", isAdded)

	blockchain.Mining()
	blockchain.Print()

	fmt.Printf("Aさん %.1f\n", blockchain.CalculateTotalAmount(walletA.BlockchainAddress()))
	fmt.Printf("Bさん %.1f\n", blockchain.CalculateTotalAmount(walletB.BlockchainAddress()))
	fmt.Printf("マイニングする人 %.1f\n", blockchain.CalculateTotalAmount(walletM.BlockchainAddress()))
}
