package main

import (
	"fmt"
	"go_blockchain/wallet"
	"log"
)

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	w := wallet.NewWallet()
	fmt.Println("PrivateKey: " + w.PrivateKeyStr())
	fmt.Println("PublicKey: " + w.PublicKeyStr())
	fmt.Println("BlockchainAddress:  " + w.BlockchainAddress())

	// Bさんに1.0コインを送る
	t := wallet.NewTransaction(w.PrivateKey(), w.PublicKey(), w.BlockchainAddress(), "B", 1.0)
	fmt.Printf("Signature %s\n", t.GenerateSignature())
}
