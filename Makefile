BINARY ?= RsaGenerator
VERSION ?= 1.0.0
DEFAULT_PORT ?= 8082


mac-app:
	@go build -gcflags=all=-trimpath="$GOPATH" -ldflags "-X main.buildVersion=$(VERSION) -s -w" -o "./mac/Rsa秘钥生成器.app/Contents/MacOS/RsaGenerator"


win-x86_64.exe:
	GOARCH=amd64 CGO_ENABLED=0 GOOS=windows go build -ldflags "-X main.buildVersion=$(VERSION) -H windowsgui"


win-386.exe:
	GOARCH=386 CGO_ENABLED=0  GOOS=windows go build -ldflags "-X main.buildVersion=$(VERSION) -H windowsgui"
