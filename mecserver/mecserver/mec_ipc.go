package mecserver

import (
	"fmt"

	"github.com/ghetzel/shmtool/shm"
)

type ipc struct {
	id        int
	segHandle *shm.Segment
	// ptr       unsafe.Pointer
	status int
}

var (
	ipc_busy        bool = false
	ipc_idle        int  = 1
	ipc_busying     int  = 2
	ipc_writing     byte = 0
	ipc_write_done  byte = 1
	ipc_write_done2 byte = 2
	ipc_cal_done    int  = 3
	ipc_cal_done2   int  = 4

	ipc_shmpool []ipc

	numOfShm int = 6
	capOfShm int = 1024000
)

// func CheckError(err error) {
// 	if err != nil {
// 		panic(err)
// 	}
// }

func createShm(numOfSum int, capacity int) ([]*shm.Segment, []int) {
	segPtrAttay := make([]*shm.Segment, numOfSum)
	shmID := make([]int, numOfSum)
	// ptr := make([]unsafe.Pointer, numOfShm)
	for ID := 0; ID < numOfSum; ID++ {
		segPtr, err := shm.Create(capacity)

		CheckError(err)
		// ptr[ID], err = segPtr.Attach()
		// fmt.Println("分别创建时的地址指针:", ptr)
		// CheckError(err)
		// err = segPtr.Detach(ptr[ID])
		// CheckError(err)
		shmID[ID] = segPtr.Id
		segPtrAttay[ID] = segPtr
		if err != nil {
			fmt.Print("分配内存出错", err)
		}
	}
	return segPtrAttay, shmID
	// print(n)
}

func initShm() {
	segPtrAttay, shmID := createShm(numOfShm, capOfShm)
	// fmt.Println("创建成功后的地址指针:", ptr)
	// ptrtemp, _ := segPtrAttay[0].Attach()
	// fmt.Println("再次attch时:", ptrtemp)
	shmpoor := make([]ipc, numOfShm)
	for i := 0; i < numOfShm; i++ {
		shmpoor[i].id = shmID[i]
		shmpoor[i].segHandle = segPtrAttay[i]
		shmpoor[i].status = ipc_idle
		// shmpoor[i].ptr = ptr[i]
	}
	//if threadNumber == 1 {
	//	ipc_shmpool = shmpoor
	//} else if threadNumber == 2 {
	//	ipc_shmpool2 = shmpoor
	//}
	ipc_shmpool = shmpoor
	// fmt.Println(ipc_shmpool)
	// ipc_shmpool[0].segHandle.Write([]byte{52, 52, 12, 14})
	// fmt.Println(getByteWithOffset(ipc_shmpool[0].segHandle, 0))
	// pt, _ := ipc_shmpool[0].segHandle.Attach()
	// fmt.Println(*(*byte)(ipc_shmpool[0].ptr))
}

func deleteShm() {
	for _, shmitem := range ipc_shmpool {
		shmitem.segHandle.Destroy()
	}
}

func getByteWithOffset(seghandle *shm.Segment, offset int) int {
	bytes, err := seghandle.ReadChunk((int64)(offset+1), (int64)(offset))
	CheckError(err)
	return int(bytes[0])
}

func findUnUsingSeg() int {
	//var i int
	var index int = -1
	for ipc_busy {
	}
	ipc_busy = true
	for i, item := range ipc_shmpool {
		if item.status == ipc_idle {
			index = i
			break
		}
	}
	// 设置Seg为忙的状态i
	ipc_shmpool[index].status = ipc_busying
	ipc_busy = false
	return index
}

func setSegFree(index int) {
	ipc_shmpool[index].status = ipc_idle

}

//func getSeg() int {
//	var index int = -1
//	for index == -1 {
//		index = findUnUsingSeg()
//	}
//	return index
//}

// func main() {
// 	initShm()
// 	fmt.Println(ipc_shmpool)
// 	defer deleteShm()
// 	//===================================================================
// 	ipc_shmpool[0].segHandle.Write([]byte{254, 253, 252, 251, 250, 249})
// 	ipc_shmpool[0].segHandle.Write([]byte{1, 2, 3, 4, 5, 6, 7})
// 	//上箭头（这样子写进去的数据是往后接着写而不覆盖============================
// 	//==================================================================
// 	ipc_shmpool[0].segHandle.Seek(0, 0)
// 	ipc_shmpool[0].segHandle.Write([]byte{8, 9, 0, 11, 12, 13})
// 	//通过调用seek函数可以重新指定写入的数据的position========================
// 	// for {
// 	// }
// 	// ptr, id := createShm(3, 1000000)
// 	// seg, err := shm.Open(65566)
// 	// if err != nil {
// 	// 	panic(err)
// 	// }
// 	// ptr, err := seg.Attach()
// 	// // ptr[0].Write([]byte{255, 111, 111, 111, 111, 111, 111, 111, 111})
// 	// // ptr[1].Write([]byte{222, 222, 222, 222, 222, 222, 222, 222, 222})
// 	// if err != nil {
// 	// 	panic(err)
// 	// }
// 	// // uptr := uintptr(ptr) + 0
// 	// // nptr_ptr := unsafe.Pointer(uptr)
// 	// fmt.Println(getByteWithOffset(ptr, 2))
// }
