# 基于UTXO实现区块链交易（第四版本的pow区块链开发）

完整教程地址:https://blog.csdn.net/Laughing_G/article/details/104054099

UTXO（Unspent Transaction Outputs）是未花费的交易输出，它是比特币交易生成及验证的一个核心概念。交易构成了一组链式结构，所有合法的比特币交易都可以追溯到前向一个或多个交易的输出，这些链条的源头都是挖矿奖励，末尾则是当前未花费的交易输出。

一、新建Trancastions.go文件，并定义结构体
[plain]
type Transaction struct {
	TXID []byte           //交易ID
	TXinputs []TXinput   //交易输入数组
	TXoutputs []TXoutput //交易输出数组
}
type TXinput struct {
	TXid []byte//引用的交易ID
	index int64//引用的output索引值
	Sig string//解锁脚本
}
type TXoutput struct {
	value float64//转账金额
	PubKeyHash string//锁定脚本
}
//设置交易id

二、序列化Transcation结构体，并设置TXID

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

三、提供创建Transcation的挖矿方法

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

四、修改block.go和blockChain.go（略，报错的文件都需要修改，可以下载源码等比）

五、修改cli.go新建getbalance命令

const AddBlockString  ="addblock"
const PrintBlockString  = "print"
const GetBlanceString  = "getbalance"
type Cli struct {
	Bc *BlockChain
}
func PrintUsage (){
	println(Usage)
}
func (cli *Cli)CheckInputLenth(){
	if len(os.Args)<2{
		fmt.Println("Invalid ARGS")
		PrintUsage()
		os.Exit(1)
	}
}
func (this *Cli)Run(){
	this.CheckInputLenth()
	AddBlocker:=flag.NewFlagSet(AddBlockString,flag.ExitOnError)
	PrintBlockChainer:=flag.NewFlagSet(PrintBlockString,flag.ExitOnError)
	getBalancer:=flag.NewFlagSet(GetBlanceString,flag.ExitOnError)
	AddBlockerParam:=AddBlocker.String("data","","AddBlock {data}")
	getBalancerParam:=getBalancer.String("address","","打印余额")
	switch os.Args[1] {
	case AddBlockString:
		//AddBlock
		err:=AddBlocker.Parse(os.Args[2:])
		if err!=nil{fmt.Println(err)}
		if AddBlocker.Parsed(){
			if *AddBlockerParam==""{PrintUsage()}else {
				this.AddBlock(*AddBlockerParam)
			}
		}
	case GetBlanceString:
		err:=getBalancer.Parse(os.Args[2:])
		if err!=nil{
			PrintUsage()
			log.Panic(err)
		}
		if getBalancer.Parsed(){
			if *getBalancerParam==""{fmt.Println(PrintUsage)}else {
				this.getBalance(*getBalancerParam)
			}
		}

六、修改command.go

func (this *Cli)getBalance(address string){
	utxos:=this.Bc.FindUTXOs(address)
	total:=0.0
	for _,utxo:=range utxos{
		total+=utxo.value
	}
	fmt.Printf("/%s/的余额为:%f\n",address,total)
}
七、在blockChain.go中新建FindCoinbaseUTXOs方法，获取UTXOs

func (this *BlockChain)FindUTXOs(address string)[]TXoutput{
	var UTXO []TXoutput
	spentOutputs := make(map[string][]int64)
	it := this.NewIterator()
	for {
		block := it.Next()
		for _, tx := range block.Transactions {
		OUTPUT:
			for i, output := range tx.TXoutputs {
				if spentOutputs[string(tx.TXID)] != nil {
					fmt.Println(spentOutputs[string(tx.TXID)])
					for _, j := range spentOutputs[string(tx.TXID)] {
						if int64(i) == j {
							continue OUTPUT
						}
					}
				}
				if output.PubKeyHash == address {
					UTXO=append(UTXO, output)
				}
			}
			if !tx.IsCoinBaseTX(tx) {
				for _, input := range tx.TXinputs {
					if input.Sig == address {
						spentOutputs[string(input.TXid)] = append(spentOutputs[string(input.TXid)], input.Index)
					}
				}
			}
		}

		if len(block.PrevHash) == 0 {
			fmt.Printf("区块遍历完成退出!")
			break
		}
	}

	return UTXO
	//
}

八、在Trancations.go中新增判断是否为挖矿交易的方法，如果是则跳过input记录前output

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


/*个人认为没有这个方法，也不会影响收益查看，因为余额虽然会在多个交易中查询UTXOS，但是挖矿交易始终是收益的源头，也就是第一条数据，而当遍历到第一条交易数据时，遍历output的循环不会在下一次执行了，所以大家可以发现就算没有添加此函数，余额也不会受影响，
但添加此函数有也添加的好处，比如会提交代码运行效率等

*/

九、创建普通交易

思路

[创建普通交易逻辑上要实现以下：（var utxos []output）

1.找到匹配的utxos，遍历余额，返回map[tanscation.TXid][outputIndex]，账户余额(resvalue)

2.判断余额与转账金额的大小，若余额小于转账金额，则返回nil，若余额大于转账金额则：

2.1新建TXinputs，记录output的ID和索引

2.2.新建txoutput，进行找零，output[amount,to]   output[resValuea - mount,from]，且返回*Transcation

实现：

1.在Transcation中实现普通交易的NewTranscations方法

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

 2.在blockChain中实现查找所需余额的FindNeedUTXOS方法

func (bc *BlockChain)FindNeedUtxos(from string,amount float64)(map[string][]uint64,float64){
	it:=bc.NewIterator()
	utxos:=make(map[string][]uint64)
	resValue:=float64(0.0)
	spentoutput :=make(map[string][]uint64)
	for {
		block:=it.Next()
		for _,transcation :=range block.Transactions{
			outputList:=[]uint64{}
			//output获取
		OUTPUT:
			for index,output :=range transcation.TXoutputs{
				if resValue>=amount{
					return utxos,resValue
				}
				if spentoutput[string(transcation.TXID)]!=nil{
					for _,value :=range spentoutput[string(transcation.TXID)]{
						if value==uint64(index){
							continue OUTPUT
						}
					}
				}
				if output.PubKeyHash==from{
					outputList=append(outputList, uint64(index))
					utxos[string(transcation.TXID)]=outputList
					resValue+=output.Value
					fmt.Printf("找到满足的金额：%f\n",output.Value)
				}
			}
			//input筛选
			if !transcation.IsCoinBaseTX(transcation){
				inputList:=[]uint64{}
				for _,input :=range transcation.TXinputs{
					if input.Sig==from{
						inputList=append(inputList, uint64(input.Index))
						spentoutput[string(input.TXid)]=inputList
					}
				}
			}
		}
		if len(block.PrevHash)==0{
			break
		}
	}
	fmt.Printf("转账结束\n")
	return utxos,resValue
}

至此，我们已经实现了UTXO的转账，特别提醒一点，一个交易中是不会存在两个同地址output（其中一个已经用过，另一个没有用过。）两个同地output，要么同时没有被用过，要么都被用过。因为区块链转账时实时，是全部转完的，即使自己有剩余，也会先拿出来，最后转给自己。所以，我们可以通过这一点，将两个utxo函数，高聚合化，但是为了节约时候，这里我就不再详细说明。

十、补充makerkleRoot生成函数，以及优化时间戳

func (this *Block)MakeMakerkleRoot()[]byte{
	//我们进行哈希拼接
	final:=[]byte{}
	for _,j :=range this.Transactions{
		final=append(final, j.TXID...)
	}
	hash:=sha256.Sum256(final)
	return hash[:]
}


//时间戳
fmt.Printf("时间戳：%s\n",time.Unix(int64(block.TimeStamp),0).Format("2006-1-2 15:04:05"))

