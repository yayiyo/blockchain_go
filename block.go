package main

import (
	"strconv"
	"bytes"
	"crypto/sha256"
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
	Hash []byte
}

//获取区块hash值
func (b *Block) SetHash() {
	//将时间戳转换成字节数组
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	//fmt.Println(timestamp)
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	//将区块信息进行hash
	hash := sha256.Sum256(headers)

	b.Hash = hash[:]
}

//创建区块
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}}

	block.SetHash()
	return block
}

//创建创世区块
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block",[]byte{})
}
