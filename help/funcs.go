package help

import (
	"fmt"
	"gearbox/types"
)

func GetApiDocsUrl(route ...types.RouteName) string {
	var _route types.RouteName
	if len(route) == 0 {
		_route = "{route_name}"
	} else {
		_route = route[0]
	}
	return fmt.Sprintf("%s/%s", DocsBaseUrl, _route)
}

func GetApiHelp(routeName types.RouteName, more ...string) string {
	var _more string
	if len(more) > 0 {
		_more = " " + more[0]
	}
	return fmt.Sprintf("see API docs for%s: %s", _more, GetApiDocsUrl(routeName))
}

func ContactSupportHelp() string {
	return "contact support"
}
