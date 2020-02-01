package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const Usage  = `
	print                   "遍历区块链           	        example:pow.exe print"
	getbalance --address    "打印当前账户余额 	        example:pow.exe getbalance --address "0x10086""
	send -from {string} -to {string} -amount {float64} -miner {string} --data {string}    "转账"
`
const PrintBlockString  = "print"
const GetBlanceString  = "getbalance"
const SendString  = "send"
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
	PrintBlockChainer:=flag.NewFlagSet(PrintBlockString,flag.ExitOnError)
	getBalancer:=flag.NewFlagSet(GetBlanceString,flag.ExitOnError)
	sender:=flag.NewFlagSet(SendString,flag.ExitOnError)
	getBalancerParam:=getBalancer.String("address","","打印余额")
	from:=sender.String("from","","发送人的地址")
	to:=sender.String("to","","接收人的地址")
	amount:=sender.Float64("amount",0.0,"转账金额")
	miner:=sender.String("miner","","矿工")
	data:=sender.String("data","","挖矿信息")
	switch os.Args[1] {
	case GetBlanceString:
		err:=getBalancer.Parse(os.Args[2:])
		if err!=nil{
			fmt.Println("输入错误，请加入参数-address 如pow.exe getbalance --address 007")
			PrintUsage()
			log.Panic(err)
		}
		if getBalancer.Parsed(){
			if *getBalancerParam==""{
				PrintUsage()
				os.Exit(1)
			}else {
				this.getBalance(*getBalancerParam)
			}
		}
	case SendString:
		err:=sender.Parse(os.Args[2:])
		if err!=nil{
			fmt.Println(err)
			PrintUsage()
		}
		if sender.Parsed(){
			if *from=="" && *to=="" && *amount ==0.0 && *miner=="" && *data==""{
				fmt.Println("参数有空 ，错误，请重新输入 例如：pow.exe send --from a --to b --amount 6.6 --miner d --data 666")
				PrintUsage()
			}else {
				this.Send(*from,*to,*amount,*miner,*data)
			}
		}
	case PrintBlockString:
		err:=PrintBlockChainer.Parse(os.Args[2:])
		if err!=nil{fmt.Println(err)}
		if PrintBlockChainer.Parsed(){
			this.PrintBlockChain()
		}
	default:
		fmt.Println("Invalid input ")
		PrintUsage()
	}
}