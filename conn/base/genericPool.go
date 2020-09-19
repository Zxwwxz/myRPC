package connBase

import (
    "errors"
    "io"
    "sync"
    "time"
)

var (
    //非法配置
    ErrInvalidConfig = errors.New("invalid pool config")
    //池已经关闭
    ErrPoolClosed    = errors.New("pool closed")
)

type factory func() (io.Closer, error)

type GenericPool struct {
    sync.Mutex                      // 多线程锁
    pool        chan io.Closer      // 池获取回收通道
    maxOpen     int                 // 池中最大资源数
    numOpen     int                 // 当前池中资源数
    minOpen     int                 // 池中最少资源数
    closed      bool                // 池是否已关闭
    maxLifetime time.Duration       // 最大断开时间
    factory     factory             // 创建连接的方法
}

func NewGenericPool(minOpen, maxOpen int, maxLifetime int64, factory factory) (*GenericPool, error) {
    if maxOpen <= 0 || minOpen > maxOpen {
        return nil, ErrInvalidConfig
    }
    p := &GenericPool{
        maxOpen:     maxOpen,
        minOpen:     minOpen,
        maxLifetime: time.Duration(maxLifetime) * time.Second,
        factory:     factory,
        pool:        make(chan io.Closer, maxOpen),
    }
    //先创建最少的
    for i := 0; i < minOpen; i++ {
        closer, err := factory()
        if err != nil {
            continue
        }
        p.numOpen++
        p.pool <- closer
    }
    return p, nil
}

//获取连接对象
func (p *GenericPool) Get() (io.Closer, error) {
    //已经关闭
    if p.closed {
        return nil, ErrPoolClosed
    }
    closer, err := p.getOrCreate()
    if err != nil {
        return nil, err
    }
    return closer, nil
}

func (p *GenericPool) getOrCreate() (io.Closer, error) {
    select {
    //可以获取
    case closer := <-p.pool:
        return closer, nil
    default:
    }
    p.Lock()
    //已经超过最大，不能在创建了
    if p.numOpen >= p.maxOpen {
        closer := <-p.pool
        p.Unlock()
        return closer, nil
    }
    // 新建连接
    closer, err := p.factory()
    if err != nil {
        p.Unlock()
        return nil, err
    }
    p.numOpen++
    p.Unlock()
    return closer, nil
}

// 释放单个资源到连接池
func (p *GenericPool) Release(closer io.Closer) error {
    //已经关闭
    if p.closed {
        return ErrPoolClosed
    }
    p.Lock()
    p.pool <- closer
    p.Unlock()
    return nil
}

// 关闭单个资源
func (p *GenericPool) Close(closer io.Closer) error {
    p.Lock()
    _ = closer.Close()
    p.numOpen--
    p.Unlock()
    return nil
}

// 关闭连接池，释放所有资源
func (p *GenericPool) Shutdown() error {
    if p.closed {
        return ErrPoolClosed
    }
    p.Lock()
    close(p.pool)
    for closer := range p.pool {
        _ = closer.Close()
        p.numOpen--
    }
    p.closed = true
    p.Unlock()
    return nil
}
