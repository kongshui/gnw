package msg

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
)

// 生成公钥和私钥
func RsaGenKey() (prvkey, pubkey []byte) {
	var (
		err error
	)

	prvPath := filepath.Join("./", "private_key.pem")
	pubPath := filepath.Join("./", "public_key.pem")
	// 生成私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println("私钥生成失败：", err)
		os.Exit(10)
	}

	// 将私钥编码为PEM格式
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPem := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	privateKeyFile, err := os.Create(prvPath)
	if err != nil {
		fmt.Println("私钥文件创建失败：", err)
		os.Exit(11)
	}
	pem.Encode(privateKeyFile, privateKeyPem)
	privateKeyFile.Close()

	// 生成公钥
	publicKey := privateKey.PublicKey

	// 将公钥编码为PEM格式
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		fmt.Println("公钥编码失败：", err)
		os.Exit(12)
	}
	publicKeyPem := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	publicKeyFile, err := os.Create(pubPath)
	if err != nil {
		fmt.Println("公钥文件创建失败：", err)
		os.Exit(13)
	}
	pem.Encode(publicKeyFile, publicKeyPem)
	publicKeyFile.Close()
	prvkey = pem.EncodeToMemory(privateKeyPem)
	pubkey = pem.EncodeToMemory(publicKeyPem)
	os.Chmod(prvPath, 0400)
	os.Chmod(pubPath, 0400)
	return
}
