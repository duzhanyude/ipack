package html

import (
	"ipack/com.ipack/buffer"
	"ipack/com.ipack/statistic"
	"strconv"
)

func StatusHtml() string {

	html := `<html>
		<head>
    		<title>连接数据展示</title>
         <style>
     
      h2{
         margin-top: 30px;
         text-align: center;
         background-color: #89bf04;
         color: #fff;
         font-weight: normal;
         padding: 15px 0
      }
	table{
		table-layout:fixed;
		width:1000px;
		word-break:break-all;
	}
   </style>
		</head>
	<body>
      <h2>连接数据展示</h2>
	<nav>
		<a href="index">首页</a> |
		<a href="timeData">实时数据</a> |
		<a href="status">连接统计</a> |
		<a href="sys">系统</a> |
	</nav>
      `
	html = html + `<div style='width:100%'><table style='margin: 0px auto;'border='1'  cellspacing='0'>
		<tr><td style='width:160px;'>抓包总数</td><td>` + strconv.Itoa(int(statistic.ReceivePackageNum)) + `</td></tr>`

	html = displayInfo(html)
	html = html + `</body></html>`
	return html
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
	content = content + "</table></div>"
	return content
}
