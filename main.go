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
}
