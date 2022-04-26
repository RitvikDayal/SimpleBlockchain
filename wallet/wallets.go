package wallet

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"io/ioutil"
	"os"
)

const wallets_file = "./tmp/wallets.dat"

type Wallets struct {
	Wallets map[string]*Wallet
}

func CreateWallets() (*Wallets, error) {
	var wallets Wallets
	wallets.Wallets = make(map[string]*Wallet)

	err := wallets.LoadWallets()
	return &wallets, err
}

func (ws *Wallets) SaveWallets() {
	var content bytes.Buffer

	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(wallets_file, content.Bytes(), 0644)
	if err != nil {
		panic(err)
	}
}

func (ws *Wallets) LoadWallets() error {
	if _, err := os.Stat(wallets_file); os.IsNotExist(err) {
		return err
	}

	var wallets Wallets

	fileContent, err := ioutil.ReadFile(wallets_file)
	if err != nil {
		panic(err)
	}

	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		panic(err)
	}

	ws.Wallets = wallets.Wallets

	return nil
}

func (ws *Wallets) GetWallet(address string) *Wallet {
	return ws.Wallets[address]
}

func (ws *Wallets) CreateWallet() string {
	wallet := NewWallet()
	// get address as string
	address := string(wallet.GetAddress())
	ws.Wallets[string(address)] = wallet
	return address
}

func (ws *Wallets) GetAddresses() []string {
	var addresses []string
	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}
	return addresses
}