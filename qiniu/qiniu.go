package qiniu

import (
    "bytes"
    "context"
    "fmt"
    "github.com/qiniu/api.v7/auth/qbox"
    "github.com/qiniu/api.v7/storage"
)

/**
 * @Author zyq
 * @Date 2023/2/23 4:28 PM
 * @Description
 **/

type Config struct {
    Key    string
    Secret string
    Host   string
    Bucket string
}

func Init(key, secret, host, bucket string) *Config {
    return &Config{
        Key:    key,
        Secret: secret,
        Host:   host,
        Bucket: bucket,
    }
}

func (c *Config) UploadBytes(data []byte, key string) (string, error) {
    upToken := c.GetToken(key)

    cfg := storage.Config{}
    formUploader := storage.NewFormUploader(&cfg)
    putExtra := storage.PutExtra{}
    ret := storage.PutRet{}
    if key == "" {
        key = c.Bucket + "/" + upToken
    }
    err := formUploader.Put(context.Background(), &ret, upToken, key, bytes.NewReader(data), int64(len(data)), &putExtra)
    if err != nil {
        return "", err
    }

    return c.Host + ret.Key, nil
}

func (c *Config) GetToken(key string) string {
    putPolicy := storage.PutPolicy{
        Scope: c.Bucket,
    }
    if key != "" {
        putPolicy.Scope = fmt.Sprintf("%s:%s", c.Bucket, key)
    }

    mac := qbox.NewMac(c.Key, c.Secret)
    return putPolicy.UploadToken(mac)
}
