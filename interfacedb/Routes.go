package interfacedb

import (
	"net/http"
)

//Route contains the routing details of a route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

//Routes contains the group of route details
type Routes []Route

var routes = Routes{
	Route{
		"Insert",
		"POST",
		"/insert",
		Insert,
	},
	Route{
		"Show",
		"GET",
		"/show/{showId}",
		Show,
	},
	Route{
		"ShowAll",
		"GET",
		"/showall",
		ShowAll,
	},
	Route{
		"Update",
		"POST",
		"/update/{updateId}",
		Update,
	},
	Route{
		"Delete",
		"POST",
		"/delete/{deleteId}",
		Delete,
	},
}
