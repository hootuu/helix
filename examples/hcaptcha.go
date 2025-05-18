package main

import (
	"fmt"
	"github.com/hootuu/helix/components/hcaptcha"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/data/hjson"
	"math/rand/v2"
	"os"
	"time"
)

type User struct {
	Name string
	Age  int
}

func main() {
	d := dict.New(&User{
		Name: "ABC",
		Age:  10,
	})
	fmt.Println(hjson.MustToString(d))
	if true {
		return
	}
	os.Setenv("HELIX_CAPTCHA_CLEAN_BEFORE_DAYS", fmt.Sprintf("%d", 0))
	os.Setenv("HELIX_CAPTCHA_CLEAN_INTERVAL", fmt.Sprintf("%d", 10*time.Second))
	helix.AfterStartup(func() {
		//fmt.Println(hcaptcha.Generate(hcaptcha.NewOptions(hcaptcha.DigitsCaptcha, 6, 0)))
		//fmt.Println(hcaptcha.Generate(hcaptcha.NewOptions(hcaptcha.LetterCaptcha, 8, 0)))
		//fmt.Println(hcaptcha.Generate(hcaptcha.NewOptions(hcaptcha.LetterUpperCaptcha, 8, 0)))
		//fmt.Println(hcaptcha.Generate(hcaptcha.NewOptions(hcaptcha.MixtureCaptcha, 10, 0)))
		test := func() {
			link := fmt.Sprintf("user_check_%d_%d", time.Now().UnixMilli(), rand.Int64())
			capCode, err := hcaptcha.New(link, hcaptcha.NewOptions(hcaptcha.DigitsCaptcha, 6, 1*time.Second))
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("capCode: ", capCode)
			ok, err := hcaptcha.Verify(link, capCode)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("ok [will true]: ", ok)
			time.Sleep(time.Duration(rand.IntN(100)) * time.Millisecond)
			//ok, err = hcaptcha.Verify(link, capCode)
			//if err != nil {
			//	fmt.Println(err)
			//	return
			//}
			//fmt.Println("ok [will false]: ", ok)
		}
		for i := 0; i < 100; i++ {
			go func() {
				for j := 0; j < 1000; j++ {
					test()
				}
			}()
		}

	})
	helix.Startup()

}
