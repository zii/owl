package method

import (
	"owl/biz/android"
	"owl/service"
)

func Device_info(md *service.Meta) (interface{}, error) {
	out, err := android.GetTLDevice()
	if err != nil {
		return nil, service.NewError(400, "DEVICE_INVALID")
	}
	return out, nil
}
