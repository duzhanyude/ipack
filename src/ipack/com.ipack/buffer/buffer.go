package buffer

import (
	"container/list"
	"ipack/com.ipack/db/leveldb"
	"net"
	"sync"
)

var clientMap sync.Map
var PackageIP sync.Map
var PackageList = new(list.List)
var WsConnect sync.Map
var GlobalDB leveldb.LevelDB

func GetClient(key string) net.Conn {
	conn, _ := clientMap.Load(key)
	if conn == nil {
		return nil
	}
	return conn.(net.Conn)
}
func SaveClient(key string, conn net.Conn) {
	clientMap.Store(key, conn)
}
func DelClient(key string) {
	conn := GetClient(key)
	if conn != nil {
		conn.Close()
	}
	clientMap.Delete(key)
}
