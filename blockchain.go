package main

import (
	"github.com/bolt"
	"log"
	"fmt"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks"

//定义区块链
type Blockchain struct {
	//指向最后一个区块的hash
	tip []byte
	db  *bolt.DB
}

//区块链迭代器
type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci:=&BlockchainIterator{bc.tip,bc.db}
	return bci
}

//返回区块链中的下一个区块
func (i *BlockchainIterator) Next() *Block {
	var block *Block

	err:=i.db.View(func(tx *bolt.Tx) error {
		b:=tx.Bucket([]byte(blocksBucket))
		encodedBlock:=b.Get(i.currentHash)
		block=DeserializeBlock(encodedBlock)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	i.currentHash = block.PrevBlockHash

	return block
}

//添加区块
func (bc *Blockchain) AddBlock(data string) {
	//获取当前区块链的最后一个区块
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("1"))

		return nil
	})
	newBlock := NewBlock(data, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err = b.Put(newBlock.Hash, newBlock.Serialize())
		err = b.Put([]byte("1"), newBlock.Hash)
		bc.tip = newBlock.Hash

		return nil
	})
}

//创建一个含有创世区块的区块链
func NewBlockchain() *Blockchain {
	var tip []byte
	//打开或新建一个数据库
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	//打开一个读写事务
	err = db.Update(func(tx *bolt.Tx) error {
		//获取数据库中的“数据表”，存在则读取键为“1”对应的值（键“1”对应的是tip），
		//如果不存在，则创建创世区块，创建Bucket，将区块保存，然后更新1键
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {
			fmt.Println("No existing blockchain found. Creating a new one ...")
			genesis := NewGenesisBlock()

			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Panic(err)
			}

			//存储创世区块
			err = b.Put(genesis.Hash, genesis.Serialize())
			if err != nil {
				log.Panic(err)
			}

			//存储指向最后一个区块的hash，定义键为1
			err = b.Put([]byte("1"), genesis.Hash)
			if err != nil {
				log.Panic(err)
			}
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("1"))
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{tip, db}
	return &bc
}
