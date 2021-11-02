package main

import (
	"bytes"
)

type TXOutput struct {
	Value int //output中的Value就是“币”
	//ScriptPubKey string //保存一个任意的用户定义的钱包地址(实际上是接受方验签, 代表受益方是谁)
	PubKeyHash []byte
}

// 以某个地址锁定交易输出
func (out *TXOutput) Lock(address string) {
	fullPayload := Base58Decode([]byte(address))
	pubkeyHash := fullPayload[1:len(fullPayload) - 4] //截取公钥hash
	out.PubKeyHash = pubkeyHash
}

// 判断该输出是否被某个公钥hash锁住
func (out *TXOutput) IsLockedWithKey(pubkeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubkeyHash) == 0
}

func NewTXOutput(value int, address string) *TXOutput {
	txo := &TXOutput{
		Value:      value,
		PubKeyHash: nil,
	}
	//Lock本质上就是对pubKeyHash进行赋值
	txo.Lock(address)
	return txo
}