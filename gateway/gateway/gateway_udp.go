package gateway

import (
	"fmt"
	"net"
)

//DownProcess ：下行接受处理函数
func DownProcess(listenHost string) {
	//listenHost := "0.0.0.0:" + gatewayDownPort
	lister, err := net.ListenPacket(udp, listenHost)
	if err != nil {
		fmt.Println("DOWNLINK:Binding Failed...err: ", err)
		return
	}
	fmt.Println("DOWNLINK:Starting binding at :", listenHost)

	defer lister.Close() //延时关闭listen
	for {
		data := make([]byte, 1024)
		//等待客户端连接
		//fmt.Println("waiting for Down comming packet...")
		//conn, err := lister.Accept()

		n, _, err := lister.ReadFrom(data)
		if err != nil {
			fmt.Println("recive failed，err: ", err)
		}
		//fmt.Println(n, remoteAddr, strings.Split(string(data), "|:|:|")[0])

		//go print_aim(aim)
		go udpDOWNProcess(data[:n], IPBytesToString(data[0:4])+":5212")
	}
}


func udpDOWNProcess(data []byte, sendHost string) {
	serverHash := data[4:24]
	context := data[24:]
	fmt.Println("DOWNLINK:" , data)
	//fmt.Println("the recv context is：", context)
	//验证报文可用性
	if packageValidCheck(context, serverHash) {
		//把报文广播出去
		conn, err := net.Dial("udp", sendHost)         //car 1.9 1.5  "192.168.1.9:5212"
		if err != nil {
			fmt.Println("发起广播报文时出错:", err)
		} else {
			_, err = conn.Write(data)
			if err != nil {
				fmt.Println("传输广播报文时出错")
			}
		}

	} else { //如果不可用
		fmt.Println("The UDP package was broken!")
	}

}


func udpSend(remoteAddr string, mess []byte) {
	conn, err := net.Dial("udp", remoteAddr)
	//defer conn.Close()

	if err == nil {
		conn.Write(mess)
	} else {
		fmt.Println("dial to", remoteAddr, "Failed")
	}
}
