package dht

import (
	"math/big"
	"fmt"
	"encoding/hex"
)

func nextId(id string) string {

	nId := big.Int{}
	nId.SetString(id, base)

	y := big.Int{}
	two := big.NewInt(2)
	one := big.NewInt(1)
	mbig := big.NewInt(m)

	y.Add(&nId, one)
	// 2^m
	two.Exp(two, mbig, nil)
	y.Mod(&y, two)

	yBytes := y.Bytes()
	yHex := fmt.Sprintf("%x", yBytes)
	return yHex
}

func prevId(id string) string {

	nId := big.Int{}
	nId.SetString(id, base)

	y := big.Int{}
	two := big.NewInt(2)
	one := big.NewInt(1)
	mbig := big.NewInt(m)

	y.Sub(&nId, one)
	// 2^m
	two.Exp(two, mbig, nil)
	y.Mod(&y, two)

	yBytes := y.Bytes()
	yHex := fmt.Sprintf("%x", yBytes)
	return yHex
}

func hexStringToByteArr(hexId string) []byte {
	var hexbytes []byte
	hexbytes, _ = hex.DecodeString(hexId)
	return hexbytes
}
