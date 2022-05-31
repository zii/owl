package method

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"owl/biz/android"
	"owl/service"

	"github.com/xuri/excelize/v2"
)

func Export(md *service.Meta, w http.ResponseWriter) {
	info := android.GetDeviceInfo()
	w.Header().Set("Content-Type", "application/vnd.ms-excel;charset=utf8")
	filename := "plist.xlsx"
	if info != nil {
		filename = fmt.Sprintf("%s.xlsx", info.Serial)
	}
	w.Header().Set("Content-Disposition", "attachment; filename="+url.QueryEscape(filename))

	plist, err := android.GetTLPlist()
	if err != nil {
		log.Println("Export err:", err)
		w.WriteHeader(500)
		return
	}

	f := excelize.NewFile()
	sh := "Sheet1"
	sheet := f.NewSheet(sh)
	f.SetActiveSheet(sheet)
	f.SetCellStr(sh, "A1", "应用名")
	f.SetCellStr(sh, "B1", "包名")
	f.SetCellStr(sh, "C1", "大小")
	f.SetCellStr(sh, "D1", "版本")
	f.SetCellStr(sh, "E1", "安装时间")
	f.SetCellStr(sh, "F1", "上次时间")
	f.SetCellStr(sh, "G1", "使用时长")
	f.SetCellStr(sh, "H1", "打开次数")
	for i, p := range plist.Plist {
		f.SetCellStr(sh, fmt.Sprintf("A%d", i+2), p.Label)
		f.SetCellStr(sh, fmt.Sprintf("B%d", i+2), p.Name)
		f.SetCellStr(sh, fmt.Sprintf("C%d", i+2), p.Size)
		f.SetCellStr(sh, fmt.Sprintf("D%d", i+2), p.Version)
		f.SetCellStr(sh, fmt.Sprintf("E%d", i+2), p.InstallTime)
		f.SetCellStr(sh, fmt.Sprintf("F%d", i+2), p.LastTimeUsed)
		f.SetCellStr(sh, fmt.Sprintf("G%d", i+2), p.TimeUsed)
		f.SetCellInt(sh, fmt.Sprintf("H%d", i+2), p.Launch)
	}
	f.SetSheetName("Sheet1", "已安装")
	sh = "已删除"
	sheet = f.NewSheet(sh)
	f.SetCellStr(sh, "A1", "已删除的包")
	for i, p := range plist.Uninstall {
		f.SetCellStr(sh, fmt.Sprintf("A%d", i+2), p)
	}
	f.Write(w)
}
