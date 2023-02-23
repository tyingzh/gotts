package main

import (
    "github.com/tyingzh/gotts"
    "github.com/tyingzh/gotts/helper"
    "github.com/tyingzh/gotts/logger"
    "github.com/tyingzh/gotts/qiniu"
    "io/ioutil"
    "log"
    "os"
)

/**
 * @Author zyq
 * @Date 2023/2/23 2:15 PM
 * @Description
 **/

type Config struct {
    Token       string           `json:"token"`
    Module      gotts.ConnModule `json:"module"`
    QiniuKey    string           `json:"qiniu_key"`
    QiniuSecret string           `json:"qiniu_secret"`
    QiniuBucket string           `json:"qiniu_bucket"`
    QiniuHost   string           `json:"qiniu_host"`
    Address     string           `json:"address"`
    Path        string           `json:"path"`
}

func GetCfg() *Config {
    f, err := ioutil.ReadFile(os.Getenv("GOPATH") + "/src/gotts/config.json")
    if err != nil {
        log.Fatal(err)
    }
    var cfg = new(Config)
    helper.JsonUnmarshalByte(f, cfg)

    if cfg.Address == "" {
        cfg.Address = ":3030"
    }
    return cfg
}

func main() {
    cfg := GetCfg()
    var ops = []gotts.Options{
        gotts.WithVoice(gotts.XiaoxiaoNeural),
        gotts.WithToken(cfg.Token),
        gotts.WithModule(cfg.Module),
    }
    if cfg.Path != "" && cfg.QiniuHost == "" {
        ops = append(ops, gotts.WithWriter(&gotts.WriterFile{Path: cfg.Path}))
    }
    if cfg.QiniuHost != "" {
        ops = append(ops, gotts.WithWriter(&gotts.WriterQiniu{Path: cfg.Path, Cfg: qiniu.Init(cfg.QiniuKey, cfg.QiniuSecret, cfg.QiniuHost, cfg.QiniuBucket)}))
    }

    logger.Debug(gotts.Init(cfg.Address, ops...))
}
