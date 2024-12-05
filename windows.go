package main

import (
	"os"

	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
)

var mainWindow *walk.MainWindow

var mainWindowWidth = 300
var mainWindowHeight = 200

func MenuBarInit() []MenuItem {
	return []MenuItem{
		Action{
			Text: "Runlog",
			OnTriggered: func() {
				OpenBrowserWeb(RunlogDirGet())
			},
		},
		Action{
			Text: "Mini Windows",
			OnTriggered: func() {
				NotifyAction()
			},
		},
		Action{
			Text: "Sponsor",
			OnTriggered: func() {
				AboutAction()
			},
		},
	}
}

var intervalNumber *walk.NumberEdit
var filterInterface *walk.ComboBox
var connectivityURL, outputFolder *walk.LineEdit

func ConsoleWidget() []Widget {
	interfaceList := InterfaceOptions()

	return []Widget{
		Label{
			Text: "Connectivity URL: ",
		},
		LineEdit{
			AssignTo: &connectivityURL,
			Text:     ConfigGet().ConnectivityURL,
			OnEditingFinished: func() {
				ConnectivityURLSave(connectivityURL.Text())
			},
		},
		Label{
			Text: "Monitor Interface: ",
		},
		ComboBox{
			AssignTo: &filterInterface,
			Model:    interfaceList,
			CurrentIndex: func() int {
				filter := ConfigGet().FilterInterface
				for i, name := range interfaceList {
					if name == filter {
						return i
					}
				}
				return 0
			},
			OnCurrentIndexChanged: func() {
				err := FilterInterfaceSave(filterInterface.Text())
				if err != nil {
					ErrorBoxAction(mainWindow, err.Error())
				}
			},
			OnBoundsChanged: func() {
				addr := ConfigGet().FilterInterface
				for i, item := range interfaceList {
					if addr == item {
						filterInterface.SetCurrentIndex(i)
						return
					}
				}
				filterInterface.SetCurrentIndex(0)
			},
		},
		Label{
			Text: "Output Folder: ",
		},
		Composite{
			Layout: HBox{MarginsZero: true},
			Children: []Widget{
				LineEdit{
					AssignTo: &outputFolder,
					Text:     ConfigGet().OutputDirectory,
					OnEditingFinished: func() {
						dir := outputFolder.Text()
						if dir != "" {
							stat, err := os.Stat(dir)
							if err != nil {
								ErrorBoxAction(mainWindow, "The server folder is not exist")
								outputFolder.SetText("")
								OutputDirectorySave("")
								return
							}
							if !stat.IsDir() {
								ErrorBoxAction(mainWindow, "The server folder is not directory")
								outputFolder.SetText("")
								OutputDirectorySave("")
								return
							}
						}
						OutputDirectorySave(dir)
					},
				},
				PushButton{
					MaxSize: Size{Width: 30},
					Text:    " ... ",
					OnClicked: func() {
						dlgDir := new(walk.FileDialog)
						dlgDir.FilePath = ConfigGet().OutputDirectory
						dlgDir.Flags = win.OFN_EXPLORER
						dlgDir.Title = "Please select a folder as output file directory"

						exist, err := dlgDir.ShowBrowseFolder(mainWindow)
						if err != nil {
							logs.Error(err.Error())
							return
						}
						if exist {
							logs.Info("select %s as output file directory", dlgDir.FilePath)
							outputFolder.SetText(dlgDir.FilePath)
							OutputDirectorySave(dlgDir.FilePath)
						}
					},
				},
			},
		},
		// Label{
		// 	Text: "Listen Address: ",
		// },
		// ComboBox{
		// 	AssignTo: &listenAddr,
		// 	CurrentIndex: func() int {
		// 		addr := ConfigGet().ListenAddr
		// 		for i, item := range interfaceList {
		// 			if addr == item {
		// 				return i
		// 			}
		// 		}
		// 		return 0
		// 	},
		// 	Model: interfaceList,
		// 	OnCurrentIndexChanged: func() {
		// 		err := ListenAddressSave(listenAddr.Text())
		// 		if err != nil {
		// 			ErrorBoxAction(mainWindow, err.Error())
		// 		} else {
		// 			BrowseURLUpdate()
		// 		}
		// 	},
		// 	OnBoundsChanged: func() {
		// 		addr := ConfigGet().ListenAddr
		// 		for i, item := range interfaceList {
		// 			if addr == item {
		// 				listenAddr.SetCurrentIndex(i)
		// 				return
		// 			}
		// 		}
		// 		listenAddr.SetCurrentIndex(0)
		// 	},
		// },
		Label{
			Text: "Loop Monitor Interval: ",
		},
		NumberEdit{
			AssignTo:    &intervalNumber,
			Value:       float64(ConfigGet().Interval),
			ToolTipText: "1~300s",
			MaxValue:    300,
			MinValue:    1,
			OnValueChanged: func() {
				err := IntervalSave(int(intervalNumber.Value()))
				if err != nil {
					ErrorBoxAction(mainWindow, err.Error())
				}
			},
		},
	}
}

func mainWindows() {
	CapSignal(CloseWindows)
	cnt, err := MainWindow{
		Title:          "Simple IP Monitor Windows " + VersionGet(),
		Icon:           ICON_Main,
		AssignTo:       &mainWindow,
		MinSize:        Size{Width: mainWindowWidth, Height: mainWindowHeight},
		Size:           Size{Width: mainWindowWidth, Height: mainWindowHeight},
		Layout:         VBox{Margins: Margins{Top: 5, Bottom: 5, Left: 5, Right: 5}},
		Font:           Font{Bold: true},
		MenuItems:      MenuBarInit(),
		StatusBarItems: StatusBarInit(),
		Children: []Widget{
			Composite{
				Layout:   Grid{Columns: 2},
				Children: ConsoleWidget(),
			},
		},
	}.Run()

	if err != nil {
		logs.Error(err.Error())
	} else {
		logs.Info("main windows exit %d", cnt)
	}

	if err := recover(); err != nil {
		logs.Error(err)
	}

	CloseWindows()
}

func CloseWindows() {
	if mainWindow != nil {
		mainWindow.Close()
		mainWindow = nil
	}
	NotifyExit()
}
