package memser

import (
	"image/color"
	"math"
	"strings"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

type Request struct {
	BgImgPath string
	FontPath  string
	FontSize  float64
	Text      string
}

func TextOnImg(input string) (string, error) {
	request := Request{
		BgImgPath: "storage/notponyal.png",
		FontPath:  "",
		FontSize:  48.0,
		Text:      input,
	}

	bgImage, err := gg.LoadImage(request.BgImgPath)
	if err != nil {
		return "", err
	}
	imgW := bgImage.Bounds().Dx()
	imgH := bgImage.Bounds().Dy()

	dc := gg.NewContext(imgW, imgH)
	dc.DrawImage(bgImage, 0, 0)

	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return "", err
	}

	face := truetype.NewFace(font, &truetype.Options{Size: 48})
	dc.SetFontFace(face)

	// on ne ponyal part
	x := float64(imgW / 2)
	y := float64((imgH / 2) + 160.0)
	maxWidth := float64(imgW) - 60.0
	dc.SetColor(color.White)
	dc.DrawStringWrapped("он не понял", x, y, 0.5, 0.5, maxWidth, 1.5, gg.AlignCenter)
	y = float64((imgH / 2) - 260.0)
	input = strings.Join(dc.WordWrap(input, maxWidth), "\n")
	dc.DrawStringWrapped(input, x, y, 0.5, 0.5, maxWidth, 1.5, gg.AlignCenter)

	pathImg := "catmeme.png"
	if err := gg.SavePNG(pathImg, dc.Image()); err != nil {
		return "", err
	}

	return pathImg, nil

}

func HoldMeme(input string) (string, error) {
	request := Request{
		BgImgPath: "storage/hold.png",
		FontPath:  "",
		FontSize:  48.0,
		Text:      input,
	}

	bgImage, err := gg.LoadImage(request.BgImgPath)
	if err != nil {
		return "", err
	}
	imgW := bgImage.Bounds().Dx()
	imgH := bgImage.Bounds().Dy()

	dc := gg.NewContext(imgW, imgH)
	dc.DrawImage(bgImage, 0, 0)

	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return "", err
	}

	face := truetype.NewFace(font, &truetype.Options{Size: 48})

	dc.SetFontFace(face)

	// on ne ponyal part
	x := float64(imgW / 2)
	y := float64((imgH / 2) - 200.0)
	maxWidth := float64(imgW) - 60.0
	dc.SetColor(color.White)
	input = strings.Join(dc.WordWrap(input, maxWidth), "\n")
	dc.DrawStringWrapped(input, x, y, 0.5, 0.5, maxWidth, 1.5, gg.AlignCenter)

	dc.SetColor(color.Black)
	for i := 0; i < 100; i++ {
		angle := float64(i) / float64(1) * 2 * math.Pi
		xOffset := math.Cos(angle) * 4
		yOffset := math.Sin(angle) * 4
		dc.DrawStringWrapped(input, x+xOffset, y+yOffset, 0.5, 0.5, maxWidth, 1.5, gg.AlignCenter)
	}

	pathImg := "holdmeme.png"
	if err := gg.SavePNG(pathImg, dc.Image()); err != nil {
		return "", err
	}

	return pathImg, nil
}

func Choice(left, right, bottom string) (output string, err error) {
	req := Request{
		BgImgPath: "./assets/chose.png",
		FontPath:  "",
		FontSize:  49.0,
		Text:      "",
	}

	// load background image as gg object
	bgImage, err := gg.LoadImage(req.BgImgPath)
	if err != nil {
		panic(err)
	}
	// measure image for next steps
	imgH := bgImage.Bounds().Dy()
	imgW := bgImage.Bounds().Dx()

	dc := gg.NewContext(imgW, imgH)
	dc.DrawImage(bgImage, 0, 0)

	// set up font
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		panic(err)
	}

	// 3 different font sizes
	// TODO make font size scaling by text volume
	face := truetype.NewFace(font, &truetype.Options{Size: 35})
	dc.SetFontFace(face)
	faceSmall := truetype.NewFace(font, &truetype.Options{Size: 25})
	faceXSmall := truetype.NewFace(font, &truetype.Options{Size: 20})

	//      dc.Rotate(gg.Radians(-10))
	dc.SetRGB(0, 0, 0)
	s := bottom
	n := 3 // "stroke" size
	for dy := -n; dy <= n; dy++ {
		for dx := -n; dx <= n; dx++ {
			if dx*dx+dy*dy >= n*n {
				// give it rounded corners
				continue
			}
			x := float64(imgW)/2 + float64(dx)
			y := float64(imgH)*7/8 + float64(dy)
			dc.DrawStringWrapped(s, x, y, 0.5, 0.5, 370, 1.5, gg.AlignCenter)
		}
	}
	dc.SetRGB(1, 1, 1)
	dc.DrawStringWrapped(s, float64(imgW)/2, float64(imgH)*7/8, 0.5, 0.5, 370, 1.5, gg.AlignCenter)

	var1 := left
	var2 := right
	dc.SetFontFace(faceSmall)
	dc.Rotate(gg.Radians(-10))
	dc.SetRGB(0, 0, 0)
	x := float64(imgW)/4 - 20.0
	y := float64(imgH)/8 - 5.0
	dc.DrawStringWrapped(var1, x, y, 0.5, 0.5, 150, 1.5, gg.AlignCenter)
	x = float64(imgW)*3/4 - 90.0
	y += 20.0
	dc.SetFontFace(faceXSmall)
	dc.DrawStringWrapped(var2, x, y, 0.5, 0.5, 150, 1.5, gg.AlignCenter)
	// TODO make sane return and error handling
	dc.SavePNG("choice.png")

	return "choice.png", nil
}
