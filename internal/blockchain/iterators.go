package blockchain

import "github.com/dgraph-io/badger"

type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

// Iterator takes our BlockChain struct and returns it as a BlockCHainIterator struct
func (chain *BlockChain) Iterator() *BlockChainIterator {
	iterator := BlockChainIterator{chain.LastHash, chain.Database}

	return &iterator
}

// Next will iterate through the BlockChainIterator
func (iterator *BlockChainIterator) Next() (*Block, error) {
	var block *Block

	err := iterator.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iterator.CurrentHash)
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			block, err = Deserialize(val)
			return err
		})
	})
	if err != nil {
		return nil, err
	}

	iterator.CurrentHash = block.PrevHash

	return block, nil
}
