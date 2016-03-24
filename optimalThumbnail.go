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

// +build linux,amd64 darwin

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
	inDir := flag.String("in", "", "input directory or file[jpg|png|gif]")    // input directory
	outDir := flag.String("out", "", "output directory") // output directory
	width := flag.Int("width", 178, "the new width")     // width
	height := flag.Int("height", 178, "the new height")  // height
	outFormat := flag.String("format", "jpg", "image format")  // image format: jpg, png, webp
	quality := flag.Int("q", 75, "image output quality") // image encode quality 0-100
	concurency := flag.Int("c", 10, "amount of threads to be used") //amount of channels used in paralel default 10


	flag.Parse()
	usage := "usage: \noptimalThumbnail -in input[file|dir] -out outputDir -width 178 -height 178 -format jpg -q 70 -c 10"

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
	//define an array for the processable files
	var files []os.FileInfo
	var input_dir = *inDir
	file_info, err := os.Stat(*inDir)
	if err != nil {
		fmt.Println(err)
		d.Close()
		os.Exit(1)
	}

	if file_info.IsDir() {
		fi, err := d.Readdir(-1)
		d.Close()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		//get the processable files:
		for _, file := range fi {
			ext := filepath.Ext(file.Name())
			if !file.IsDir() && (ext == ".png" || ext == ".jpg" || ext == ".gif") {
				files = append(files, file)
			}
		}
	} else {
		//d is a file not a directory, but a file, append it
		files = append(files, file_info)
		input_dir = filepath.Dir(*inDir)
	}

	//count the files
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
			processFiles(files_chunk, input_dir, *outDir, *width, *height, *quality, *outFormat)

			percent := float64(end) / float64(files_count)
			fmt.Printf("\r%v-%v/%v (%v%%)", i, end, files_count, int(percent * 100))

			i += max_channels+1
			if files_count <= i {
				is_more = false
			}
		}

	} else {
		for _, f := range files {
			fmt.Printf("%v\n", f.Name())
		}
		processFiles(files, input_dir, *outDir, *width, *height, *quality, *outFormat)
	}


	end_time := time.Now()
	fmt.Printf("Processing time : %v \n", end_time.Sub(start_time))
}

