package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

//出块奖励
const subsidy = 10

type TXInput struct {
	TXid []byte //储存该交易所引用的交易id
	Vout int //Vout保存交易中的一个output的索引
	ScriptSig string //由于没有实现address，保存一个任意的用户定义的钱包地址(实际上是输入方验签, 用于验证该交易由谁发出)
}

//检查地址是否能启动事务输入
func (input *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return input.ScriptSig == unlockingData
}

type TXOutput struct {
	Value int //output中的Value就是“币”
	ScriptPubKey string //保存一个任意的用户定义的钱包地址(实际上是接受方验签, 代表受益方是谁)
}

//判断地址是否可以解锁输出
func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}

type Transaction struct {
	ID []byte //交易编号
	Vin []TXInput //输入
	Vout []TXOutput //输出
}

//检查交易是否为coinbase交易
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].TXid) == 0 && tx.Vin[0].Vout == -1
}

//设置交易ID，交易ID实际上就是交易的Hash
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer //开辟内存
	var hash [32]byte
	encoder := gob.NewEncoder(&encoded)  //编码对象
	err := encoder.Encode(tx) //编码
	if err != nil {
		log.Panic(err)
	}
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

//创建铸币交易
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("给%s的挖矿奖励", to)
	}
	txin := TXInput{
		TXid:      []byte{},
		Vout:      -1, //-1代表是coinbase交易
		ScriptSig: data,
	}
	txout := TXOutput{
		Value:        subsidy,
		ScriptPubKey: to,
	}
	tx := Transaction{
		ID:   nil,
		Vin:  []TXInput{txin},
		Vout: []TXOutput{txout},
	}
	return &tx
}

func NewUTXOTransaction(from, to string, amount int, bc *BlockChain) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	//查询最小的可用于支付的UTXO
	accumulateRewards, validOutputs := bc.FindSpendableOutputs(from, amount)
	if accumulateRewards < amount {
		log.Panic("交易金额不足")
	}
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}
		for _, out := range outs {
			input := TXInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}
	outputs = append(outputs, TXOutput{amount, to})
	if accumulateRewards > amount {
		//找零
		outputs = append(outputs, TXOutput{accumulateRewards - amount, from})
	}
	tx := Transaction{nil, inputs, outputs}
	tx.SetID()
	return &tx
}