package service

import (
	"encoding/json"
	"fmt"
	"github.com/LonelySnail/monkey/agent"
	"github.com/LonelySnail/monkey/module"
	"github.com/LonelySnail/monkey/util"
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
		return BOOL, util.BoolToBytes(v2), nil
	case int32:
		return INT, util.Int32ToBytes(v2), nil
	case int64:
		return LONG, util.Int64ToBytes(v2), nil
	case float32:
		return FLOAT,util.Float32ToBytes(v2), nil
	case float64:
		return DOUBLE,util.Float64ToBytes(v2), nil
	case []byte:
		return BYTES, v2, nil
	case map[string]interface{}:
		bytes, err := util.MapToBytes(v2)
		if err != nil {
			return MAP, nil, err
		}
		return MAP, bytes, nil
	case map[string]string:
		bytes, err := util.MapToBytesString(v2)
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

func Bytes2Args(app module.IDefaultApp, argsType string, arg []byte) (interface{}, error) {
	switch argsType {
	case NULL:
		return nil, nil
	case STRING:
		return string(arg), nil
	case BOOL:
		return util.BytesToBool(arg), nil
	case INT:
		return util.BytesToInt32(arg), nil
	case LONG:
		return util.BytesToInt64(arg), nil
	case FLOAT:
		return util.BytesToFloat32(arg), nil
	case DOUBLE:
		return util.BytesToFloat64(arg), nil
	case BYTES:
		return arg, nil
	case MAP:
		mps, errs := util.BytesToMap(arg)
		if errs != nil {
			return nil, errs
		}
		return mps, nil
	case MAPSTR:
		mps, errs := util.BytesToMapString(arg)
		if errs != nil {
			return nil, errs
		}
		return mps, nil

	default:
		return agent.NewGateSession(app,arg)

	}
	return nil, fmt.Errorf("Bytes2Args [%s] not registered", argsType)
}