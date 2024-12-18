package userService

import (
	"fmt"
	"github.com/mojocn/base64Captcha"
	"image/color"
	"strings"
)

var digitDriver = &base64Captcha.DriverString{
	Height:          32,
	Width:           100,
	NoiseCount:      0,
	ShowLineOptions: 2,
	Length:          4,
	Source:          "abcdefghijklmnopqrstuvwxyz",
	BgColor:         &color.RGBA{R: 3, G: 102, B: 214, A: 125},
	Fonts:           []string{"wqy-microhei.ttc"},
}

var store = base64Captcha.DefaultMemStore

// CaptchaGenerate 生成验证码
func CaptchaGenerate() (string, string, string, error) {
	driver := digitDriver.ConvertFonts()
	b := base64Captcha.NewCaptcha(driver, store)
	id, b64s, _, err := b.Generate()
	hcode := store.Get(id, false)
	if err != nil {
		fmt.Println("Error generating captcha:", err)
		return "", "", "", err
	}
	return id, b64s, hcode, nil
}

// GetCodeAnswer 验证验证码
func GetCodeAnswer(id, code string) bool {
	return store.Verify(id, strings.ToLower(code), false)
}
