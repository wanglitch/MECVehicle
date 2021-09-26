package mecserver

import (
	"encoding/json"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var zkservers string
var conn *zk.Conn

type zklistener func()

func Testlisten() {
	fmt.Println("listen test")
}

//获取本机的ip地址
func localIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
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

func Register(path, port string) {
	ip := localIP()
	if ip == "" {
		panic("get local ip error")
	}
	zkKey := ip + ":" + port
	value := make(map[string]interface{})
	value["host"] = ip
	p, _ := strconv.Atoi(port)
	value["port"] = p
	zkVal, _ := json.Marshal(value)
	go func() {
		for {
			var exist bool
			var err error
			registerPath := path + "/" + zkKey
			exist, _, err = conn.Exists(registerPath)
			if err != nil {
				fmt.Println("Exist Err:", err)
				time.Sleep(time.Second)
				continue
			}

			if !exist {
				//flag = zk.FlagEphemeral
				//flags有4种取值：
				//0:永久，除非手动删除
				//zk.FlagEphemeral = 1:短暂，session断开则改节点也被删除，session维持时间为zk.Connect的第二个参数
				//zk.FlagSequence  = 2:会自动在节点后面添加序号
				//3:Ephemeral和Sequence，即，短暂且自动添加序号
				//zk.WorldACL(zk.PermAll)//控制访问权限模式
				//var flag int32 =1
				_, err := conn.Create(registerPath, zkVal, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
				if err != nil {
					fmt.Println("Create Err:", err)
					time.Sleep(time.Second)
					continue
				}
			}
			time.Sleep(10 * time.Second)
		}
	}()
}

func Getip() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return (ipnet.IP.String())
			}
		}
	}
	return "未检测到ip地址"
}

func Callback(event zk.Event) {
	fmt.Println(">>>>>>>>>>>>>>>>>>>")
	fmt.Println("path:", event.Path)
	fmt.Println("type:", event.Type.String())
	fmt.Println("state:", event.State.String())
	fmt.Println("<<<<<<<<<<<<<<<<<<<")
}

//删除注册的节点
func Unregister(registerPath string) {
	if conn != nil {
		exist, _, err := conn.Exists(registerPath)
		if exist {
			err = conn.Delete(registerPath, 0)
			if err != nil {
				fmt.Println("Delete ZK Node:", err)
			}
		}
		//conn.Close()//如果节点是自动注册的，则使用该操作，直接Close掉conn即可。
	}
}

//获取zk结点中的配置信息，比如超时时间等
func GetConf(path string) (string, string, error) {
	ret, _, err1 := conn.Get(path)
	if err1 != nil {
		fmt.Println("Get Err:", err1)
		return "", "", err1
	}
	return path, string(ret), nil
}

//获取zk结点下，注册的机器的ip信息
func GetBatchConf(path string) (map[string]string, error) {
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

//获取所有结点信息 并存入map
func GetAllConf(path string) (map[string]string, error) {
	ret := make(map[string]string)
	keys, _, err1 := conn.Children(path)
	if err1 != nil {
		return ret, err1
	}
	for _, key := range keys {
		_, conf, _ := GetConf("/" + key)
		ret[key] = conf
	}
	return ret, nil
}

func ZKOperateTest() {
	fmt.Printf("ZKOperateTest\n")

	var path = "/"

	children, _, err := conn.Children(path)
	if err != nil {
		return
	}

	for _, p := range children {

		nodesname, conf, _ := GetConf("/" + p)
		fmt.Println("conf: ", nodesname, ":", conf)

	}
}

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
