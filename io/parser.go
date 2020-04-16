package io

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
)

func ToJson(i interface{}) ([]byte, error) {
	return json.Marshal(i)
}

func ToYaml(i interface{}) ([]byte, error) {
	return yaml.Marshal(i)
}

func FromJson(b []byte,i interface{}) error {
	return json.Unmarshal(b, i)
}

func FromYaml(b []byte,i interface{}) error {
	return yaml.Unmarshal(b, i)
}

