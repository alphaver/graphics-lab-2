package winapi

import (
	"errors"
	"github.com/lxn/win"
	"strings"
)

const (
	BIF_RETURNONLYFSDIRS = 0x01
	BIF_EDITBOX          = 0x10
	BIF_NEWDIALOGSTYLE   = 0x40
)

func MessageBox(message string, caption string, flags uint32) {
	sysMessage := win.SysAllocString(message)
	defer win.SysFreeString(sysMessage)
	sysCaption := win.SysAllocString(caption)
	defer win.SysFreeString(sysCaption)

	win.MessageBox(0, sysMessage, sysCaption, flags)
}

func OpenDirectory() (dirPath string, err error) {
	selectedDir := win.SysAllocString(strings.Repeat(" ", win.MAX_PATH))
	defer win.SysFreeString(selectedDir)
	title := win.SysAllocString("Browse for folders")
	defer win.SysFreeString(title)

	browseInfo := &win.BROWSEINFO{
		HwndOwner: win.HWND(0),
		PidlRoot:  0,
		LpszTitle: title,
		UlFlags:   BIF_RETURNONLYFSDIRS | BIF_EDITBOX | BIF_NEWDIALOGSTYLE,
	}

	result := win.SHBrowseForFolder(browseInfo)
	if result == 0 {
		return "", nil
	}
	if !win.SHGetPathFromIDList(result, selectedDir) {
		return "", errors.New("can't get the directory path")
	}
	win.CoTaskMemFree(result)
	return win.UTF16PtrToString(selectedDir), nil
}