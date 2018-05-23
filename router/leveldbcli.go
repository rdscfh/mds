package router

import (
	"github.com/syndtr/goleveldb/leveldb"
)

var db *leveldb.DB
var err error

func init() {
	var err error
	db, err = leveldb.OpenFile("./db", nil)
	if err != nil {
		panic(err)
	}
}
func DBGet(key string) (string, error) {
	v, err := db.Get([]byte(key), nil)
	if err != nil {
		return "", err
	}
	return string(v), err
}

func DBSet(key string, val []byte) error {
	_, err = db.Get([]byte(key), nil)
	if err != nil {
		return db.Put([]byte(key), val, nil)
	}
	return err
}

func DBatchSet(key string, val map[string][]byte) error {
	batch := new(leveldb.Batch)
	for k, v := range val {
		rk := key + "_" + k
		batch.Put([]byte(rk), v)
	}
	return db.Write(batch, nil)
}

func DBHGet(key string, filed string) (string, error) {
	return DBGet(key + "_" + filed)
}
