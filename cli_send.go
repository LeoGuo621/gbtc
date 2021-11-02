package main

import (
	"fmt"
	"log"
)

func (cli *CLI) send(from, to string, amount int) {
	if !ValidateAddress(from) {
		log.Panic("invalid from address")
	}
	if !ValidateAddress(to) {
		log.Panic("invalid to address")
	}
	bc := NewBlockChain()
	defer bc.db.Close()
	tx := NewUTXOTransaction(from, to, amount, bc)
	bc.MineBlock([]*Transaction{tx})
	fmt.Println("交易成功")
}