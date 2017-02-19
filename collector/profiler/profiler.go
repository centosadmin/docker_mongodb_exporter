package profiler

import (
	"sync"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Profiler struct {
	Session *mgo.Session
	tailers map[string]*Tailer
	wg      sync.WaitGroup
}

func New(sess *mgo.Session) *Profiler {
	return &Profiler{
		Session: sess,
		tailers: make(map[string]*Tailer),
	}
}

func (p *Profiler) Start() {
	tailers := make(map[string]*Tailer)
	sess := conn.GetSession()
	dbs := getDBs(sess)
	p.wg.Add(len(dbs))
	for _, db := range dbs {
		if db == "admin" || db == "local" {
			continue
		}
		if hasSystemProfile(sess, db) {
			tailers[db] = tailer.New(conn.GetConn(), db)
			go tailers[db].Tail(p.wg)
		}
	}
	p.wg.Wait()
}
