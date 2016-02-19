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
	"strings"
	"github.com/BPerlakiH/optimalThumbnail/process"
)

import _ "image/jpeg"
import _ "image/gif"

func processFiles(files []os.FileInfo, inDir string, outDir string, width int, height int, quality int, outFormat string) {
	files_count := len(files)
	//go parallel on all files at the same time:
	channels := make(chan int, files_count)

	for index := range files {
		filename := files[index].Name()
		// fmt.Println(filename)
		go func(filename string) {
			// Combine the directory path with the filename to get 
			// the full path of the image
			basename := filepath.Base(filename)
			outputFileName := strings.TrimSuffix(basename, filepath.Ext(basename)) + "." + outFormat

			fullImagePath := filepath.Join(inDir, filename)
			fullImageOutPath := filepath.Join(outDir, outputFileName)
			// fmt.Println(fullImageOutPath)

			process.ProcessFile(fullImagePath, fullImageOutPath, width, height, quality)
			//open the channel
			channels <- 1
		}(filename)
	}
	//close all channels
	for i := 0; i < files_count; i++ {
		<-channels
	}
}

func main() {

	start_time := time.Now()

	// Read the CMD options
	inDir := flag.String("in", "", "input directory")    // input directory
	outDir := flag.String("out", "", "output directory") // output directory
	width := flag.Int("width", 178, "the new width")     // width
	height := flag.Int("height", 178, "the new height")  // height
	outFormat := flag.String("format", "jpg", "image format")  // image format: jpg, png, webp
	quality := flag.Int("q", 75, "image output quality") // image encode quality 0-100
	concurency := flag.Int("c", 10, "amount of threads to be used") //amount of channels used in paralel default 10


	flag.Parse()
	usage := "usage: \noptimalThumbnail -in inputDir -out outputDir -width 178 -height 178 -format jpg -q 70 -c 10"

	if *inDir == "" || *outDir == "" {
		log.Fatal(usage)
	}

	if *outFormat != "jpg" && *outFormat != "png" && *outFormat != "webp" {
		fmt.Printf("invalid image format (%v), it can be: [jpg|png|webp]\n", *outFormat)
		log.Fatal(usage)
	}

	// Print the cmd options
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
	max_channels := *concurency
	
	if max_channels < files_count {
		//go parallel on the channels:
		is_more := true
		i := 0
		for is_more == true {
			end := i + max_channels
			end = int(math.Min(float64(end), float64(files_count)))
			files_chunk := files[i:end]
			processFiles(files_chunk, *inDir, *outDir, *width, *height, *quality, *outFormat)

			percent := float64(end) / float64(files_count)
			fmt.Printf("\r%v-%v/%v (%v%%)", i, end, files_count, int(percent * 100))

			i += max_channels+1
			if files_count <= i {
				is_more = false
			}
			// time.Sleep(500 * time.Millisecond)
		}

		
		// for i := 0; i <= channel_count; i++ {
		// 	go func(i int) {
		// 		min := i * files_per_channel
		// 		channel_ends := min + files_per_channel
		// 		// fmt.Printf("min: %v, channel_ends: %v, files_count: %v\n", min, channel_ends, files_count)
		// 		max := int(math.Min(float64(channel_ends), float64(files_count)))
		// 		// fmt.Printf("%v - %v\n", min, max)
		// 		for k := min; k < max; k++ {
		// 			filename := files[k].Name()
		// 			percent := float64(k) / float64(files_count)
		// 			fmt.Printf("\r%v %v/%v (%v%%)", filename, k, files_count, int(percent * 100))
		// 			fullImagePath := filepath.Join(*inDir, filename)

		// 			basename := filepath.Base(filename)
		// 			outputFileName := strings.TrimSuffix(basename, filepath.Ext(basename)) + "." + *outFormat

		// 			fullImageOutPath := filepath.Join(*outDir, outputFileName)
		// 			process.ProcessFile(fullImagePath, fullImageOutPath, *width, *height, *quality)
		// 		}
		// 		//open the channel
		// 		channels <- 1
		// 	}(i)
		// 	//close the channel
		// 	<-channels
		// }
	} else {
		for _, f := range files {
			fmt.Printf("%v\n", f.Name())
		}
		processFiles(files, *inDir, *outDir, *width, *height, *quality, *outFormat)
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

