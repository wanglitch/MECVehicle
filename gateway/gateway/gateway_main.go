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
	fmt.Println("Starting the service...")
	//初始化
	// 更新时间
	updateDua := time.Second * 10
	// 查询路径
	pathRoot := "/"
	//status,aim :=upd(conn)
	mecNodes, _ := GetAllConf(pathRoot)
	time.Sleep(time.Second * 2)

	go TCPprocess(mecNodes, "192.168.1.7:8888" , "192.168.1.7:8889")   //shost:listen  shost2:send
	go TCPprocess(mecNodes, "192.168.1.7:8898" , "192.168.1.7:8899")
	go DownProcess("0.0.0.0:5211", "192.168.1.9:5212")
	go DownProcess("0.0.0.0:5221", "192.168.1.10:5222")

	//开启定时器更新服务器信息
	for true {
		//更新服务器信息 （sec Second per update）
		mecNodes, _ = GetAllConf(pathRoot)

		if Debug == true {
			fmt.Println("Running & MecNodes now:", mecNodes)
		}

		time.Sleep(updateDua)
	}
}


