package wallet

import (
	"bytes"
	// "crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/gob"
	"errors"
	"gochain/constcoe"
	"gochain/utils"
	"os"
)

type Wallet struct{
	PrivateKey []byte
	PublicKey []byte
}

func NewKeyPair()(ecdsa.PrivateKey, []byte){
	curve := elliptic.P256()
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	utils.Handle(err)
	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	return *privateKey, publicKey
}



func NewWallet() *Wallet{
	privateKey, publicKey := NewKeyPair()
	privateKeyBytes, err := x509.MarshalECPrivateKey(&privateKey)
	utils.Handle(err)
	wallet := Wallet{privateKeyBytes, publicKey}
	return &wallet
}

func (w *Wallet) EcdsaPrivateKey() ecdsa.PrivateKey{
	ecdsaPrivateKey, err := x509.ParseECPrivateKey(w.PrivateKey)
	utils.Handle(err)
	return *ecdsaPrivateKey
}


func (w *Wallet) Address() []byte {
	pubHash := utils.PublicKeyHash(w.PublicKey)
	return utils.PubHash2Address(pubHash)
}


func (w *Wallet) Save() {
	filename := constcoe.Wallets + string(w.Address()) + ".wlt"
	var content bytes.Buffer
	// gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(w)
	utils.Handle(err)
	err = os.WriteFile(filename, content.Bytes(), 0644)
	utils.Handle(err)
}

func LoadWallet(address string) *Wallet {
	filename := constcoe.Wallets + address + ".wlt"
	if !utils.FileExists(filename) {
		
		utils.Handle(errors.New("no wallet with such address"))
	}
	var w Wallet
	// gob.Register(elliptic.P256())
	fileContent, err := os.ReadFile(filename)
	utils.Handle(err)
	decoder := gob.NewDecoder(bytes.NewBuffer(fileContent))
	err = decoder.Decode(&w)
	utils.Handle(err)
	return &w
}