package hcaptcha

import (
	"fmt"
	"math/rand"
	"strings"
)

const (
	gDigitLib       = "023456789"
	gUpperLetterLib = "ABCDEFGHIJKLMNPQRSTUVWXYZ"
	gLowerLetterLib = "abcdefghijkmnpqrstuvwxyz"
	gLetterLib      = gUpperLetterLib + gLowerLetterLib
	gMixLib         = gLetterLib + gDigitLib
)

func Generate(opt *Options) Captcha {
	lib := ""
	switch opt.Type {
	case DigitsCaptcha:
		lib = gDigitLib
	case LetterCaptcha:
		lib = gLetterLib
	case LetterUpperCaptcha:
		lib = gUpperLetterLib
	case LetterLowerCaptcha:
		lib = gLowerLetterLib
	case MixtureCaptcha:
		lib = gMixLib
	default:
		return NilCaptcha
	}
	libLen := uint32(len(lib))
	var buf strings.Builder
	for i := 0; i < opt.Length; i++ {
		x := rand.Uint32()
		n := x % libLen
		d := lib[n]
		buf.WriteString(fmt.Sprintf("%c", d))
	}
	return buf.String()
}
