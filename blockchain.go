package main

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
)

//TODO: 持久化
//当前目录下的数据库文件
const dbFile = "blockchain.db"
const blockBucket = "blocks"
const genesisCoinbaseData = "一生不弱于人 --Leo"

type BlockChain struct {
	//blocks []*Block // 区块指针数组
	//TODO: 持久化，重写区块链结构
	//tip意为"末梢", 这里记录链中最新一个区块的hash
	tip []byte
	db *bolt.DB
}


/*
func (bc *BlockChain) AddBlock(data string) {

	//prevBlock := bc.blocks[len(bc.blocks) - 1]
	//newBlock := NewBlock(data, prevBlock.Hash)
	//bc.blocks = append(bc.blocks, newBlock)

	//TODO: 持久化。1.获取上一区块hash 2.生成新区块 3.存入数据库 4.更新区块链的last区块hash
	var lastHash []byte //上一区块hash
	err := bc.db.View(func(tx *bolt.Tx) error {
		block := tx.Bucket([]byte(blockBucket))
		lastHash = block.Get([]byte("last"))
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	newBlock := NewBlock(data, lastHash)
	err = bc.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket)) //取出区块所在的数据
		err := bucket.Put(newBlock.Hash, newBlock.Serialize()) // 向数据库重存入数据
		if err != nil {
			log.Panic(err)
		}
		err = bucket.Put([]byte("last"), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}
		bc.tip = newBlock.Hash
		return nil
	})
}

 */

//TODO: 挖矿
func (bc *BlockChain) MineBlock(transactions []*Transaction)  {
	//lashHash保存当前数据库最新区块的hash
	var lastHash []byte
	for _, tx := range transactions {
		if bc.VerifyTransaction(tx) != true {
			log.Panic("invalid tx")
		}
	}

	err := bc.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		lastHash = bucket.Get([]byte("last"))
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	//创建新区块
	newBlock := NewBlock(transactions, lastHash)
	err = bc.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		err := bucket.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}
		err = bucket.Put([]byte("last"), newBlock.Hash)
		bc.tip = newBlock.Hash
		return nil
	})
}

//TODO:查找一个address所有未使用输出的交易 (此方法相对复杂)
func (bc *BlockChain) FindUnspentTransactions(pubKeyHash []byte) []Transaction {
	var unspentTXs []Transaction
	//某个交易已经花费的交易输出，构建tx->VOutIdx的map
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()
	for {
		//从新区块向旧区块遍历
		block := bci.Next()
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)
		Outputs:
			for outIdx, out := range tx.Vout {
				//若该交易输出已被花费，直接跳过此输出
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				//这里找到了一个未被花费的交易输出
				//如果可以被pubKeyHash解锁，就是属于pubKeyHash的utxo
				if out.IsLockedWithKey(pubKeyHash) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}
			//维护spentTXOs
			//首先排除coinbase交易
			if tx.IsCoinbase() == false {
				for _, in := range tx.Vin {
					//若该交易的输入可由pubKeyHash解锁，说明pubKeyHash对应的address已使用过
					if in.UsesKey(pubKeyHash) {
						inTXID := hex.EncodeToString(in.TXid)
						spentTXOs[inTXID] = append(spentTXOs[inTXID], in.Vout)
					}
				}
			}
		}
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	return unspentTXs
}
//TODO: 返回一个address所有未使用的交易输出
func (bc *BlockChain) FindUTXO(pubKeyHash []byte) []TXOutput {
	var UTXOs []TXOutput
	//先找到所有交易
	unspentTransactions := bc.FindUnspentTransactions(pubKeyHash)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Vout {
			//可解锁代表是address的资产
			if out.IsLockedWithKey(pubKeyHash) {
				UTXOs = append(UTXOs, out)
			}
		}
	}
	return UTXOs
}

//查找可以执行的最小交易
func (bc *BlockChain) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTXs := bc.FindUnspentTransactions(pubKeyHash)
	accumulated := 0 //累计金额
	Work:
		for _, tx := range unspentTXs {
			txID := hex.EncodeToString(tx.ID)
			for outIdx, out := range tx.Vout {
				if out.IsLockedWithKey(pubKeyHash) && accumulated < amount {
					accumulated += out.Value //金额累加
					unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
					if accumulated >= amount {
						break Work
					}
				}
			}
		}

	return accumulated, unspentOutputs
}

func (bc *BlockChain) Iterator() *BlockChainIterator {
	bci := &BlockChainIterator{
		currentHash: bc.tip,
		db:          bc.db,
	}
	return bci
}



//判断数据库是否存在
func dbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

func NewBlockChain() *BlockChain {
	//return &BlockChain{[]*Block{NewGenesisBlock()}}
	if dbExists() == false {
		fmt.Println("数据库不存在, 请创建数据库")
		os.Exit(1)
	}
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	//处理数据更新
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket)) //按照名称打开数据库的桶
		tip = bucket.Get([]byte("last"))
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := BlockChain{tip, db}
	return &bc
}

//创建一个区块链存入数据库
func CreateBlockchain(address string) *BlockChain {
	if dbExists() {
		fmt.Println("数据库已存在，无需创建")
		os.Exit(1)
	}
	var tip []byte
	coinbaseTX := NewCoinbaseTX(address, genesisCoinbaseData)
	genesis := NewGenesisBlock(coinbaseTX) //创建创世区块
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucket([]byte(blockBucket))
		if err != nil {
			log.Panic(err)
		}
		err = bucket.Put(genesis.Hash, genesis.Serialize())
		if err != nil {
			log.Panic(err)
		}
		err = bucket.Put([]byte("last"), genesis.Hash)
		if err != nil {
			log.Panic(err)
		}
		tip = genesis.Hash
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	bc := BlockChain{tip, db}
	return &bc
}

// 交易签名
func (bc *BlockChain) SignTransaction(tx *Transaction, privateKey ecdsa.PrivateKey)  {
	prevTXs := make(map[string]Transaction)
	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.TXid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}
	tx.Sign(privateKey, prevTXs)
}

func (bc *BlockChain) FindTransaction(ID []byte) (Transaction, error) {
	bci := bc.Iterator()
	for {
		block := bci.Next()
		for _, tx := range block.Transactions {
			if bytes.Compare(tx.ID, ID) == 0 {
				return *tx, nil
			}
		}
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	return Transaction{}, errors.New("invalid transaction ID")
}

func (bc *BlockChain) VerifyTransaction(tx *Transaction) bool {
	prevTXs := make(map[string]Transaction)
	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.TXid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}
	return tx.Verify(prevTXs)
}