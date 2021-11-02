package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const walletfile = "wallet.dat" //钱包文件

type Wallets struct {
	Wallets map[string]*Wallet //以钱包地址作为key
}

// 创建一个钱包或抓取已存在的钱包
func NewWallets() (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)
	err := wallets.LoadFromFile()
	return &wallets, err
}

//创建一个钱包
func (ws *Wallets) CreateWallet() string {
	wallet := NewWallet()//创建钱包
	address := fmt.Sprintf("%s",wallet.GetAddress())
	ws.Wallets[address] = wallet//保存钱包
	return address
}

//读钱包文件
func (ws *Wallets) LoadFromFile() error {
	myWalletFile := walletfile
	// 若该文件不存在则返回错误
	if _, err := os.Stat(myWalletFile); os.IsNotExist(err) {
		return err
	}
	// 读取钱包内容
	fileContent, err := ioutil.ReadFile(myWalletFile)
	if err != nil {
		log.Panic(err)
	}
	// 解析二进制数据
	var wallets Wallets
	//注册加密算法
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panic(err)
	}
	ws.Wallets = wallets.Wallets
	return nil
}

//写钱包文件
func (ws *Wallets) SaveToFile() {
	var content bytes.Buffer
	myWalletFile := walletfile
	// 注册加密算法, 选用P256椭圆曲线
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&content) // 生成编码器
	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}
	err = ioutil.WriteFile(myWalletFile, content.Bytes(), 0644) // r-4, w-2, x-1
	if err != nil {
		log.Panic(err)
	}
}
// 抓取所有钱包地址
func (ws *Wallets) GetAddresses() []string {
	var addresses []string
	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}
	return addresses
}

//抓取单个钱包
func (ws *Wallets) GetWallet(address string) Wallet {
	return *ws.Wallets[address]
}