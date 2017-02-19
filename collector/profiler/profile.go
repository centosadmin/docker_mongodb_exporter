package profiler

import (
	"time"
)

type Profile struct {
	Timestamp      time.Time `bson:"ts,omitempty"`
	Operation      string    `bson:"op,omitempty"`
	Namespace      string    `bson:"ns,omitempty"`
	DocsExamined   float64   `bson:"docsExamined,omitempty"`
	KeysExamined   float64   `bson:"keysExamined,omitempty"`
	KeyUpdates     float64   `bson:"keyUpdates,omitempty"`
	NReturned      float64   `bson:"nreturned,omitempty"`
	NumYield       float64   `bson:"numYield,omitempty"`
	WriteConflicts float64   `bson:"writeConflicts,omitempty"`
	Millis         float64   `bson:"millis,omitempty"`
}
