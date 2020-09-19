package connBase

import "io"

type Pool interface {
    Get() (io.Closer, error) // 获取资源
    Release(io.Closer) error     // 释放资源
    Close(io.Closer) error       // 关闭资源
    Shutdown() error             // 关闭池
}