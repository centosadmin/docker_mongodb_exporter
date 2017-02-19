package tailer

import (
	"time"
)

type Profile struct {
	Operation      string    `bson:"op,omitempty"`
	Namespace      string    `bson:"ns,omitempty"`
	Timestamp      time.Time `bson:"ts,omitempty"`
	NumYield       float64   `bson:"numYield,omitempty"`
	NReturned      float64   `bson:"nreturned,omitempty"`
	Millis         float64   `bson:"millis,omitempty"`
	WriteConflicts float64   `bson:"writeConflicts,omitempty"`
	KeyUpdates     float64   `bson:"keyUpdates,omitempty"`
	DocsExamined   float64   `bson:"docsExamined,omitempty"`
	KeysExamined   float64   `bson:"keysExamined,omitempty"`
}
