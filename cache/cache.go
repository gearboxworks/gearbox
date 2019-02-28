package cache

import (
	"encoding/json"
	"fmt"
	"gearbox/only"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

const ErrCacheMiss = "cache key '%s' not found"

type Cache struct {
	Dir string
}

type Wrapper struct {
	Expires string `json:"expires"`
	Data    []byte `json:"data"`
}

func NewCache(dir string) *Cache {
	return &Cache{
		Dir: dir,
	}
}

func (me *Cache) Close(f *os.File) {
	_ = f.Close()
}

func (me *Cache) Get(key string) (data []byte, err error) {
	for range only.Once {
		fn := fmt.Sprintf("%s/%s.json", me.Dir, key)
		var f *os.File
		f, err = os.Open(fn)
		if err != nil {
			pe, ok := err.(*os.PathError)
			if !ok {
				break
			}
			if pe.Err == syscall.ENOENT && pe.Op == "open" {
				err = fmt.Errorf(ErrCacheMiss, key)
				break
			}
			break
		}
		defer me.Close(f)
		var b []byte
		b, err = ioutil.ReadFile(fn)
		if err != nil {
			break
		}
		w := Wrapper{}
		err = json.Unmarshal(b, &w)
		if err != nil {
			break
		}
		expires, err := time.Parse(time.RFC3339, w.Expires)
		if err != nil {
			break
		}
		if expires.Before(time.Now()) {
			_ = os.Remove(fn)
			err = fmt.Errorf(ErrCacheMiss, key)
			break
		}

		data = w.Data
	}
	return data, err
}

func (me *Cache) Set(key string, b []byte, duration string) (err error) {
	for range only.Once {
		dur, err := time.ParseDuration(duration)
		if err != nil {
			break
		}
		w := &Wrapper{
			Expires: time.Now().Add(dur).Format(time.RFC3339),
			Data:    b,
		}
		b, err := json.Marshal(w)
		if err != nil {
			break
		}
		f := fmt.Sprintf("%s/%s.json", me.Dir, key)
		d := filepath.Dir(f)
		if !dirExists(d) {
			err = os.Mkdir(filepath.Dir(f), 0777)
			if err != nil {
				break
			}
		}
		err = ioutil.WriteFile(f, b, 0777)
		if err != nil {
			break
		}
	}
	return err
}

func dirExists(d string) bool {
	_, err := os.Stat(d)
	return !os.IsNotExist(err)
}
