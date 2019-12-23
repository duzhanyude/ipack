package leveldb

import (
	"github.com/syndtr/goleveldb/leveldb"
	"log"
)

var Db *leveldb.DB

type LevelDB struct {
	DbName string
}

func (l *LevelDB) Init() {
	db, err := leveldb.OpenFile(l.DbName, nil)
	if err != nil {
		log.Println(err)
	}
	Db = db
}
func (l *LevelDB) Save(k, v string) error {
	return Db.Put([]byte(k), []byte(v), nil)
}
func (l *LevelDB) Get(k string) string {
	data, _ := Db.Get([]byte(k), nil)
	return string(data)
}
func (l *LevelDB) Del(k string) error {
	return Db.Delete([]byte(k), nil)
}
