package functions

import (
	"encoding/json"
	"strings"
)

func MaskingDataFromBytes(from []byte, maskingString string) (to []byte) {

	var f interface{}
	json.Unmarshal(from, &f)

	v, ok := f.(map[string]interface{})
	if ok {

		sfields := strings.Split(maskingString, ",")
		for _, sfield := range sfields {

			layers := strings.Split(sfield, ".")
			temp := &v

			for i, layer := range layers {

				if i == len(layers)-1 {
					_, exist := (*temp)[layer]
					if exist {
						(*temp)[layer] = "-MASKED-"
					}
				} else {
					_, exist := (*temp)[layer]
					if exist {
						temp2, okk := (*temp)[layer].(map[string]interface{})
						if okk {
							temp = &temp2
						} else {
							break
						}
					} else {
						break
					}
				}
			}
		}
	}

	rb, _ := json.Marshal(v)
	return rb
}
