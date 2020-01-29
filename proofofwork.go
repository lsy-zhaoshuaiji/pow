package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

type ProofOfWork struct {
	block *Block
	target *big.Int
}
func NewProofOfWork(block *Block)*ProofOfWork{
	pow:=ProofOfWork{
		block:block,
	}
	//挖矿难度值
	diffculty:="0000100000000000000000000000000000000000000000000000000000000000"
	tmp:=big.Int{}
	tmp.SetString(diffculty,16)
	pow.target=&tmp
	pow.block.Difficulty=pow.target.Uint64()
	return &pow
}
func (this *ProofOfWork)run()([]byte,uint64){
	var nonce uint64
	var hash [32]byte
	for {
		//1.拼接区块数据
		byteList:=[][]byte{
			IntToByte(this.block.Version),
			this.block.PrevHash,
			this.block.MerkelRoot,
			IntToByte(this.block.TimeStamp),
			IntToByte(this.block.Difficulty),
			IntToByte(nonce),
			this.block.Data,
		}
		blockinfo:=bytes.Join(byteList,[]byte{})
		//2.将拼接后的数据进行哈希256运算
		hash=sha256.Sum256(blockinfo)
		//3.将哈希数据转为big.int类型
		tmp:=big.Int{}
		tmp.SetBytes(hash[:])
		//4.比较
		status:=tmp.Cmp(this.target)
		if status==-1{
			fmt.Printf("挖矿成功，难度为:%d,   当前哈希为：%x,  nonce为：%d\n",this.block.Difficulty,hash,nonce)
			break;
		}else {
			//fmt.Println(nonce)
			nonce++
		}
	}
	return hash[:],nonce
}
