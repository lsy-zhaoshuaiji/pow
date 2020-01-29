package main

import (
	"flag"
	"fmt"
	"os"
)

const Usage  = `
	AddBlock --data  	"add block to blockChain" 	example:./block AddBlock {DATA}
	PrintBlockChain     "print all blockChain data"
`
const AddBlockString  ="AddBlock"
const PrintBlockString  = "PrintBlockChain"
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
	AddBlockerParam:=AddBlocker.String("data","","AddBlock {data}")
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
	case PrintBlockString:
		//PrintBlockChain
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