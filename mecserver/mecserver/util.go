package mecserver

import (
	"bytes"
	"encoding/binary"
)

//CheckError .
//检查错误并进行合适的处理
func CheckError(err error) int {
	if err != nil {
		panic(err)
	}
	return 0
}

//整形转换成字节
func IntToBytes(n int) []byte {
	x := int32(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.LittleEndian, x)
	return bytesBuffer.Bytes()
}
func ByteToInt8(b []byte) int8 {
	bytesBuffer := bytes.NewBuffer(b)

	var x int8
	binary.Read(bytesBuffer, binary.LittleEndian, &x)

	return int8(x)
}
func ByteToInt16(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int16
	binary.Read(bytesBuffer, binary.LittleEndian, &x)

	return int(x)
}

//字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	binary.Read(bytesBuffer, binary.LittleEndian, &x)

	return int(x)
}
