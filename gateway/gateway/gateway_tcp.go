package gateway

import (
	"fmt"
	"io/ioutil"
	"net"
)

func TCPprocess(aim map[string]string, shost string, shost2 string) {
	// fmt.Println("服务器开始监听...")
	//1.tcp表示使用网络协议是tcp
	//2.0.0.0.0:8888表示在本地监听8888端口
	//shost = "192.168.1.7:8888"      //192.168.17.128
	lister, err := net.Listen("tcp", shost)
	if err != nil {
		fmt.Println("UPLINK:Binding Failed...err: ", err)
		return
	}
	fmt.Println("UPLINK:Starting binding at :", shost)

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
		//go print_aim(aim)
		go process(conn, aim, shost2)

	}
	//fmt.Printf("lister=%v\n",lister)
}

func process(con net.Conn, aim map[string]string, shost2 string) {
	//循环接收客户端发送的数据
	defer con.Close() //关闭con
	var buf []byte

	//创建一个新的切片
	//buf := make([]byte, 1000000)
	//con.Read(buf)
	//1.等待客户端通过con发送信息
	//2.如果客户端没有write[发送]，协程就会阻塞于此
	//fmt.Printf("服务器等待客户端 %s 发送信息\n", con.RemoteAddr().String())
	//n, err := con.Read(buf)

	//buf, n := tcprecv(con)
	buf, n := tcphandler(con)

	fmt.Println("Gateway的tcp收到了", n, "个字节")

	ioutil.WriteFile(`temp`, buf[9:n], 0666)
	//ioutil.WriteFile(`1`, buf[0:n], 0666)

	if n != 0 {

		bestMEC := selectBestMECServer(aim)
		bestMEC = "0.0.0.0"       //"192.168.2.240"

		if bestMEC != "" {
			//fmt.Println("the bestMEC is:", bestMEC)
			//发送任务
			//fmt.Println(string(data))
			//data := append(sourceIP, byteSum)
			//data = append(data, context...)

			context := buf[0:n]
			//tcpsend(bestMEC+":"+mecServerPort, context, n)
			tcpsend(shost2, context, n)
			//tcpsendfromfile(bestMEC+":"+mecServerPort)

		} else {
			fmt.Println("No bestMEC has find!")
		}
	} else {
		fmt.Println("something wrong happend while recving!")
	}
}


func print_aim(m map[string]string) {
	fmt.Println("当前的aim:", m)

}
