package helper

import (
    "encoding/json"
    "github.com/tyingzh/gotts/logger"
    "unsafe"
)

/**
 * @Author zyq
 * @Date 2023/2/23 11:26 AM
 * @Description
 **/

const (
    InitDateTime = "2006-01-02 15:04:05"
)

func SafeGo(f func()) {
    go func() {
        defer RecoverPanic()
        f()
    }()
}

// RecoverPanic 恢复panic
func RecoverPanic() {
    err := recover()
    if err != nil {
        logger.Sugar.Error(err)

        //buf := make([]byte, 2048)
        //n := runtime.Stack(buf, false)
        //logger.Sugar.Error(fmt.Sprintf("%s", buf[:n]))
    }
}

func JsonMarshal(v interface{}) string {
    bytes, err := json.Marshal(v)
    if err != nil {
        logger.Sugar.Error("json序列化：", err)
    }
    return Bytes2str(bytes)
}

func JsonUnmarshal(str string, v interface{}) {
    err := json.Unmarshal(Str2bytes(str), v)
    if err != nil {
        //logger.Sugar.Error("json反序列化：", err)
    }
    return
}

func JsonMarshalByte(v interface{}) []byte {
    bytes, err := json.Marshal(v)
    if err != nil {
        logger.Sugar.Error("json序列化：", err)
    }
    return bytes
}

func JsonUnmarshalByte(src []byte, v interface{}) {
    err := json.Unmarshal(src, v)
    if err != nil {
        //logger.Sugar.Error("json反序列化：", err)
    }
    return
}

func Str2bytes(s string) []byte {
    x := (*[2]uintptr)(unsafe.Pointer(&s))
    h := [3]uintptr{x[0], x[1], x[1]}
    return *(*[]byte)(unsafe.Pointer(&h))
}

func Bytes2str(b []byte) string {
    return *(*string)(unsafe.Pointer(&b))
}
