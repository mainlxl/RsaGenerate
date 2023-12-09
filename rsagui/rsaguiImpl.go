package rsagui

import (
	"RsaGenerator/rsautils"
	"errors"
	"fmt"
	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
	"os"
	"os/user"
	"path"
)

//::private::
type TMainWindowFields struct {
}

func (f *TMainWindow) OnFormCreate(sender vcl.IObject) {
	f.OutText.ReadOnly()
	u, _ := user.Current()
	f.OutPath.SetText(path.Join(u.HomeDir, "Desktop"))
}

func (f *TMainWindow) OnPathSelectBtnClick(sender vcl.IObject) {
	dialog := vcl.NewSelectDirectoryDialog(f.TForm)
	dialog.SetOptions(types.NewSet(types.OfAllowMultiSelect, types.OfFileMustExist, types.OfPathMustExist, types.OfShowHelp))
	dialog.SetFilter("所有文件|*.pem")
	dialog.SetTitle("请选择输出目录")
	dialog.SetInitialDir(f.OutPath.Text())
	if dialog.Execute() {
		vcl.ShowMessage(fmt.Sprintf("私钥路径: %s", dialog.FileName()))
	}
}

//通过私钥获取公钥
func (f *TMainWindow) OnObtainPrivateInfoBtnClick(sender vcl.IObject) {
	dialog := vcl.NewOpenDialog(f.TForm)
	dialog.SetOptions(types.NewSet(types.OfAllowMultiSelect, types.OfFileMustExist, types.OfPathMustExist, types.OfShowHelp))
	dialog.SetFilter("所有文件|*.pem")
	dialog.SetTitle("请选择私钥")
	dialog.SetInitialDir(f.OutPath.Text())
	if dialog.Execute() {
		key, err := rsautils.ParsePrivateKey(dialog.FileName())
		if err != nil {
			vcl.ShowMessage(err.Error())
			return
		} else {
			path, err := ensurePath(f.OutPath.Text())
			if err != nil {
				vcl.ShowMessage(err.Error())
				return
			}
			info, err := rsautils.OutputPem(key, path)
			f.OutText.Lines().Clear()
			f.OutText.Lines().Add(info)
			f.OutText.SetSelStart(0)
			if err != nil {
				vcl.ShowMessage(err.Error())
			}
		}
	}
}
func (f *TMainWindow) OnGenerateBtnClick(sender vcl.IObject) {
	path, err := ensurePath(f.OutPath.Text())
	if err != nil {
		vcl.ShowMessage(err.Error())
		return
	}
	key, err := rsautils.Generate()
	if err != nil {
		vcl.ShowMessage(err.Error())
		return
	} else {
		info, err := rsautils.OutputPem(key, path)
		f.OutText.Lines().Clear()
		f.OutText.Lines().Add(info)
		f.OutText.SetSelStart(0)
		if err != nil {
			vcl.ShowMessage(err.Error())
		}
	}
}

func ensurePath(path string) (string, error) {
	info, err := os.Stat(path)
	if err == nil {
		if info.IsDir() {
			return path, nil
		} else {
			return "", errors.New("输出路径不是文件夹")
		}
	} else if os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755) // 创建路径
		if err != nil {
			return "", fmt.Errorf("创建路径失败:", err)
		} else {
			return path, nil
		}
	} else {
		return path, fmt.Errorf("发生错误:", err)
	}
}

func (f *TMainWindow) OnOutTextChange(sender vcl.IObject) {

}
