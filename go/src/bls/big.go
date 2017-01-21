package bls

import "math/big"

func Decimal2Big(s string) (x big.Int) {
	x.SetString(s,10)
	return
}
func Hex2Big(s string) (x big.Int) {
	x.SetString(s,16)
	return
}
func Bytes2Big(b []byte) (x big.Int) {
	// big endian
	x.SetBytes(b)
	return
}
