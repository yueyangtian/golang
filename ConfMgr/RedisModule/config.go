package RedisModule
import (
    "github.com/go-redis/redis"
    "time"
)

type ConfRedis struct {
    ServiceName            string
    ServiceHost            string
    ConfKey                string
    ConfType               string
    ConfData               string
}

type  ClientConf struct {
    redis_addr              string
    redis_password          string
    redis_db                int
}

var Client *redis.Client
var CrdbConf *ClientConf


func init() {
    CrdbConf = new(ClientConf)
    CrdbConf.redis_addr = "127.0.0.1:6379"
    CrdbConf.redis_password = ""
    CrdbConf.redis_db = 0

    Client = redis.NewClient(&redis.Options {
        Addr:               CrdbConf.redis_addr,
        Password:           CrdbConf.redis_password,
        DB:                 CrdbConf.redis_db,
        DialTimeout:        10*time.Second,
        ReadTimeout:        30*time.Second,
        WriteTimeout:       30*time.Second,
        PoolSize:           10,
        PoolTimeout:        30*time.Second })
}
