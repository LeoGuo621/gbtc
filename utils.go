package main

import (
	"bytes"
	"encoding/binary"
	"log"
)

// 整数转化为16进制
func IntToHex(num int64) []byte {
	buf := new(bytes.Buffer) //开辟内存储存字节集
	err := binary.Write(buf, binary.BigEndian, num) //num转化字节集写入buf
	if err != nil {
		log.Panic(err)
	}
	return buf.Bytes()
}

// ReverseBytes reverses a byte array
func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}