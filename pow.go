package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

var (
	maxNonce = math.MaxInt64
)

const targetBits = 12

type ProofOfWork struct {
	block *Block
	target *big.Int //存储特定hash值的一个整数
}

//创建工作量证明的挖矿对象
func NewProofOfWork(block *Block) *ProofOfWork {
	//初始化目标整数
	target := big.NewInt(1)
	//左移位，设置难度值
	target.Lsh(target, uint(256 - targetBits))
	pow := &ProofOfWork{block, target}
	return pow
}

//准备数据进行挖矿计算
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join([][]byte{
		pow.block.PrevBlockHash, //父区块hash
		pow.block.HashTransaction(), //当前数据
		IntToHex(pow.block.Timestamp), //16进制的时间戳
		IntToHex(int64(targetBits)), //16进制的难度值
		IntToHex(int64(nonce)), //保存工作量证明的nonce
	}, []byte{})
	return data
}

//开始挖矿
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0
	//fmt.Printf("当前挖矿计算的区块数据: %s", pow.block.Data)
	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data) //计算hash
		fmt.Printf("\r%x", hash) //打印显示hash
		hashInt.SetBytes(hash[:]) //获取要对比的数据

		// Cmp compares x and y and returns:
		//
		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y
		//
		if hashInt.Cmp(pow.target) == -1 {
			break
		}else {
			nonce++
		}
	}
	fmt.Printf("\n\n")
	// nonce为puzzle答案, hash为区块哈希
	return nonce, hash[:]
}

//挖矿校验
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int
	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
	isValid := hashInt.Cmp(pow.target) == -1
	return isValid
}