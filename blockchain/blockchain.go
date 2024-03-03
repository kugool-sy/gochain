package blockchain

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"gochain/blockchain/transaction"
	"gochain/constcoe"
	"gochain/utils"
	"runtime"

	"github.com/dgraph-io/badger"
)

type BlockChain struct {
	LastHash []byte
	// Blocks []*Block
	Database *badger.DB
}

type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func InitBlockChain(address []byte) *BlockChain {
	var lastHash []byte

	if utils.FileExists(constcoe.BCFile) {
		fmt.Println("blockchain already exists")
		runtime.Goexit()
	}

	opts := badger.DefaultOptions(constcoe.BCPath)
	opts.Logger = nil

	db, err := badger.Open(opts)
	utils.Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		genesis := GenesisBlock(address)
		fmt.Println("Genesis Block Created")
		err = txn.Set(genesis.Hash, genesis.Serialize())
		utils.Handle(err)
		err = txn.Set([]byte("lh"), genesis.Hash)
		utils.Handle(err)
		err = txn.Set([]byte("ogprevhash"), genesis.PrevHash)
		utils.Handle(err)
		lastHash = genesis.Hash
		return err

	})

	utils.Handle(err)
	blockChain := BlockChain{lastHash, db}
	return &blockChain
}

func ContinueBlockChain() *BlockChain {
	if !utils.FileExists(constcoe.BCFile) {
		fmt.Println("No blockchain found, please create one first")
		runtime.Goexit()
	}

	var lastHash []byte
	opts := badger.DefaultOptions(constcoe.BCPath)
	opts.Logger = nil
	db, err := badger.Open((opts))
	utils.Handle(err)

	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		utils.Handle(err)
		item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
		utils.Handle(err)
		return err

	})
	utils.Handle(err)
	chain := BlockChain{lastHash, db}
	return &chain

}

func (bc *BlockChain) Iterator() *BlockChainIterator {
	iterator := BlockChainIterator{bc.LastHash, bc.Database}
	return &iterator
}

func (iterator *BlockChainIterator) Next() *Block {

	var block *Block

	err := iterator.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iterator.CurrentHash)
		utils.Handle(err)
		err = item.Value(func(val []byte) error {
			block = DeSerializeBlock(val)
			return nil
		})

		utils.Handle(err)
		return err
	})
	utils.Handle(err)
	iterator.CurrentHash = block.PrevHash
	return block
}

func (bc *BlockChain) BackOgPrevHash() []byte {
	var ogPrevHash []byte

	err := bc.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("ogprevhash"))
		utils.Handle(err)
		err = item.Value(func(val []byte) error {
			ogPrevHash = val
			return nil
		})
		utils.Handle(err)
		return err
	})

	utils.Handle(err)
	return ogPrevHash
}

func (bc *BlockChain) AddBlock(newBlock *Block) {
	var lastHash []byte

	err := bc.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		utils.Handle(err)
		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
		utils.Handle(err)
		return err
	})

	utils.Handle(err)

	if !bytes.Equal(newBlock.PrevHash, lastHash) {
		fmt.Println("This block is out of age")
		runtime.Goexit()
	}

	err = bc.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		utils.Handle(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)
		bc.LastHash = newBlock.Hash
		return err

	})

	utils.Handle(err)

}

func (bc *BlockChain) FindUnspentTransactions(address []byte) []transaction.Transaction {
	var unSpentTxs []transaction.Transaction
	spentTxs := make(map[string][]int) // can't use type []byte as key value

	// 倒序遍历所有区块
	iter := bc.Iterator()

all:
	for {
		block := iter.Next()

		// 遍历所有交易
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		IterOutputs:
			for outIdx, output := range tx.Outputs {

				// 检测是否是utxo
				if spentTxs[txID] != nil {
					for _, spentOut := range spentTxs[txID] {
						if outIdx == spentOut {
							continue IterOutputs
						}
					}
				}

				if output.ToAddressRight(address) {
					unSpentTxs = append(unSpentTxs, *tx)
				}
			}

			if !tx.IsBase() {
				for _, input := range tx.Inputs {
					if input.FromAddressRight(address) {
						inTxID := hex.EncodeToString(input.TxID)
						spentTxs[inTxID] = append(spentTxs[inTxID], input.OutIdx)
					}
				}
			}
			if bytes.Equal(block.PrevHash, bc.BackOgPrevHash()) {
				break all
			}
		}
	}

	return unSpentTxs

}

func (bc *BlockChain) FindUTXOs(address []byte) (int, map[string]int) {
	unspentOuts := make(map[string]int)
	unspentTxs := bc.FindUnspentTransactions(address)
	accumulated := 0

work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Outputs {
			if out.ToAddressRight(address) {
				unspentOuts[txID] = outIdx
				accumulated += out.Value

				continue work
			}
		}
	}

	return accumulated, unspentOuts
}

func (bc *BlockChain) FindSpendableOutputs(address []byte, amount int) (int, map[string]int) {
	unspentOuts := make(map[string]int)
	unspentTxs := bc.FindUnspentTransactions(address)
	accumulated := 0

work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)
		for outIdx, out := range tx.Outputs {
			if out.ToAddressRight(address) && accumulated < amount {
				accumulated += out.Value
				unspentOuts[txID] = outIdx

				if accumulated >= amount {
					break work
				}

				continue work
			}
		}
	}

	return accumulated, unspentOuts

}

func (bc *BlockChain) CreateTransaction(from, to []byte, amount int) (*transaction.Transaction, bool) {
	var inputs []transaction.TxInput
	var outputs []transaction.TxOutput

	acc, validOutput := bc.FindSpendableOutputs(from, amount)
	if acc < amount {
		fmt.Println("not enough coin")
		return &transaction.Transaction{}, false
	}

	for txid, outIdx := range validOutput {
		txID, err := hex.DecodeString(txid)
		utils.Handle(err)
		input := transaction.TxInput{TxID: txID, OutIdx: outIdx, FromAddress: from}
		inputs = append(inputs, input)
	}

	outputs = append(outputs, transaction.TxOutput{Value: amount, ToAddress: to})
	if acc > amount {
		outputs = append(outputs, transaction.TxOutput{Value: acc - amount, ToAddress: from})
	}
	tx := transaction.Transaction{ID: nil, Inputs: inputs, Outputs: outputs}
	tx.SetID()

	return &tx, true
}

// func (bc *BlockChain) Mine(txs []*transaction.Transaction) {
// 	bc.AddBlock(txs)
// }
