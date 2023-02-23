### go-tts

##### Text To Speech：文字转语音

### 实现
#### golang + websocket

### module
#### 已实现 Bing/Azure两种websocket方式

#### Bing

````
    logger.Debug(gotts.Init("127.0.0.1:3030",
        gotts.WithWriter(&gotts.WriterFile{Path: "/Users/tsying/Work/go/src/gozny/log/tts/"}),
        gotts.WithVoice(gotts.YunyangNeural),
        gotts.WithToken("token"),
        gotts.WithModule(gotts.ConnModuleBing)))
````

#### Azure

````
    logger.Debug(gotts.Init("127.0.0.1:3030",
        gotts.WithWriter(&gotts.WriterFile{Path: "/Users/tsying/Work/go/src/gozny/log/tts/"}),
        gotts.WithVoice(gotts.YunyangNeural),
        gotts.WithToken("token"),
        gotts.WithModule(gotts.ConnModuleAzure)))
````

### 语音输出方式
#### 可通过重定义接口 Writer方式实现多种语音输出方式 
