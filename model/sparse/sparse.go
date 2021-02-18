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
	memTable map[string][]byte
	blockSize int
	metricReadIndexElapsed time.Duration
	metricReadDataElapsed time.Duration
}

func (s *sparse) Write(data []byte) {
	seg := input.GetRootAsTraceSegmentRequest(data, 0)
	endpoint := new(input.Field)
	if !seg.Fields(endpoint, 3) {
		log.Fatalf("failed to load segmentID")
	}
	key := string(endpoint.Value())
	var buffer []byte
	if b, ok := s.memTable[key]; ok {
		buffer = append(b, seg.SpansBytes()...)
	} else {
		buffer = seg.SpansBytes()
	}
	if len(buffer) > s.blockSize * 1024 {
		k := append([]byte(key), []byte(strconv.FormatInt(time.Now().UnixNano(), 10))...)
		s.db.Write(k, buffer)
		s.index.Write(append(seg.TraceID(), k...), nil)
		s.memTable[key] = nil
	} else {
		s.memTable[key] = buffer
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
	for key, buffer := range s.memTable {
		s.db.Write(append([]byte(key), []byte(strconv.FormatInt(time.Now().UnixNano(), 10))...), buffer)
	}
	s.db.Close()
}

func newSparse(blockSize int) model.Model {
	newDB := db.NewDB(blockSize)
	indexDB := db.NewDB(blockSize)
	return &sparse{db: &newDB, index: &indexDB, memTable: make(map[string][]byte, 10), blockSize: blockSize}
}
