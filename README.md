# optimalThumbnail
Optimal thumbnail creator written in golang
based on:
https://github.com/GHamrouni/OptimalCrop

Processes multiple files on 12 channels, and won't run in to the "too many files open" error.


Installation
------------
```bash
$ go get github.com/BPerlakiH/optimalThumbnail
```

Included libraries:
"github.com/nfnt/resize"

Usage
-----

optimalThumbnail -in your_input_folder_path -out your_output_folder_path -width 154 -height 154 -format jpg -q 85 -c 50

flags:
format - (optional) output format of the image, it can be jpg, png or webp [default jpg]
q - (optional) quality of the image encoding 0-100 [default 75]
c - (optional) concurency, the amount of go rutines launched at once [default 10] 