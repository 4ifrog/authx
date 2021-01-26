package redis

import (
	"bytes"
	"encoding/gob"
)

func gobEncodedBytes(val interface{}) (*bytes.Buffer, error) {
	// Use native gob encoding for the fastest serialization.
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(val); err != nil {
		return nil, err
	}

	return &buf, nil
}

func gobDecodedBytes(data []byte, val interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(val); err != nil {
		return err
	}

	return nil
}
