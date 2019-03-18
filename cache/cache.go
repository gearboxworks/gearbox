package cache

import (
	"encoding/json"
	"fmt"
	"gearbox/only"
	"gearbox/stat"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

const ErrCacheMiss = "cache key not found"
const ErrExpired = "cache key expired"

type Cache struct {
	Dir     string
	Disable bool
}

type Wrapper struct {
	Expires string `json:"expires"`
	Data    string `json:"data"`
}

func NewCache(dir string) *Cache {
	return &Cache{
		Dir: dir,
	}
}

func (me *Cache) Close(f *os.File) {
	_ = f.Close()
}

func (me *Cache) GetCacheFilepath(key string) string {
	return filepath.FromSlash(fmt.Sprintf("%s/%s.json", me.Dir, key))
}

func (me *Cache) VerifyCacheFile(key string) (fp string, status stat.Status) {
	var f *os.File
	var err error
	for range only.Once {
		fp := me.GetCacheFilepath(key)
		f, err = os.Open(fp)
		if err != nil {
			pe, ok := err.(*os.PathError)
			if !ok {
				break
			}
			if pe.Err == syscall.ENOENT && pe.Op == "open" {
				err = fmt.Errorf(ErrCacheMiss)
				break
			}
			break
		}
	}
	me.Close(f)
	if err != nil {
		var msg string
		if err.Error() == ErrCacheMiss {
			msg = fmt.Sprintf("no cache file for key '%s'", key)
		} else {
			msg = fmt.Sprintf("cannot open cache file for key '%s'", key)
		}
		status = stat.NewFailStatus(&stat.Args{
			Error:   err,
			Message: msg,
		})
	}
	return fp, status
}

func (me *Cache) Get(key string) (data []byte, ok bool, status stat.Status) {
	for range only.Once {
		if me.Disable {
			break
		}
		var fp string
		fp, status = me.VerifyCacheFile(key)
		if status.IsError() {
			break
		}
		var b []byte
		b, err := ioutil.ReadFile(fp)
		if err != nil {
			status = stat.NewFailStatus(&stat.Args{
				Error:   err,
				Message: fmt.Sprintf("could not read file '%s'", fp),
				Help:    fmt.Sprintf("ensure you have permissions to read '%s'", fp),
			})
			break
		}
		w := Wrapper{}
		err = json.Unmarshal(b, &w)
		if err != nil {
			status = stat.NewFailStatus(&stat.Args{
				Error:   err,
				Message: fmt.Sprintf("could not unmarshal JSON in file '%s'", fp),
				Help:    fmt.Sprintf("try deleting the files your cache at '%s'", filepath.Dir(fp)),
			})
			break
		}
		data = []byte(w.Data)
		expires, err := time.Parse(time.RFC3339, w.Expires)
		if err != nil {
			status = stat.NewFailStatus(&stat.Args{
				Error:   err,
				Message: fmt.Sprintf("failed to calculate cache expiration for file '%s'", fp),
				Help:    fmt.Sprintf("try deleting the files your cache at '%s'", filepath.Dir(fp)),
			})
			break
		}
		if expires.Before(time.Now()) {
			//_ = os.Remove(fp)
			status = stat.NewStatus(&stat.Args{
				Failed:     false,
				Error:      fmt.Errorf(ErrExpired),
				Message:    fmt.Sprintf("cache expired for key '%s'", key),
				HttpStatus: http.StatusOK,
			})
			break
		}
		ok = true
	}
	return data, ok, status
}

func (me *Cache) Set(key string, b []byte, duration string) (status stat.Status) {
	for range only.Once {
		dur, err := time.ParseDuration(duration)
		if err != nil {
			break
		}
		w := &Wrapper{
			Expires: time.Now().Add(dur).Format(time.RFC3339),
			Data:    string(b),
		}
		b, err := json.Marshal(w)
		if err != nil {
			status = stat.NewFailStatus(&stat.Args{
				Error:   err,
				Message: fmt.Sprintf("could not marshal JSON to cache key '%s'", key),
				Help:    "this should never happen, so try rebooting. Or contacting support",
			})
			break
		}
		fp := me.GetCacheFilepath(key)
		d := filepath.Dir(fp)
		if !dirExists(d) {
			err = os.Mkdir(filepath.Dir(fp), 0777)
			if err != nil {
				status = stat.NewFailStatus(&stat.Args{
					Error:   err,
					Message: fmt.Sprintf("unable to create cache directory '%s'", d),
					Help:    fmt.Sprintf("ensure you have permissions to '%s'", filepath.Dir(d)),
				})
				break
			}
		}
		err = ioutil.WriteFile(fp, b, 0777)
		if err != nil {
			status = stat.NewFailStatus(&stat.Args{
				Error:   err,
				Message: fmt.Sprintf("unable to write to cache file '%s'", fp),
				Help:    fmt.Sprintf("ensure you have permissions to '%s'", filepath.Dir(d)),
			})
			break
		}
	}
	return status
}

func dirExists(d string) bool {
	_, err := os.Stat(d)
	return !os.IsNotExist(err)
}
