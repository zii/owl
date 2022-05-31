package base

import (
	"bytes"
	"encoding/json"
)

// json.Marshal和json.Unmarshal处理int64有bug
func JsonEncode(v interface{}) ([]byte, error) {
	b := new(bytes.Buffer)
	e := json.NewEncoder(b)

	// https://studygolang.com/articles/14473?fr=sidebar
	e.SetEscapeHTML(false) // fix: encoding 时可能会自动追加 '\n'
	err := e.Encode(v)
	if err != nil {
		return nil, err
	}

	blob := b.Bytes()
	if len(blob) > 0 && blob[len(blob)-1] == '\n' {
		blob = blob[:len(blob)-1]
	}
	return blob, nil
}

func MustJsonEncode(v interface{}) []byte {
	b, e := JsonEncode(v)
	Raise(e)
	return b
}

func JsonDecode(data []byte, v interface{}) error {
	e := json.NewDecoder(bytes.NewBuffer(data))
	e.UseNumber()
	if err := e.Decode(v); err != nil {
		return err
	}

	return nil
}

func MustJsonDecode(data []byte, v interface{}) {
	e := json.NewDecoder(bytes.NewBuffer(data))
	e.UseNumber()
	if err := e.Decode(v); err != nil {
		Raise(err)
	}
}

func JsonString(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
