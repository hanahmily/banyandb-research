package dense

import (
	"encoding/binary"

	db "github.com/hanahmily/banyandb-research"
	"github.com/hanahmily/banyandb-research/api/input"
	"github.com/hanahmily/banyandb-research/model"
)

var _ model.Model = &dense{}

type dense struct {
	db *db.DB
}

func (s *dense) Write(data []byte) {
	seg := input.GetRootAsTraceSegmentRequest(data, 0)
	traceID := seg.TraceID()
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(seg.StartTime()))
	key := append(traceID, b...)
	s.db.Write(key, seg.SpansBytes())
}

func (s dense) Get(traceID string) {
	s.db.Read(traceID, false, nil)
}

func (s *dense) Finish() {
	s.db.Close()
}

func newDense(blockSize int) model.Model {
	newDb:= db.NewDB(blockSize, db.CompressionAlgorithm_Snappy)
	return &dense{db: &newDb}
}

var _ model.Model = &compactedDense{}

type compactedDense struct {
	db *db.DB
}

func (c *compactedDense) Write(data []byte) {
	seg := input.GetRootAsTraceSegmentRequest(data, 0)
	traceID := seg.TraceID()
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(seg.StartTime()))
	key := append(traceID, b...)
	key = append(key, seg.SpansBytes()...)
	c.db.Write(key, nil)
}

func (c *compactedDense) Get(traceID string) {
	c.db.Read(traceID, true, nil)
}

func (c *compactedDense) Finish() {
	c.db.Close()
}

func newCompactedDense(blockSize int) model.Model {
	newDb:= db.NewDB(blockSize, db.CompressionAlgorithm_Snappy)
	return &compactedDense{db: &newDb}
}

