package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"log"
)

const version = byte(0x00)   //钱包版本
const addressChecksumLen = 4 //地址校验和长度

type Wallet struct {
	PrivateKey ecdsa.PrivateKey // 解锁钱包的权限
	PublicKey []byte //即收款地址
}

func NewWallet() *Wallet {
	private, public := newKeyPair()
	wallet := Wallet{
		PrivateKey: private,
		PublicKey:  public,
	}
	return &wallet
}

// 创建公私钥对
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	// 椭圆曲线加密，选用P256椭圆曲线方案
	curve := elliptic.P256()
	// 需要传入曲线与一个随机数生成器
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	//利用私钥推导出公钥
	publicKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, publicKey
}

// 生成校验和
func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])
	return secondSHA[:addressChecksumLen]
}

//公钥哈希处理，先后经过SHA256与RIPEMD160
func HashPubKey(pubkey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubkey)
	R160Hasher := ripemd160.New()
	_, err := R160Hasher.Write(publicSHA256[:]) //传入sha256结果并加密
	if err != nil {
		log.Panic(err)
	}
	publicR160Hash := R160Hasher.Sum(nil)
	return publicR160Hash
}

//获取钱包地址
func (w Wallet) GetAddress() []byte {
	pubKeyHash := HashPubKey(w.PublicKey)
	versionPayload := append([]byte{version}, pubKeyHash...) // 拼接钱包版本号与公钥hash
	checksum := checksum(versionPayload) //获取校验和
	fullPayload := append(versionPayload, checksum...) //拼接钱包版本号、公钥hash、校验和
	address := Base58Encode(fullPayload)
	return address
}

//校验钱包地址
func ValidateAddress(address string) bool {
	fullPayload := Base58Decode([]byte(address))
	actualChecksum := fullPayload[len(fullPayload) -addressChecksumLen:] //获取解码出的校验和
	version := fullPayload[0]                                            // 取得钱包版本
	pubKeyHash := fullPayload[1:len(fullPayload) -addressChecksumLen]
	targetChecksum := checksum(append([]byte{version}, pubKeyHash...)) // 计算版本号+pubKeyHash的校验和
	return bytes.Compare(actualChecksum, targetChecksum) == 0 // 对比是否相等
	
}