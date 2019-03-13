package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	pb "gopkg.in/cheggaaa/pb.v1"
)

const RefreshRate = time.Millisecond * 100

// WriteCounter counts the number of bytes written to it. It implements to the io.Writer
// interface and we can pass this into io.TeeReader() which will report progress on each
// write cycle.
type WriteCounter struct {
	n   int // bytes read so far
	bar *pb.ProgressBar
}

func NewWriteCounter(total int) *WriteCounter {
	b := pb.New(total)
	b.SetRefreshRate(RefreshRate)
	b.ShowTimeLeft = true
	b.ShowSpeed = true
	b.SetUnits(pb.U_BYTES)

	return &WriteCounter{
		bar: b,
	}
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	wc.n += len(p)
	wc.bar.Set(wc.n)
	return wc.n, nil
}

func (wc *WriteCounter) Start() {
	wc.bar.Start()
}

func (wc *WriteCounter) Finish() {
	wc.bar.Finish()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Provide a url to download.")
		os.Exit(1)
	}

	fileUrl := os.Args[1]
	err := DownloadFile(fileUrl, path.Base(fileUrl))
	if err != nil {
		panic(err)
	}
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory. We pass an io.TeeReader
// into Copy() to report progress on the download.
func DownloadFile(url string, filename string) error {

	// Create the file, but give it a tmp file extension, this means we won't overwrite a
	// file until it's downloaded, but we'll remove the tmp extension once downloaded.
	out, err := os.Create(filename + ".tmp")
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fsize, _ := strconv.Atoi(resp.Header.Get("Content-Length"))

	// Create our progress reporter and pass it to be used alongside our writer
	counter := NewWriteCounter(fsize)
	counter.Start()

	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
	if err != nil {
		return err
	}

	counter.Finish()

	err = os.Rename(filename+".tmp", filename)
	if err != nil {
		return err
	}

	return nil
}
