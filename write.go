package gotts

import (
    "errors"
    "fmt"
    "github.com/tyingzh/gotts/qiniu"
    "os"
)

/**
 * @Author zyq
 * @Date 2023/2/20 5:04 PM
 * @Description
 **/

type WriterFile struct {
    Path string
}

func (w *WriterFile) Write(reqId string, body []byte) (string, error) {
    if w.Path == "" {
        w.Path = "/Users/tsying/Work/py/skytts/python_cli_demo/"
    }
    filename := fmt.Sprintf("%s%s.mp3", w.Path, reqId)
    if err := os.WriteFile(filename, body, os.ModePerm); err != nil {
        return "", err
    }
    return filename, nil
}

type WriterQiniu struct {
    Cfg  *qiniu.Config
    Path string
}

func (w *WriterQiniu) Write(reqId string, body []byte) (string, error) {
    if w.Cfg == nil {
        return "", errors.New("七牛未配置")
    }
    if w.Path == "" {
        w.Path = "tts/"
    }
    return w.Cfg.UploadBytes(body, fmt.Sprintf("%s%s.mp3", w.Path, reqId))
}
