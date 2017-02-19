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

type Connection struct {
	sync.Mutex
	Session *mgo.Session
}

func NewConnection(uri string) (*Connection, error) {
	var err error
	conn := &Connection{}
	conn.Lock()
	conn.Session, err = mgo.Dial(uri)
	conn.Unlock()
	return conn, err
}

func (c *Connection) GetConn() *mgo.Session {
	c.Lock()
	defer c.Unlock()
	return c.Session.Copy()
}

func (c *Connection) Close() {
	c.Session.Close()
}

func StartProfilerListener() {
	conn, err := NewConnection("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	tailers := make(map[string]*Tailer)
	sess := conn.GetConn()
	dbs := getDBs(sess)

	var wg sync.WaitGroup
	wg.Add(len(dbs))
	for _, db := range dbs {
		if db == "admin" || db == "local" {
			continue
		}
		if hasSystemProfile(sess, db) {
			tailers[db] = tailer.New(conn.GetConn(), db)
			go tailers[db].Tail(wg)
		}
	}
	wg.Wait()
}
