package proto

import "fmt"

type Success struct {
	Ok     bool        `json:"ok"`
	Result interface{} `json:"result"`
}

type Fail struct {
	Ok          bool        `json:"ok"`
	ErrorCode   int         `json:"error_code"`
	Description string      `json:"description"`
	Data        interface{} `json:"data,omitempty"`
}

func (f *Fail) Error() string {
	return fmt.Sprintf("%d %s", f.ErrorCode, f.Description)
}
