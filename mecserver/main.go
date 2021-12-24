package main

import "mecserver1.7/mecserver"

//func fromIPtoMessage(IPNumber []byte) (writeDone byte, calDone int) {
//	for i, v := range IPNumber {
//		if i == 3 {
//			return v, 1
//		}
//	}
//}

func main() {
	//var a  = "192.168.1.11"
	//var aa byte = 1
	//b := []byte(a)
	//fmt.Println(b)
	////writeDone,_ := fromIPtoMessage(b)
	////fmt.Println(writeDone)
	//fmt.Println(b[len(b)-1:])
	//fmt.Println(reflect.TypeOf(aa))
	mecserver.RegularRgister()
}
