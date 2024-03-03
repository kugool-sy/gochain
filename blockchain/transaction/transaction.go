package transaction

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"gochain/constcoe"
	"gochain/utils"
)

type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

func (tx *Transaction) TxHash() []byte {
	var encoded bytes.Buffer
	var hash [32]byte

	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(tx)
	utils.Handle(err)

	hash = sha256.Sum256(encoded.Bytes())
	return hash[:]
}

func (tx *Transaction) SetID() {
	tx.ID = tx.TxHash()
}

func BaseTx(toAddress []byte) *Transaction {
	txInput := TxInput{[]byte{}, -1, []byte{}}
	txOutput := TxOutput{constcoe.InitCoin, toAddress}

	tx := Transaction{[]byte("This is base transaction!"), []TxInput{txInput}, []TxOutput{txOutput}}
	return &tx
}

func (tx *Transaction) IsBase() bool {
	return len(tx.Inputs) == 1 && tx.Inputs[0].OutIdx == -1
}
