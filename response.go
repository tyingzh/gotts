package gotts

import (
    "github.com/tyingzh/gotts/helper"
    "net/http"
)

/**
 * @Author zyq
 * @Date 2023/2/20 5:35 PM
 * @Description
 **/

type Resp struct {
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}

func Fail(w http.ResponseWriter, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusBadRequest)
    _, _ = w.Write(helper.JsonMarshalByte(Resp{
        Message: message,
    }))
}

func Success(w http.ResponseWriter, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    _, _ = w.Write(helper.JsonMarshalByte(Resp{
        Message: "成功",
        Data:    data,
    }))
}
