package main

import (
	"fmt"
	"time"
)

func (this *Cli)PrintBlockChain(){
	Iterator:=this.Bc.NewIterator()
	HeightblockChain:=0
	for{
		HeightblockChain++
		block:=Iterator.Next()
		if len(block.PrevHash)==0{
			Iterator.Restore()
			break
		}

	}

	for{
		block:=Iterator.Next()
		fmt.Printf("=======当前区块高度:%d======\n",HeightblockChain)
		fmt.Printf("当前哈希：%x\n",block.Hash)
		fmt.Printf("上一级哈希：%x\n",block.PrevHash)
		//fmt.Printf("交易信息：%s\n",block.Transactions[0].TXinputs[0].Sig)
		fmt.Printf("梅克尔根：%x\n",block.MerkelRoot)
		fmt.Printf("时间戳：%s\n",time.Unix(int64(block.TimeStamp),0).Format("2006-1-2 15:04:05"))

		HeightblockChain--
		if len(block.PrevHash)==0{
			break
		}
	}
}
func (this *Cli)getBalance(address string){
	fmt.Println("began")
	utxos:=this.Bc.FindUTXOs(address)
	total:=0.0
	for _,utxo:=range utxos{
		total+=utxo.Value
	}
	fmt.Printf("/%s/的余额为:%f\n",address,total)
}
func (this *Cli)Send(from ,to string,amount float64,miner,data string){
	//1.创建挖矿交易,可以理解为挖矿交易就是为了争夺普通交易的上传权
	CoinBaseTranscation:=NewCoinBaseTX(miner,data)
	//2.创建普通交易
	//func NewTranscations(from, to string,amount float64,bc *BlockChain )*Transaction{
	tx:=NewTranscations(from,to,amount,this.Bc)
	if tx==nil{
		fmt.Println("余额不足，请检查钱包额余2")
	}
	//3.调用Addblock（可以理解为，将交易数据保存在数据库，因为Addblock中就是存放Transcations的过程）
	//func (bc *BlockChain) AddBlock(txs []*Transaction) {
	this.Bc.AddBlock([]*Transaction{CoinBaseTranscation,tx})

}
func (this *Cli)NewWallet(){
	wallet:=NewWallet()
	fmt.Printf("%v\n",wallet.PrivateKey)
	fmt.Printf("%v\n",wallet.PublicKey)
}