package api

import (
	"encoding/base64"
	"math/rand"
	"time"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"

	"github.com/hanahmily/banyandb-research/api/input"
	v3common "github.com/hanahmily/banyandb-research/api/pb/skywalking/network/common/v3"
	v3 "github.com/hanahmily/banyandb-research/api/pb/skywalking/network/language/agent/v3"
)

const spanSize = 20 
const SegmentVariants = 10

type sqlCollection []string

var endpointNames = make([]string, SegmentVariants, SegmentVariants)
var sqlTmpls = make([]sqlCollection, 0, SegmentVariants)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}


func init() {
	for i := 0; i < SegmentVariants; i++ {
		length := rand.Intn(50 - 5) + 5
		endpointName := make([]byte, length * 2)
		rand.Read(endpointName)
		endpointNames[i] = string(endpointName)
		sqlCol := make(sqlCollection, 0)
		for i := 0; i < 20; i++ {
			sqlCol = append(sqlCol, randStringBytes(1000))
		}
		sqlTmpls = append(sqlTmpls, sqlCol)
	}
}

func GenerateInput(isProxy bool) ([]string, string) {
	result := make([]string, 0, SegmentVariants)
	b := flatbuffers.NewBuilder(0)

	traceID := uuid.NewString()
	for i := 0; i < SegmentVariants; i++ {
		b.Reset()
		segmentID := uuid.NewString()
		fieldsPoss := createFields(
			b,
			segmentID,
			"produtpage",
			"a12ff60b-5807-463b_1",
			endpointNames[i],
			"a12ff60b-5807-463b_1",
			"500",
			"true",
		)
		input.TraceSegmentRequestStartFieldsVector(b, 7)
		for _, pos := range fieldsPoss {
			b.PrependUOffsetT(pos)
		}

		fieldsPos := b.EndVector(7)
		traceIDPos := b.CreateByteString([]byte(traceID))

		segment := v3.SegmentObject{
			TraceId: traceID,
			TraceSegmentId: segmentID,
			ServiceInstance: "a12ff60b-5807-463b_1",
			Service: "productpage",
			Spans: generateSpans(isProxy, i),
		}
		data, _ := proto.Marshal(&segment)
		//log.Info("segment size:%s", humanize.Bytes(uint64(len(data))))
		spanPos := b.CreateByteVector(data)

		input.TraceSegmentRequestStart(b)
		input.TraceSegmentRequestAddTraceID(b, traceIDPos)
		now := time.Now().UnixNano()
		input.TraceSegmentRequestAddStartTime(b, now)
		input.TraceSegmentRequestAddEndTime(b, now + int64(500 * time.Millisecond))
		input.TraceSegmentRequestAddFields(b, fieldsPos)

		input.TraceSegmentRequestAddSpans(b, spanPos)
		b.Finish(input.TraceSegmentRequestEnd(b))
		result = append(result, base64.StdEncoding.EncodeToString(b.FinishedBytes()))
	}
	return result, traceID
}

func generateSpans(isProxy bool, segIdx int) []*v3.SpanObject {
	spanId := 0
	spans := make([]*v3.SpanObject, 0)
	spans = append(spans, &v3.SpanObject{
		SpanId: int32(spanId),
		ParentSpanId: 0,
		StartTime: time.Now().UnixNano(),
		EndTime: time.Now().UnixNano(),
		OperationName: "",
		SpanType: v3.SpanType_Entry,
		SpanLayer: v3.SpanLayer_Http,
		ComponentId: 11,
		IsError: false,
		Tags: []*v3common.KeyStringValuePair{
			{
				Key: "url",
				Value: "http://test.com/vvv",
			},
			{
				Key: "param",
				Value: randStringBytes(20),
			},
			{
				Key: "status_code",
				Value: "200",
			},
		},
	})

	if !isProxy {
		for i := 0; i < spanSize; i++ {
			parentID := spanId
			spanId++
			spans = append(spans, &v3.SpanObject{
				SpanId: int32(spanId),
				ParentSpanId: int32(parentID),
				StartTime: time.Now().UnixNano(),
				EndTime: time.Now().UnixNano(),
				OperationName: "",
				SpanType: v3.SpanType_Exit,
				SpanLayer: v3.SpanLayer_Database,
				ComponentId: 21,
				IsError: false,
				Peer: "10.0.0.112:8978",
				Tags: []*v3common.KeyStringValuePair{
					{
						Key: "db.statement",
						Value: sqlTmpls[segIdx][i],
					},
					{
						Key:   "db.bind_vars",
						Value: randStringBytes(20),
					},
					{
						Key: "db.type",
						Value: "sql",
					},
					{
						Key: "db.instance",
						Value: "10.0.0.11:1234",
					},
				},
			})
		}
	}
	parentID := spanId
	spanId++
	spans = append(spans, &v3.SpanObject{
		SpanId: int32(spanId),
		ParentSpanId: int32(parentID),
		StartTime: time.Now().UnixNano(),
		EndTime: time.Now().UnixNano(),
		OperationName: "",
		SpanType: v3.SpanType_Exit,
		SpanLayer: v3.SpanLayer_Http,
		ComponentId: 11,
		IsError: false,
		Peer: "10.0.0.112:8978",
		Tags: []*v3common.KeyStringValuePair{
			{
				Key: "url",
				Value: "http://test.com/vvv",
			},
			{
				Key: "method",
				Value: "GET",
			},
			{
				Key: "param",
				Value: randStringBytes(20),
			},
			{
				Key: "status_code",
				Value: "200",
			},
		},
	})
	return spans
}

func createFields(b *flatbuffers.Builder, values ...string) []flatbuffers.UOffsetT {
	result := make([]flatbuffers.UOffsetT, len(values))
	for _, v := range values {
		vp := b.CreateString(v)
		input.FieldStart(b)
		input.FieldAddValue(b, vp)
		result = append(result, input.FieldEnd(b))
	}
	return result
}
