package main

import (
	"fmt"
	"strconv"
)

func (cli *CLI) showBlockChain() {
	bc := NewBlockChain()
	defer bc.db.Close()
	bci := bc.Iterator()
	for {
		block := bci.Next()
		fmt.Printf("父区块哈希: %x\n", block.PrevBlockHash)
		fmt.Printf("当前哈希: %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("pow校验: %s\n", strconv.FormatBool(pow.Validate()))
		for _, tx := range block.Transactions {
			fmt.Println(tx)
		}
		fmt.Println()

		// 从后往前遍历，找到了创世区块，终止循环
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}
