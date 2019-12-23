package html

import (
	"ipack/com.ipack/statistic"
	"strconv"
	"strings"
)

var to string
var i int

func IndexHtml() string {
	i = 2
	to = ""
	ipInfo := make(map[string]int)

	html := `
	<html>
		<head>
    		<title>首页</title>
			 <link href="https://cdn.bootcss.com/vis/4.21.0/vis.min.css" rel="stylesheet">
    		<style type="text/css">
        		#mynetwork {
            		width: 100%;
            		height: 400px;
            		border: 1px solid lightgray;
        	}
    </style>
         <style>
      h2{
         margin-top: 30px;
         text-align: center;
         background-color:#89bf04;
         color: #fff;
         font-weight: normal;
         padding: 15px 0
      }
     
   	</style>
	</head>
	<body>
	<h2>首页</h2>
	<nav>
	<a href="index">首页</a> |
	<a href="timeData">实时数据</a> |
	<a href="status">连接统计</a> |
	<a href="sys">系统</a> |
	</nav>
	
	<div id="mynetwork"></div>
	<script src="https://cdn.bootcss.com/vis/4.21.0/vis.min.js"></script>
<script type="text/javascript">
    // create an array with nodes
    var nodes = new vis.DataSet([
        {id: 1, label: '本机电脑\r\r` + strconv.Itoa(int(statistic.ReceivePackageNum)) + `'},
	`
	for elem := statistic.NetCardList.Front(); elem != nil; elem = elem.Next() {
		net := strings.Split(elem.Value.(string), "@@@")
		html = html + `{id: ` + strconv.Itoa(int(i)) + `, label: '` + net[0] + "\\r" + net[1] + `'},`
		to = to + `{from: ` + strconv.Itoa(int(i)) + `, to: 1,label:''},`
		ipInfo[net[1]] = i
		i++
	}
	statistic.PIP.Range(func(k, v interface{}) bool {
		value := ipInfo[k.(string)]
		if value == 0 {
			html = html + `{id: ` + strconv.Itoa(int(i)) + `, label: '` + k.(string) + `'},`
			ipInfo[k.(string)] = i
		}
		i++
		return true
	})
	statistic.FromToIP.Range(func(k, v interface{}) bool {
		ft := strings.Split(k.(string), "-")
		f := ipInfo[ft[0]]
		d := ipInfo[ft[1]]
		to = to + `{from: ` + strconv.Itoa(int(f)) + `, to: ` + strconv.Itoa(int(d)) + `,label:'` + strconv.Itoa(v.(int)) + `'},`
		return true
	})

	html = html + `
    ]);
    // create an array with edges
    var edges = new vis.DataSet([
        ` + to + `
    ]);
 
    // create a network
    var container = document.getElementById('mynetwork');
 
    // provide the data in the vis format
    var data = {
        nodes: nodes,
        edges: edges
    };
    var options = {
		 configure: {
			enabled: true,
			filter: 'nodes,edges',
			container: undefined,
        showButton: true
    },
	edges:{
		arrows:{
			to:{enabled:true}	
		}
		}
	};
 
    // initialize your network!
    var network = new vis.Network(container, data, options);
</script>
</body>
</html>
`

	return html
}
