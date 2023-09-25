package serializers

import (
	"bytes"
	"encoding/gob"
	"serialization/internal/logger"
)

type NativeSerializer struct {
	data *Data
}

func (s *NativeSerializer) PrepareData() {
	s.data = &data
}

func (s *NativeSerializer) Serialize() []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	err := encoder.Encode(s.data)
	logger.FailOnError("native encode failed", err)

	return buffer.Bytes()
}

func (s *NativeSerializer) Deserialize(raw []byte) {
	var buffer bytes.Buffer
	buffer.Write(raw)

	decoder := gob.NewDecoder(&buffer)

	err := decoder.Decode(s.data)
	logger.FailOnError("native decode failed", err)
}

func (s *NativeSerializer) CheckResult() {
	s.data.Compare(&data)
}
