package http

import (
	"capture/com.capture/buffer"
	"capture/com.capture/service/http/handler"
	"capture/com.capture/service/http/html"
	"capture/com.capture/statistic"
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type HttpServer struct {
	Listener string
}

var HtmlMap = make(map[string]func() string)

//开启服务
func (h *HttpServer) StartHttp() {
	HtmlMap["/index"] = html.IndexHtml
	//HtmlMap["/timeData"]=html.IndexHtml
	HtmlMap["/sys"] = html.SystemHtml

	http.HandleFunc("/", h.index)
	http.HandleFunc("/status", h.displayStatu)

	http.Handle("/getOnline", websocket.Handler(handler.Online))
	err := http.ListenAndServe(":"+h.Listener, nil)
	if err != nil {
		log.Fatal("ListenAndService", err)
	}
}

func (h *HttpServer) displayStatu(w http.ResponseWriter, r *http.Request) {
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
		"<tr><td>抓包总数</td><td>%d</td></tr>"

	fmt.Fprintf(w, displayInfo(callBack), statistic.ReceivePackageNum) //这个写入到w的是输出到客户端的
}
func (h *HttpServer) index(w http.ResponseWriter, r *http.Request) {
	f := HtmlMap[r.URL.Path]
	if f != nil {
		fmt.Fprintf(w, f())
	} else {
		fmt.Fprintf(w, strings.ReplaceAll(html.TimeMontorData(), "${ip}", r.Host))
	}

}
func displayInfo(content string) string {
	if buffer.PackageList.Len() > 0 {
		content = content + "<tr><td>连接总数</td><td>" + strconv.Itoa(buffer.PackageList.Len()) + "</td></tr>"
		for elem := buffer.PackageList.Front(); elem != nil; elem = elem.Next() {
			s, _ := buffer.PackageIP.Load(elem.Value.(string))
			//fileObj.Write([]byte(elem.Value.(string)+"@@@"+s))
			content = content + "<tr><td>" + elem.Value.(string) + "</td><td>" + s.(string) + "</td></tr>"
		}
	}
	content = content + "</table>"
	return content
}
