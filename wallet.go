package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
)

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey []byte
}
func NewWallet()*Wallet{
	curve:=elliptic.P256()
	privateKey,err:=ecdsa.GenerateKey(curve,rand.Reader)
	if err!=nil{fmt.Println(err)}
	publicKeyOrign:=privateKey.PublicKey
	publicKey:=append(publicKeyOrign.X.Bytes(),publicKeyOrign.Y.Bytes()...)
	return &Wallet{privateKey,publicKey}
}