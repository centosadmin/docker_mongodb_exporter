package profiler

import (
	"strconv"
	"sync"

	"github.com/percona/mongodb_exporter/collector/profiler/tailer"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func getDBs(conn *mgo.Session) []string {
	var databases []string
	res := struct {
		Databases []struct {
			Name string `bson:"name"`
		} `bson:"databases"`
	}{}
	err := conn.DB("admin").Run(bson.D{{"listDatabases", 1}}, &res)
	if err == nil {
		for _, db := range res.Databases {
			databases = append(databases, db.Name)
		}
	}
	return databases
}

func getCollections(conn *mgo.Session, db string) []string {
	var collections []string
	res := struct {
		Cursor struct {
			FirstBatch []struct {
				Name string `bson:"name"`
			} `bson:"firstBatch"`
		} `bson:"cursor"`
	}{}
	err := conn.DB(db).Run(bson.D{{"listCollections", 1}}, &res)
	if err == nil {
		for _, coll := range res.Cursor.FirstBatch {
			collections = append(collections, coll.Name)
		}
	}
	return collections
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
	var has bool
	for _, coll := range getCollections(conn, db) {
		if coll == "system.profile" {
			has = true
			break
		}
	}
	return has
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

func main() {
	conn, err := NewConnection("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	tailers := make(map[string]*tailer.Tailer)
	done := make(chan bool)
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
			go tailers[db].Tail(done, wg)
		}
	}
	wg.Wait()
	<-done
}
