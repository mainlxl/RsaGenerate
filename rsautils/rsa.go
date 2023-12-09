package rsautils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

var BitsStr = "2048"

type OutLins interface {
	Add(S string) int32
}

func Generate() (*rsa.PrivateKey, error) {
	// 生成RSA密钥对
	atoi, err := strconv.Atoi(BitsStr)
	if err != nil {
		return nil, err
	}
	privateKey, err := rsa.GenerateKey(rand.Reader, atoi)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func ParsePrivateKey(inputText string) (*rsa.PrivateKey, error) {
	pemData, err := os.ReadFile(inputText)
	if err != nil {
		return nil, fmt.Errorf("读取PEM文件时发生错误：", err)
	}

	// 解码PEM块
	block, _ := pem.Decode(pemData)
	if err != nil {
		return nil, fmt.Errorf("无法解码PEM块:", err)
	}

	// 解析RSA私钥
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("PKCS#8&PKCS#1解码失败:", err)
		}
		// 将解析后的私钥对象转换为ecdsa.PrivateKey类型
		ecdsaKey, ok := privateKey.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("无法转换为ECDSA私钥类型")
		}
		return ecdsaKey, nil
	}
	return privateKey, nil

}

func OutputPem(privateKey *rsa.PrivateKey, profix string) (string, error) {

	// 获取公钥
	publicKey := &privateKey.PublicKey

	info, err := os.OpenFile(profix+"/秘钥详情.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
	defer info.Close()

	var outInfoLog = &strings.Builder{}
	d := privateKey.D.Text(16)
	fmt.Fprintf(outInfoLog, "私钥指数(16): %s\n", d)
	fmt.Fprintf(outInfoLog, "私钥长度: %d\n", len(d))
	n := privateKey.N.Text(16)
	fmt.Fprintf(outInfoLog, "\n模数(16): %s\n", n)
	fmt.Fprintf(outInfoLog, "模数长度: %d\n", len(n))
	e := privateKey.E
	fmt.Fprintf(outInfoLog, "\n公钥指数(16进制): %s\n\n", strconv.FormatInt(int64(e), 16))

	info.Write([]byte(outInfoLog.String()))

	// 将公钥转换为DER编码
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		fmt.Fprintf(outInfoLog, "将公钥转换为DER编码时发生错误：", err)
		return outInfoLog.String(), err
	}

	// 创建一个PEM块，将DER编码的公钥放入其中
	publicKeyBlock := pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}

	// 将PEM块写入文件
	publicKeyFile, err := os.Create(profix + "/publickey.pem")
	if err != nil {
		return outInfoLog.String(), fmt.Errorf("创建公钥文件时发生错误：", err)
	}

	defer publicKeyFile.Close()

	err = pem.Encode(publicKeyFile, &publicKeyBlock)
	if err != nil {
		return outInfoLog.String(), fmt.Errorf("写入公钥文件时发生错误：", err)
	}

	publicKeyBase64, _ := os.ReadFile(publicKeyFile.Name())
	fmt.Fprintf(info, "\n%s\n", string(publicKeyBase64))
	fmt.Fprintf(outInfoLog, "\n公钥已保存到%s文件中", publicKeyFile.Name())
	fmt.Fprintf(outInfoLog, "\n%s\n", string(publicKeyBase64))

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)

	privateKeyBlock := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	// 将PEM块写入文件
	privateKeyFile, err := os.Create(profix + "/privatekey.pem")
	if err != nil {
		return outInfoLog.String(), fmt.Errorf("创建公钥文件时发生错误：", err)
	}
	defer privateKeyFile.Close()

	err = pem.Encode(privateKeyFile, &privateKeyBlock)
	if err != nil {
		return outInfoLog.String(), fmt.Errorf("写入私钥文件时发生错误：", err)
	}

	privateKeyBase64, _ := os.ReadFile(privateKeyFile.Name())
	fmt.Fprintf(info, "\n%s\n", string(privateKeyBase64))
	fmt.Fprintf(outInfoLog, "\n私钥已保存到%s文件中", privateKeyFile.Name())
	fmt.Fprintf(outInfoLog, "\n%s\n", string(privateKeyBase64))

	if runtime.GOOS == "windows" {
		// 如果是Windows系统，将"\r\n"替换为"\n"
		processedString := strings.Replace(outInfoLog.String(), "\n", "\n\n", -1)
		return processedString, nil
	} else {
		// 对于其他系统，不做处理，直接输出原始字符串
		return outInfoLog.String(), nil
	}
}
