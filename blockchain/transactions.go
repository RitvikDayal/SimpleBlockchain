package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
)

type Transaction struct {
	ID      []byte
	Inputs  []TXInput
	Outputs []TXOutput
}

type TXInput struct {
	Txid      []byte
	Vout      int
	Signature string
}

type TXOutput struct {
	Value      int
	PubKeyHash string
}

func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	HandleError(err)

	hash = sha256.Sum256(encoded.Bytes())

	tx.ID = hash[:]
}

func CoinbaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to %s", to)
	}

	txin := TXInput{[]byte{}, -1, data}
	txout := TXOutput{10, to}

	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.SetID()
	return &tx
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].Txid) == 0 && tx.Inputs[0].Vout == -1
}

func (in *TXInput) CanUnlockIn(data string) bool {
	return in.Signature == data
}

func (out *TXOutput) CanBeUnlockedWith(data string) bool {
	return out.PubKeyHash == data
}

func (bc *Blockchain) FindUTXO(address string) []TXOutput {
	var UTXOs []TXOutput

	UTXs := bc.UnspentTransactions(address)

	for _, tx := range UTXs {
		for _, out := range tx.Outputs {
			if out.CanBeUnlockedWith(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

func (bc *Blockchain) FindSOuts(address string, amount int) (int, map[string][]int) {
	var UOuts = make(map[string][]int)

	UTXs := bc.UnspentTransactions(address)
	accumulated := 0

Worker:
	for _, tx := range UTXs {
		txId := hex.EncodeToString(tx.ID)
		for outIdx, out := range tx.Outputs {
			if out.CanBeUnlockedWith(address) && accumulated < amount {
				accumulated += out.Value
				UOuts[txId] = append(UOuts[txId], outIdx)
			}
			if accumulated >= amount {
				break Worker
			}
		}
	}

	return accumulated, UOuts
}

func NewTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	accumulated, UOuts := bc.FindSOuts(from, amount)
	if accumulated < amount {
		panic("Not enough funds")
	}

	for txId, outs := range UOuts {
		txId, err := hex.DecodeString(txId)
		HandleError(err)
		for _, outIdx := range outs {
			input := TXInput{[]byte(txId), outIdx, from}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, TXOutput{amount, to})
	if accumulated > amount {
		outputs = append(outputs, TXOutput{accumulated - amount, from})
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetID()

	return &tx
}
