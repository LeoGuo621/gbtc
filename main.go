package main

func main() {
	//bc := NewBlockChain()
	//bc.AddBlock("11111")
	////bc.AddBlock("22222")
	////bc.AddBlock("33333")
	//for _, block := range bc.blocks {
	//	fmt.Printf("PrevHash: %x\n", block.PrevBlockHash)
	//	fmt.Printf("Data: %s\n", block.Data)
	//	fmt.Printf("Hash: %x\n", block.Hash)
	//	pow := NewProofOfWork(block)
	//	fmt.Printf("pow校验: %s\n", strconv.FormatBool(pow.Validate()))
	//	fmt.Println()
	//}

	//blockchain := CreateBlockchain("leo")
	//defer blockchain.db.Close() //延迟关闭数据库
	//cli := CLI{blockchain}
	cli := CLI{}
	cli.Run()
}
