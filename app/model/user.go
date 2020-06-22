package model

import (
	couchdb "github.com/leesper/couchdb-golang"
)

var couchDB *couchdb.Database

func init() {
	btDB, err := couchdb.NewDatabase("http://localhost:5984/flashlight")
	if err != nil {
		panic(err)
	}
}
