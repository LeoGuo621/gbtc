package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"
)

//出块奖励
const subsidy = 1000

type Transaction struct {
	ID []byte //交易编号
	Vin []TXInput //输入
	Vout []TXOutput //输出
}

// 序列化交易
func (tx Transaction) Serialize() []byte {
	var encoded bytes.Buffer
	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	return encoded.Bytes()
}

// 反序列化交易
func DeserializeTX(data []byte) Transaction {
	var transaction Transaction
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&transaction)
	if err != nil {
		log.Panic(err)
	}
	return transaction
}

//对交易进行hash
func (tx *Transaction) Hash() []byte {
	var hash [32]byte
	txCopy := *tx
	txCopy.ID = []byte{}
	hash = sha256.Sum256(txCopy.Serialize())
	return hash[:]
}

//签名交易
func (tx *Transaction) Sign(privateKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	if tx.IsCoinbase() {
		return //铸币交易无需签名
	}
	for _, vin := range tx.Vin {
		//如果该交易的输入交易id为nil则代表先前交易有误
		if prevTXs[hex.EncodeToString(vin.TXid)].ID == nil {
			//TODO: 校验
			log.Panic("交易输入有误")
		}
	}
	txCopy := tx.TrimmedCopy()
	for inID, vin := range txCopy.Vin {
		//获取vin引用的先前交易
		prevTX := prevTXs[hex.EncodeToString(vin.TXid)]
		//设定vin的签名为空与公钥
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKey = prevTX.Vout[vin.Vout].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Vin[inID].PubKey = nil

		//dataToSign := fmt.Sprintf("%x\n", txCopy)
		r, s, err := ecdsa.Sign(rand.Reader, &privateKey, txCopy.ID)
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)
		//签名
		tx.Vin[inID].Signature = signature
	}
}

//获取交易在需要签名时的缺失签名等字段的副本
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TXInput
	var outputs []TXOutput
	for _, vin := range tx.Vin {
		inputs = append(inputs, TXInput{vin.TXid, vin.Vout, nil, nil})
	}
	for _, vout := range tx.Vout {
		outputs = append(outputs, TXOutput{vout.Value, vout.PubKeyHash})
	}
	txCopy := Transaction{
		ID:   tx.ID,
		Vin:  inputs,
		Vout: outputs,
	}
	return txCopy
}

func (tx Transaction) String() string {
	var lines []string
	lines = append(lines, fmt.Sprintf("Transaction %x\n", tx.ID))
	for i, input := range tx.Vin {
		lines = append(lines, fmt.Sprintf("Input %d", i))
		lines = append(lines, fmt.Sprintf("TXid %x", input.TXid))
		lines = append(lines, fmt.Sprintf("Vout %d", input.Vout))
		lines = append(lines, fmt.Sprintf("Signature %d", input.Signature))
		lines = append(lines, fmt.Sprintf("PubKey %d", input.PubKey))
	}
	for i, output := range tx.Vout {
		lines = append(lines, fmt.Sprintf("Out %d", i))
		lines = append(lines, fmt.Sprintf("Value %d", output.Value))
		lines = append(lines, fmt.Sprintf("PubKeyHash %d", output.PubKeyHash))
	}
	return strings.Join(lines, "\n")
}

//签名认证
func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}
	for _, vin := range tx.Vin {
		//如果该交易的输入交易id为nil则代表先前交易有误
		if prevTXs[hex.EncodeToString(vin.TXid)].ID == nil {
			log.Panic("交易输入有误")
		}
	}
	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for inID, vin := range tx.Vin {
		prevTX := prevTXs[hex.EncodeToString(vin.TXid)]
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKey = prevTX.Vout[vin.Vout].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Vin[inID].PubKey = nil

		r, s := big.Int{}, big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])

		x, y := big.Int{}, big.Int{}
		keyLen := len(vin.PubKey)
		x.SetBytes(vin.PubKey[:(keyLen / 2)])
		y.SetBytes(vin.PubKey[(keyLen / 2):])

		//dataToVerify := fmt.Sprintf("%x\n", txCopy)
		rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
		if ecdsa.Verify(&rawPubKey, txCopy.ID, &r, &s) == false {
			return false
		}
		//txCopy.Vin[inID].PubKey = nil
	}
	return true
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
	//若data为空，将他设置成随机值
	if data == "" {
		//data = fmt.Sprintf("给%s的挖矿奖励", to)
		randData := make([]byte, 20)
		_, err := rand.Read(randData)
		if err != nil {
			log.Panic(err)
		}
		data = fmt.Sprintf("%x", randData)
	}
	txin := TXInput{
		TXid:      []byte{},
		Vout:      -1, //-1代表是coinbase交易
		Signature: nil,
		PubKey:    []byte(data),
	}
	txout := *NewTXOutput(subsidy, to)
	tx := Transaction{
		ID:   nil,
		Vin:  []TXInput{txin},
		Vout: []TXOutput{txout},
	}
	tx.ID = tx.Hash()
	return &tx
}

// TODO: 钱包功能下的转账交易
func NewUTXOTransaction(from, to string, amount int, bc *BlockChain) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	wallets, err := NewWallets()
	if err != nil {
		log.Panic(err)
	}
	wallet := wallets.GetWallet(from)
	//获取公钥哈哈希
	pubKeyHash := HashPubKey(wallet.PublicKey)
	//查询最小的可用于支付的UTXO
	accumulateRewards, validOutputs := bc.FindSpendableOutputs(pubKeyHash, amount)
	if accumulateRewards < amount {
		log.Panic("交易金额不足")
	}
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}
		for _, out := range outs {
			//注意input中的签名一开始为nil
			input := TXInput{txID, out, nil, wallet.PublicKey}
			inputs = append(inputs, input)
		}
	}
	outputs = append(outputs, *NewTXOutput(amount, to))
	if accumulateRewards > amount {
		//找零
		outputs = append(outputs, *NewTXOutput(accumulateRewards - amount, from))
	}
	tx := Transaction{nil, inputs, outputs}
	tx.SetID() //设置交易hash
	bc.SignTransaction(&tx, wallet.PrivateKey)
	return &tx
}