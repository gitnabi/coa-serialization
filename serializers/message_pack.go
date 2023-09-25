package serializers

import (
	"serialization/internal/logger"

	msgpack "github.com/shamaton/msgpack/v2"
)

type MessagePackSerializer struct {
	data *Data
}

func (s *MessagePackSerializer) PrepareData() {
	s.data = &data
}

func (s *MessagePackSerializer) Serialize() []byte {
	bytes, err := msgpack.Marshal(s.data)
	logger.FailOnError("msgpack marshal failed", err)
	return bytes
}

func (s *MessagePackSerializer) Deserialize(raw []byte) {
	err := msgpack.Unmarshal(raw, s.data)
	logger.FailOnError("msgpack unmarshal failed", err)
}

func (s *MessagePackSerializer) CheckResult() {
	s.data.Compare(&data)
}
