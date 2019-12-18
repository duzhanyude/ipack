package html

func TimeMontorData() string {
	timeHtml := `<html>
		<head>
    		<title>在线实时数据</title>
         <style>
     
      h2{
         margin-top: 30px;
         text-align: center;
         background-color: #89bf04;
         color: #fff;
         font-weight: normal;
         padding: 15px 0
      }
      #chat{
         text-align: center;
        
      }
      #win{
         margin-top: 20px;
         text-align: center;
      }
      #sse{
         margin-top: 10px;
         text-align: center;
      }
      #sse button{
         background-color: #009688;
         color: #fff;
         height: 40px;
         border: 0;
         border-radius: 3px 3px;
         padding-left: 10px;
         padding-right: 10px;
         cursor: pointer;
      }
		textarea {
			background-color: #000;
			color:#fff;
			font-size:14px;
			font-weight: bold;
			padding:10px;
		}
   </style>
		</head>
	<body>
	<script type="text/javascript">
    var sock = null;
    // var wsuri = "wss://127.0.0.1:8080"; //本地的地址 是可以改变的哦
     var wsuri = "ws://${ip}/getOnline"; //本地的地址 是可以改变的哦


    window.onload = socketMessage

	function socketMessage () {
        //可以看到客户端JS，很容易的就通过WebSocket函数建立了一个与服务器的连接sock，当握手成功后，会触发WebScoket对象的onopen事件，告诉客户端连接已经成功建立。客户端一共绑定了四个事件。
        console.log("开始了 onload");

        sock = new WebSocket(wsuri);
        //建立连接后触发
        sock.onopen = function() {
            console.log(" 建立连接后触发 connected to " + wsuri);
        }
        // 关闭连接时候触发
        sock.onclose = function(e) {
            //console.log("关闭连接时候触发 connection closed (" + e.code + ")");
			//window.location.reload();
			sock.close()
			socketMessage()
        }
        // 收到消息后触发
        sock.onmessage = function(e) {
            //console.log("收到消息后触发 message received: " + e.data);
			var history= document.getElementById("history")
			if(history.value.length>5000){
				history.value = history.value.substring(0,5000)
			}
			history.value = history.value+ e.data
			history.scrollTop = history.scrollHeight
        }
        //发生错误的时候触发
        sock.onerror=function (e) {
			//console.log(e)
            console.log("发生错误时候触发"+wsuri)
        }
    };

</script>
      <h2>实时数据</h2>
	<nav>
		<a href="index">首页</a> |
		<a href="timeData">实时数据</a> |
		<a href="sys">系统</a> |
	</nav>
      <div id="chat">
           <textarea id="history" cols="198" rows="37" readonly></textarea>
      </div>
</body>
</html>`
	return timeHtml
}
