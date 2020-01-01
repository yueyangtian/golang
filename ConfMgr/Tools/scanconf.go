package Tools

import (
	redis "ConfMgr/RedisModule"
	ini "github.com/Unknwon/goconfig"
	"encoding/json"
)

func ScanIniFile(filename string)(jstr []byte, err error) {
	cfg, err := ini.LoadConfigFile(filename)
	if err != nil {
		return
	}

	section_map := make(map[string]interface{})
	sessionList := cfg.GetSectionList()
	for _, section := range sessionList {
		ini_map := make(map[string]interface{})
		keyList := cfg.GetKeyList(section)
		for _, key := range keyList {
			val, err := cfg.GetValue(section,key)
			if err != nil {
				continue
			}
			ini_map[key] = val
		}
		section_map[section] = ini_map
	}

	str, err := json.Marshal(section_map)
	if err != nil {
		return nil, err
	}
	jstr = []byte(str)
	return jstr, nil
}
func UpdateInRedis(hostip, servicename, conftype , data string)(err error) {
	var conf redis.ConfRedis

	conf.ServiceHost = hostip
	conf.ServiceName = servicename
	conf.ConfType = "ini"
	conf.ConfData = data

	str, err := json.Marshal(conf)
	if err != nil {
		return err
	}
	conf_new, err := redis.Parser(str)
	if err != nil {
		return err
	}
	err = redis.UpdateConf(conf_new["ConfKey"].(string),conf_new)
	if err != nil {
		return err
	}
	return nil
}