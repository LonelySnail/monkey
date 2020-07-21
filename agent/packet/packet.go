package packet

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"unsafe"
)

const (
	CONNECT = iota
	PING
	DATA
	DISCONNECT
)
const maxLength = 65535

/*
 -----------------------------------------------
|												|
|byte	|  byte|2byte|body										|
 -----------------------------------------------

*/

var packetErr = errors.New("packet error")

type Pack struct {
	Type    byte
	length  uint16
	Payload []byte
}

func GetMagic() byte {
	return 0x80
}

func readInt(r *bufio.Reader) (int, error) {
	buf := make([]byte, 2)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return 0, err
	}
	return int(binary.BigEndian.Uint16(buf)), nil
}

func SliceByteToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func StringToSliceByte(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func Packet(typ byte, payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})

	if err := binary.Write(buf, binary.BigEndian, GetMagic()); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.BigEndian, typ); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, uint16(len(payload))); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, payload); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func UnPacket(r *bufio.Reader) (*Pack, error) {
	pack := new(Pack)
	magic, err := r.ReadByte()
	if err != nil {
		return pack, err
	}
	if magic != GetMagic() {
		return pack, packetErr
	}
	typ, err := r.ReadByte()
	if err != nil {
		return pack, err
	}
	pack.Type = typ
	n, err := readInt(r)
	if err != nil {
		return pack, err
	}
	if n > maxLength {
		return pack, packetErr
	}
	payload := make([]byte, n)
	_, err = io.ReadFull(r, payload)
	if err != nil {
		return pack, err
	}
	pack.Payload = payload

	return pack, nil
}
