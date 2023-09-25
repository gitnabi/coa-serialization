package serializers

import (
	"serialization/internal/logger"

	"encoding/xml"
)

type XmlSerializer struct {
	data *Data
}

func (s *XmlSerializer) PrepareData() {
	s.data = &data
}

func (s *XmlSerializer) Serialize() []byte {
	bytes, err := xml.Marshal(s.data)
	logger.FailOnError("xml marshal failed", err)
	return bytes
}

func (s *XmlSerializer) Deserialize(raw []byte) {
	tmp := Data{}
	// not everything is erased, some values are appended to current ones
	err := xml.Unmarshal(raw, &tmp)
	logger.FailOnError("xml unmarshal failed", err)
	s.data = &tmp
}

func (s *XmlSerializer) CheckResult() {
	s.data.Compare(&data)
}
