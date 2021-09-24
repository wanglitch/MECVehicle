package gateway

import (
	"fmt"
	"net"
	"os"
)

// func tcpsend(aim string,mess []byte]){
// 	conn,err := net.Dial("tcp",aim)
// 	defer conn.Close()
// 	if err != nil{
// 		fmt.Println("客户端 dial err：",err)
// 		return
// 	}
// 		_, err = conn.Write(mess)
// 		if err != nil{
// 			fmt.Println("conn Write err:",err)
// 		}
// }
//conn,err := net.Dial("tcp","127.0.0.1:8888")

// func tcpsend(remoteAddr string, mess []byte) {
// 	conn, err := net.Dial("tcp", remoteAddr)
// 	defer conn.Close()
// 	if err != nil {
// 		fmt.Println("客户端 dial err:", err)
// 		return
// 	}
// 	_, err = conn.Write(mess)
// 	if err != nil {
// 		fmt.Println("conn Write err:", err)
// 	}
// }

// tcp receive message
func tcphandler(con net.Conn) ([]byte, int) {
	var allbuf []byte
	var m int = 0
	for{
		//创建一个新的切片
		buf := make([]byte, 1024)
		n, err := con.Read(buf)
		allbuf = append(allbuf, buf[:n]...)
		m = m + n
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
		}
	}
	defer con.Close()
	return allbuf, m
}


func tcprecv(conn net.Conn) ([]byte, int) {
	f, _ := os.Create("temp")
	returndata := make([]byte, 1000000)
	buf := make([]byte, 4096)
	index := 0
	// var n int
	var length int = 0
	for {
		n, err := conn.Read(buf)
		f.Write(buf[:n])
		// fmt.Println("阶段性", n)
		if n == 0 {
			// fmt.Println("文件读取完毕")
			break
		}
		length = length + n
		if err != nil {
			fmt.Println("conn.Read err:", err)
			return nil, 0
		}
		copy(returndata[index*4096:index*4096+n], buf[0:n])
		index = index + 1
		// file.Write(buf[:n])
	}
	defer f.Close()
	return returndata, length
}

func tcpsend(remoteAddr string, data []byte, len int) {
	conn, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		fmt.Println("tcp dial to", remoteAddr, "failed")
	}
	defer conn.Close()

	//打开要传输的文件
	// file, err := os.Open(filepath)``
	// if err != nil {
	// 	fmt.Println("os.Open err", err)
	// 	return
	// }
	var i int = 0
	buf := make([]byte, 1024)  //4096
	segsum := len / 1024
	//循环读取文件内容，写入远程连接
	for i = 0; i < segsum; i++ {
		// n, err := file.Read(buf)
		// if err == io.EOF {
		// 	fmt.Println("文件读取完毕")
		// 	return
		// }
		copy(buf, data[i*1024:1024*(i+1)])
		_, err = conn.Write(buf[:1024])
		if err != nil {
			fmt.Println("conn.Write err:", err)
			return
		}
	}
	copy(buf, data[segsum*1024:len])
	_, err = conn.Write(buf[:len-segsum*1024])
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return
	}

}

//func tcpresend(remoteAddr string, data []byte) {
//	conn, err := net.Dial("tcp", remoteAddr)
//	if err != nil {
//		fmt.Println("tcp dial to", remoteAddr, "failed")
//	}
//
//
//	var count int64
//	for {
//		buf := make([]byte, 1024)
//		//读取文件内容
//		n, err := f.Read(buf)
//		conn.Write(buf[:n])
//		if err != nil {
//			if err.Error() == "EOF" {
//				break
//			}
//		}
//
//		count += int64(n)
//	}
//
//	defer conn.Close()
//	defer f.Close()
//	return
//}


func Tcpsendfromfile(remoteAddr string) {
	conn, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		fmt.Println("tcp dial to", remoteAddr, "failed")
	}
	//打开要传输的文件
	f, err := os.Open("/home/wan/goWorkPlace/src/gateway1.7/1")

	if err != nil {
		fmt.Println("os.Open err", err)
		return
	}
	var count int64
	for {
		buf := make([]byte, 1024)
		//读取文件内容
		n, err := f.Read(buf)
		conn.Write(buf[:n])
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
		}

		count += int64(n)
		}

	defer conn.Close()
	defer f.Close()
	return
}
