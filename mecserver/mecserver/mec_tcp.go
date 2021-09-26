package mecserver

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
)

func TCPprocess(listenHost string, sendHost string) {
	fmt.Println("服务器开始监听...")
	//1.tcp表示使用网络协议是tcp
	//2.0.0.0.0:8888表示在本地监听8888端口
	//lister, err := net.Listen("tcp", "0.0.0.0:8889")
	lister, err := net.Listen("tcp", listenHost)
	if err != nil {
		fmt.Println("监听失败...err: ", err)
		return
	}
	defer lister.Close() //延时关闭listen

	//循环等待客户端连接
	for {
		//等待客户端连接
		//fmt.Println("等待客户端连接")
		conn, err := lister.Accept()
		if err != nil {
			fmt.Printf("连接Accept() 失败，err: ", err)
		} else {
			fmt.Printf("Accept() suc conn=%v,客户端IP=%v\n", conn, conn.RemoteAddr().String())
		}
		go Process(conn, sendHost)
	}
	//fmt.Printf("lister=%v\n",lister)
}

func Process(con net.Conn, sendHost string) {
	//循环接收客户端发送的数据
	//defer con.Close() //关闭con
	//var context = []byte{}
	var sourceIP []byte

	//创建一个新的切片
	var buf []byte
	//buf := make([]byte, 1000000)
	//con.Read(buf)
	//1.等待客户端通过con发送信息
	//2.如果客户端没有write[发送]，协程就会阻塞于此
	//fmt.Printf("服务器等待客户端 %s 发送信息\n", con.RemoteAddr().String())

	//buf, n := tcprecv(con)
	buf, n := tcpHandler(con)
	fmt.Println("MEC的tcp收到了", n, "个字节")

	//context := buf[0:n]
	//3.服务器显示客户端信息
	//fmt.Printf("收到了客户端（IP：%v）%d 个字节数据",con.RemoteAddr().String(),n)
	//fmt.Printf("收到了客户端 %s 数据:%s ", con.RemoteAddr().String(), string(buf[:n]))
	sourceIP = buf[0:4]
	context := buf[4:n]

	ioutil.WriteFile(`temp`, buf[9:n], 0666)

	taskdata := context
	calResult, err := calculateTask(taskdata, IPBytesToString(sourceIP))
	//calResult, err := calculateTask(buf)

	callength := len(calResult)
	data := make([]byte, callength+24)
	if err != nil {
		fmt.Println("err occured and the err is :", err)
	} else {
		//fmt.Println("the calculate result is:", string(calResult))
		//data = []byte(fmt.Sprintf("%x", calSHA1SUM([]byte(calResult))) + delimiter + seqNo + delimiter + "leng" + delimiter + string(calResult) + delimiter + sourceIP + "|:|:|")
		//fmt.Println(string(data))
		//data := append(sourceIP, calSHA1SUM(calResult)...)
		// data := calSHA1SUM(append(sourceIP, calResult...))
		//data = append(data, calResult...)
		// fmt.Println("the cal result is:", (calResult))
		// fmt.Println("the sourceIP is:", sourceIP)
		// fmt.Println("the calSUM:", calSHA1SUM(calResult))

		copy(data[0:4], sourceIP)
		copy(data[4:24], calSHA1SUM(calResult))
		copy(data[24:24+callength], calResult)
		fmt.Println("DOWNLINK:the sending data is:", data)
		//udpSend(strings.Split(con.RemoteAddr().String(), ":")[0]+":"+gwDownPort, data)
		//udpSend("192.168.1.7"+":"+gwDownPort, data)
		udpSend(sendHost, data)
	}
}

// tcp receive message
func tcpHandler(con net.Conn) ([]byte, int) {
	var allbuf []byte
	var m int = 0
	for {
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

//将BytesIP转为StringIP
func IPBytesToString(b []byte) string {
	var buf bytes.Buffer
	for i, v := range b {
		t := strconv.FormatInt(int64(v), 10)
		if i < 3 {
			buf.WriteString(t + ".")
		} else {
			buf.WriteString(t)
		}
	}
	return buf.String()
}

func tcprecv(conn net.Conn) ([]byte, int) {
	f, _ := os.Create("temp")
	returndata := make([]byte, 100000)
	buf := make([]byte, 4096)
	index := 0
	// var n int
	var length int = 0
	for {
		n, err := conn.Read(buf)
		f.Write(buf[:n])
		//fmt.Println("阶段性收到了", n, "个字节.")
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
	// file, err := os.Open(filepath)
	// if err != nil {
	// 	fmt.Println("os.Open err", err)
	// 	return
	// }
	var i int
	buf := make([]byte, 4096)
	segsum := len / 4096
	//循环读取文件内容，写入远程连接
	for i = 0; i < segsum; i++ {
		// n, err := file.Read(buf)
		// if err == io.EOF {
		// 	fmt.Println("文件读取完毕")
		// 	return
		// }
		copy(buf, data[i*4096:4096*(i+1)])
		_, err = conn.Write(buf[:4096])
		if err != nil {
			fmt.Println("conn.Write err:", err)
			return
		}
	}
	copy(buf, data[i*4096:len-segsum*4096])
}
