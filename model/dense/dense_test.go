package dense

import (
	"encoding/base64"
	"log"
	"testing"
	"time"

	"github.com/hanahmily/banyandb-research/api"
	"github.com/hanahmily/banyandb-research/model"
)

func BenchmarkDenseIndex(b *testing.B) {
	benchmarks := []struct {
		name string
		model model.Model
		isProxy bool
	}{
		{
			"512-bytes Database dense index",
			newDense(512),
			false,
		},
		{
			"1024-bytes Database dense index",
			newDense(1024),
			false,
		},
		//{
		//	"512-bytes Database compacted dense index",
		//	newCompactedDense(512),
		//	false,
		//},
		//{
		//	"1024-bytes Database compacted dense index",
		//	newCompactedDense(1024),
		//	false,
		//},
		{
			"512-bytes Proxy dense index",
			newDense(512),
			true,
		},
		{
			"1024-bytes Proxy dense index",
			newDense(1024),
			true,
		},
		//{
		//	"512-bytes Proxy compacted dense index",
		//	newCompactedDense(512),
		//	true,
		//},
		//{
		//	"1024-bytes Proxy compacted dense index",
		//	newCompactedDense(1024),
		//	true,
		//},
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
			log.Println("getting by trace id elapsed: ", elapsed)
			defer bm.model.Finish()
		})
	}
}

