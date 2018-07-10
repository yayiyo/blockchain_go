package main

import "crypto/sha256"

//定义Merkle结构体
type MerkleTree struct {
	RootNode *MerkleNode
}

//定义Merkle节点
type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}

//创建MerkleTree，传入的数据是多个未经Hash的交易
func NewMerkleTree(data [][]byte) *MerkleTree {
	var nodes []MerkleNode

	//判断交易个数是否是偶数，如果不是，将最后一个交易复制，添加到data中
	if len(data)%2 != 0 {
		data = append(data, data[len(data)-1])
	}

	//遍历所有交易，将交易Hash之后作为数据创建叶子节点
	for _, datum := range data {
		node := NewMerkleNode(nil, nil, datum)
		nodes = append(nodes, *node)
	}

	//如果又四个交易，生成一个根节点需要进行两轮生成树
	for i := 0; i < len(data)/2; i++ {
		var newLevel []MerkleNode

		for j := 0; j < len(nodes); j += 2 {
			node := NewMerkleNode(&nodes[j], &nodes[j+1], nil)
			newLevel = append(newLevel, *node)
		}

		nodes = newLevel
	}

	mTree := MerkleTree{&nodes[0]}

	return &mTree
}

//创建MerkleNode
func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	mNode := MerkleNode{}

	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		mNode.Data = hash[:]
	} else {
		prevHashes := append(left.Data, right.Data...)
		hash := sha256.Sum256(prevHashes)
		mNode.Data = hash[:]
	}

	mNode.Left = left
	mNode.Right = right

	return &mNode
}
