package android

import (
	"strconv"
	"strings"
)

type PropValue string

func (p PropValue) Int() int {
	i, _ := strconv.Atoi(string(p))
	return i
}

func (p PropValue) String() string {
	return string(p)
}

type Props map[string]PropValue

func (ps Props) Get(k string) PropValue {
	return ps[k]
}

func (ps Props) Operator() string {
	var keys = []string{"gsm.operator.alpha", "gsm.operator.alpha0", "gsm.operator.alpha1"}
	for _, k := range keys {
		s := ps.Get(k).String()
		s = strings.Trim(s, ",")
		if s != "" {
			return s
		}
	}
	return ""
}
