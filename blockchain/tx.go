package blockchain

type TXInput struct {
	Txid      []byte
	Vout      int
	Signature string
}

type TXOutput struct {
	Value      int
	PubKeyHash string
}

func (in *TXInput) CanUnlockIn(data string) bool {
	return in.Signature == data
}

func (out *TXOutput) CanBeUnlockedWith(data string) bool {
	return out.PubKeyHash == data
}
