package main

import (
	"flag"
)

var myFlagSet = flag.NewFlagSet("test", flag.ExitOnError)
var stringFlag = myFlagSet.String("abc", "default value", "help mesage")

//func main() {
//	fmt.Println(os.Args[2:])
//	err:=myFlagSet.Parse(os.Args[2:])
//	if myFlagSet.Parsed(){
//		fmt.Println(err,"22222222222222")
//		fmt.Println(*stringFlag,"333333333333333")
//		args := myFlagSet.Args()
//		for i := range args {
//			fmt.Println(i, myFlagSet.Arg(i))
//		}
//	}
//
//}