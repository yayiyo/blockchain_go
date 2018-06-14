package main

import (
	"log"
	"encoding/hex"
	"fmt"
	"bytes"
	"encoding/gob"
	"crypto/sha256"
)

//定义挖出创世区块的奖励
const subsidy = 10

//定义交易
type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

//判断是否是coinbase
func (tx Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

//设置交易的ID
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

//交易输入
type TXInput struct {
	Txid      []byte
	Vout      int
	ScriptSig string
}

//交易输出
type TXOutput struct {
	Value        int
	ScriptPubkey string
}

//锁定
func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}

//解锁
func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubkey == unlockingData
}

//创建一个coinbase交易，用于产生创世区块的币
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txin := TXInput{[]byte{}, -1, data}
	txout := TXOutput{subsidy, to}
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.SetID()

	return &tx
}

//创建交易
func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	acc, validOutputs := bc.FindSpendableOutputs(from, amount)

	if acc < amount {
		log.Panic("ERROR:Not enough funds")
	}

	//列出所有输入
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

	//列出所有输出
	outputs = append(outputs, TXOutput{amount, to})
	if acc > amount {
		outputs = append(outputs, TXOutput{acc - amount, from})
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetID()

	return &tx
}
