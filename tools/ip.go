package tools

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/viper"
)

func GetLocation(ip string) string {
	var address = "已关闭位置获取"
	if viper.GetBool("settings.public.isLocation") {
		if ip == "127.0.0.1" || ip == "localhost" {
			return "内部IP"
		}
		resp, err := http.Get("https://restapi.amap.com/v3/ip?ip=" + ip + "&key=3fabc36c20379fbb9300c79b19d5d05e")
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		s, err := ioutil.ReadAll(resp.Body)

		m := make(map[string]string)

		err = json.Unmarshal(s, &m)
		if err != nil {
			fmt.Println("Umarshal failed:", err)
		}
		if m["province"] == "" {
			return "未知位置"
		}
		address = m["province"] + "-" + m["city"]
	}
	return address
}
