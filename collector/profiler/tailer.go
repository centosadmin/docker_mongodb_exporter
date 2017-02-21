package profiler

import (
	"sync"
	"time"

	"github.com/golang/glog"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Tailer struct {
	Database string
	Session  *mgo.Session
	LastTs   time.Time
	stats    map[string]*Stats
}

func NewTailer(sess *mgo.Session, db string) *Tailer {
	return &Tailer{
		Database: db,
		Session:  sess,
		LastTs:   time.Now(),
		stats:    make(map[string]*Stats),
	}
}

func (t *Tailer) SystemProfile() *mgo.Collection {
	return t.Session.DB(t.Database).C("system.profile")
}

func (t *Tailer) GetStats(name string) *Stats {
	if _, ok := t.stats[name]; ok {
		return t.stats[name].Get()
	}
	return &Stats{}
}

func (t *Tailer) SetStats(name string, stats *Stats) {
	if _, ok := t.stats[name]; ok {
		t.stats[name].Lock()
		defer t.stats[name].Unlock()
		t.stats[name] = stats
	}
	t.stats[name] = stats
}

func (t *Tailer) tailQuery() bson.M {
	selfNs := t.Database + ".system.profile"
	return bson.M{
		"ts": bson.M{"$gt": t.LastTs},
		"ns": bson.M{"$ne": selfNs},
		"$or": []bson.M{
			bson.M{"millis": bson.M{"$gt": 0}},
			bson.M{"docsExamined": bson.M{"$gt": 0}},
			bson.M{"numYield": bson.M{"$gt": 0}},
		},
	}
}

func (t *Tailer) Tail(wg sync.WaitGroup) {
	defer wg.Done()
	glog.Infof("Starting profiler tailing on database '%s'\n", t.Database)
	iter := t.SystemProfile().Find(t.tailQuery()).Tail(5 * time.Second)
	for {
		profile := &Profile{}
		for iter.Next(profile) {
			t.LastTs = profile.Timestamp
			stats := t.GetStats(profile.Operation)
			stats.Count += 1
			stats.DocsExamined += profile.DocsExamined
			stats.KeysExamined += profile.KeysExamined
			stats.KeyUpdates += profile.KeyUpdates
			stats.Millis += profile.Millis
			stats.NumYields += profile.NumYield
			stats.WriteConflicts += profile.WriteConflicts
			t.SetStats(profile.Operation, stats)
			time.Sleep(10 * time.Millisecond)
		}
		if iter.Err() != nil {
			glog.Errorf("Profiler tailing error on database '%s': %s\n", t.Database, iter.Err().Error())
			break
		}
		if iter.Timeout() {
			time.Sleep(50 * time.Millisecond)
			continue
		}
		time.Sleep(50 * time.Millisecond)
		glog.V(2).Infof("Restarting profiling tailing on '%s' due to timeout\n", t.Database)
		iter = t.SystemProfile().Find(t.tailQuery()).Tail(5 * time.Second)
	}
	iter.Close()
}
