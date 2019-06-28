package main

import (
	"fmt"
	"github.com/martinlindhe/imgcat/lib"
	"gopkg.in/alecthomas/kingpin.v2"
	"image"
	"os"
	"strconv"
	"strings"
)

type file struct {
	FileName string
}

// exists reports whether the named file or directory exists.
func exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func (i *file) Set(value string) error {

	if !exists(value) {
		return fmt.Errorf("'%s' does not exist", value)
	}
	i.FileName = value
	return nil
}

func (i *file) String() string {
	return i.FileName
}

func imageList(s kingpin.Settings) (target *[]string) {
	file := new(file)
	s.SetValue(file)
	split(file.FileName)
	return
}

func split(src string) (target *[]string) {
	index := strings.Index(src, ".")
	name := src[:index]
	suffix := src[index+1:]
	fIn, _ := os.Open(src)
	defer fIn.Close()
	origin, fm, err := image.Decode(fIn)
	width := origin.Bounds().Dx() / 4
	height := origin.Bounds().Dy() / 4
	target = &[]string{}
	for i := 0; i < 16; i++ {
		x := i % 4
		y := i >> 2
		sx := x * width
		ex := (x + 1) * width
		sy := y * height
		ey := (y + 1) * height
		dst := name + strconv.Itoa(i) + "." + suffix
		*target = append(*target, dst)
		fmt.Println("src=", src, " dst=", dst)
		fOut, _ := os.Create(dst)
		defer fOut.Close()

		if err != nil {
			panic(err)
			return
		}
		err = imgcat.Clip(origin, fm, fOut, sx, sy, ex, ey, 0)
		if err != nil {
			panic(err)
		}
	}
	return
}

func main() {

	//setting := kingpin.Arg("file", "Image file to show.").Required()
	//verbose := kingpin.Flag("verbose", "Verbose mode.").Short('v').Bool()
	//
	//// support -h for --help
	//kingpin.CommandLine.HelpFlag.Short('h')
	//kingpin.Parse()

	files := split(os.Args[1])

	//if len(*files) > 1 && *verbose {
	//	heading = true
	//}

	out, err := os.OpenFile("/dev/tty", os.O_WRONLY, 0)
	for i, file := range *files {
		if i&3 == 0 {
			fmt.Fprintln(out)
		}
		err = imgcat.CatFile(file, out)
		if err != nil {
			fmt.Println(err)
		}
	}
}
