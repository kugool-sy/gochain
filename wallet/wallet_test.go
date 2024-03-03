package wallet

import (
	// "bytes"
	"fmt"
	"testing"
)

func TestNewKeyPair(t *testing.T){
	fmt.Println("test ok")
	t.Log("hhdhdhdhdhahsdfa")
	a,b:=NewKeyPair()
	fmt.Print(a,b)

}

func TestSave(t *testing.T){
	newWallet := NewWallet()
	fmt.Println(string(newWallet.Address()))
	fmt.Printf("private key: %x\n", newWallet.PrivateKey)
	newWallet.Save()
}

func TestLoadWallet(t *testing.T) {
	wallet:=LoadWallet("19LuoGo4M4xvDcxzmCD3DTie5vjbzL2nvr")
	fmt.Printf("private key: %x\n", wallet.PrivateKey)
	fmt.Println(string(wallet.Address()))
}


// func TestSave2(t *testing.T){
// 	newWallet := NewWallet2()
// 	fmt.Println(string(newWallet.Address()))
// 	newWallet.Save()
// }


// func TestLoadWallet2(t *testing.T) {
// 	wallet2:=LoadWallet2("1M85M2Bp5GnFiC63nkpPbTiajhMccziGdw")
// 	fmt.Println(string(wallet2.Address()))
// }
