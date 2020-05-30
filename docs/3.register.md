# 三、服务注册

##服务注册和发现原理
    服务端，客户端，注册中心
    选型：go最好consul和etcd
        服务健康检查，kv存储，支持watch
        consul：多数据中心，raft，cp
        zookeeper：paxos，cp
        etcd：raft，cp
        euerka：ap

##选项模式
##组件设计