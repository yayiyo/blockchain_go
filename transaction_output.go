package main

import (
	"bytes"
)

//交易输出
type TXOutput struct {
	Value      int
	PubkeyHash []byte
}

func (out *TXOutput) Lock(address []byte) {
	pubKeyHash := Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	out.PubkeyHash = pubKeyHash
}

//检查提供的公钥是否用于锁定输出
func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubkeyHash, pubKeyHash) == 0
}

//创建TXOutput
func NewTXOutput(value int, address string) *TXOutput {
	txo:=&TXOutput{value,nil}
	txo.Lock([]byte(address))

	return txo
}
