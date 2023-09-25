package serializers

import (
	"serialization/internal/logger"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

type ProtobufSerializer struct {
	data *structpb.Struct
}

func (s *ProtobufSerializer) PrepareData() {
	values := []*structpb.Value{}
	for _, cloth := range data.Clothes {
		value := &structpb.Value{
			Kind: &structpb.Value_StringValue{
				StringValue: cloth,
			},
		}
		values = append(values, value)
	}
	s.data = &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"name": {
				Kind: &structpb.Value_StringValue{
					StringValue: data.Name,
				},
			},
			"bank_balance": {
				Kind: &structpb.Value_NumberValue{
					NumberValue: float64(data.BankBalance),
				},
			},
			"weight": {
				Kind: &structpb.Value_NumberValue{
					NumberValue: data.Weight,
				},
			},
			"clothes": {
				Kind: &structpb.Value_ListValue{
					ListValue: &structpb.ListValue{
						Values: values,
					},
				},
			},
			"is_alive": {
				Kind: &structpb.Value_BoolValue{
					BoolValue: data.IsAlive,
				},
			},
		},
	}
}

func (s *ProtobufSerializer) Serialize() []byte {
	bytes, err := proto.Marshal(s.data)
	logger.FailOnError("protobuf marshal failed", err)
	return bytes
}

func (s *ProtobufSerializer) Deserialize(raw []byte) {
	err := proto.Unmarshal(raw, s.data)
	logger.FailOnError("protobuf unmarshal failed", err)
}

func (s *ProtobufSerializer) CheckResult() {
	clothes := []string{}
	for _, value := range s.data.Fields["clothes"].GetListValue().Values {
		clothes = append(clothes, value.GetStringValue())
	}

	data.Compare(&Data{
		Name:        s.data.Fields["name"].GetStringValue(),
		BankBalance: int64(s.data.Fields["bank_balance"].GetNumberValue()),
		Weight:      s.data.Fields["weight"].GetNumberValue(),
		Clothes:     clothes,
		IsAlive:     s.data.Fields["is_alive"].GetBoolValue(),
	})
}
