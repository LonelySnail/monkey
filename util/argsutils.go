package util

import (
	"encoding/binary"
	"encoding/json"
	"math"
)

var (
	NULL    = "null"    //nil   null
	BOOL    = "bool"    //bool
	INT     = "int"     //int
	LONG    = "long"    //long64
	FLOAT   = "float"   //float32
	DOUBLE  = "double"  //float64
	BYTES   = "bytes"   //[]byte
	STRING  = "string"  //string
	MAPSTR  = "mapstr"  //map[string]string{}
	MAP     = "map"     //map[string]interface{}
	Session = "session"
)

// BoolToBytes bool->bytes
func BoolToBytes(v bool) []byte {
	var buf = make([]byte, 1)
	if v {
		buf[0] = 1
	} else {
		buf[0] = 0
	}
	return buf
}

// BytesToBool bytes->bool
func BytesToBool(buf []byte) bool {
	var data bool = buf[0] != 0
	return data
}

// Int32ToBytes Int32ToBytes
func Int32ToBytes(i int32) []byte {
	var buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(i))
	return buf
}

// BytesToInt32 BytesToInt32
func BytesToInt32(buf []byte) int32 {
	return int32(binary.BigEndian.Uint32(buf))
}

// Int64ToBytes Int64ToBytes
func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

// BytesToInt64 BytesToInt64
func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}

// Float32ToBytes Float32ToBytes
func Float32ToBytes(float float32) []byte {
	bits := math.Float32bits(float)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)

	return bytes
}

// BytesToFloat32 BytesToFloat32
func BytesToFloat32(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)

	return math.Float32frombits(bits)
}

// Float64ToBytes Float64ToBytes
func Float64ToBytes(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)

	return bytes
}

// BytesToFloat64 BytesToFloat64
func BytesToFloat64(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)

	return math.Float64frombits(bits)
}

// MapToBytes MapToBytes
func MapToBytes(jmap map[string]interface{}) ([]byte, error) {
	bytes, err := json.Marshal(jmap)
	return bytes, err
}

// BytesToMap BytesToMap
func BytesToMap(bytes []byte) (map[string]interface{}, error) {
	v := make(map[string]interface{})
	err := json.Unmarshal(bytes, &v)

	return v, err
}

// MapToBytesString MapToBytesString
func MapToBytesString(jmap map[string]string) ([]byte, error) {
	bytes, err := json.Marshal(jmap)
	return bytes, err
}

// BytesToMapString BytesToMapString
func BytesToMapString(bytes []byte) (map[string]string, error) {
	v := make(map[string]string)
	err := json.Unmarshal(bytes, &v)

	return v, err
}