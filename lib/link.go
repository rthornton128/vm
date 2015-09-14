package vm

import "fmt"

func toAddress(b []byte) int16 {
	fmt.Println("toaddr")
	return int16(b[0]<<8) | int16(b[1])
}

func toBytes(a uint16) (byte, byte) {
	return byte(a >> 8), byte(a & 0xff)
}
