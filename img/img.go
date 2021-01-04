package img

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"os/exec"
	"strconv"
)

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
func checksum(img image.Image, len, width, x, y int) int {
	sum := 0
	rect := img.Bounds()
	pt := rect.Max
	xlim := min(x+len, pt.X)
	ylim := min(y+width, pt.Y)
	for i := x; i < xlim; i++ {
		for j := y; j < ylim; j++ {
			r, g, b, _ := img.At(i, j).RGBA()
			// fmt.Println(r, g, b)
			k := int(int(r) + int(g) + int(b))
			// fmt.Println(k)
			sum += k
			sum ^= 0xF1021012
			// sum ^= 0xFFFFFFFF
		}
	}
	return sum
}

// MakeChecksum makes checksum based on len and width
func MakeChecksum(img image.Image, len, width int, res *[]int, channel chan bool) {
	ans := make([]int, 0)
	rect := img.Bounds()
	pt := rect.Max

	// fmt.Println(pt)
	for i := 0; i < pt.X; i = i + len {
		for j := width; j < pt.Y; j = j + width {
			// fmt.Println(i, j)
			ans = append(ans, checksum(img, len, width, i, j))
			// break
		}
		// break
	}
	*res = ans
	channel <- true
}

// ExtractFrames given a video file extracts frames in certain intervals.
func ExtractFrames(filename string, interval, width, height, startime int) []image.Image {
	ans := make([]image.Image, 0)
	fmt.Printf("Extracting frames of size %dx%d\n", width, height)
	i := startime
	for {
		cmd := exec.Command("ffmpeg", "-accurate_seek", "-ss", strconv.Itoa(i), "-i", filename, "-s", fmt.Sprintf("%dx%d", width, height), "-vframes", "1", "-f", "image2", "pipe:1")
		buffer, _ := cmd.Output()
		fmt.Printf("Frame on %d sec , Size = %d\n", i, len(buffer))
		img, _, err := image.Decode(bytes.NewReader(buffer))
		if err != nil {
			// log.Fatal(err)
			fmt.Printf("Last Frame Size was 0, So end of file reached\n")
			break
		}
		ans = append(ans, img)
		i += interval
	}
	return ans
}

// SaveImages saves images in temp/dir/img_xxx.jpg
func SaveImages(frames []image.Image, dir string) {
	n := len(frames)
	folder := os.TempDir() + "/" + dir
	// fmt.Println(folder)
	_, err := os.Stat(folder)
	if err == nil {
		err = os.RemoveAll(folder)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = os.Mkdir(folder, 0755)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < n; i++ {
		filename := folder + "/img_" + strconv.Itoa(i) + ".jpg"
		f, err := os.Create(filename)
		if err != nil {
			log.Fatal(err)
		}
		jpeg.Encode(f, frames[i], nil)
	}
}
