package gateway

import (
	"bytes"
	"crypto/sha1"
	"fmt"
)

func calSHA1SUM(context []byte) []byte {
	hashObj := sha1.New()
	hashObj.Write([]byte(context))
	clienthash := hashObj.Sum(nil)
	return clienthash
}

//serverHash|:|:|seqNo|:|:|packageLength|:|:|Context
func packageValidCheck(context []byte, serverHash []byte) bool {
	var result bool = false
	clienthash := calSHA1SUM(context)

	// if Debug == true {
	// 	fmt.Println("the context is :", string(context))
	// 	// fmt.Println("the clienthash is :", string(clienthash))
	// 	fmt.Println("tht serverhash is :", serverHashString)
	// 	fmt.Println("the clienthash is :", clienthashString)
	// 	//fmt.Printf("the clienthash is :%x\n", clienthash)
	// 	fmt.Println("the compare result is :", serverHashString == clienthashString)
	// }

	//fmt.Println(context)
	//fmt.Printf("%x\n", clienthash)

	// fmt.Println(context)
	//fmt.Println(clienthash)
	//fmt.Println(serverHash)
	if 0 == bytes.Compare(clienthash, serverHash) {
		result = true
	} else {
		result = false
	}
	//fmt.Println(context)
	//fmt.Printf("%x\n", clienthash)

	return result
}

//TestpackageValidCheck ：测试报文验证函数可用性
func TestpackageValidCheck(context []byte, serverHash []byte) {
	if packageValidCheck(context, serverHash) {
		fmt.Println("the Test Succeeded")
	} else {
		fmt.Println("the Test Failed")
	}
}
