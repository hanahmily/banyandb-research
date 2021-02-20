package sparse

import (
	"log"
	"strconv"
	"time"

	db "github.com/hanahmily/banyandb-research"
	"github.com/hanahmily/banyandb-research/api/input"
	"github.com/hanahmily/banyandb-research/model"
)

var _ model.Model = &sparse{}
type sparse struct {
	db *db.DB
	index *db.DB
	memTable map[string]*memTable
	blockSize int
	metricReadIndexElapsed time.Duration
	metricReadDataElapsed time.Duration
}

type memTable struct {
	q []byte
	startTime int64
}

func (m *memTable) getID(key string) []byte {
	return append([]byte(key), []byte(strconv.FormatInt(m.startTime, 10))...)
}


func (s *sparse) Write(data []byte) {
	seg := input.GetRootAsTraceSegmentRequest(data, 0)
	endpoint := new(input.Field)
	if !seg.Fields(endpoint, 3) {
		log.Fatalf("failed to load segmentID")
	}
	key := string(endpoint.Value())
	m, ok := s.memTable[key]
	if ok {
		m.q = append(m.q, seg.SpansBytes()...)
	} else {
		m = &memTable{q: seg.SpansBytes(), startTime: seg.StartTime()}
		s.memTable[key] = m
	}
	s.index.Write(append(seg.TraceID(), m.getID(key)...), nil)
	if len(m.q) > s.blockSize * 1024 {
		s.db.Write(m.getID(key), m.q)
		delete(s.memTable, key)
	}
}

func (s *sparse) Get(traceID string) {
	eps := make([]string, 0)
	t1 := time.Now()
	s.index.Read(traceID, true, func(key, val []byte) {
		endpoint := key[len(traceID):]
		eps = append(eps, string(endpoint))
		
	})
	elapsed := time.Since(t1)
	s.metricReadIndexElapsed = s.metricReadIndexElapsed + elapsed
	t2 := time.Now()
	for _, ep := range eps {
		s.db.Read(ep, false, nil)
	}
	elapsed1 := time.Since(t2)
	s.metricReadDataElapsed = s.metricReadDataElapsed + elapsed1
}

func (s *sparse) Finish() {
	log.Printf("querying index elapsed: %v, data elapsed: %v", s.metricReadIndexElapsed, s.metricReadDataElapsed)
	for key, m := range s.memTable {
		s.db.Write(m.getID(key), m.q)
	}
	s.index.Close()
	s.db.Close()
}

func newSparse(blockSize int) model.Model {
	newDB := db.NewDB(blockSize)
	indexDB := db.NewDB(blockSize)
	return &sparse{db: &newDB, index: &indexDB, memTable: make(map[string]*memTable, 10), blockSize: blockSize}
}
