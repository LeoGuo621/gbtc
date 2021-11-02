package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type CLI struct {
	blockchain *BlockChain
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("createWallet 创建钱包")
	fmt.Println("listAddresses 显示所有地址")
	fmt.Println("getBalance -address 根据地址查询金额")
	fmt.Println("createBlockChain 根据地址创建区块链")
	fmt.Println("send -from FROM_ADDR -to TO_ADDR -amount AMOUNT 转账")
	fmt.Println("addBlock 向区块链增加区块")
	fmt.Println("showBlockChain 显示区块链")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) Run() {
	cli.validateArgs()

	createWalletCMD := flag.NewFlagSet("createWallet", flag.ExitOnError)
	listAddressesCMD := flag.NewFlagSet("listAddresses", flag.ExitOnError)

	showBlockChainCMD := flag.NewFlagSet("showBlockChain", flag.ExitOnError)
	getBalanceCMD := flag.NewFlagSet("getBalance", flag.ExitOnError)
	createBlockChainCMD := flag.NewFlagSet("createBlockChain", flag.ExitOnError)
	sendCMD := flag.NewFlagSet("send", flag.ExitOnError)

	getBalanceAddress := getBalanceCMD.String("address", "", "查询地址")
	createBlockChainAddress := createBlockChainCMD.String("address", "", "创世地址")

	sendFrom := sendCMD.String("from", "", "转出方")
	sendTo := sendCMD.String("to", "", "接收方")
	sendAmount := sendCMD.Int("amount", 0, "金额")

	switch os.Args[1] {
	case "createWallet":
		err := createWalletCMD.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "listAddresses":
		err := listAddressesCMD.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCMD.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getBalance":
		err := getBalanceCMD.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createBlockChain":
		err := createBlockChainCMD.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "showBlockChain":
		err := showBlockChainCMD.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if showBlockChainCMD.Parsed() {
		cli.showBlockChain()
	}

	if createWalletCMD.Parsed() {
		cli.createWallet()
	}

	if listAddressesCMD.Parsed() {
		cli.listAddresses()
	}

	if sendCMD.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCMD.Usage()
			os.Exit(1)
		}
		cli.send(*sendFrom, *sendTo, *sendAmount)
	}

	if getBalanceCMD.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCMD.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceAddress)
	}

	if createBlockChainCMD.Parsed() {
		if *createBlockChainAddress == "" {
			createBlockChainCMD.Usage()
			os.Exit(1)
		}
		cli.createBlockchain(*createBlockChainAddress)
	}
}