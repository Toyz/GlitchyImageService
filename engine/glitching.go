package engine

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"strings"

	glitch "github.com/sugoiuguu/go-glitch"
)

func gifImage(file io.Reader, expressions []string) (error, *bytes.Buffer, image.Rectangle) {
	var bounds image.Rectangle
	buff := new(bytes.Buffer)
	lGif, err := gif.DecodeAll(file)

	bounds = lGif.Image[0].Bounds()
	if err != nil {
		return err, nil, image.Rectangle{}
	}

	out := lGif
	for _, expression := range expressions {
		expr, err := glitch.CompileExpression(expression)
		if err != nil {
			return err, nil, bounds
		}

		newImage, err := expr.JumbleGIFPixels(out)
		if err != nil {
			out = nil
			return err, nil, bounds
		}
		out = newImage
		newImage = nil
	}

	err = gif.EncodeAll(buff, out)
	if err != nil {
		return err, nil, bounds
	}

	return nil, buff, bounds
}

func ProcessImage(file io.Reader, mime string, expressions []string) (error, *bytes.Buffer, image.Rectangle) {
	buff := new(bytes.Buffer)
	var bounds image.Rectangle
	switch strings.ToLower(mime) {
	case "image/gif":
		err, by, rect := gifImage(file, expressions)
		bounds = rect

		if err != nil {
			return err, nil, bounds
		}
		buff = by
		break
	default:
		img, _, err := image.Decode(file)
		if err != nil {
			return err, nil, image.Rectangle{}
		}
		bounds = img.Bounds()

		out := img
		for _, expression := range expressions {
			expr, err := glitch.CompileExpression(expression)
			if err != nil {
				return err, nil, bounds
			}

			newImage, err := expr.JumblePixels(out)
			if err != nil {
				out = nil
				return err, nil, bounds
			}
			out = newImage
			newImage = nil
		}

		switch strings.ToLower(mime) {
		case "image/png":
			png.Encode(buff, out)
			break
		case "image/jpg", "image/jpeg":
			jpeg.Encode(buff, out, nil)
			break
		}
	}

	return nil, buff, bounds
}
