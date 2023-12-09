:: go get github.com/akavel/rsrc
:: go install github.com/akavel/rsrc
rsrc.exe -ico logo.ico -manifest manifest.xml -arch amd64 -o resource.syso
copy resource.syso ..\resource.syso
copy liblcl.dll .\build\liblcl.dll
go build -o ".\build\Rsa.exe" -ldflags "-X main.buildVersion=1.0.0 -H windowsgui" ..\ 