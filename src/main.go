package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jasonlvhit/gocron"
	"github.com/parnurzeal/gorequest"
	//"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"time"
	"github.com/shirou/gopsutil/host"
	"math"
)

//Config test
type Config struct {
	Server string `json:"server"`
	Node   string `json:"node"`
}

func uploadStatus(config Config) {
	fmt.Println("Upload status...")
	v, _ := mem.VirtualMemory()
	h, _ := host.Info()

	//c, _ := cpu.Percent(0, false)
	//d, _ := disk.Usage("c:")

	usedMemory := v.Used / (1024 * 1024)
	// fmt.Printf("Used: %v\n", usedMemory)
	// fmt.Printf("HD: %v GB  Free: %v GB Usage:%f\n", d.Total/1024/1024/1024, d.Free/1024/1024/1024, d.UsedPercent)

	days := int32(math.Floor(float64(h.Uptime) / (60 * 60 * 24)))
	hours := int32(math.Floor(float64(h.Uptime - uint64(days * 3600 * 24)) / 3600))
	minutes := int32(math.Floor(float64(h.Uptime - uint64(days * 3600 * 24) - uint64(hours * 3600)) / 60))
	seconds := int32(math.Floor(float64(h.Uptime - uint64(days * 3600 * 24) - uint64(hours * 3600) - uint64(minutes * 60))))
	upTimeFormat := fmt.Sprintf("<span class='label label-primary'>%dd</span> <span class='label label-success'>%dh</span> <span class='label label-info'>%dm</span> <span class='label label-default'>%ds</span>", days, hours, minutes, seconds)

	hostFormat := fmt.Sprintf("<div><b>Name</b>: %s</div>" +
		" <div style='margin-top:20px;'><b>Uptime</b>: %s</div>", h.Hostname, upTimeFormat)

	jsonBody := fmt.Sprintf(`{"nodeCode":"%s", "nodeStatus":{"status":"up","metrics":[` +
		`{"title":"Host","type":"text","value": "%s","width":0.25},` +
		`{"title":"Mem","type":"line","value": {"x":"%s","y":%d,"xinterval":60,"xcount":15},"width":0.75}` +
		`]}}`, config.Node, hostFormat, time.Now().Format("15:04"), usedMemory)

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
