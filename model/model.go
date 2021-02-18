package model

type Model interface {
	Write(data []byte)
	Get(traceID string)
	Finish()
}
