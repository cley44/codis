package utils

import (
	"encoding/json"
	"fmt"
)

func PrintJSON(in ...interface{}) {
	output := []interface{}{}

	for _, item := range in {
		j, err := json.Marshal(item)
		if err != nil {
			panic(err)
		}

		output = append(output, string(j))
	}

	fmt.Println(output...) //nolint:forbidigo
}

func PrintJSONIndent(in ...interface{}) {
	output := []interface{}{}

	for _, item := range in {
		j, err := json.MarshalIndent(item, "", "  ")
		if err != nil {
			panic(err)
		}

		output = append(output, string(j))
	}

	fmt.Println(output...) //nolint:forbidigo
}
