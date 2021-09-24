package gateway

import (
	"fmt"
	"strings"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

var zkservers string
var conn *zk.Conn

type zklistener func()

func Testlisten() {
	fmt.Println("listen test")
}

//GetConf : 获取zk结点中的配置信息，比如超时时间等
func GetConf(conn *zk.Conn, path string) (string, string, error) {

	ret, _, err1 := conn.Get(path)
	if err1 != nil {
		fmt.Println("Get Err:", err1)
		return "", "", err1
	}
	return path, string(ret), nil
}

//GetBatchConf : 获取zk结点下，注册的机器的ip信息
func GetBatchConf(path string) (map[string]string, error) {
	var hosts = []string{"vehicle.sidx.fun:52181"}
	conn, _, err := zk.Connect(hosts, time.Second*5)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ret := make(map[string]string)
	keys, _, err1 := conn.Children(path)
	if err1 != nil {
		return ret, err1
	}
	for _, key := range keys {
		newPath := path + "/" + key
		val, _, err := conn.Get(newPath)
		if err != nil {
			continue
		}
		ret[key] = string(val)
	}
	return ret, nil
}

//GetAllConf : 获取所有结点信息 并存入map
func GetAllConf(path string) (map[string]string, error) {
	ret := make(map[string]string)
	var hosts = []string{"localhost:2181"} //"vehicle.sidx.fun:52181"
	conn, _, err := zk.Connect(hosts, time.Second*5, zk.WithLogInfo(false))
	if err != nil {
		panic(err)
	}

	keys, _, err := conn.Children(path)
	if err != nil {
		panic(err)
	}
	for _, key := range keys {
		_, conf, _ := GetConf(conn, "/"+key)
		ret[key] = conf
	}
	conn.Close()
	return ret, nil
}

// 首先调用此方法设置zk地址
// 初始化zk server
func SetZkServers(servers string, listeners ...zklistener) {
	if zkservers != servers {
		zkservers = servers
		initConn()
		for _, listener := range listeners {
			listener() //定义当设置zk server时候需要通知的回调函数，做一些操作
		}
	}
}

func initConn() {
	servers := strings.Split(zkservers, ",")
	//注意，第二个参数为创建之后session维持的时间，因为session消
	//失之后，才会将注册的ip地址从zk上摘除，所以不能太长，否则影响服务
	//正常功能，一般为1s.
	//option := zk.WithEventCallback(Callback)
	connZK, _, err := zk.Connect(servers, time.Second*5) //测试所以写的60s
	if err != nil {
		panic(err)
	}
	if conn != nil {
		conn.Close()
	}
	conn = connZK
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
