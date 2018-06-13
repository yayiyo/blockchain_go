package main

//定义区块链
type Blockchain struct {
	blocks []*Block
}

//添加区块
func (bc *Blockchain) AddBlock(data string) {
	//获取当前区块链的最后一个区块
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, newBlock)
}

//创建一个含有创世区块的区块链
func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}
