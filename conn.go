package gotts

import (
    "github.com/gorilla/websocket"
    "sync/atomic"
)

/**
 * @Author zyq
 * @Date 2023/2/21 10:47 AM
 * @Description
 **/

type Conn struct {
    ID string
    *websocket.Conn
    done   chan struct{} // conn done
    closed int32         // 已结束
}
type ConnAzure struct {
    Conn // Exceeded maximum websocket connection duration(> 1200000ms)
}

type ConnBing struct {
    Conn
}

type IConn interface {
    Read(s *Server)
    Write(reqId, text, voice string) error
    Init(token string)
    Heartbeat(s *Server)
    GetID() string
}

type ConnModule string

const (
    ConnModuleAzure ConnModule = "azure"
    ConnModuleBing  ConnModule = "bing"
)

func (c *Conn) isClosed() bool {
    return atomic.LoadInt32(&c.closed) == 1
}

func (c *Conn) setClosed() {
    atomic.StoreInt32(&c.closed, 1)
}
func (c *Conn) GetID() string {
    return c.ID
}

func NewConn(module ConnModule) IConn {
    switch module {
    case ConnModuleAzure:
        return &ConnAzure{
            Conn: Conn{
                ID:   GetUUID(),
                done: make(chan struct{}),
            },
        }
    default:
        return &ConnBing{
            Conn: Conn{
                ID:   GetUUID(),
                done: make(chan struct{}),
            },
        }
    }
}

func GetConn(v interface{}, module ConnModule) IConn {
    switch module {
    case ConnModuleAzure:
        vv, _ := v.(*ConnAzure)
        return vv
    default:
        vv, _ := v.(*ConnBing)
        return vv
    }
}
