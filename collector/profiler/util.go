package profiler

import (
	"sync"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func getDBs(conn *mgo.Session) []string {
	var dbs []string
	dbs, _ = conn.DatabaseNames()
	return dbs
}

func getCollections(conn *mgo.Session, db string) []string {
	var colls []string
	colls, _ = conn.DB(db).CollectionNames()
	return colls
}

func getProfilingLevel(conn *mgo.Session, db string) (int, int) {
	res := struct {
		Was    int `bson:"was"`
		SlowMs int `bson:"slowms"`
	}{}
	err := conn.DB(db).Run(bson.D{{"profile", -1}}, &res)
	if err != nil {
		return 0, 0
	}
	return res.Was, res.SlowMs
}

func hasSystemProfile(conn *mgo.Session, db string) bool {
	for _, coll := range getCollections(conn, db) {
		if coll == "system.profile" {
			return true
		}
	}
	return false
}
