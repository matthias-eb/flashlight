package model

import (
	"encoding/json"
	"fmt"

	couchdb "github.com/leesper/couchdb-golang"
)

type indexesDB struct {
	Indexes indexDB `json:"indexes"`
}

type indexDB struct {
	Type string   `json:"type"`
	Def  defIndex `json:"def"`
}

type defIndex struct {
	Fields []fieldsIndex `json:"fields"`
}

type fieldsIndex struct {
	Timestamp string `json:"timestamp"`
}

var couchDB *couchdb.Database

const databaseName = "flashlight"

func init() {
	var err error
	var indexesDB []indexDB
	var timestampIndexable = false

	server, err := couchdb.NewServer("http://localhost:5984")
	if err != nil {
		panic(err)
	}

	auth_Token, err := server.Login("fishface", "followthefish")
	if err != nil {
		panic(err)
	}

	if !server.Contains(databaseName) {
		fmt.Printf("Creating new Database..\n")
		err = server.VerifyToken(auth_Token)
		if err != nil {
			panic(err)
		}
		// Fails due to being unauthorized, even after logging in... But why? No Use for the Auth_token foundet, needs to be in cookies somehow..
		couchDB, err = server.Create(databaseName)
		if err != nil {
			panic(err)
		}
	} else {
		couchDB, err = server.Get(databaseName)
		if err != nil {
			panic(err)
		}
	}

	indexes, err := couchDB.GetIndex()
	if err != nil {
		panic(err)
	}

	indexBytes, err := json.Marshal(indexes["indexes"])
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(indexBytes, &indexesDB)
	if err != nil {
		panic(err)
	}

	for _, indexDB := range indexesDB {
		if indexDB.Type == "json" && indexDB.Def.Fields[0].Timestamp == "asc" {
			timestampIndexable = true
		}

	}

	// This does not work and will onlythrow an internal server error
	if !timestampIndexable {
		_, _, err = couchDB.PutIndex([]string{"timestamp"}, "", "")
		if err != nil {
			panic(err)
		}
	}

	if err != nil {
		panic(err)
	}
}
