pow.exe send -from 张三 -to 李四 -amount 10 -miner 班长 -data "张三转李四10"
pow.exe send -from 张三 -to 王五 -amount 20 -miner 班长 -data "张三转王五20"

pow.exe send -from 王五 -to 李四 -amount 2 -miner 班长 -data "王五转李四2"
pow.exe send -from 王五 -to 李四 -amount 3 -miner 班长 -data "王五转李四3"
pow.exe send -from 王五 -to 张三 -amount 5 -miner 班长 -data "王五转张三5"


pow.exe send -from 李四 -to 赵六 -amount 14 -miner 班长 -data "李四转赵六14"