package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
)

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int
}

func (b *Block) DeriveHash() {
	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})
	hash := sha256.Sum256(info)
	b.Hash = hash[:]
}

func (b *Block) Serialize() ([]byte, error) {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)
	if err != nil {
		return nil, err
	}

	return res.Bytes(), nil
}

func Deserialize(data []byte) (*Block, error) {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)
	if err != nil {
		return nil, err
	}

	return &block, nil
}

func CreateBlock(data string, prevHash []byte, difficulty uint) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash, 0}
	// Don't forget to add the 0 at the end for the nonce!
	pow := NewProofOfWork(block, difficulty)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}
