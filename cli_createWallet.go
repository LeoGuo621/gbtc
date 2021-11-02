package main

import "fmt"

func (cli *CLI) createWallet() {
	wallets, _ := NewWallets()
	address := wallets.CreateWallet()
	wallets.SaveToFile()
	fmt.Println("新生成的钱包地址为", address)
}