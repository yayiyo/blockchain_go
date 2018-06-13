package main

import (
	"time"
)

//定义一个区块
type Block struct {
	//当前时间戳
	Timestamp int64
	//交易信息
	Data []byte
	//前一个区块的Hash值
	PrevBlockHash []byte
	//当前区块的Hash值
	Hash  []byte
	Nonce int
}

//创建区块
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash
	block.Nonce = nonce
	return block
}

//创建创世区块
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}
