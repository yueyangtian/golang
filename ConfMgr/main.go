package main
import (
    redis "ConfMgr/RedisModule"
    tools "ConfMgr/Tools"
    "github.com/fatih/structs"
    "encoding/json"
    "fmt"
)

func TestUpdate() {
   var conf redis.ConfRedis

   conf.ServiceHost = "127.0.0.1"
   conf.ServiceName = "test"
   conf.ConfKey = "test:127.0.0.1"
   conf.ConfType = "json"
   conf.ConfData = "{\"type_1\":\"test\", \"type_2\":\"test\"}"
   conf_map := structs.Map(&conf)
   var conf_list []string
   for k,v := range conf_map {
       conf_list = append(conf_list,k,v.(string))
   }

   redis.UpdateConf(conf_map["ConfKey"].(string),conf_map)
}
func TestParser() {
   var conf redis.ConfRedis

   conf.ServiceHost = "127.0.0.1"
   conf.ServiceName = "test4"
   conf.ConfType = "ini"
   conf.ConfData = "{\"type_1\":\"test\", \"type_2\":\"test\", \"type_3\":{\"index\": \"1\"}}"

   str, err := json.Marshal(conf)
   if err != nil {
       fmt.Printf("%s\n", err.Error())
       return
   }
   conf_new, err := redis.Parser(str)
   if err != nil {
       fmt.Printf("%s\n", err.Error())
       return
   }
   for k,v := range conf_new {
       fmt.Printf("%s:%s\n", k, v.(string))
   }
   redis.UpdateConf(conf_new["ConfKey"].(string),conf_new)

}
func TestGetConf() {
   s,_ := redis.GetConf("metadataService:10.12.23.177")
   fmt.Println(string(s))
}
func TestTools() {
    s, _ := tools.ScanIniFile("/opt/tvu/metadatad/conf/metadatad.conf")
    tools.UpdateInRedis("10.12.23.177/18","metadataService", "ini", string(s))
    fmt.Println(string(s))
}
func main() {
    //TestParser()
    //TestGetConf()
    TestTools()
}
