package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

type CLI struct {
	blockchain *BlockChain
}

func (cli *CLI) createBlockchain(address string) {
	bc := CreateBlockchain(address)
	bc.db.Close()
	fmt.Println("创建成功, 创世者: ", address)
}

func (cli *CLI) getBalance(address string) {
	bc := NewBlockChain()
	defer bc.db.Close()
	balance := 0
	UTXOs := bc.FindUTXO(address)
	for _, out := range UTXOs {
		balance += out.Value
	}
	fmt.Printf("地址%s的余额为%d\n", address, balance)
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("getbalance -address 根据地址查询金额")
	fmt.Println("createblockchain 根据地址创建区块链")
	fmt.Println("send -from FROMADDR -to TOADDR -amount AMOUNT 转账")
	fmt.Println("addblock 向区块链增加区块")
	fmt.Println("showchain 显示区块链")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

//func (cli *CLI) addBlock(data string) {
//	cli.blockchain.AddBlock(data) //二次调用增加区块
//	fmt.Println("增加区块成功")
//}
func (cli *CLI) send(from, to string, amount int) {
	bc := NewBlockChain()
	defer bc.db.Close()
	tx := NewUTXOTransaction(from, to, amount, bc)
	bc.MineBlock([]*Transaction{tx})
	fmt.Println("交易成功")
}

func (cli *CLI) showBlockChain() {
	//bci := cli.blockchain.Iterator() //创建循环迭代器
	//for {
	//	block := bci.Next()
	//	fmt.Printf("父区块哈希: %x\n", block.PrevBlockHash)
	//	fmt.Printf("数据: %s\n", block.Data)
	//	fmt.Printf("当前哈希: %x\n", block.Hash)
	//	pow := NewProofOfWork(block)
	//	fmt.Printf("pow校验: %s\n", strconv.FormatBool(pow.Validate()))
	//	fmt.Println()
	//
	//	// 从后往前遍历，找到了创世区块，终止循环
	//	if len(block.PrevBlockHash) == 0 {
	//		break
	//	}
	//}

	bc := NewBlockChain()
	defer bc.db.Close()
	bci := bc.Iterator()
	for {
		block := bci.Next()
		fmt.Printf("父区块哈希: %x\n", block.PrevBlockHash)
		fmt.Printf("当前哈希: %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("pow校验: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
		// 从后往前遍历，找到了创世区块，终止循环
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}

func (cli *CLI) Run() {
	cli.validateArgs()
	//addblockcmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	showchaincmd := flag.NewFlagSet("showchain", flag.ExitOnError)
	getbalancecmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createblockchaincmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendcmd := flag.NewFlagSet("send", flag.ExitOnError)

	//addBlockData := addblockcmd.String("data", "", "Block Data")
	getbalanceaddress := getbalancecmd.String("address", "", "查询地址")
	createblockchainaddress := createblockchaincmd.String("address", "", "创世地址")
	sendfrom := sendcmd.String("from", "", "转出方")
	sendto := sendcmd.String("to", "", "接收方")
	sendamount := sendcmd.Int("amount", 0, "金额")
	switch os.Args[1] {
	//case "addblock":
	//	err := addblockcmd.Parse(os.Args[2:])
	//	if err != nil {
	//		log.Panic(err)
	//	}
	case "send":
		err := sendcmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getbalance":
		err := getbalancecmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createblockchaincmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "showchain":
		err := showchaincmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	//if addblockcmd.Parsed() {
	//	if *addBlockData == "" {
	//		addblockcmd.Usage()
	//		os.Exit(1)
	//	}else {
	//		cli.addBlock(*addBlockData)
	//	}
	//}

	if showchaincmd.Parsed() {
		cli.showBlockChain()
	}
	if sendcmd.Parsed() {
		if *sendfrom == "" || *sendto == "" || *sendamount <= 0 {
			sendcmd.Usage()
			os.Exit(1)
		}
		cli.send(*sendfrom, *sendto, *sendamount)
	}

	if getbalancecmd.Parsed() {
		if *getbalanceaddress == "" {
			getbalancecmd.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getbalanceaddress)
	}

	if createblockchaincmd.Parsed() {
		if *createblockchainaddress == "" {
			createblockchaincmd.Usage()
			os.Exit(1)
		}
		cli.createBlockchain(*createblockchainaddress)
	}
}