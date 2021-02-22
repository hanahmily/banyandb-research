package compacted

import (
	"encoding/base64"
	"log"
	"testing"
	"time"

	db "github.com/hanahmily/banyandb-research"
	"github.com/hanahmily/banyandb-research/api"
	"github.com/hanahmily/banyandb-research/model"
)

func BenchmarkCompactedIndex(b *testing.B) {
	benchmarks := []struct {
		name string
		model model.Model
		isProxy bool
	}{
		{
			"512-bytes Database compacted index",
			newCompacted(512, db.CompressionAlgorithm_Snappy),
			false,
		},
		{
			"1024-bytes Database compacted index",
			newCompacted(1024, db.CompressionAlgorithm_Snappy),
			false,
		},
		{
			"512-bytes Proxy compacted index",
			newCompacted(512, db.CompressionAlgorithm_Snappy),
			true,
		},
		{
			"1024-bytes Proxy compacted index",
			newCompacted(1024, db.CompressionAlgorithm_Snappy),
			true,
		},
		{
			"512-bytes Database compacted index",
			newCompacted(512, db.CompressionAlgorithm_LZ4),
			false,
		},
		{
			"1024-bytes Database compacted index",
			newCompacted(1024, db.CompressionAlgorithm_LZ4),
			false,
		},
		{
			"512-bytes Proxy compacted index",
			newCompacted(512, db.CompressionAlgorithm_LZ4),
			true,
		},
		{
			"1024-bytes Proxy compacted index",
			newCompacted(1024, db.CompressionAlgorithm_LZ4),
			true,
		},
		{
			"512-bytes Database compacted index",
			newCompacted(512, db.CompressionAlgorithm_ZSTD),
			false,
		},
		{
			"1024-bytes Database compacted index",
			newCompacted(1024, db.CompressionAlgorithm_ZSTD),
			false,
		},
		{
			"512-bytes Proxy compacted index",
			newCompacted(512, db.CompressionAlgorithm_ZSTD),
			true,
		},
		{
			"1024-bytes Proxy compacted index",
			newCompacted(1024, db.CompressionAlgorithm_ZSTD),
			true,
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			var traceSize int
			if bm.isProxy {
				traceSize = 900000
			} else {
				traceSize = 20000
			}
			traces := make([]string, 0, traceSize)
			for i := 0; i < traceSize; i++ {
				data, traceID := api.GenerateInput(bm.isProxy)
				traces = append(traces, traceID)
				for j := 0; j < api.SegmentVariants; j++ {
					bytes, err := base64.StdEncoding.DecodeString(data[j % len(data)])
					if err != nil {
						log.Fatalf("failed to decode data")
					}
					bm.model.Write(bytes)
				}
			}
			t1 := time.Now()
			for i := 0; i < traceSize; i++ {
				bm.model.Get(traces[i])
			}
			elapsed := time.Since(t1)
			e := elapsed.Nanoseconds() / int64(traceSize)
			log.Printf("getting by trace id elapsed: %v, %v per query ", elapsed,  time.Duration(e))
			defer bm.model.Finish()
		})
	}
}
