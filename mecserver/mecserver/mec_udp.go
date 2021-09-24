package mecserver

import (
	"fmt"
	"net"
)

func udpSend(remoteAddr string, mess []byte) {
	conn, err := net.Dial("udp", remoteAddr)
	//defer conn.Close()

	if err == nil {
		conn.Write(mess)
	} else {
		fmt.Println("dial to", remoteAddr, "Failed")
	}
}
