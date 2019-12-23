package httpServer

import (
	"fmt"
	"log"
	"net/http"
)

type HttpServer struct {
	Num  *int32
	Flag *string
}

//开启服务
func (h *HttpServer) StartHttp() {
	http.HandleFunc("/", h.sayHelloName)
	err := http.ListenAndServe(":9091", nil)
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

		"<tr><td>总包</td><td>%d</td></tr>" +
		"<tr><td>连接</td><td>%s</td></tr>" +
		"</table>"

	fmt.Fprintf(w, callBack, *h.Num, *h.Flag) //这个写入到w的是输出到客户端的
}
