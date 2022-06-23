package services

import (
	"encoding/json"
	"fmt"
)

func Println(data map[string]interface{}) {
	bs, _ := json.Marshal(data)
	fmt.Println(string(bs))
}
