package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
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
	Sig string//解锁脚本
}
type TXoutput struct {
	Value float64//转账金额
	PubKeyHash string//锁定脚本
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
	input:=TXinput{[]byte{},-1,address}
	//3.在output中，金额为btc常量，reward{12.5}，锁定脚本为address
	output:=TXoutput{reward,address}
	//4.创建Transcation交易，并设置TXid
	tx:=Transaction{[]byte{},[]TXinput{input},[]TXoutput{output}}
	//通过SetHash方法创建交易ID
	tx.SetHash()
	return &tx
}
func NewTranscations(from, to string,amount float64,bc *BlockChain )*Transaction{
	utxos,resValue:=bc.FindNeedUtxos(from,amount)
	if resValue<amount{
		fmt.Println("余额不足，请检查钱包额余")
		return nil
	}else {
		//额余充足，进行转账和找零
		var inputs []TXinput
		var outputs []TXoutput
		for id,indexList :=range utxos{
			for _,index :=range indexList{
				inputs=append(inputs, TXinput{[]byte(id),int64(index),from})
			}
		}
		outputs=append(outputs,TXoutput{amount,to})
		//找零
		outputs=append(outputs,TXoutput{resValue-amount,from})
		tx:=Transaction{
			TXID:[]byte{},
			TXinputs:inputs,
			TXoutputs:outputs,
		}
		tx.SetHash()
		return &tx
	}
}
//4.根据交易调整程序