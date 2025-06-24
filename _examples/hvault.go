package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/hootuu/helix/components/hvault"
	"github.com/hootuu/helix/helix"
)

func main() {

	helix.AfterStartup(func() {
		srcData := []byte("example.hvault")
		enData, err := hvault.Encrypt(srcData)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(hex.EncodeToString(enData))
		deData, err := hvault.Decrypt(enData)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(deData))
		if !bytes.Equal(srcData, deData) {
			fmt.Println("err: !bytes.Equal(srcData, deData)")
			return
		}
		pwd := []byte("12345678")
		pwdEnData, err := hvault.EncryptWithPwd(srcData, pwd)
		if err != nil {
			fmt.Println(err)
			return
		}
		pwdDeData, err := hvault.DecryptWithPwd(pwdEnData, pwd)
		if err != nil {
			fmt.Println(err)
			return
		}
		if !bytes.Equal(srcData, pwdDeData) {
			fmt.Println("err: !bytes.Equal(srcData, pwdDeData)")
			return
		}
		fmt.Println(string(pwdDeData))
	})
	helix.Startup()
}
