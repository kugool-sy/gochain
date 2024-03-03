package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"gochain/blockchain/transaction"
	"gochain/utils"
	"time"
)

type Block struct {
	Timestamp int64
	Hash      []byte
	PrevHash  []byte
	Target    []byte
	Nonce     int64
	// Data      []byte
	Transactions []*transaction.Transaction
}

func (b *Block) SetHash() {
	information := bytes.Join([][]byte{utils.ToHexInt(b.Timestamp), b.PrevHash, b.Target, utils.ToHexInt(b.Nonce), b.BackTransactionSummary()}, []byte{})
	hash := sha256.Sum256(information)
	b.Hash = hash[:]
}

func (block Block) ShowInfo() {
	fmt.Println("======================================")
	fmt.Printf("prev block hash: %x\n", block.PrevHash)
	fmt.Printf("block hash: %x\n", block.Hash)
	fmt.Printf("block time: %d\n", block.Timestamp)
	// fmt.Printf("block data: %s\n", block.Data)
	fmt.Printf("block target: %x\n", block.Target)
	fmt.Printf("block nonce: %d\n", block.Nonce)
	fmt.Println("======================================")

}

func CreateBlock(prevhash []byte, txs []*transaction.Transaction) *Block {
	block := Block{time.Now().Unix(), []byte{}, prevhash, []byte{}, 0, txs}
	block.Target = block.GetTarget()
	block.Nonce = block.FindNonce()
	block.SetHash()
	return &block
}

func GenesisBlock(owner []byte) *Block {
	// tx := transaction.BaseTx([]byte("kugool"))
	tx := transaction.BaseTx(owner)
	return CreateBlock([]byte("kugool is awesome"), []*transaction.Transaction{tx})
}

func (b *Block) BackTransactionSummary() []byte {
	txIDs := make([][]byte, 0)
	for _, txID := range b.Transactions {
		txIDs = append(txIDs, txID.ID)
	}

	return bytes.Join(txIDs, []byte{})
}

func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(b)
	utils.Handle(err)
	return res.Bytes()
}

func DeSerializeBlock(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	utils.Handle(err)
	return &block
}
