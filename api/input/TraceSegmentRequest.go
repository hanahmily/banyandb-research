// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package input

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type TraceSegmentRequest struct {
	_tab flatbuffers.Table
}

func GetRootAsTraceSegmentRequest(buf []byte, offset flatbuffers.UOffsetT) *TraceSegmentRequest {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &TraceSegmentRequest{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *TraceSegmentRequest) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *TraceSegmentRequest) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *TraceSegmentRequest) TraceID() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *TraceSegmentRequest) StartTime() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *TraceSegmentRequest) MutateStartTime(n int64) bool {
	return rcv._tab.MutateInt64Slot(6, n)
}

func (rcv *TraceSegmentRequest) EndTime() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *TraceSegmentRequest) MutateEndTime(n int64) bool {
	return rcv._tab.MutateInt64Slot(8, n)
}

func (rcv *TraceSegmentRequest) Fields(obj *Field, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *TraceSegmentRequest) FieldsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *TraceSegmentRequest) Tags(obj *Tag, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *TraceSegmentRequest) TagsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *TraceSegmentRequest) Spans(j int) byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetByte(a + flatbuffers.UOffsetT(j*1))
	}
	return 0
}

func (rcv *TraceSegmentRequest) SpansLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *TraceSegmentRequest) SpansBytes() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *TraceSegmentRequest) MutateSpans(j int, n byte) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateByte(a+flatbuffers.UOffsetT(j*1), n)
	}
	return false
}

func TraceSegmentRequestStart(builder *flatbuffers.Builder) {
	builder.StartObject(6)
}
func TraceSegmentRequestAddTraceID(builder *flatbuffers.Builder, traceID flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(traceID), 0)
}
func TraceSegmentRequestAddStartTime(builder *flatbuffers.Builder, startTime int64) {
	builder.PrependInt64Slot(1, startTime, 0)
}
func TraceSegmentRequestAddEndTime(builder *flatbuffers.Builder, endTime int64) {
	builder.PrependInt64Slot(2, endTime, 0)
}
func TraceSegmentRequestAddFields(builder *flatbuffers.Builder, fields flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(fields), 0)
}
func TraceSegmentRequestStartFieldsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func TraceSegmentRequestAddTags(builder *flatbuffers.Builder, tags flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(4, flatbuffers.UOffsetT(tags), 0)
}
func TraceSegmentRequestStartTagsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func TraceSegmentRequestAddSpans(builder *flatbuffers.Builder, spans flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(5, flatbuffers.UOffsetT(spans), 0)
}
func TraceSegmentRequestStartSpansVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(1, numElems, 1)
}
func TraceSegmentRequestEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
