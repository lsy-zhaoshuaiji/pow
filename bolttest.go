package main

import (
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
)

func test1(){
	db, err := bolt.Open("testDb.db", 0400, nil)//第二个参数为权限
	if err != nil {
		fmt.Println(" bolt Open err :", err)
		return
	}
	defer db.Close()
	//创建bucket
	//View参数为一个函数类型，是一个事务
	err = db.View(func(tx *bolt.Tx) error {
		//打开一个bucket
		b1 := tx.Bucket([]byte("bucket1"))
		//没有这个bucket
		if b1 == nil {
			return errors.New("bucket do not exist!")
		}
		v1 := b1.Get([]byte("key1"))
		v2 := b1.Get([]byte("key2"))
		v3 := b1.Get([]byte("key3"))

		fmt.Printf("v1:%s\n", string(v1))
		fmt.Printf("v2:%s\n", string(v2))
		fmt.Printf("v3:%s\n", string(v3))
		return nil
	})

	if err != nil {
		fmt.Printf("db.View err:", err)
	}

	return

}
