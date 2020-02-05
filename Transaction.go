package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"math/big"
	"strings"
)

const reward  = 50
//1.定义交易结构
type Transaction struct {
	TXID []byte           //交易ID
	TXinputs []TXinput   //交易输入数组
	TXoutputs []TXoutput //交易输出数组
}
type TXinput struct {
	TXid []byte//引用的交易ID
	Index int64//引用的output索引值
	//Sig string//解锁脚本
	//签名,由r、s拼接成的hash
	Sigrnature []byte
	//公钥，由X Y拼接的公钥
	PublicKey  []byte
}
type TXoutput struct {
	Value float64//转账金额
	//PubKeyHash string//锁定脚本
	//公钥的哈希
	PublickHash []byte
}
func (TX *TXoutput)Lock(address string){
	//1.base58解码
	TX.PublickHash=GetPublicHash(address)
}
func NewTXoutput(value float64,address string)*TXoutput{
	TX:=TXoutput{
		Value:value,
	}
	TX.Lock(address)
	return &TX
}
//设置交易id
func (tx *Transaction)SetHash(){
	var buffer bytes.Buffer
	encoder:=gob.NewEncoder(&buffer)
	err:=encoder.Encode(tx)
	if err!=nil {
		fmt.Println(err)
	}
	hash:=sha256.Sum256(buffer.Bytes())
	tx.TXID=hash[:]
}
//2.提供创建交易方法（挖矿交易）
func (tx *Transaction)IsCoinBaseTX(txs *Transaction)bool{
	//1.TXid为空
	//2.index为-1
	//3.只有一个Input
	if len(txs.TXinputs) ==1{
		if len(txs.TXinputs[0].TXid)==0 && txs.TXinputs[0].Index==-1{
			return true
		}
	}
	return false
}
//3.创建挖矿交易
func NewCoinBaseTX(address string,data string) *Transaction{
	//1.挖矿交易只有一个Input和一个output
	//2.在input时，TXid为空，index为-1，解锁脚本为：矿池地址
	input:=TXinput{[]byte{},-1,nil,[]byte(data)}
	//3.在output中，金额为btc常量，reward{12.5}，锁定脚本为address
	output:=NewTXoutput(reward,address)
	//4.创建Transcation交易，并设置TXid
	tx:=Transaction{[]byte{},[]TXinput{input},[]TXoutput{*output}}
	//通过SetHash方法创建交易ID
	tx.SetHash()
	return &tx
}
func NewTranscations(from, to string,amount float64,bc *BlockChain )*Transaction{

	//1. 创建交易之后要进行数字签名->所以需要私钥->打开钱包"NewWallets()"
	ws := NewWallets()

	//2. 找到自己的钱包，根据地址返回自己的wallet
	wallet := ws.WalletsMap[from]
	if wallet == nil {
		fmt.Printf("没有找到该地址的钱包，交易创建失败!\n")
		return nil
	}

	//3. 得到对应的公钥，私钥
	pubKey := wallet.PubKey
	privateKey := wallet.PrivateKey //稍后再用

	//传递公钥的哈希，而不是传递地址
	pubKeyHash := Newripe160Hash(pubKey)

	//1. 找到最合理UTXO集合 map[string][]uint64
	utxos, resValue := bc.FindNeedUTXOs(pubKeyHash, amount)
	if resValue<amount{
		fmt.Println("余额不足，请检查钱包额余")
		return nil
	}else {
		//额余充足，进行转账和找零
		var inputs []TXinput
		var outputs []TXoutput
		for id,indexList :=range utxos{
			for _,index :=range indexList{
				inputs=append(inputs, TXinput{[]byte(id),int64(index),nil,pubKey})
			}
		}
		output:=NewTXoutput(amount,to)
		outputs=append(outputs,*output)
		//找零
		output1:=NewTXoutput(resValue-amount,from)
		outputs=append(outputs,*output1)
		tx:=Transaction{
			TXID:[]byte{},
			TXinputs:inputs,
			TXoutputs:outputs,
		}
		tx.SetHash()
		bc.SignTransaction(&tx,privateKey)
		return &tx
	}
}
func (tx *Transaction)Verify(prevTXs map[string]Transaction)bool{
	if !tx.IsCoinBaseTX(tx){
		return true
	}
	txCopy:=tx.TrimmedCopy()
	for index,input :=range tx.TXinputs{
		prevTranscation:=prevTXs[string(input.TXid)]
		if len(prevTranscation.TXID)==0{
			log.Panic("error......")
		}
		txCopy.TXinputs[index].PublicKey=prevTranscation.TXoutputs[input.Index].PublickHash
		txCopy.SetHash()
		txCopy.TXinputs[index].PublicKey=nil
		originSignature:=txCopy.TXID
		signature:=input.Sigrnature
		publicKey:=input.PublicKey
		r :=big.Int{}
		s :=big.Int{}
		r.SetBytes(signature[:len(signature)/2])
		s.SetBytes(signature[len(signature)/2:])
		x :=big.Int{}
		y :=big.Int{}
		x.SetBytes(publicKey[:len(publicKey)/2])
		y.SetBytes(publicKey[len(publicKey)/2:])
		originPublicKey:=ecdsa.PublicKey{elliptic.P256(),&x,&y}
		//func Verify(pub *PublicKey, hash []byte, r, s *big.Int) bool {
		if !ecdsa.Verify(&originPublicKey,originSignature,&r,&s){
			return false
		}
	}
	return true
}
func (tx *Transaction)Sign(privateKey *ecdsa.PrivateKey, prevTXs map[string]Transaction){
	if tx.IsCoinBaseTX(tx){
		return
	}

	txCopy:=tx.TrimmedCopy()
	for index,input :=range txCopy.TXinputs{
		prevTranscation:=prevTXs[string(input.TXid)]
		if len(prevTranscation.TXID)==0{
			log.Panic("交易错误，..........")
		}
		txCopy.TXinputs[index].PublicKey=prevTranscation.TXinputs[input.Index].PublicKey
		txCopy.SetHash()
		txCopy.TXinputs[index].PublicKey=nil
		r,s,err:=ecdsa.Sign(rand.Reader,privateKey,txCopy.TXID)
		if err!=nil{log.Panic(err)}
		signature:=append(r.Bytes(),s.Bytes()...)
		tx.TXinputs[index].Sigrnature=signature
	}
}
func (tx *Transaction)TrimmedCopy()Transaction{
	var inputs[]TXinput
	var outputs[]TXoutput
	for _,input :=range tx.TXinputs{
		inputs=append(inputs, TXinput{input.TXid,input.Index,nil,nil})
	}
	for _,output :=range tx.TXoutputs{
		outputs=append(outputs, output)
	}
	return Transaction{tx.TXID,inputs,outputs}
}
func (tx Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction %x:", tx.TXID))

	for i, input := range tx.TXinputs {

		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TXID:      %x", input.TXid))
		lines = append(lines, fmt.Sprintf("       Out:       %d", input.Index))
		lines = append(lines, fmt.Sprintf("       Signature: %x", input.Sigrnature))
		lines = append(lines, fmt.Sprintf("       PubKey:    %x", input.PublicKey))
	}

	for i, output := range tx.TXoutputs{
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:  %f", output.Value))
		lines = append(lines, fmt.Sprintf("       Script: %x", output.PublickHash))
	}

	return strings.Join(lines, "\n")
}
//