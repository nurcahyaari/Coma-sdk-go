package comasdkgo

import "encoding/json"

type Message struct {
	Data json.RawMessage `json:"data"`
}
