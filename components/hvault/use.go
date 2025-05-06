package hvault

import (
	"context"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/helix"
)

func Encrypt(src []byte) ([]byte, error) {
	return doEncrypt(src)
}

func Decrypt(src []byte) ([]byte, error) {
	return doDecrypt(src)
}

func DecryptWithPwd(src []byte, pwdBytes []byte) ([]byte, error) {
	return doDecryptWithPwd(src, pwdBytes)
}

func EncryptWithPwd(src []byte, pwdBytes []byte) ([]byte, error) {
	return doEncryptWithPwd(src, pwdBytes)
}

func init() {
	helix.Use(helix.BuildHelix(
		"hvault",
		func() (context.Context, error) {
			return nil, zplt.HelixPgDB().PG().AutoMigrate(&VaultM{})
		},
		func(ctx context.Context) {

		},
	))
}
