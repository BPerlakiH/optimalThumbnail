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
	"image/png"
	"image/draw"
	"log"
	"os"
	"path/filepath"
	"github.com/BPerlakiH/optimalThumbnail/optimal"
)

// Resize and crop an image
func ProcessFile(inputFile string, outputFile string, width int, height int) {

	// Only three image formats are supported (png, jpg, gif)
	if filepath.Ext(inputFile) == ".png" ||
		filepath.Ext(inputFile) == ".jpg" ||
		filepath.Ext(inputFile) == ".gif" {

		imagefile, err := os.Open(inputFile)
		// TODO: implement a saner way to handle
		// oprn errors
		if err != nil {
			fmt.Printf("Open error \n" + inputFile)
			log.Println(err) // Don't use log.Fatal to exit
			os.Exit(1)
		} else {
			// fmt.Printf("Openned: " + inputFile + "\n")
		}

		
		// Decode the image.
		m, _, err := image.Decode(imagefile)
		imagefile.Close()

		if err != nil {
			fmt.Printf("Decode error \n")
			log.Println(err)
			os.Exit(1)
		}

		b := m.Bounds()

		// All images are converted to the NRGBA type
		rgbaImage := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
		draw.Draw(rgbaImage, rgbaImage.Bounds(), m, b.Min, draw.Src)

		// Perform an optimal resize with 4 iterations  
		m2 := optimal.OptimalResize(rgbaImage, width, height, 8)

		fo, err := os.Create(outputFile)

		if err != nil {
			fmt.Printf("create file error \n")
			panic(err)
		}
		
		w := bufio.NewWriter(fo)
		png.Encode(w, m2)
		w.Flush()
		fo.Close()
	}
}