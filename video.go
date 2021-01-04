package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	"log"
	"os"
	_ "reflect"
	"strconv"
	"strings"

	"./img"

	"github.com/jung-kurt/gofpdf"
	_ "golang.org/x/image/bmp"
)

const threshold float64 = 0.9

var dir string = "SavedImages"

// returns true if they are similar above a certain threshold
func comp(a, b []int) bool {
	n := len(a)
	cnt := 0
	for i := 0; i < n; i++ {
		if a[i] == b[i] {
			cnt++
		}
	}
	val := float64(float64(cnt) / float64(n))
	return val > threshold
}
// Uses images to make pdf 
func makePdf(n int, output string) {
	pdf := gofpdf.New("L", "pt", "A3", "")
	pdf.SetFont("Arial", "B", 16)
	width := 1200
	height := 800
	folder := os.TempDir() + "/" + dir
	for i := 0; i < n; i++ {
		filename := folder + "/img_" + strconv.Itoa(i) + ".jpg"
		pdf.ImageOptions(filename, 0, 0, float64(width), float64(height), true, gofpdf.ImageOptions{}, 0, "")
	}
	err := pdf.OutputFileAndClose(output)
	if err != nil {
		log.Fatal(err)
	}
	err = os.RemoveAll(folder)
	if err != nil {
		log.Fatal(err)
	}
}
func main() {
	var filename string
	var outfile string
	var interval int
	var resolution string
	var boxlen int
	var startime int
	flag.StringVar(&filename, "p", "//", "input file path such as: abc/d/e.mp4")
	flag.StringVar(&outfile, "o", "out.pdf", "output file name such as : example.pdf")
	flag.IntVar(&interval, "i", 60, "Interval in seconds such as: 10")
	flag.StringVar(&resolution, "s", "800x600", "Resolution such as : 1920x1080")
	flag.IntVar(&boxlen, "c", 100, "Length of box for checksum(Affects Performace for small values and accuracy for large): 100")
	flag.IntVar(&startime, "ss", 0, "Start time in seconds such as: 10")
	flag.Parse()
	resolarr := strings.Split(resolution, "x")
	width, err := strconv.Atoi(resolarr[0])
	if err != nil && width <= 0 {
		fmt.Println("Enter a valid width")
		return
	}
	height, err := strconv.Atoi(resolarr[1])
	if err != nil && height <= 0 {
		fmt.Println("Enter a valid height")
		return
	}
	fmt.Println(filename)
	if filename == "//" {
		fmt.Println("Pls use arg -h to see all valid options")
		return
	}
	if interval <= 0 {
		fmt.Println("Pls input a positive length of interval")
		return
	}
	frames := img.ExtractFrames(filename, interval, width, height, startime)
	size := len(frames)
	if size == 0 {
		fmt.Printf("File doesn't Exist , pls enter valid file\n")
		return
	}
	fmt.Println(size)
	ans := make([][]int, size)
	arrchan := make([]chan bool, size)
	for i := range arrchan {
		arrchan[i] = make(chan bool)
	}
	for i, v := range frames {
		go img.MakeChecksum(v, boxlen, boxlen, &ans[i], arrchan[i])
	}
	res := make([]bool, size)
	res[0] = true

	for i := 0; i < size; i++ {
		<-arrchan[i]
		if i > 0 {
			res[i] = !comp(ans[i], ans[i-1])
		}
	}
	uniqFrames := make([]image.Image, 0)
	for i := 0; i < size; i++ {
		if res[i] {
			uniqFrames = append(uniqFrames, frames[i])
		}
	}
	// fmt.Println(res)
	img.SaveImages(uniqFrames, dir)
	makePdf(len(uniqFrames), outfile)
}
