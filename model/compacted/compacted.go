package compacted

import (
	"bytes"

	db "github.com/hanahmily/banyandb-research"
	"github.com/hanahmily/banyandb-research/api/input"
	"github.com/hanahmily/banyandb-research/model"
)

var _ model.Model = &compacted{}
type compacted struct {
	db *db.DB
	blockSize int
	currentTraceID []byte
	buf []byte
}

func (s *compacted) Write(data []byte) {
	seg := input.GetRootAsTraceSegmentRequest(data, 0)
	if s.currentTraceID == nil {
		s.currentTraceID = seg.TraceID()
	}
	if bytes.Compare(s.currentTraceID, seg.TraceID()) == 0 {
		s.buf = append(s.buf, seg.SpansBytes()...)
	} else {
		s.db.Write(s.currentTraceID, s.buf)
		s.buf = seg.SpansBytes()
		s.currentTraceID = seg.TraceID()
	}
}

func (s *compacted) Get(traceID string) {
	s.db.Read(traceID, false, nil)
}

func (s *compacted) Finish() {
	if len(s.buf) > 0 {
		s.db.Write(s.currentTraceID, s.buf)
	}
	s.db.Close()
}

func newCompacted(blockSize int, algorithm db.CompressionAlgorithm) model.Model {
	newDB := db.NewDB(blockSize, algorithm)
	return &compacted{db: &newDB, blockSize: blockSize}
}
