package blockchain

import (
	"fmt"
	"gochain/utils"
)

func (bc *BlockChain) RunMine() {
	transactionPool := CreateTransactionPool()
	candidateBlock := CreateBlock(bc.LastHash, transactionPool.PubTx)

	if candidateBlock.ValidatePoW() {
		bc.AddBlock(candidateBlock)
		err := RemoveTransactionPool()
		utils.Handle(err)
		return
	} else {
		fmt.Println("Block has invalid nonce.")
		return
	}
}
