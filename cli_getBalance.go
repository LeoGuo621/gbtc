package main

import (
	"fmt"
	"log"
)

func (cli *CLI) getBalance(address string)  {
	if !ValidateAddress(address) {
		log.Panic("invalid address")
	}
	bc := NewBlockChain()
	defer bc.db.Close()
	balance := 0
	//通过地址获取pubKeyHash
	// fullPayload = version + pubKeyHash + checksum
	fullPayload := Base58Decode([]byte(address))
	pubKeyHash := fullPayload[1:len(fullPayload) - 4]
	UTXOs := bc.FindUTXO(pubKeyHash)
	for _, out := range UTXOs {
		balance += out.Value
	}
	fmt.Printf("地址%s的余额为%d\n", address, balance)
}