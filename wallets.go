package main

import (
	"btcutil/base58"
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
)

const walletFile  = "Wallets.dat"
type Wallets struct {
	WalletsMap map[string]*Wallet
}
func NewWallets()*Wallets{
	var walletsObj Wallets
	walletsObj.WalletsMap=make(map[string]*Wallet)
	walletsObj.LoadFile()
	return &walletsObj
}
func (ws *Wallets)CreateWallets()string{
	walletObj:=NewWallet()
	address:=walletObj.NewAdress()
	ws.WalletsMap[address]=walletObj
	ws.SaveTofile()
	return address
}
func (ws *Wallets)SaveTofile(){
	var buffer bytes.Buffer
	gob.Register(elliptic.P256())
	encoder:=gob.NewEncoder(&buffer)
	err:=encoder.Encode(ws)
	if err!=nil{fmt.Println(err,)}
	//func WriteFile(filename string, data []byte, perm os.FileMode) error {
	err=ioutil.WriteFile(walletFile,buffer.Bytes(),0600)
	if err!=nil{fmt.Println(err,"mmmmmmmmmmmm")}
}
func (ws *Wallets)LoadFile(){
	_, err := os.Stat(walletFile)
	if os.IsNotExist(err) {
		//ws.WalletsMap = make(map[string]*Wallet)
		return
	}
	//func ReadFile(filename string) ([]byte, error) {
	content,err:=ioutil.ReadFile(walletFile)
	if err!=nil{fmt.Println(err)}
	gob.Register(elliptic.P256())
	decoder:=gob.NewDecoder(bytes.NewReader(content))
	var walletsObj Wallets
	err=decoder.Decode(&walletsObj)
	if err!=nil{fmt.Println(err)}
	ws.WalletsMap=walletsObj.WalletsMap
}
func (ws *Wallets)PrintAdress()[]string{
	var addList []string
	for address,_ :=range ws.WalletsMap{
		addList=append(addList, address)
	}
	return addList
}
func GetPublicHash(address string)[]byte{
	payLoad:=base58.Decode(address)
	payLoad=payLoad[1:len(payLoad)-4]
	//hash:=Newripe160Hash(payLoad)
	return payLoad
	//比特币钱包生成publicKey(X,Y生成的)------>hash:=ripemd160Hash加密 ----->verion+hash------->hash2:=[hash256[verion+hash]][:4]+verion+hash---->base58
	//反解   base58  ---->ressultByte[1:len(resultByte)-4](publickHash)
}