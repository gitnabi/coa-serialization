package serializers

import (
	"log"
	"math"
	"reflect"
)

type Data struct {
	Name        string
	BankBalance int64
	Weight      float64
	Clothes     []string
	IsAlive     bool
}

var data = Data{
	Name:        "name",
	BankBalance: math.MaxInt32,
	Weight:      100.500,
	Clothes:     []string{"qwer", "bebe", "1324", "", "|^_^|"},
	IsAlive:     false,
}

// func GetDataSize() int {
// 	result := int(unsafe.Sizeof(d))
// 	result += len(d.Name)
// 	for _, cloth := range d.Clothes {
// 		result += len(cloth)
// 	}
// 	return result
// }

func (d *Data) Compare(other *Data) {
	if d.Name != other.Name {
		log.Fatalf("d.Name != other.Name[%s != %s]", d.Name, other.Name)
	}
	if d.BankBalance != other.BankBalance {
		log.Fatalf("d.BankBalance != other.BankBalance[%d != %d]", d.BankBalance, other.BankBalance)
	}
	if d.Weight != other.Weight {
		log.Fatalf("d.Weight != other.Weight[%f != %f]", d.Weight, other.Weight)
	}
	if d.IsAlive != other.IsAlive {
		log.Fatalf("d.IsAlive != other.IsAlive[%t != %t]", d.IsAlive, other.IsAlive)
	}
	if !reflect.DeepEqual(d.Clothes, other.Clothes) {
		log.Fatalf("d.Clothes != other.Clothes[ %v != %v ]", d.Clothes, other.Clothes)
	}
}

type Serializer interface {
	PrepareData()
	Serialize() []byte
	Deserialize([]byte)
	CheckResult()
}
