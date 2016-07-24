package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jasonlvhit/gocron"
	"github.com/parnurzeal/gorequest"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

//Config test
type Config struct {
	Server string `json:"server"`
	Node   string `json:"node"`
}

func uploadStatus(config Config) {
	v, _ := mem.VirtualMemory()
	c, _ := cpu.Percent(0, false)
	d, _ := disk.Usage("c:")

	fmt.Println("c:", c)
	fmt.Printf("Total: %v, Free:%v, UsedPercent:%f\n", v.Total, v.Free, v.UsedPercent)
	fmt.Printf("HD: %v GB  Free: %v GB Usage:%f\n", d.Total/1024/1024/1024, d.Free/1024/1024/1024, d.UsedPercent)

	jsonBody := fmt.Sprintf(`{"nodeCode":"%s", "nodeStatus":{"status":"up","metrics":[`+
		`{"title":"Mem","type":"gauge","value": %f,"width":0.25}`+
		`,{"title":"CPU","type":"line","value": {"x":"","y":%f,"xcount":40},"width":0.75}`+
		`]}}`, config.Node, v.UsedPercent, c[0])

	gorequest.New().Post(config.Server + "/api/status/upload/").Send(jsonBody).End()
}

func main() {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	configuration := Config{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("config: %#v", configuration)
	gocron.Every(1).Seconds().Do(uploadStatus, configuration)
	<-gocron.Start()
}
