package gateway

import (
	"fmt"
	"time"
)

const Debug = true

//type gw struct {mecServerPort
//	zkConn *(zk.Conn)
//}

//vehicle.sidx.fun:52181

//Gateway ： gateway运行的主函数
func Gateway() {
	//初始化
	// SetZkServers("localhost:2181", Testlisten) //vehicle.sidx.fun:52181
	fmt.Println("Starting the service...")

	//status,aim :=upd(conn)
	mecNodes, _ := GetAllConf("/")

	go TCPprocess(mecNodes) //shost:"192.168.1.7:8888"listen  shost2: "192.168.1.7:8889"send
	go DownProcess("0.0.0.0:5211")

	//开启定时器更新服务器信息
	for true {
		//更新服务器信息 （sec Second per update）
		mecNodes, _ = GetAllConf("/")

		if Debug == true {
			fmt.Println("Running & MecNodes now:", mecNodes)
		}

		time.Sleep(time.Second * 10)
	}
}
