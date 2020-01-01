package RedisModule
import (
    "fmt"
    "errors"
    "encoding/json"
    "reflect"
    "strings"
    "github.com/go-redis/redis"
    "github.com/fatih/structs"
)

func Parser(conf_str []byte)(conf map[string]interface{}, err error) {
    var struct_conf ConfRedis
    err = json.Unmarshal(conf_str, &struct_conf)
    if err != nil {
        return nil, err
    }

    struct_conf.ConfKey = struct_conf.ServiceName + ":" + struct_conf.ServiceHost
    conf_map := structs.Map(&struct_conf)
    return conf_map, nil
}
func UpdateConf(key string, conf map[string]interface{}) (err error) {
    version_key := "version:" + key
    fn := func(tx *redis.Tx) error {
        _, err := tx.Get(version_key).Result()
        if err != nil && err != redis.Nil {
            return err
        }

        _,err = tx.Pipelined(func(pipe redis.Pipeliner) error {
            err = pipe.HMSet(key,conf).Err()
            if err != nil {
                return err
            }
            err = pipe.Incr(version_key).Err()
            if err != nil {
                return err
            }
            return nil
        })
        return err
    }
    err = Client.Watch(fn, key, version_key)
    return
}
func json2ini(jstr string)(ini_str []byte, err error) {
    var inistr string
    jsonMap := make(map[string]interface{})
    err = json.Unmarshal([]byte(jstr), &jsonMap)
    if err != nil {
        return nil, err
    }
    for k,v := range jsonMap {
        if reflect.TypeOf(k).String() != "string" {
            err = errors.New("configure index error")
            return nil, err
        }
        if reflect.TypeOf(v).String() == "map[string]interface {}" {
            second_map := reflect.ValueOf(v).Interface().(map[string]interface{})
            inistr = inistr + "[" + strings.Trim(k," ") + "]\n"
            for k1,v1 := range(second_map) {
                if reflect.TypeOf(v1).String() == "map[string]interface {}" {
                    err = errors.New("configure index error")
                    return nil, err
                }
                inistr = inistr + strings.Trim(k1," ") + " = " + strings.Trim(v1.(string), " ") + "\n"

                fmt.Println(reflect.TypeOf(k1).String(),":", reflect.TypeOf(v1).String())
            }
        } else {
            inistr = inistr + strings.Trim(k," ") + " = " + strings.Trim(v.(string), " ") + "\n"
        }
    }
    ini_str =  []byte(inistr)
    return
}
func GetConf(key string)(conf_str []byte, err error) {
    m, err := Client.HGetAll(key).Result()
    if err != nil {
        return nil, err
    }
    if _, ok := m["ConfType"]; !ok {
        err = errors.New("hash index empty")
        return nil, err
    }
    if _, ok := m["ConfData"]; !ok {
        err = errors.New("hash index empty")
        return nil, err
    }

    Type := m["ConfType"]
    Data := m["ConfData"]
    if Type == "ini" {
        conf_str, err = json2ini(Data)
    } else if Type == "json" {
        conf_str = []byte(Data)
    } else {
        err = errors.New("not support this format")
        conf_str = nil
    }
    return conf_str, nil

}
