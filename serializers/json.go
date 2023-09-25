package serializers

import (
	"encoding/json"
	"serialization/internal/logger"
)

type JsonSerializer struct {
	data *Data
}

func (s *JsonSerializer) PrepareData() {
	s.data = &data
}

func (s *JsonSerializer) Serialize() []byte {
	bytes, err := json.Marshal(s.data)
	logger.FailOnError("json marshal failed", err)
	return bytes
}

func (s *JsonSerializer) Deserialize(raw []byte) {
	err := json.Unmarshal(raw, s.data)
	logger.FailOnError("json unmarshal failed", err)
}

func (s *JsonSerializer) CheckResult() {
	s.data.Compare(&data)
}
