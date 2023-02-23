package gotts

import (
    "fmt"
    "github.com/gorilla/websocket"
    "github.com/tyingzh/gotts/helper"
    "github.com/tyingzh/gotts/logger"
    "log"
    "net/http"
    "strings"
    "time"
)

/**
 * @Author zyq
 * @Date 2023/2/23 8:57 AM
 * @Description
 **/

func (c *ConnAzure) Init(token string) {
    if token == "" {
        token = "TKB3oaPu1fUpQCBsuD-jEUHgZrTIZyIgBys9FoNe9sQ4giv5aeEmMuPmlpV4z1duB2jw7GuesvLmGUIRir"
    }
    //"eastus.api.speech.microsoft.com"
    u := fmt.Sprintf("wss://eastus.api.speech.microsoft.com/cognitiveservices/websocket/v1?TrafficType=AzureDemo&X-ConnectionId=%s", c.ID)

    //logger.Debugf("connecting to %s", u.String())
    header := http.Header{}
    header.Add("Origin", "https://azure.microsoft.com")

    var err error
    wsConn, _, err := websocket.DefaultDialer.Dial(u, header)
    if err != nil {
        log.Fatal("dial:", err)
    }
    c.Conn.Conn = wsConn
}

func (c *ConnAzure) Write(reqId, text, voice string) error {
    var payload1 = `{"context":{"system":{"name":"SpeechSDK","version":"1.12.1-rc.1","build":"JavaScript","lang":"JavaScript","os":{"platform":"Browser/Linux x86_64","name":"Mozilla/5.0 (X11; Linux x86_64; rv:78.0) Gecko/20100101 Firefox/78.0","version":"5.0 (X11)"}}}}`
    var message1 = "Path : speech.config\r\nX-RequestId: " + reqId + "\r\nX-Timestamp: " + GetXTime() + "\r\nContent-Type: application/json\r\n\r\n" + payload1

    err := c.WriteMessage(websocket.TextMessage, []byte(message1))
    if err != nil {
        return err
    }

    var payload2 = `{"synthesis":{"audio":{"metadataOptions":{"sentenceBoundaryEnabled":false,"wordBoundaryEnabled":false},"outputFormat":"audio-16khz-32kbitrate-mono-mp3"}}}`
    var message2 = "Path : synthesis.context\r\nX-RequestId: " + reqId + "\r\nX-Timestamp: " + GetXTime() + "\r\nContent-Type: application/json\r\n\r\n" + payload2

    err = c.WriteMessage(websocket.TextMessage, []byte(message2))
    if err != nil {
        return err
    }

    if voice == "" {
        voice = XiaoxiaoNeural
    }

    //var payload_3 = '<speak xmlns="http://www.w3.org/2001/10/synthesis" xmlns:mstts="http://www.w3.org/2001/mstts" xmlns:emo="http://www.w3.org/2009/10/emotionml" version="1.0" xml:lang="en-US"><voice name="' + voice + '"><mstts:express-as style="General"><prosody rate="'+spd+'%" pitch="'+ptc+'%">'+ msg_content +'</prosody></mstts:express-as></voice></speak>'
    var payload3 = `<speak xmlns="http://www.w3.org/2001/10/synthesis" xmlns:mstts="http://www.w3.org/2001/mstts" xmlns:emo="http://www.w3.org/2009/10/emotionml" version="1.0" xml:lang="en-US">
    <voice name="` + voice + `">
        <prosody rate="0%" pitch="0%">` + text + `</prosody>
    </voice>
</speak>`
    // <voice>可批量</voice><voice>可批量</voice>

    var message3 = "Path: ssml\r\nX-RequestId: " + reqId + "\r\nX-Timestamp: " + GetXTime() + "\r\nContent-Type: application/ssml+xml\r\n\r\n" + payload3
    err = c.WriteMessage(websocket.TextMessage, []byte(message3))
    if err != nil {
        return err
    }
    return nil
}

func (c *ConnAzure) Read(s *Server) {
    defer close(c.done)

    const (
        endRespPath = "Path:turn.end"
        needle      = "Path:audio\r\n"
        requestId   = "X-RequestId:"
    )
    for {
        _, message, err := c.ReadMessage()
        if err != nil {
            logger.Debugf("%s close with read err:%+v", c.ID, err)
            return
        }
        msg := string(message)
        rIdx := strings.Index(msg, requestId)
        if rIdx < 0 {
            continue
        }
        rIdx += len(requestId)
        reIdx := strings.Index(msg, "\r\n")
        if reIdx < rIdx {
            continue
        }
        rId := msg[rIdx:reIdx]
        v, _ := s.rm.Load(rId)
        vv, _ := v.(string)
        if strings.Contains(msg, endRespPath) {
            u, _ := s.writer.Write(rId, []byte(vv))
            s.rm.Store(rId, u)
            s.cancelRequest(rId)
            continue
        }
        startIdx := strings.Index(msg, needle)
        if startIdx < 0 {
            continue
        }
        startIdx += len(needle)
        vv += msg[startIdx:]
        s.rm.Store(rId, vv)
    }
}

func (c *ConnAzure) Heartbeat(s *Server) {
    helper.SafeGo(func() {
        pingTicket := time.NewTicker(time.Second * 30)  // max=60s
        sendTicket := time.NewTicker(30 * time.Second)  // max=300s
        closeTicket := time.NewTicker(time.Minute * 15) // max=1200s
        // 阻塞
        for {
            select {
            case <-pingTicket.C: // 每30s ping一次
                helper.SafeGo(func() {
                    if err := c.WriteMessage(websocket.PingMessage, []byte{}); err != nil { // 超过1分钟不ping断链
                        //return
                    }
                })
            case <-sendTicket.C: // 每180s 执行一次
                helper.SafeGo(func() {
                    err := c.Write(GetUUID(), "叮咚", "")
                    if err != nil {
                        logger.Debugf("%s ping failed: %+v", c.ID, err)
                    }
                })
                logger.Debugf("%s test 30 seconds", c.ID)
            case <-closeTicket.C: // 每600s 执行一次
                s.Conn.Delete(c.ID)
                if c.isClosed() {
                    return
                }
                s.start <- struct{}{}
                logger.Debugf("%s c closing", c.ID)
                select {
                case <-time.After(time.Second * 5):
                    if !c.isClosed() {
                        c.setClosed()
                        _ = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
                        logger.Debugf("%s c closed", c.ID)
                    }
                    return
                }
            case <-c.done:
                s.Conn.Delete(c.ID)
                if !c.isClosed() {
                    c.setClosed()
                    s.start <- struct{}{}
                    logger.Debugf("%s c closed", c.ID)
                }
                return
            }
        }
    })
}
