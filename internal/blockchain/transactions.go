package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/ADTEULAA2023/tacochain/pkg"
)

type Transaction struct {
	ID      []byte
	Inputs  []TransactionInput
	Outputs []TransactionOutput
}

//TxOutput represents a transaction in the blockchain
//For Example, I sent you 5 coins. Value would == 5, and it would have my unique PubKey
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]

}

//CoinbaseTx is the function that will run when someone on a node succesfully "mines" a block. The reward inside as it were.
func CoinbaseTx(toAddress, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Coins to %s", toAddress)
	}
	//Since this is the "first" transaction of the block, it has no previous output to reference.
	//This means that we initialize it with no ID, and it's OutputIndex is -1
	txIn := TransactionInput{[]byte{}, -1, data}
	//txOut will represent the amount of tokens(reward) given to the person(toAddress) that executed CoinbaseTx
	txOut := TransactionOutput{nil, "", toAddress} // You can see it follows {value, PubKey}

	tx := Transaction{nil, []TransactionInput{txIn}, []TransactionOutput{txOut}}

	return &tx

}
func (tx *Transaction) IsCoinbase() bool {
	//This checks a transaction and will only return true if it is a newly minted "coin"
	// Aka a Coinbase transaction
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}

func NewTransaction(from, to string, data []byte, chain *BlockChain) (*Transaction, error) {
	var inputs []TransactionInput
	var outputs []TransactionOutput

	validOutputs, err := chain.FindSpendableOutputs(from, data)
	if err != nil {
		return nil, err
	}

	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			return nil, err
		}

		for _, out := range outs {
			input := TransactionInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	private, public, encryptedData, err := pkg.EncryptTransactionData(data)
	if err != nil {
		return nil, err
	}

	outputs = append(outputs, TransactionOutput{encryptedData, public, to})
	log.Println("==== private key:", private)
	log.Println("==== public key:", public)

	tx := Transaction{nil, inputs, outputs}
	tx.SetID()

	return &tx, nil
}
