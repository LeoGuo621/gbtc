package main

import (
	"fmt"
	"log"
)

func (cli *CLI) createBlockchain(address string)  {
	if !ValidateAddress(address) {
		log.Panic("invalid address")
	}
	bc := CreateBlockchain(address)
	bc.db.Close()
	fmt.Println("创建区块链成功")
}
