package transaction

import "bytes"

type TxOutput struct {
	Value     int
	ToAddress []byte
}

type TxInput struct {
	TxID        []byte
	OutIdx      int
	FromAddress []byte
}

func (txIn *TxInput) FromAddressRight(address []byte) bool {
	return bytes.Equal(txIn.FromAddress, address)
}

func (txOut *TxOutput) ToAddressRight(address []byte) bool {
	return bytes.Equal(txOut.ToAddress, address)
}
