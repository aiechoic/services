package encoding

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
)

var (
	GobSerializer  Serializer = &gobSerializer{}
	JSONSerializer Serializer = &jsonSerializer{}
)

type Serializer interface {
	Serialize(v any) ([]byte, error)
	Deserialize(data []byte, v any) error
}

type jsonSerializer struct{}

func (j *jsonSerializer) Serialize(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (j *jsonSerializer) Deserialize(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

type gobSerializer struct{}

func (g *gobSerializer) Serialize(v any) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(v)
	return buf.Bytes(), err
}

func (g *gobSerializer) Deserialize(data []byte, v any) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(v)
}
