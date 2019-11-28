package init

import (
	"com.capture/conf"
	"com.capture/dispatch"
	"com.capture/filter"
	"com.capture/log"
	"com.capture/message"
	packHandler "com.capture/package"
	"com.capture/service/http"
	"fmt"
)

func init() {
	//初始化日志配置
	logConf := log.LogConf{conf.LogFile}
	logConf.InitLog()

	client := message.Client{*conf.DesHost, *conf.En}
	config := conf.Conf{client}

	var message message.Message
	message = &client

	var f filter.Filter
	f = &filter.ContentFilter{*conf.ContentFilters}

	dis := dispatch.Dispatch{config, message, f}

	server := http.HttpServer{"9091"}
	go server.StartHttp()

	s := packHandler.Pack{*conf.PackageFilter, dis}

	go s.InitCapture()

	bannar()
}
func bannar() {
	fmt.Println("")
	fmt.Println("***************************************************************************")
	fmt.Println("***************************************************************************")
	fmt.Println("********************** welcome use capture utils **************************")
	fmt.Println("************************* start  success! *********************************")
	fmt.Println("***************************************************************************")
	fmt.Println("")
}
