package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
)

const pngSide = 16

// pngPalette цвета для генерации визуально различных изображений
var pngPalette = []color.RGBA{
	{R: 220, G: 60, B: 60, A: 255},
	{R: 60, G: 200, B: 90, A: 255},
	{R: 70, G: 90, B: 220, A: 255},
	{R: 230, G: 190, B: 40, A: 255},
}

// GeneratePNG создаёт корректное одноцветное PNG-изображение и возвращает его байты;
// цвет выбирается по индексу, поэтому разные индексы дают разные файлы
func GeneratePNG(index int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, pngSide, pngSide))
	fill := pngPalette[index%len(pngPalette)]
	draw.Draw(img, img.Bounds(), &image.Uniform{C: fill}, image.Point{}, draw.Src)

	var buffer bytes.Buffer
	if err := png.Encode(&buffer, img); err != nil {
		panic(fmt.Sprintf("Не удалось сгенерировать PNG: %s", err.Error()))
	}

	return buffer.Bytes()
}
