package main

/*
区块的创建分为7步
1.定义结构体
2.创建区块
3.生成哈希
4.定义区块链结构体
5.生成区块链并添加创世
6.生成创世块块
7.添加其他区块
*/

func main(){
	bc:=NewBlockChain()
	cli:=Cli{bc}
	cli.Run()
}
