package mecserver

import (
	"fmt"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

func Createnodes() {
	var path = "/work1"
	var data = []byte("work1")
	// permission
	var acls = zk.WorldACL(zk.PermAll)
	//_, _, _, err := conn.ExistsW(path)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	// create
	p, err_create := conn.Create(path, data, 0, acls)
	if err_create != nil {
		fmt.Println(err_create)
		return
	}
	fmt.Println("created:", p)
}

func Createmarknodes(path string, data []byte) {
	//var path = "/1"
	//var data = []byte("1")
	// permission
	var acls = zk.WorldACL(zk.PermAll)
	_, _, _, err := conn.ExistsW(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	// create
	p, err_create := conn.Create(path, data, zk.FlagEphemeral, acls)
	if err_create != nil {
		fmt.Println(err_create)
		return
	}
	fmt.Println("created:", p)
}

func CreateChildnodes() {

	// try create child node
	var root_path = "/work1"
	//_, _, _, err := conn.ChildrenW(root_path)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	cur_time := time.Now().Format //.Unix() //"child"
	ch_path := fmt.Sprintf("%s/ch_%d", root_path, cur_time)
	var acls = zk.WorldACL(zk.PermAll)
	_, err := conn.Create(ch_path, []byte("1"), zk.FlagEphemeral, acls)
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Printf("ch_path: %s create\n", p)
}

func RegularRgister() {
	SetZkServers("localhost:2181", Testlisten) //vehicle.sidx.fun:52181
	//_, _, _, err := conn.ChildrenW("/work1")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}

	go taskServerInit("127.0.0.1:5215", "127.0.0.1:5214", 1)
	go taskServerInit("127.0.0.1:5225", "127.0.0.1:5224", 2)
	go TCPprocess("0.0.0.0:8889", "192.168.1.7:5211", 1)
	go TCPprocess("0.0.0.0:8899", "192.168.1.7:5221", 2)

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
