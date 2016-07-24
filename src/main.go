package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jasonlvhit/gocron"
	"github.com/parnurzeal/gorequest"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

//Config test
type Config struct {
	Server string `json:"server"`
	Node   string `json:"node"`
}

func uploadStatus(config Config) {
	fmt.Println("Upload status...")
	v, _ := mem.VirtualMemory()
	c, _ := cpu.Percent(0, false)
	//d, _ := disk.Usage("c:")

	// fmt.Printf("Total: %v, Free:%v, UsedPercent:%f\n", v.Total, v.Free, v.UsedPercent)
	// fmt.Printf("HD: %v GB  Free: %v GB Usage:%f\n", d.Total/1024/1024/1024, d.Free/1024/1024/1024, d.UsedPercent)

	jsonBody := fmt.Sprintf(`{"nodeCode":"%s", "nodeStatus":{"status":"up","metrics":[`+
		`{"title":"Mem","type":"gauge","value": %f,"width":0.25}`+
		`,{"title":"CPU","type":"line","value": {"x":"","y":%f,"xcount":40},"width":0.75}`+
		`]}}`, config.Node, v.UsedPercent, c[0])

	gorequest.New().Post(config.Server + "/api/status/upload/").Send(jsonBody).End()
}

func main() {
	configPath := "config.json"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}
	if _, err := os.Stat(configPath); os.IsExist(err) {
		fmt.Println("config file: " + configPath + " not exist!")
		return
	}

	file, _ := os.Open(configPath)
	decoder := json.NewDecoder(file)
	configuration := Config{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("Read config file error: ", err)
		return
	}

	if configuration.Server == "" {
		fmt.Println("Server can't be empty")
		return
	}
	if configuration.Node == "" {
		fmt.Println("Node can't be empty")
		return
	}

	gocron.Every(1).Seconds().Do(uploadStatus, configuration)
	<-gocron.Start()
}
