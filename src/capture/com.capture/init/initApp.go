package init

import (
	"capture/com.capture/buffer"
	"capture/com.capture/conf"
	"capture/com.capture/db/leveldb"
	"capture/com.capture/dispatch"
	"capture/com.capture/filter"
	"capture/com.capture/log"
	"capture/com.capture/message"
	packHandler "capture/com.capture/package"
	"capture/com.capture/service/http"
	"fmt"
)

func init() {
	//初始化日志配置
	logConf := log.LogConf{conf.LogFile}
	go logConf.InitLog()
	//初始化数据库
	go initDB()

	client := message.Client{*conf.DesHost, *conf.En}
	config := conf.Conf{client}

	var messC message.Message
	messC = &client

	var messW message.Message
	messW = &message.WebSocketClient{}

	resgiterM := message.RegisterMessage{}

	resgiterM.Register(messC)
	resgiterM.Register(messW)

	var f filter.Filter
	f = &filter.ContentFilter{*conf.ContentFilters}

	dis := dispatch.Dispatch{config, resgiterM, f, *conf.ContentLogFile}

	go dis.WriteFile()

	server := http.HttpServer{*conf.WebPort}
	go server.StartHttp()

	s := packHandler.Pack{*conf.PackageFilter, dis}

	go s.InitCapture()

	bannar()
}

//初始化数据库
func initDB() {
	db := leveldb.LevelDB{"db"}
	db.Init()
	buffer.GlobalDB = db
	//db.Save("kry","你好")
	//fmt.Println(db.Get("kry"))
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
