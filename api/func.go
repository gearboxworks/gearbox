package api

import (
	"fmt"
	"github.com/labstack/echo"
)

func optionsHandler(ctx echo.Context) error {
	return nil
}

func GetApiDocsUrl(route ...RouteName) string {
	var _route RouteName
	if len(route) == 0 {
		_route = "{route_name}"
	} else {
		_route = route[0]
	}
	return fmt.Sprintf("%s/%s", DocsBaseUrl, _route)
}

func GetApiHelp(routeName RouteName, more ...string) string {
	var _more string
	if len(more) > 0 {
		_more = " " + more[0]
	}
	return fmt.Sprintf("see API docs for%s: %s", _more, GetApiDocsUrl(routeName))
}

func ContactSupportHelp() string {
	return "contact support"
}

func noop(x ...interface{}) {}
