package test

import (
	"gearbox/api"
	"gearbox/gearbox"
	"gearbox/status"
	"gearbox/status/is"
	"gearbox/test/includes"
	"net/http"
	"testing"
	"time"
)

func TestHostApiResponses(t *testing.T) {
	gb, sts := includes.GearboxConstructAndInitialize(t)
	if status.IsError(sts) {
		t.Error(sts.Message())
	}
	go gb.GetHostApi().Start()
	defer gb.GetHostApi().Stop()
	time.Sleep(time.Second * 1)
	ha := gb.GetHostApi()

	for methodName, methodMap := range ha.GetMethodMap() {
		if methodName != http.MethodGet {
			continue
		}
		t.Run(string(methodName), func(t *testing.T) {
			for routeName, _ := range methodMap {
				t.Run(string(routeName), func(t *testing.T) {
					sts = testResource(t, ha, routeName)
					if is.Error(sts) {
						t.Error(sts.Message())
					}
				})
			}
		})
	}
}

func testResource(t *testing.T, ha gearbox.HostApi, rn api.RouteName) (sts status.Status) {
	//
	//for range only.Once {
	//	var url api.UriTemplate
	//	var values api.ValuesFuncValues
	//	values, sts = testValues(ha, rn)
	//	if status.IsError(sts) {
	//		break
	//	}
	//	if values == nil {
	//		values = make(api.ValuesFuncValues, 1)
	//		values[0] = make(api.ValueFuncVarsValues, 1)
	//	}
	//	for i := range values[0] {
	//		var vars api.UriTemplateVars
	//		vars, sts = ha.GetUriTemplateVars(rn, values, i)
	//		if status.IsError(sts) {
	//			break
	//		}
	//		url, sts = ha.GetUrl(rn, vars)
	//		if status.IsError(sts) {
	//			break
	//		}
	//		hc := &test.HttpClient{}
	//		sts = hc.GET(string(url))
	//		if status.IsError(sts) {
	//			break
	//		}
	//		var body []byte
	//		body, sts = hc.GetBody()
	//		if status.IsError(sts) {
	//			break
	//		}
	//		var fn string
	//		if vars == nil {
	//			fn = string(rn)
	//		} else {
	//			fn = fmt.Sprintf("%s/%s", rn, strings.Join(vars.Values(), "/"))
	//		}
	//		gf := test.NewGoldFile(fn)
	//		var expected []byte
	//		expected, sts = gf.Read()
	//		if status.IsError(sts) {
	//			break
	//		}
	//		if gf.DoRunTest() && !bytes.Equal(expected, body) {
	//			opts := jsondiff.DefaultConsoleOptions()
	//			diff, s := jsondiff.Compare(expected, body, &opts)
	//			t.Error(fmt.Sprintf("no match; %s: %s", s, diff.String()))
	//		}
	//		sts = gf.Write(body)
	//	}
	//}
	return sts
}
func testValues(ha gearbox.HostApi, rn api.RouteName) (values api.ValuesFuncValues, sts status.Status) {
	//for range only.Once {
	//	var valuesFunc api.ValuesFunc
	//	valuesFunc, sts = ha.GetValuesFunc(rn)
	//	if status.IsError(sts) {
	//		break
	//	}
	//	if valuesFunc == nil {
	//		break
	//	}
	//	values, sts = valuesFunc()
	//	if len(values) == 0 {
	//		values = api.ValuesFuncValues{
	//			api.ValueFuncVarsValues{},
	//		}
	//	}
	//}
	return values, sts
}
