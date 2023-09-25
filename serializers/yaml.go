package serializers

import (
	"serialization/internal/logger"

	yaml "gopkg.in/yaml.v3"
)

type YamlSerializer struct {
	data *Data
}

func (s *YamlSerializer) PrepareData() {
	s.data = &data
}

func (s *YamlSerializer) Serialize() []byte {
	bytes, err := yaml.Marshal(s.data)
	logger.FailOnError("yaml marshal failed", err)
	return bytes
}

func (s *YamlSerializer) Deserialize(raw []byte) {
	err := yaml.Unmarshal(raw, s.data)
	logger.FailOnError("yaml unmarshal failed", err)
}

func (s *YamlSerializer) CheckResult() {
	s.data.Compare(&data)
}
