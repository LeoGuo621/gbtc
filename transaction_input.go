package main

import "bytes"

type TXInput struct {
	TXid []byte //储存该交易所引用的交易id
	Vout int //Vout保存交易中的一个output的索引
	Signature []byte //签名
	PubKey []byte // 公钥
}

// 检测输入是否是合法的
func (input *TXInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := HashPubKey(input.PubKey)
	return bytes.Compare(lockingHash, pubKeyHash) == 0
}