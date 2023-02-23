package gotts

import (
    "github.com/tyingzh/gotts/logger"
    "os"
    "testing"
)

/**
 * @Author zyq
 * @Date 2023/2/20 5:09 PM
 * @Description
 **/

func TestTTS(t *testing.T) {
    t.Run("StartBing", func(t *testing.T) {
        logger.Debug(Init("127.0.0.1:3030", WithWriter(&WriterFile{Path: os.Getenv("GOPATH") + "/src/gotts/log/"}), WithVoice(YunyangNeural), WithToken("token")))
    })
    t.Run("StartAzure", func(t *testing.T) {
        logger.Debug(Init("127.0.0.1:3030", WithWriter(&WriterFile{Path: os.Getenv("GOPATH") + "/src/gotts/log/"}), WithVoice(YunyangNeural), WithToken("token"), WithModule(ConnModuleAzure)))
    })
}
