package main

import "fmt"

func (this *Cli)AddBlock(data string){
	this.Bc.AddBlock(data)
}
func (this *Cli)PrintBlockChain(){
	//TODO
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
		fmt.Printf("交易信息：%s\n",block.Data)
		HeightblockChain--
		if len(block.PrevHash)==0{
			break
		}
	}
}
