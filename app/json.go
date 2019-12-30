// json.go kee > 2019/12/17

package app

import (
	"encoding/json"
	"fmt"
)

type Json map[string]interface{}

func (Json) Marshal(inter interface{}) (str []byte) {
	str, _ = json.Marshal(inter)
	return
}

func (d Json) Unmarshal(str []byte) Json {
	json.Unmarshal(str, &d)
	return d
}

func (d Json) Conver(inter interface{}) Json {
	return d.Unmarshal(d.Marshal(inter))
}

func (d Json) JsonUnmarshal(jStr string, out interface{}) {
	json.Unmarshal([]byte(jStr), &out)
}

func (d Json) ToStruct(out interface{}) {
	dStr := d.ToString()
	json.Unmarshal([]byte(dStr), &out)
}

func (d Json) ToString() string {
	return string(d.Marshal(d))
}

func (d Json) Get(name string) interface{} {
	v := d[name]
	typeof := fmt.Sprintf("%T", v)
	if "[]uint8" == typeof {
		return UbTS(v.([]uint8))
	}
	return v
}

// []uint8 to
func UbTS(bs []uint8) string {
	bt := []byte{}
	for _, b := range bs {
		bt = append(bt, byte(b))
	}
	return string(bt)
}
