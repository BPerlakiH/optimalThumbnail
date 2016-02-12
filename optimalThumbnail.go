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

package main

import (
	"fmt"
	"flag"
	"log"
	"os"
	"path/filepath"
	"time"
	"math"
	"github.com/BPerlakiH/optimalThumbnail/process"
)

import _ "image/jpeg"
import _ "image/gif"

func main() {

	start_time := time.Now()

	// Read the CMD options
	inDir := flag.String("in", "", "input directory")    // input directory
	outDir := flag.String("out", "", "output directory") // output directory
	width := flag.Int("width", 128, "the new width")     // width
	height := flag.Int("height", 128, "the new height")  // height

	flag.Parse()

	if *inDir == "" || *outDir == "" {
		log.Fatal("usage: \n imageResizer -in inputDir -out outputDir -width 128 -height 128")
	}

	// Print the cmd options

	fmt.Printf("image resize daemon \n")

	fmt.Printf("Input:  %s \n", *inDir)
	fmt.Printf("Output: %s \n", *outDir)
	fmt.Printf("Width:  %d \n", *width)
	fmt.Printf("Height: %d \n", *height)

	d, err := os.Open(*inDir)
	if err != nil {
		fmt.Println(err)
		d.Close()
		os.Exit(1)
	}

	fi, err := d.Readdir(-1)
	d.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	
	//get the processable files:
	var files []os.FileInfo
	for _, file := range fi {
		ext := filepath.Ext(file.Name())
		if !file.IsDir() && (ext == ".png" || ext == ".jpg" || ext == ".gif") {
			files = append(files, file)
		}
	}
	files_count := len(files)


	//BRAKE DOWN TO CHANELS
	channel_count := 12
	
	if channel_count < files_count {
		//go parallel on the channels:
		channels := make(chan int, channel_count)
		files_per_channel := int(math.Floor(float64(files_count) / float64(channel_count)))

		for i := 0; i < channel_count; i++ {
			go func(i int) {
				min := i * files_per_channel
				max := int(math.Min(float64((i+1) * files_per_channel), float64(files_count)))
				for k := min; k < max; k++ {
					filename := files[k].Name()
					percent := float64(k) / float64(files_count)
					fmt.Printf("\r%v %v/%v (%v%%)", filename, k, files_count, int(percent * 100))
					fullImagePath := filepath.Join(*inDir, filename)
					fullImageOutPath := filepath.Join(*outDir, filepath.Base(filename))
					process.ProcessFile(fullImagePath, fullImageOutPath, *width, *height)
				}
				//open the channel
				channels <- 1
			}(i)
			//close the channel
			<-channels
		}
	} else {
		//go parallel on all files at the same time:
		channels := make(chan int, files_count)
		for index := 0; index < files_count; index++ {
			filename := files[index].Name()
			fmt.Println(filename)
			go func(filename string) {
				// Combine the directory path with the filename to get 
				// the full path of the image
				fullImagePath := filepath.Join(*inDir, filename)
				fullImageOutPath := filepath.Join(*outDir, filepath.Base(filename))

				process.ProcessFile(fullImagePath, fullImageOutPath, *width, *height)
				//open the channel
				channels <- 1
			}(filename)
		}
		//close all channels
		for i := 0; i < files_count; i++ {
			<-channels
		}
	}

	//LINEAR ONE
	// for index := 0; index < files_count-1; index++ {
	// 	filename := files[index].Name()
	// 	// Combine the directory path with the filename to get 
	// 	// the full path of the image
	// 	fullImagePath := filepath.Join(*inDir, filename)
	// 	fullImageOutPath := filepath.Join(*outDir, filepath.Base(filename))
	// 	process.ProcessFile(fullImagePath, fullImageOutPath, *width, *height)
	// }



	end_time := time.Now()
	fmt.Printf("Processing time : %v \n", end_time.Sub(start_time))
}
