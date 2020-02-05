package main

import (
	"btcutil/base58"
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PubKey []byte
}
func NewWallet()*Wallet{
	curve:=elliptic.P256()
	privateKey,err:=ecdsa.GenerateKey(curve,rand.Reader)
	if err!=nil{fmt.Println(err)}
	publicKeyOrign:=privateKey.PublicKey
	public:=append(publicKeyOrign.X.Bytes(),publicKeyOrign.Y.Bytes()...)
	return &Wallet{privateKey,public}
}
func (this *Wallet)NewAdress()string{
	//1.获取pubKey
	pubKey:=this.PubKey
	//2.获取ripemd160哈希
	rip160HashValue:=Newripe160Hash(pubKey)
	//3.将ripeHash与version进行拼接
	version:=byte(00)
	payload:=append([]byte{version},rip160HashValue...)
	//4.拷贝一份payload做hash后截取前4个字节
	CheckCode:=CheckCode(payload)
	//5.再次拼接
	payload=append(payload,CheckCode...)
	//6.做base58
	address:=base58.Encode(payload)
	return address
}
func CheckCode(data []byte)[]byte{
	//两次hash
	hash1:=sha256.Sum256(data)
	hash2:=sha256.Sum256(hash1[:])
	return hash2[:4]
}
func Newripe160Hash(pubKey []byte)[]byte{
	hash:=sha256.Sum256(pubKey)
	//创建编码器
	ripe:=ripemd160.New()
	//写入
	_,err:=ripe.Write(hash[:])
	if err!=nil {
		fmt.Println(err)
	}
	//生成哈希
	rip160HashValue:=ripe.Sum(nil)
	return rip160HashValue

}
func IsValidAddress(address string)bool{
	payLoad:=base58.Decode(address)
	if len(payLoad)<4{
		fmt.Println("钱包地址长度太短")
		return false
	}
	Code1:=payLoad[len(payLoad)-4:]
	payLoad=payLoad[:len(payLoad)-4]
	Code2:=CheckCode(payLoad)
	return bytes.Equal(Code2,Code1)
}