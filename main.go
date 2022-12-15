package main

import (
	"bufio"
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/nfnt/resize"
	"github.com/schollz/progressbar/v3"
)

func main() {

	for {
		fmt.Print("Path to compress images->")
		scanner1 := bufio.NewScanner(os.Stdin)
		var path string
		if scanner1.Scan() {
			path = strings.Trim(scanner1.Text(), "\"")
		}
		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Println("path not found")
		} else {
			allfiles := listFiles(path)
			fmt.Printf("found %d files are you sure ? (y/n)", len(allfiles))
			scanner2 := bufio.NewScanner(os.Stdin)
			if scanner2.Scan() {
				if scanner2.Text() == "y" {
					bar := progressbar.Default(int64(len(allfiles)))
					for i := 0; i < len(allfiles); i++ {
						makesmall(allfiles[i])
						bar.Add(1)
					}
					fmt.Println("DONE")
				} else {
					os.Exit(0)
				}
			}
		}
	}

}

func makesmall(path string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("error open file")
		return
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		fmt.Println("cannot decode file")
		file.Close()
		return
	}
	file.Close()

	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(2000, 0, img, resize.Lanczos3)

	out, err := os.Create(path)
	if err != nil {
		fmt.Println("cannot create file")
		out.Close()
		return
	}
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, m, nil)
}

func listFiles(desination string) []string {
	var files []string

	filess, err := ioutil.ReadDir(desination)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range filess {
		if !file.IsDir() {
			files = append(files, desination+"/"+file.Name())
		}
	}

	return files

}
