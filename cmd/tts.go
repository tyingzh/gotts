package main

import (
    "github.com/tyingzh/gotts"
    "github.com/tyingzh/gotts/logger"
)

/**
 * @Author zyq
 * @Date 2023/2/23 2:15 PM
 * @Description
 **/

func main() {
    logger.Debug(gotts.Init("127.0.0.1:3030",
        gotts.WithWriter(&gotts.WriterFile{Path: "/Users/tsying/Work/go/src/gotts/log/"}),
        gotts.WithVoice(gotts.YunyangNeural),
        gotts.WithToken("token"),
        gotts.WithModule(gotts.ConnModuleAzure)))
}
