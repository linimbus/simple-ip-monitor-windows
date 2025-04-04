package main

import (
	"os"

	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
)

func IconLoadFromBox(filename string, size walk.Size) *walk.Icon {
	body, err := Asset(filename)
	if err != nil {
		logs.Error(err.Error())
		return walk.IconApplication()
	}
	dir := DEFAULT_HOME + "\\icon\\"
	_, err = os.Stat(dir)
	if err != nil {
		err = os.MkdirAll(dir, 644)
		if err != nil {
			logs.Error(err.Error())
			return walk.IconApplication()
		}
	}
	filepath := dir + filename
	err = SaveToFile(filepath, body)
	if err != nil {
		logs.Error(err.Error())
		return walk.IconApplication()
	}
	icon, err := walk.NewIconFromFileWithSize(filepath, size)
	if err != nil {
		logs.Error(err.Error())
		return walk.IconApplication()
	}
	return icon
}

var ICON_Main *walk.Icon
var ICON_Status *walk.Icon

var ICON_Max_Size = walk.Size{
	Width: 128, Height: 128,
}

var ICON_Min_Size = walk.Size{
	Width: 16, Height: 16,
}

func IconInit() {
	ICON_Main = IconLoadFromBox("main.ico", ICON_Max_Size)
	ICON_Status = IconLoadFromBox("status.ico", ICON_Min_Size)
}
