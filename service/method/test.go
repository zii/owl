package method

import (
	"owl/biz/android"
	"owl/service"
)

func Test(md *service.Meta) (interface{}, error) {
	id := md.Get("id").Int()
	if id == 1 {
		android.GetTLPlist()
	} else if id == 2 {
		android.TestLabelIcon()
	} else if id == 3 {
		android.TestLS()
	}
	return id, nil
}
