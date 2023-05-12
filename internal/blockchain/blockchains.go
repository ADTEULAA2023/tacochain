package blockchain

import (
	"encoding/hex"
	"fmt"
	"os"
	"runtime"

	"github.com/dgraph-io/badger"
)

const (
	dbPath      = "./tmp/blocks"
	dbFile      = "./tmp/blocks/MANIFEST"          // This can be used to verify that the blockchain exists
	genesisData = "First Transaction from Genesis" // This is arbitrary data for our genesis data
)

//BlockChain is an array of block pointers
type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

func DBexists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

//ContinueBlockChain will be called to append to an existing blockchain
func ContinueBlockChain(address string) (*BlockChain, error) {
	if !DBexists() {
		fmt.Println("No blockchain found, please create one first")
		runtime.Goexit()
	}

	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)
	opts.Logger = nil
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(txn *badger.Txn) error {
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
		return nil, err
	}

	chain := BlockChain{lastHash, db}
	return &chain, nil
}

//InitBlockChain will be what starts a new blockChain
func InitBlockChain(address string) (*BlockChain, error) {
	var lastHash []byte

	if DBexists() {
		fmt.Println("blockchain already exists")
		runtime.Goexit()
	}

	opts := badger.DefaultOptions(dbPath)
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(txn *badger.Txn) error {

		cbtx := CoinbaseTx(address, genesisData)
		genesis := Genesis(cbtx)
		fmt.Println("Genesis Created")
		serialize, err := genesis.Serialize()
		if err != nil {
			return err
		}

		err = txn.Set(genesis.Hash, serialize)
		if err != nil {
			return err
		}

		err = txn.Set([]byte("lh"), genesis.Hash)

		lastHash = genesis.Hash

		return err

	})

	if err != nil {
		return nil, err
	}

	blockchain := BlockChain{lastHash, db}
	return &blockchain, nil
}

// AddBlock Will add a Block type unit to a blockchain
func (chain *BlockChain) AddBlock(transactions []*Transaction) error {
	var lastHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
	})
	if err != nil {
		return err
	}

	newBlock := CreateBlock(transactions, lastHash)
	return chain.Database.Update(func(transaction *badger.Txn) error {
		serialized, err := newBlock.Serialize()
		if err != nil {
			return err
		}
		err = transaction.Set(newBlock.Hash, serialized)
		if err != nil {
			return err
		}

		err = transaction.Set([]byte("lh"), newBlock.Hash)
		chain.LastHash = newBlock.Hash
		return err
	})
}

func (chain *BlockChain) FindUnspentTransactions(address string) ([]Transaction, error) {
	var unspentTxs []Transaction

	spentTXOs := make(map[string][]int)

	iter := chain.Iterator()

	for {
		block, err := iter.Next()
		if err != nil {
			return nil, err
		}

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Outputs {
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				if out.CanBeUnlocked(address) {
					unspentTxs = append(unspentTxs, *tx)
				}
			}
			if !tx.IsCoinbase() {
				for _, in := range tx.Inputs {
					if in.CanUnlock(address) {
						inTxID := hex.EncodeToString(in.ID)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Out)
					}
				}
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}
	return unspentTxs, nil
}

func (chain *BlockChain) FindUTXO(address string) ([]TransactionOutput, error) {
	var UTXOs []TransactionOutput
	unspentTransactions, err := chain.FindUnspentTransactions(address)
	if err != nil {
		return nil, err
	}

	for _, tx := range unspentTransactions {
		for _, out := range tx.Outputs {
			if out.CanBeUnlocked(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs, nil
}

func (chain *BlockChain) FindSpendableOutputs(address string, data []byte) (map[string][]int, error) {
	unspentOuts := make(map[string][]int)
	unspentTxs, err := chain.FindUnspentTransactions(address)
	if err != nil {
		return nil, err
	}

	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)
		for outIdx, out := range tx.Outputs {
			if out.CanBeUnlocked(address) {
				unspentOuts[txID] = append(unspentOuts[txID], outIdx)
			}
		}
	}
	return unspentOuts, nil
}
