package http

import (
	"fmt"
	"golang.org/x/net/websocket"
	"ipack/com.ipack/service/http/handler"
	"ipack/com.ipack/service/http/html"
	"log"
	"net/http"
	"strings"
)

type HttpServer struct {
	Listener string
}

var HtmlMap = make(map[string]func() string)

//开启服务
func (h *HttpServer) StartHttp() {
	HtmlMap["/index"] = html.IndexHtml
	HtmlMap["/sys"] = html.SystemHtml
	HtmlMap["/status"] = html.StatusHtml
	http.HandleFunc("/", h.index)

	http.Handle("/getOnline", websocket.Handler(handler.Online))
	err := http.ListenAndServe(":"+h.Listener, nil)
	if err != nil {
		log.Fatal("ListenAndService", err)
	}
}
func (h *HttpServer) index(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	f := HtmlMap[r.URL.Path]
	if f != nil {
		fmt.Fprintf(w, f())
	} else {
		fmt.Fprintf(w, strings.ReplaceAll(html.TimeMontorData(), "${ip}", r.Host))
	}
}
