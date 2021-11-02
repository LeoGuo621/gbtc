package main

import (
	"github.com/boltdb/bolt"
	"log"
)

//区块迭代器
type BlockChainIterator struct {
	currentHash []byte
	db *bolt.DB
}

//根据迭代器取下一个区块
func (bci *BlockChainIterator) Next() *Block {
	var block *Block
	err := bci.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		encodedBlock := bucket.Get(bci.currentHash) //抓取当前迭代器所在位置的区块数据
		block = Deserialize(encodedBlock) //解码
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	bci.currentHash = block.PrevBlockHash
	return block
}