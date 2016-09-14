package main

import "uniagent/common/net/restful"

var routes = []restful.Route{
	{
		"GET",
		"/parameters/path/:user",
		PathParameters,
	},
	{
		"Get",
		"/parameters/query",

		QueryParameters,
	},
	{
		"PUT",
		"/handling/request",
		HandlingRequest,
	},
}
