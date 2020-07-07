package model

import (
	couchdb "github.com/leesper/couchdb-golang"
)

var couchDB *couchdb.Database

const databaseName = "flashlight"

func init() {
	var err error

	server, err := couchdb.NewServer("http://localhost:5984")
	if err != nil {
		panic(err)
	}

	if !server.Contains(databaseName) {
		couchDB, err = server.Create(databaseName)
	} else {
		couchDB, err = server.Get(databaseName)
	}
	if err != nil {
		panic(err)
	}
}
