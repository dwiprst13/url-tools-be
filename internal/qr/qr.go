package qr

import (
	"bytes"
	"image"
	"image/png"
	"os"

	"github.com/skip2/go-qrcode"
	"golang.org/x/image/draw"
)

func GenerateQRWithLogo(url string, size int, logoPath string) ([]byte, error) {
	qrCode, err := qrcode.New(url, qrcode.Medium)
	if err != nil {
		return nil, err
	}
	qrCode.DisableBorder = false
	qrImg := qrCode.Image(size)

	logoFile, err := os.Open(logoPath)
	if err != nil {
		return nil, err
	}
	defer logoFile.Close()

	logoImg, _, err := image.Decode(logoFile)
	if err != nil {
		return nil, err
	}

	logoBounds := logoImg.Bounds()
	logoW := logoBounds.Dx()
	logoH := logoBounds.Dy()
	aspect := float64(logoW) / float64(logoH)

	maxLogoW := size / 3
	targetW := maxLogoW
	targetH := int(float64(maxLogoW) / aspect)

	dstLogo := image.NewRGBA(image.Rect(0, 0, targetW, targetH))
	draw.NearestNeighbor.Scale(dstLogo, dstLogo.Bounds(), logoImg, logoBounds, draw.Over, nil)

	offset := image.Pt((qrImg.Bounds().Dx()-targetW)/2, (qrImg.Bounds().Dy()-targetH)/2)
	outImg := image.NewRGBA(qrImg.Bounds())
	draw.Draw(outImg, qrImg.Bounds(), qrImg, image.Point{}, draw.Src)
	draw.Draw(outImg, dstLogo.Bounds().Add(offset), dstLogo, image.Point{}, draw.Over)

	var buf bytes.Buffer
	encoder := png.Encoder{CompressionLevel: png.BestCompression}
	if err := encoder.Encode(&buf, outImg); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
