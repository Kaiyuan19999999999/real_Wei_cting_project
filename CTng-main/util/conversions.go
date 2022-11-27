package util

import (
	"crypto"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"
	"math"
	"reflect"

	"github.com/Workiva/go-datastructures/bitarray"
)

func BitsToBytes(ba bitarray.BitArray) []byte {
	it := ba.Blocks()
	byteArr1 := make([]byte, 0)
	byteArr2 := make([]byte, 8)
	for it.Next() {
		_, currBlock := it.Value()
		binary.LittleEndian.PutUint64(byteArr2, uint64(currBlock))
		byteArr1 = append(byteArr1, byteArr2...)
	}
	return byteArr1
}

func BytesToBits(bytes []byte) bitarray.BitArray {
	ba := bitarray.NewBitArray(uint64(len(bytes)*8), false)
	for i, b := range bytes {
		for j := 0; j < 8; j++ {
			if int(b)&int(math.Pow(2, float64(j))) == int(math.Pow(2, float64(j))) {
				ba.SetBit(uint64(i*8 + j))
			}
		}
	}
	return ba
}

func PEM2PrivKey(s string) crypto.PrivateKey {
	p, _ := pem.Decode([]byte(s))
	if p == nil {
		panic("no PEM block found in " + s)
	}

	// Try various different private key formats one after another.
	if rsaPrivKey, err := x509.ParsePKCS1PrivateKey(p.Bytes); err == nil {
		return *rsaPrivKey
	}
	if pkcs8Key, err := x509.ParsePKCS8PrivateKey(p.Bytes); err == nil {
		if reflect.TypeOf(pkcs8Key).Kind() == reflect.Ptr {
			pkcs8Key = reflect.ValueOf(pkcs8Key).Elem().Interface()
		}
		return pkcs8Key
	}

	return nil
}

func PEM2PK(s string) crypto.PublicKey {
	p, _ := pem.Decode([]byte(s))
	if p == nil {
		panic("no PEM block found in " + s)
	}
	pubKey, _ := x509.ParsePKIXPublicKey(p.Bytes)
	if pubKey == nil {
		panic("public key not parsed from " + s)
	}
	return pubKey
}
