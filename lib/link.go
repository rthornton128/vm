package vm

func toAddress(b []byte) uint16 {
	return uint16(b[0])<<8 | uint16(b[1])
}

func toBytes(a uint16) []byte {
	return []byte{byte(a >> 8), byte(a & 0xff)}
}
