package util

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/LonelySnail/monkey/module"
	"math"
	"reflect"
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

func ArgsTypeAnd2Bytes(arg interface{}) (string, []byte, error) {
	if arg == nil {
		return NULL, nil, nil
	}
	switch v2 := arg.(type) {
	case []uint8:
		return BYTES, v2, nil
	}
	switch v2 := arg.(type) {
	case nil:
		return NULL, nil, nil
	case string:
		return STRING, []byte(v2), nil
	case bool:
		return BOOL, BoolToBytes(v2), nil
	case int32:
		return INT, Int32ToBytes(v2), nil
	case int64:
		return LONG, Int64ToBytes(v2), nil
	case float32:
		return FLOAT,Float32ToBytes(v2), nil
	case float64:
		return DOUBLE,Float64ToBytes(v2), nil
	case []byte:
		return BYTES, v2, nil
	case map[string]interface{}:
		bytes, err := MapToBytes(v2)
		if err != nil {
			return MAP, nil, err
		}
		return MAP, bytes, nil
	case map[string]string:
		bytes, err := MapToBytesString(v2)
		if err != nil {
			return MAPSTR, nil, err
		}
		return MAPSTR, bytes, nil

	default:
		bytes,err := json.Marshal(arg)
		if err != nil {
			return Session, nil, err
		}
		return Session,bytes,err
	}
	return "", nil, fmt.Errorf("Args2Bytes [%s] not registered structure type", reflect.TypeOf(arg))
}

func Bytes2Args(app module.IDefaultApp, argsType string, args []byte) (interface{}, error) {
	switch argsType {
	case NULL:
		return nil, nil
	case STRING:
		return string(args), nil
	case BOOL:
		return BytesToBool(args), nil
	case INT:
		return BytesToInt32(args), nil
	case LONG:
		return BytesToInt64(args), nil
	case FLOAT:
		return BytesToFloat32(args), nil
	case DOUBLE:
		return BytesToFloat64(args), nil
	case BYTES:
		return args, nil
	case MAP:
		mps, errs := BytesToMap(args)
		if errs != nil {
			return nil, errs
		}
		return mps, nil
	case MAPSTR:
		mps, errs := BytesToMapString(args)
		if errs != nil {
			return nil, errs
		}
		return mps, nil

	default:

		//a := new(module.Message)
		//err := json.Unmarshal(args,a)
		//return a,err
	}
	return nil, fmt.Errorf("Bytes2Args [%s] not registered", argsType)
}

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