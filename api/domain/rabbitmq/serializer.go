package rabbitmq

import (
	"encoding/json"

	jsoniter "github.com/json-iterator/go"
)

var encoder jsoniter.API

func initSerializer() {
	// similar to jsoniter.ConfigCompatibleWithStandardLibrary
	// adding TagKey
	encoder = jsoniter.Config{
		EscapeHTML:             true,
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
		TagKey:                 "amqp",
	}.Froze()
}

func serialize(input interface{}) ([]byte, error) {
	return json.Marshal(input)
	// return encoder.Marshal(input)
}

func unserialize(input []byte, output interface{}) error {
	return json.Unmarshal(input, output)
	// return encoder.Unmarshal(input, output)
}
