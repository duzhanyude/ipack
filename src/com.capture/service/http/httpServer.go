package http

import (
	"com.capture/buffer"
	"com.capture/statistic"
	"fmt"
	"log"
	"net/http"
)

type HttpServer struct {
	Listerner string
}

//开启服务
func (h *HttpServer) StartHttp() {
	http.HandleFunc("/", h.sayHelloName)
	err := http.ListenAndServe(":"+h.Listerner, nil)
	if err != nil {
		log.Fatal("ListenAndService", err)
	}
}
func (h *HttpServer) sayHelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	//fmt.Println(r.Form)
	//fmt.Println("path", r.URL.Path)
	//fmt.Println("scheme", r.URL.Scheme)
	//fmt.Println(r.Form["url_long"])
	//for k, v := range r.Form {
	//	fmt.Println("key:", k)
	//	fmt.Println("val:", strings.Join(v, ""))
	//}
	callBack := "<table style='margin: 40px auto;'border='1'  cellspacing='0'>" +
		"<tr><td>总包</td><td>%d</td></tr>"

	fmt.Fprintf(w, displayInfo(callBack), statistic.ReceivePackageNum) //这个写入到w的是输出到客户端的
}

func displayInfo(content string) string {
	if len(buffer.PackageIP) > 0 {
		totalMap := buffer.PackageIP
		for k, v := range totalMap {
			content = content + "<tr><td>" + k + "</td><td>" + v + "</td></tr>"
		}
	}
	content = content + "</table>"
	return content
}
