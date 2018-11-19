package main

import (
	"flag"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	localappdata      = os.Getenv("LOCALAPPDATA")
	localWallpaperDir = localappdata + `\Packages\Microsoft.Windows.ContentDeliveryManager_cw5n1h2txyewy\LocalState\Assets\`

	saveDir   = flag.String("save", ".", "save wallpaper directory")
	todayFlag = flag.Bool("today", false, "save today wallpaper only")

	sizeLimit = int64(100 * 1000) //100KB
)

func main() {
	flag.Parse()
	dir := filepath.Dir(localWallpaperDir)
	today := time.Now().AddDate(0, 0, -1)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		src := filepath.Join(dir, file.Name())
		dst := *saveDir + "\\" + file.Name() + ".jpg"

		if file.Size() < sizeLimit {
			continue
		}
		if *todayFlag && file.ModTime().Before(today) { // file < now -1
			continue
		}
		d, err := getImageDimension(src)
		if err != nil {
			log.Error(err)
			continue
		}
		if !(d.x == 1920 && d.y == 1080) {
			continue
		}
		if err := copyfile(src, dst); err != nil {
			log.Error(err)
		}
	}
}

func copyfile(srcName string, dstName string) error {
	src, err := os.Open(srcName)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(dstName)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}
	return nil
}

type dim struct {
	x int
	y int
}

func getImageDimension(name string) (dim, error) {
	var d dim
	file, err := os.Open(name)
	if err != nil {
		return d, err
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		return d, err
	}
	b := img.Bounds()
	return dim{x: b.Max.X, y: b.Max.Y}, nil
}
