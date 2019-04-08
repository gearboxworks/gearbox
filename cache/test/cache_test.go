package test

import (
	"encoding/json"
	"fmt"
	"gearbox/cache"
	"gearbox/only"
	"gearbox/status"
	"gearbox/status/is"
	"gearbox/test/mock"
	"gearbox/types"
	"gearbox/util"
	"gearbox/validate"
	"net/http"
	"testing"
)

const cachekey = "foo"
const noCacheFile = "no cache file for key 'foo'"
const cacheKeyNotFound = "cache key not found"
const fooRetrieved = "cache retrieved for key 'foo'"
const fooExpired = "cache expired for key 'foo'"
const cacheKeyExpired = "cache key expired"

func TestCache(t *testing.T) {

	oss := mock.NewOsSupport(t)
	c := cache.NewCache(oss.GetCacheDir())

	//t.Run("SetGet", func(t *testing.T) {
	//	sts := c.Initialize()
	//	if status.IsError(sts) {
	//		t.Error(sts.Message())
	//	}
	//})

	t.Run("GetCacheFilepath()", func(t *testing.T) {
		fp := c.GetCacheFilepath(cachekey)
		got := util.ExtractRelativePath(fp, oss.GetCacheDir())
		wanted := types.RelativePath(fmt.Sprintf("/%s.json", cachekey))
		if got != wanted {
			t.Errorf("Wanted: %s, Got: %s", got, wanted)
		}
	})

	t.Run("Clear()", func(t *testing.T) {
		sts := c.Clear(cachekey)
		if is.Error(sts) {
			t.Errorf("Unable to clear cache '%s'", cachekey)
		}
	})

	t.Run("Get(does-not-exist)", func(t *testing.T) {
		_ = c.Clear(cachekey)
		data, sts := getCacheValue(t, c, &validate.Args{
			MustFail: true,
		})
		if is.Success(sts) {
			t.Errorf("error expected for '%s'; got: %s", cachekey, sts.Message())
		}
		if len(data) > 0 {
			t.Errorf("data returned from cleared cache '%s': %s", cachekey, string(data))
		}
	})

	t.Run("Set()", func(t *testing.T) {
		_ = c.Clear(cachekey)
		data := newTestStruct(nil)
		b, err := json.Marshal(data)
		if err != nil {
			t.Error("unable to marshal test data")
		}
		sts := c.Set(cachekey, b, "1s")
		if sts == nil {
			t.Error("sts returned nil")

		} else if sts.IsError() {
			t.Errorf("sts.Set() failed: %s", sts.Message())

		} else if sts.Message() != fmt.Sprintf("cache set for key '%s'", cachekey) {
			t.Errorf("sts.Set() returned unexpected success message: %s", sts.Message())

		}
	})

	t.Run("Set(),Get()", func(t *testing.T) {
		_ = c.Clear(cachekey)
		data := newTestStruct(nil)
		b, err := json.Marshal(data)
		if err != nil {
			t.Error("unable to marshal test data")
		}
		sts := c.Set(cachekey, b, "10m")
		if is.Error(sts) {
			t.Errorf("unable to set cache '%s'", cachekey)
		}
		b, sts = getCacheValue(t, c, &validate.Args{
			MustSucceed: true,
		})
		if is.Error(sts) {
			t.Errorf("unable to get cache '%s': %s", cachekey, sts.Message())
		}
		matchValues(t, data, b)
	})

	t.Run("Expired", func(t *testing.T) {
		for range only.Once {
			_ = c.Clear(cachekey)
			data := newTestStruct(nil)
			b, err := json.Marshal(data)
			if err != nil {
				t.Error("unable to marshal test data")
			}
			sts := c.Set(cachekey, b, "-1s")
			if is.Error(sts) {
				t.Errorf("unable to set cache '%s'", cachekey)
			}
			b, sts = getCacheValue(t, c, &validate.Args{
				MustExpire: true,
			})
			if sts == nil {
				t.Error("sts; got: nil, wanted: not nil")
				break
			}
			if len(b) == 0 {
				t.Error("b; got: \"\", wanted non-empty")
			}
			if is.Error(sts) {
				t.Error("sts.IsError(); got: false, wanted: true")
			}
			if sts.Message() != fooExpired {
				t.Errorf("sts.Message(); got: '%s', wanted: '%s'", fooExpired, sts.Message())
			}
			if sts.Cause().Error() != cacheKeyExpired {
				t.Errorf("sts.Message(); got: '%s', wanted: '%s'", cacheKeyExpired, sts.Message())
			}
			if sts.HttpStatus() != http.StatusOK {
				t.Errorf("sts.HttpStatus() got: %d, wanted: %d", sts.HttpStatus(), http.StatusOK)
			}
			matchValues(t, data, b)
		}
	})

}

func matchValues(t *testing.T, data *testStruct, b []byte) {
	data2 := &testStruct{}
	err := json.Unmarshal(b, &data2)
	if err != nil {
		t.Error("unable to unmarshal retrieved cache value")
	}
	if data.Name != data2.Name {
		t.Error("retrieved cache name does not match name")
	}
	if data.Number != data2.Number {
		t.Error("retrieved cache number does not match name")
	}
	if data.MapStringInt != nil {
		t.Error("retrieved cache Map not nil")
	}

}

type testStruct struct {
	Name         string
	Number       int
	MapStringInt map[string]int
}

func newTestStruct(m map[string]int) *testStruct {
	return &testStruct{
		Name:   "Parent",
		Number: 10,
	}
}

func getCacheValue(t *testing.T, c *cache.Cache, args *validate.Args) (data []byte, sts status.Status) {
	var ok bool
	data, ok, sts = c.Get(cachekey)
	if sts == nil {
		t.Error("sts; got: nil, wanted: not nil")
	}
	for range only.Once {
		switch {
		case args.MustFail:
			if len(data) > 0 {
				t.Errorf("data; got: %s, wanted: \"\"", string(data))
			}
			if ok {
				t.Error("ok; got true, wanted: false")
			}
			if is.Success(sts) {
				t.Error("sts.IsSuccess(); got: false, wanted: true")
			}
			if sts == nil {
				break
			}
			if sts.HttpStatus() != http.StatusInternalServerError {
				t.Errorf("sts.HttpStatus() got: %d, wanted: %d", sts.HttpStatus(), http.StatusInternalServerError)
			}
			if sts.Message() != noCacheFile {
				t.Errorf("sts.Message(); got: '%s', wanted: '%s'", noCacheFile, sts.Message())
			}
			if sts.Cause() == nil {
				t.Error("sts.Message(); got: nil, wanted: not nil")
			}
			if sts.Cause().Error() != cacheKeyNotFound {
				t.Errorf("sts.Cause().Error() '%s'", sts.Cause().Error())
			}

		case args.MustSucceed:
			if len(data) == 0 {
				t.Error("data; got: \"\", wanted non-empty")
			}
			if !ok {
				t.Error("ok; got: false, wanted: true")
			}
			if is.Error(sts) {
				t.Error("sts.IsError(); got: true, wanted: false")
			}
			if sts == nil {
				break
			}
			if sts.HttpStatus() != http.StatusOK {
				t.Errorf("sts.HttpStatus() got: %d, wanted: %d", sts.HttpStatus(), http.StatusOK)
			}
			if sts.Message() != fooRetrieved {
				t.Errorf("sts.Message(); got: '%s', wanted: '%s'", fooRetrieved, sts.Message())
			}
			if sts.Cause() != nil {
				t.Errorf("sts.Message(); got: '%s', wanted: nil", sts.Cause().Error())
			}

		case args.MustExpire:
			if len(data) == 0 {
				t.Error("data; got: \"\", wanted non-empty")
			}
			if ok {
				t.Error("ok; got: true, wanted: false")
			}
			if is.Error(sts) {
				t.Error("sts.IsError(); got: true, wanted: false")
			}
			if sts == nil {
				break
			}
			if sts.HttpStatus() != http.StatusOK {
				t.Errorf("sts.HttpStatus() got: %d, wanted: %d", sts.HttpStatus(), http.StatusOK)
			}
			if sts.Message() != fooExpired {
				t.Errorf("sts.Message(); got: '%s', wanted: '%s'", fooExpired, sts.Message())
			}
			if sts.Cause() == nil {
				t.Error("sts.Message(); got: nil, wanted: not nil")
			}
			if sts.Cause().Error() != cacheKeyExpired {
				t.Errorf("sts.Message(); got: '%s', wanted: '%s'", sts.Cause().Error(), cacheKeyExpired)
			}

		}

	}
	return data, sts
}
