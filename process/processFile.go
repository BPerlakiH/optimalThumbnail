/*
Copyright (c) 2016, Balazs Perlaki-Horvath <ba.perlaki@gmail.com>

Permission to use, copy, modify, and/or distribute this software for any purpose
with or without fee is hereby granted, provided that the above copyright notice
and this permission notice appear in all copies.

THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH
REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND
FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT,
INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS
OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER
TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF
THIS SOFTWARE.
*/


package process

import (
	"bufio"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"image/draw"
	"log"
	"os"
	"path/filepath"
	"github.com/BPerlakiH/optimalThumbnail/optimal"
	"github.com/chai2010/webp"
)

// import _ "image/jpeg"
// import _ "image/png"

// Resize and crop an image
func ProcessFile(inputFile string, outputFile string, width int, height int, quality int) {

	if inputFile == "" || outputFile == "" {
		return
	}

	imagefile, err := os.Open(inputFile)
	// TODO: implement a saner way to handle
	// open errors
	if err != nil {
		fmt.Printf("Open error \n" + inputFile)
		log.Println(err) // Don't use log.Fatal to exit
		return
	}

	
	// Decode the image.
	m, _, err := image.Decode(imagefile)
	imagefile.Close()

	if err != nil {
		fmt.Printf("Decode error \n")
		log.Println(err)
		return
	}

	b := m.Bounds()

	// All images are converted to the NRGBA type
	rgbaImage := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(rgbaImage, rgbaImage.Bounds(), m, b.Min, draw.Src)

	// Perform an optimal resize with 4 iterations  
	m2 := optimal.OptimalResize(rgbaImage, width, height, 4)

	fo, err := os.Create(outputFile)

	if err != nil {
		fmt.Printf("create file error \n")
		fo.Close()
		return
	}

	writer := bufio.NewWriter(fo)

	switch filepath.Ext(outputFile) {
		case ".png":
			png.Encode(writer, m2)
			break
		case ".webp":
			webp.Encode(writer, m2, &webp.Options{Lossless: false, Quality: float32(quality)})
			break
		default: //default to jpg
			jpeg.Encode(writer, m2, &jpeg.Options{Quality: quality})
			break
	}
	writer.Flush()
	fo.Close()
}