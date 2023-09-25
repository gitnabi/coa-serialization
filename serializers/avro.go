package serializers

import (
	"serialization/internal/logger"

	"github.com/hamba/avro"
)

type AvroSerializer struct {
	data   *Data
	schema avro.Schema
}

func (s *AvroSerializer) PrepareData() {
	var err error
	s.schema, err = avro.Parse(`{
		"type": "record",
		"name": "serialization",
		"namespace": "org.hamba.avro",
		"fields" : [
			{"name": "Name", "type": "string"},
			{"name": "BankBalance", "type": "long"},
			{"name": "Weight", "type": "double"},
			{"name": "Clothes", "type": {"type":"array", "items": "string"}},
			{"name": "IsAlive", "type": "boolean"}
		]
	}`)
	logger.FailOnError("avro parse schema failed", err)

	s.data = &data
}

func (s *AvroSerializer) Serialize() []byte {
	bytes, err := avro.Marshal(s.schema, s.data)
	logger.FailOnError("avro marshal failed", err)
	return bytes
}

func (s *AvroSerializer) Deserialize(raw []byte) {
	err := avro.Unmarshal(s.schema, raw, s.data)
	logger.FailOnError("avro unmarshal failed", err)
}

func (s *AvroSerializer) CheckResult() {
	s.data.Compare(&data)
}
