package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
)

// Proof of Work Algorithm for the block chain.
/*
Algorithm:
1. Take data from block.
2. Create counter(nonce) which starts at 0.
3. Create a hash of the data and the nonce.
4. Check if the hash meets certain requirements.
	Requirements:
		1. The first few characters of the hash should be 0.
		2. The hash should be less than a certain number of characters.
*/

const Difficulty = 20

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func NewProofOfWork(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))

	pow := &ProofOfWork{block, target}

	return pow
}

func IntToHex(n int64) []byte {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.BigEndian, n)
	if err != nil {
		panic(err)
	}

	return buffer.Bytes()
}

func (pow *ProofOfWork) InitData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.Data,
			IntToHex(pow.Block.Timestamp.Unix()),
			IntToHex(int64(nonce)),
			IntToHex(int64(Difficulty)),
		},
		[]byte{},
	)

	return data
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte

	nonce := 0

	for nonce < math.MaxInt64 {
		data := pow.InitData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)

		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.Target) == -1 {
			fmt.Printf("\r%x", hash)
			break
		} else {
			nonce++
		}
	}
	fmt.Println()

	return nonce, hash[:]
}

// Validation of Proof of Work.
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.InitData(pow.Block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(pow.Target) == -1
}
