package blockchain

import "github.com/dgraph-io/badger"

type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func (chain *BlockChain) Iterator() *BlockChainIterator {
	iterator := BlockChainIterator{chain.LastHash, chain.Database}

	return &iterator
}

func (iterator *BlockChainIterator) Next() (*Block, error) {
	var block *Block

	err := iterator.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iterator.CurrentHash)
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			block, err = Deserialize(val)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
		return err
	})
	if err != nil {
		return nil, err
	}

	iterator.CurrentHash = block.PrevHash

	return block, nil
}
