package main

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

//4. 引入区块链
//2. BlockChain结构重写
//
//使用数据库代替数组

type BlockChain struct {
	//定一个区块链数组
	//blocks []*Block
	db *bolt.DB

	tail []byte //存储最后一个区块的哈希
}

const blockChainDb = "blockChain.db"
const blockBucket = "blockBucket"

//5. 定义一个区块链
func NewBlockChain(address string) *BlockChain {
	//return &BlockChain{
	//	blocks: []*Block{genesisBlock},
	//}

	//最后一个区块的哈希， 从数据库中读出来的
	var lastHash []byte

	//1. 打开数据库
	db, err := bolt.Open(blockChainDb, 0600, nil)
	//defer db.Close()

	if err != nil {
		log.Panic("打开数据库失败！")
	}

	//将要操作数据库（改写）
	_ = db.Update(func(tx *bolt.Tx) error {
		//2. 找到抽屉bucket(如果没有，就创建）
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			//没有抽屉，我们需要创建
			bucket, err = tx.CreateBucket([]byte(blockBucket))
			if err != nil {
				log.Panic("创建bucket(b1)失败")
			}

			//创建一个创世块，并作为第一个区块添加到区块链中
			genesisBlock := GenesisBlock(address)
			//fmt.Printf("genesisBlock :%s\n", genesisBlock)

			//3. 写数据
			//hash作为key， block的字节流作为value，尚未实现
			_ = bucket.Put(genesisBlock.Hash, genesisBlock.Serialize())
			_ = bucket.Put([]byte("LastHashKey"), genesisBlock.Hash)
			lastHash = genesisBlock.Hash

			////这是为了读数据测试，马上删掉,套路!
			//blockBytes := bucket.Get(genesisBlock.Hash)
			//block := Deserialize(blockBytes)
			//fmt.Printf("block info : %s\n", block)

		} else {
			lastHash = bucket.Get([]byte("LastHashKey"))
		}

		return nil
	})

	return &BlockChain{db, lastHash}
}

//定义一个创世块
func GenesisBlock(address string) *Block {
	coinbase := NewCoinBaseTX(address, "Go一期创世块，老牛逼了！")
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

//5. 添加区块
func (bc *BlockChain) AddBlock(txs []*Transaction) {
	for _,tx :=range txs{
		status:=bc.VerifyTranscation(tx)
		if !status{
			return
		}
	}
	//如何获取前区块的哈希呢？？
	db := bc.db         //区块链数据库
	lastHash := bc.tail //最后一个区块的哈希

	_ = db.Update(func(tx *bolt.Tx) error {

		//完成数据添加
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			log.Panic("bucket 不应该为空，请检查!")
		}

		//a. 创建新的区块
		block := NewBlock(txs, lastHash)

		//b. 添加到区块链db中
		//hash作为key， block的字节流作为value，尚未实现
		_ = bucket.Put(block.Hash, block.Serialize())
		_ = bucket.Put([]byte("LastHashKey"), block.Hash)

		//c. 更新一下内存中的区块链，指的是把最后的小尾巴tail更新一下
		bc.tail = block.Hash

		return nil
	})
}

//找到指定地址的所有的utxo
func (bc *BlockChain) FindUTXOs(publicHash []byte) []TXoutput {
	var UTXO []TXoutput

	txs := bc.FindUTXOTransactions(publicHash)

	for _, tx := range txs {
		for _, output := range tx.TXoutputs {
			if bytes.Equal(publicHash,output.PublickHash) {
				UTXO = append(UTXO, output)
			}
		}
	}

	return UTXO
}

//根据需求找到合理的utxo
func (bc *BlockChain) FindNeedUTXOs(senderPubKeyHash []byte, amount float64) (map[string][]uint64, float64) {
	//找到的合理的utxos集合
	utxos := make(map[string][]uint64)
	var calc float64
	txs := bc.FindUTXOTransactions(senderPubKeyHash)
	for _, tx := range txs {
		for i, output := range tx.TXoutputs {
			if bytes.Equal(senderPubKeyHash,output.PublickHash) {

				if calc < amount {
					utxos[string(tx.TXID)] = append(utxos[string(tx.TXID)], uint64(i))
					calc += output.Value

					//加完之后满足条件了，
					if calc >= amount {
						fmt.Printf("找到了满足的金额：%f\n", calc)
						return utxos, calc
					}
				} else {
					fmt.Printf("不满足转账金额,当前总额：%f， 目标金额: %f\n", calc, amount)
				}
			}
		}
	}

	return utxos, calc
}

func (bc *BlockChain) FindUTXOTransactions(senderPubKeyHash []byte) []*Transaction {
	var txs []*Transaction //存储所有包含utxo交易集合
	//我们定义一个map来保存消费过的output，key是这个output的交易id，value是这个交易中索引的数组
	//map[交易id][]int64
	spentOutputs := make(map[string][]int64)

	//创建迭代器
	it := bc.NewIterator()

	for {
		//1.遍历区块
		block := it.Next()

		//2. 遍历交易
		for _, tx := range block.Transactions {
			//fmt.Printf("current txid : %x\n", tx.TXID)

		OUTPUT:
			//3. 遍历output，找到和自己相关的utxo(在添加output之前检查一下是否已经消耗过)
			//	i : 0, 1, 2, 3
			for i, output := range tx.TXoutputs {
				if spentOutputs[string(tx.TXID)] != nil {
					for _, j := range spentOutputs[string(tx.TXID)] {
						//[]int64{0, 1} , j : 0, 1
						if int64(i) == j {
							//fmt.Printf("111111")
							//当前准备添加output已经消耗过了，不要再加了
							continue OUTPUT
						}
					}
				}
				//这个output和我们目标的地址相同，满足条件，加到返回UTXO数组中
				if  bytes.Equal(output.PublickHash,senderPubKeyHash) {
					//fmt.Printf("222222")
					//UTXO = append(UTXO, output)

					//!!!!!重点
					//返回所有包含我的outx的交易的集合
					txs = append(txs, tx)

					//fmt.Printf("333333 : %f\n", UTXO[0].Value)
				} else {
					//fmt.Printf("333333")
				}
			}

			//如果当前交易是挖矿交易的话，那么不做遍历，直接跳过

			if !tx.IsCoinBaseTX(tx) {
				//4. 遍历input，找到自己花费过的utxo的集合(把自己消耗过的标示出来)
				for _, input := range tx.TXinputs {
					//判断一下当前这个input和目标（李四）是否一致，如果相同，说明这个是李四消耗过的output,就加进来
					if bytes.Equal(Newripe160Hash(input.PublicKey),senderPubKeyHash) {
						//spentOutputs := make(map[string][]int64)
						//indexArray := spentOutputs[string(input.TXid)]
						//indexArray = append(indexArray, input.Index)
						spentOutputs[string(input.TXid)] = append(spentOutputs[string(input.TXid)], input.Index)
						//map[2222] = []int64{0}
						//map[3333] = []int64{0, 1}
					}
				}
			} else {
				//fmt.Printf("这是coinbase，不做input遍历！")
			}
		}

		if len(block.PrevHash) == 0 {
			fmt.Printf("区块遍历完成退出!")
			break
		}
	}

	return txs
}
func (bc *BlockChain) SignTransaction(tx *Transaction, privateKey *ecdsa.PrivateKey){
	prevTXs := make(map[string]Transaction)
	for _,input :=range tx.TXinputs{
		Tid:=input.TXid
		tx,err:=bc.FindTransactionByTXid(Tid)
		if err!=nil{
			continue
		}
		prevTXs[string(input.TXid)] = tx
	}
	tx.Sign(privateKey, prevTXs)
}
func (bc *BlockChain) FindTransactionByTXid(TXid []byte)(Transaction,error){
	it:=bc.NewIterator()
	for {
		block:=it.Next()
		for _,tx :=range block.Transactions{
			if bytes.Equal(tx.TXID,TXid){
				return *tx,nil
			}
		}
		if len(block.PrevHash)==0{
			break
		}
	}
	return Transaction{},errors.New("not find ")
}
func (bc *BlockChain) VerifyTranscation(tx *Transaction)bool{
	if tx.IsCoinBaseTX(tx){
		return true
	}
	prevTXs := make(map[string]Transaction)
	for _,input :=range tx.TXinputs{
		Tid:=input.TXid
		tx,err:=bc.FindTransactionByTXid(Tid)
		if err!=nil{
			continue
		}
		prevTXs[string(input.TXid)] = tx
	}
	return tx.Verify(prevTXs)
}