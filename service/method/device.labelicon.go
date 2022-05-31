package method

import (
	"owl/biz/android"
	"owl/biz/proto"
	"owl/service"
)

func Device_labelicon(md *service.Meta) (interface{}, error) {
	pname := md.Get("name").String()
	label, icon := android.GetLabelIcon(pname)
	out := &proto.LabelIcon{
		Label: label,
		Icon:  icon,
	}
	return out, nil
}
