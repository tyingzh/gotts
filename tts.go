package gotts

import (
    "context"
    "errors"
    "github.com/google/uuid"
    "github.com/gorilla/websocket"
    "github.com/tyingzh/gotts/helper"
    "github.com/tyingzh/gotts/logger"
    "log"
    "net/http"
    "os"
    "os/signal"
    "strings"
    "sync"
    "sync/atomic"
    "syscall"
    "time"
)

/**
 * @Author zyq
 * @Date 2023/2/20 11:30 AM
 * @Description
eg:
Init("127.0.0.1:3030, "token", nil")
 **/

type Server struct {
    Conn           sync.Map       // websocket conn Exceeded maximum websocket connection duration
    rm             sync.Map       // request map
    addr           string         // websocket addr
    token          string         // azure websocket token
    interrupt      chan os.Signal // service signal
    writer         Writer         // writer 方式
    voice          string
    success        int32
    failed         int32
    maxConnectTime int64         // 最大连接时间 15分钟
    exit           chan struct{} // 退出
    start          chan struct{}
    last           int64
    module         ConnModule
}

type Writer interface {
    Write(reqId string, body []byte) (string, error)
}

func Init(addr string, ops ...Options) error {
    if addr == "" {
        return errors.New("address is empty")
    }
    op := GetOption(ops...)

    if op.w == nil {
        op.w = new(WriterFile)
    }

    c := &Server{
        addr:      addr,
        token:     op.token,
        interrupt: make(chan os.Signal, 1),
        writer:    op.w,
        voice:     op.voice,
        exit:      make(chan struct{}),
        start:     make(chan struct{}),
        module:    op.module,
    }
    signal.Notify(c.interrupt, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
    c.Start()
    return nil
}

func (s *Server) Start() {
    http.HandleFunc("/send", s.send)
    http.HandleFunc("/stat", s.stat)

    helper.SafeGo(func() {
        s.start <- struct{}{}
        logger.Debug("start: ", s.addr)
        log.Fatal(http.ListenAndServe(s.addr, nil))
    })

    // 阻塞
    for {
        select {
        case <-s.start: // 启动tts连接
            helper.SafeGo(s.newConn)
        case <-s.interrupt:
            logger.Debug("interrupt")
            s.Conn.Range(func(key, value any) bool {
                conn, ok := value.(*Conn)
                if !ok {
                    s.Conn.Delete(key)
                    return true
                }
                s.Conn.Delete(key)
                if conn.isClosed() {
                    return true
                }
                // Cleanly close the connection by sending a close message and then
                // waiting (with timeout) for the server to close the connection.
                err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
                if err != nil {
                    logger.Debug("write close:", err)
                    return true
                }
                select {
                case <-conn.done:
                case <-time.After(time.Second):
                }
                return true
            })
            return
        case <-s.exit:
            logger.Debug("exit server with no connection")
            return
        }
    }
}

func (s *Server) newConn() {
    conn := NewConn(s.module)
    conn.Init(s.token)
    s.Conn.Store(conn.GetID(), conn)
    logger.Debugf("%s conn[%s] creating", s.module, conn.GetID())
    conn.Heartbeat(s)
    conn.Read(s)
}

func (s *Server) getConn() IConn {
    var c IConn
    s.Conn.Range(func(key, value any) bool {
        conn := GetConn(value, s.module)
        if conn == nil {
            s.Conn.Delete(key)
            return true
        }
        c = conn
        return true
    })
    return c
}

func (s *Server) send(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query()
    var (
        text  = query.Get("text")
        token = query.Get("token")
        voice = query.Get("voice")
    )
    if text == "" {
        Fail(w, "输入内容不能为空")
        return
    }
    if s.token != "" && token != s.token {
        Fail(w, "请输入正确的token")
        return
    }
    if voice == "" {
        voice = s.voice
    }
    conn := s.getConn()
    if conn == nil {
        Fail(w, "服务未启动")
        return
    }
    res, err := s.handleSend(conn, text, voice)
    if err != nil {
        Fail(w, err.Error())
        return
    }
    Success(w, res)
}

func (s *Server) handleSend(conn IConn, text, voice string) (string, error) {
    text = ReplaceText(text) // replace
    var reqId = GetUUID()
    s.rm.Store(reqId, "")

    if err := conn.Write(reqId, text, voice); err != nil {
        return "", err
    }

    atomic.StoreInt64(&s.last, time.Now().Unix())

    ctx, cancel := context.WithTimeout(context.WithValue(context.Background(), reqId, ""), time.Second*30)
    cancelKey := reqId + "_cancel"
    s.rm.Store(cancelKey, cancel)
    defer cancel()
    select {
    case <-ctx.Done():
        switch ctx.Err() {
        case context.DeadlineExceeded:
            s.rm.Delete(reqId)
            atomic.AddInt32(&s.failed, 1)
            return "", errors.New("timeout")
        default:
            s.rm.Delete(cancelKey)
            v, _ := s.rm.Load(reqId)
            vv, _ := v.(string)
            s.rm.Delete(reqId)
            atomic.AddInt32(&s.success, 1)
            return vv, nil
        }
    }
}

func (s *Server) cancelRequest(reqId string) {
    v, ok := s.rm.Load(reqId + "_cancel")
    if !ok {
        return
    }
    vv, ok := v.(context.CancelFunc)
    if !ok {
        return
    }
    vv()
}

func (s *Server) countConn() int {
    var n int
    s.Conn.Range(func(key, value any) bool {
        conn := GetConn(value, s.module)
        if conn == nil {
            s.Conn.Delete(key)
            return true
        }
        n++
        return false
    })
    return n
}

type StatResp struct {
    Success   int32  `json:"success"`
    Fail      int32  `json:"fail"`
    LastTime  string `json:"last_time"`
    ConnCount int    `json:"conn_count"`
}

func (s *Server) stat(w http.ResponseWriter, r *http.Request) {
    Success(w, StatResp{
        Success:   atomic.LoadInt32(&s.success),
        Fail:      atomic.LoadInt32(&s.failed),
        LastTime:  time.Unix(atomic.LoadInt64(&s.last), 0).Format(helper.InitDateTime),
        ConnCount: s.countConn(),
    })
}

func GetXTime() string {
    return time.Now().Format("2006-01-02T15:04:05.000Z")
}

func GetUUID() string {
    return strings.ToUpper(strings.Replace(uuid.NewString(), "-", "", -1))
}

func ReplaceText(text string) string {
    for k, v := range map[string]string{
        "覃": "秦",
    } {
        if strings.Contains(text, k) {
            text = strings.Replace(text, k, v, -1)
        }
    }
    return text
}
