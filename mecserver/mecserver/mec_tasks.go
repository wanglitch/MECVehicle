package mecserver

import (
	"fmt"
	"net"
)

//Tasks : 定义MEC服务器可计算的任务大类
type Tasks struct {
	taskType int
}

func taskServerInit(listenHost string, sendHost string) {
	initShm()

	// 开启监听
	//listenHost := "127.0.0.1:5215"
	lister, err := net.ListenPacket("udp", listenHost)
	CheckError(err)
	defer lister.Close() //延时关闭listen
	// 唤醒计算服务器
	data := string(numOfShm) + "|:|:|"
	for i := 0; i < numOfShm; i++ {
		data = data + fmt.Sprintf("%d", ipc_shmpool[i].id) + "|:|:|"
	}
	fmt.Println("UPLINK: ", []byte(data))
	//udpSend("127.0.0.1:5214", []byte(data))
	udpSend(sendHost, []byte(data))
	//等待接受到确认信息
	dataByte := make([]byte, 1024)

	_, _, err = lister.ReadFrom(dataByte)
	if err != nil {
		fmt.Println("recive failed, err: ", err)
	}
	//监听服务器开启成功！！！！！

}

func taskEncode() {}
func taskDecode() {}

// 在shm中写入数据，之后读取计算结果
func calculateTask(taskdata []byte, IPNumber string) (output []byte, err error) {
	//数据处理
	//arg := taskdata
	//tasktype := taskdata[0]
	// 提取任务并写入数据
	taskparam := taskdata[0:] // 提取数据
	//taskparamLen := len(taskparam) //提取任务长度
	shmIndex := findUnUsingSeg()

	shm := ipc_shmpool[shmIndex].segHandle
	shm.Seek(0, 0)
	writeData := append([]byte{ipc_writing}, taskparam...)
	ipc_shmpool[shmIndex].segHandle.Write(writeData)
	shm.Seek(0, 0)
	if IPNumber == "192.168.1.9" {
		ipc_shmpool[shmIndex].segHandle.Write([]byte{ipc_write_done})
		for getByteWithOffset(ipc_shmpool[shmIndex].segHandle, 0) != ipc_cal_done {
			// fmt.Println(ipc_shmpool[shmIndex].ptr)
		}
	} else if IPNumber == "192.168.1.10" {
		ipc_shmpool[shmIndex].segHandle.Write([]byte{ipc_write_done2})
		for getByteWithOffset(ipc_shmpool[shmIndex].segHandle, 0) != ipc_cal_done2 {
			// fmt.Println(ipc_shmpool[shmIndex].ptr)
		}
	}
	// 提取任务计算结果
	readData := make([]byte, 100)
	shm.Seek(0, 0)
	_, err = shm.Read(readData)
	CheckError(err)
	// 释放shm
	setSegFree(shmIndex)
	// 处理计算结果
	resultlen := readData[1:3]
	resultlenInt := ByteToInt16(resultlen)
	resultData := readData[3 : 3+resultlenInt]
	output = resultData

	//// fmt.Println("接受到的数据的长度为：", len(taskparam)
	//// taskParamStr := string(taskparam)
	////fmt.Println(taskdata)
	////运行脚本=========================================================
	////if tasktype == IntToBytes(taskImageDetect)[0] {
	////command := "/home/magizio/Code/py/taskImageDetect.py"
	//// cmd := exec.Command(command)
	//// output, err = cmd.Output()
	//// if err != nil {
	//// 	fmt.Println("命令执行出错拉！", err)
	//// }
	//
	//// os.Remove("/home/magizio/Code/py/tmp.pic")
	////fmt.Println(string(output))
	////output = taskdata
	//// }
	//// err = nil
	////fmt.Println(string(output))
	//output = resultData
	return
}
