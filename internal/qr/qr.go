
package qr

import (
	"bytes"
	// "image"
	"image/color"
	"image/png"
	// "os"

	"github.com/fogleman/gg"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/skip2/go-qrcode"
	// "golang.org/x/image/draw"
)

func parseHexColor(s string) color.Color {
	c, _ := colorful.Hex(s)
	r, g, b := c.RGB255()
	return color.RGBA{r, g, b, 255}
}

func GenerateQRWithStyle(url string, size int, label string, colorStr string) ([]byte, error) {
	qrCode, err := qrcode.New(url, qrcode.Medium)
	if err != nil {
		return nil, err
	}
	qrCode.DisableBorder = false
	fgColor := parseHexColor(colorStr)
	qrCode.ForegroundColor = fgColor

	qrImg := qrCode.Image(size)

	labelHeight := 0
	if label != "" {
		labelHeight = 40 
	}

	margin := 20
	qrW := qrImg.Bounds().Dx()
	qrH := qrImg.Bounds().Dy()

	canvasW := qrW + margin*2
	canvasH := qrH + margin*2 + labelHeight

	dc := gg.NewContext(canvasW, canvasH)

	dc.SetColor(color.White)
	dc.Clear()

	dc.SetColor(fgColor)
	dc.SetLineWidth(8)
	dc.DrawRectangle(4, 4, float64(canvasW-8), float64(canvasH-8))
	dc.Stroke()

	dc.DrawImageAnchored(qrImg, canvasW/2, (canvasH/2)-2, 0.5, 0.5)

	if label != "" {
		if err := dc.LoadFontFace("D:/Projek/Go Gank/url-tools-be/assets/arial.ttf", 36); err == nil {
			r, g, b, _ := fgColor.RGBA()
			dc.SetRGB(float64(r)/65535, float64(g)/65535, float64(b)/65535)
			dc.DrawStringAnchored(label, float64(canvasW/2), float64(canvasH-40), 0.5, 0.5)
		}
	}

	var buf bytes.Buffer
	png.Encode(&buf, dc.Image())

	return buf.Bytes(), nil
}
