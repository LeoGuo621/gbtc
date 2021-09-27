package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	Timestamp int64
	//Data []byte // 交易数据
	Transactions  []*Transaction
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

/*
// 设定Block对象hash
func (b *Block) SetHash() {
	// 处理当前的时间，转化位10进制字符串，再变为字节
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	// 组合需要hash的数据
	header := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(header)
	b.Hash = hash[:]
}
 */

// 创建区块
func NewBlock(transactions []*Transaction, prevBlockHash []byte) *Block {
	// 取得一个区块初始化之后的指针
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Transactions:  transactions,
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
		Nonce:         0,
	}
	//block.SetHash()
	//TODO: proof of work 计算hash
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Hash = hash
	block.Nonce = nonce
	return block
}
// 创世区块
func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

//TODO: 持久化

//转化区块对象为字节集，可以写入文件
func (block *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err)
	}
	return res.Bytes()
}

//读取文件，读到二进制字节集，二进制字节集转化为对象
func Deserialize(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}

//把所有交易叠加在一起取hash
func (block *Block) HashTransaction() []byte {
	var txHashes [][]byte
	var txHash [32]byte
	for _, tx := range block.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:]
}