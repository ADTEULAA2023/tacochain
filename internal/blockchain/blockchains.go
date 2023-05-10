package blockchain

import (
	"fmt"

	"github.com/dgraph-io/badger"
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

func (chain *BlockChain) AddBlock(data string, difficulty uint) error {
	var lastHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
		if err != nil {
			return err
		}
		return err
	})
	if err != nil {
		return err
	}

	newBlock := CreateBlock(data, lastHash, difficulty)

	err = chain.Database.Update(func(transaction *badger.Txn) error {
		serialization, err := newBlock.Serialize()
		if err != nil {
			return err
		}

		err = transaction.Set(newBlock.Hash, serialization)
		if err != nil {
			return err
		}

		err = transaction.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash
		return err
	})
	if err != nil {
		return err
	}

	return nil
}

func Genesis(difficulty uint) *Block {
	return CreateBlock("Genesis", []byte{}, difficulty)
}

func InitBlockChain(dbPath string, difficulty uint) (*BlockChain, error) {
	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)
	opts.Truncate = true

	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	//part 1 finished

	err = db.Update(func(txn *badger.Txn) error {
		// "lh" stand for last hash
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			fmt.Println("No existing blockchain found")
			genesis := Genesis(difficulty)
			fmt.Println("Genesis proved")
			serialization, err := genesis.Serialize()
			if err != nil {
				return err
			}

			err = txn.Set(genesis.Hash, serialization)
			if err != nil {
				return err
			}
			lastHash = genesis.Hash
			return txn.Set([]byte("lh"), genesis.Hash)
			//part 2/3 finished
		} else {
			item, err := txn.Get([]byte("lh"))
			if err != nil {
				return err
			}
			err = item.Value(func(val []byte) error {
				lastHash = val
				return nil
			})
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	blockchain := BlockChain{lastHash, db}
	return &blockchain, nil
	//that's everything!
}
