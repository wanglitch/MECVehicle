package gateway

import (
	"fmt"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

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
	var hosts = []string{"localhost:2181"}                      //"vehicle.sidx.fun:52181"
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
