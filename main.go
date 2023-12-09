// 由res2go IDE插件自动生成。
package main

import (
	"RsaGenerator/rsagui"
	"fmt"
	"github.com/ying32/govcl/vcl"
)

var buildVersion string = "1.0.0"

func main() {
	vcl.Application.SetScaled(true)
	vcl.Application.SetTitle("RsaGenerate")
	vcl.Application.Initialize()
	vcl.Application.SetMainFormOnTaskBar(true)
	vcl.Application.CreateForm(&rsagui.MainWindow)
	rsagui.MainWindow.Label2.SetCaption(fmt.Sprintf("作者: Mainli 版本: %s", buildVersion))
	vcl.Application.Run()
}
