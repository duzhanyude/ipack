package conf

import (
	"flag"
	"ipack/com.ipack/message"
)

var Num *int32

var F = flag.Int("f", 0, " 1 on 0 off flag server")

var TurnOff = flag.Int("on", 0, "log 2 receive 3 send 4 all")

var PackageFilter = flag.String("pfilter", "", "capture package filter")

//var PackageFilter = flag.String("pfilter", "tcp and dst port 80", "capture package filter")

//var DesHost = flag.String("d", "222.180.198.138:18729", "send host more;")
var DesHost = flag.String("d", "", "send host more;")

var ServicePort = flag.String("listener", "80", "listener port")

var WebPort = flag.String("web", "9091", "web port")

var ContentFilters = flag.String("cfilter", "", "filter content")

var En = flag.String("en", "0", " 1 Encrypt 2 decode")

var LogFile = flag.String("logFile", "", "logFile path")

var ContentLogFile = flag.String("clogFile", "statLog.txt", "contentlog path")

type Conf struct {
	Cli message.Client
}

func init() {
	flag.Parse() //暂停获取参数
}
