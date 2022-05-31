package method

import (
	"owl/biz/android"
	"owl/service"
)

func Device_plist(md *service.Meta) (interface{}, error) {
	out, err := android.GetTLPlist()
	if err != nil {
		return nil, service.NewError(400, "PLIST_INVALID")
	}
	return out, nil
}
