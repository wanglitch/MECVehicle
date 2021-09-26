package mecserver

import (
	"fmt"
	"time"
)

func RegularRgister() {
	SetZkServers("localhost:2181", Testlisten) //vehicle.sidx.fun:52181
	//_, _, _, err := conn.ChildrenW("/work1")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}

	go taskServerInit("127.0.0.1:5215", "127.0.0.1:5214", 1)
	go taskServerInit("127.0.0.1:5225", "127.0.0.1:5224", 2)
	go TCPprocess("0.0.0.0:8889", "192.168.1.7:5211")

	//Createnodes()
	//CreateChildnodes()

	path := GetConnMACAddress("0.0.0.0")
	ip := Getip()
	Createmarknodes("/"+path, []byte(ip))

	for {
		//Register("/work1", "5555") //注册的路径和端口

		//nodesname, conf, _ := GetConf("/" + path)
		//fmt.Println("conf: ", nodesname[1:], ":", conf)

		//ZKOperateTest()
		allConf, _ := GetAllConf("/")
		fmt.Println("batchConf:", allConf)

		time.Sleep(time.Second * 10)
	}

}
