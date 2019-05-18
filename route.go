package main

import (
	"encoding/json"
	"github.com/leeif/mercury/house"
	"github.com/leeif/mercury/utils"
	"net/http"
	"strings"
)

var route *Route

type RouteFunc func(w http.ResponseWriter, r *http.Request)

func init() {
	route = &Route{
		Get:  make(map[string]RouteFunc),
		Post: make(map[string]RouteFunc),
		WS:   make(map[string]RouteFunc),
	}
	route.routeAPI()
	route.routeWS()
}

type Route struct {
	Get  map[string]RouteFunc
	Post map[string]RouteFunc
	WS   map[string]RouteFunc
}

func (route *Route) r(w http.ResponseWriter, r *http.Request) {
	var routeFunc RouteFunc
	path := r.URL.Path
	utils.Debug("path : %s", r.URL.Path)
	// ws
	routeFunc = route.WS[path]
	if routeFunc != nil {
		routeFunc(w, r)
		return
	}
	//api
	switch strings.ToLower(r.Method) {
	case "get":
		routeFunc = route.Get[path]
	case "post":
		routeFunc = route.Post[path]
	}

	if routeFunc != nil {
		routeFunc(w, r)
	} else {
		route.responseError(http.StatusNotFound, "not found", w)
	}
}

func (route *Route) routeAPI() {
	route.Get["/api/token"] = func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		body := make(map[string]interface{})
		if id == "" {
			route.responseError(http.StatusBadRequest, "need id in url query", w)
			return
		}
		body["token"] = house.GetHouse().GetToken(id)
		route.responseOK(body, w)
	}

	route.Post["/api/room/add"] = func(w http.ResponseWriter, r *http.Request) {
		room    := r.URL.Query().Get("room")
		members := strings.Split(r.URL.Query().Get("member"), "-")
		house.GetHouse().RoomAdd(room, members)
		route.responseOK(nil, w)
	}

	route.Post["/api/room/delete"] = func(w http.ResponseWriter, r *http.Request) {

	}
}

func (route *Route) routeWS() {
	route.WS["/ws/connect"] = func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "" {
			route.responseError(http.StatusBadRequest, "need a token", w)
			return
		}
		member := house.GetHouse().GetMemberFromToken(token)
		if member != nil {
			member.GenerateConnection(w, r)
		}
	}
}

func (route *Route) responseOK(body interface{}, w http.ResponseWriter) {
	res := make(map[string]interface{})
	res["status"] = "ok"
	if body != nil {
		res["body"] = body
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if b, err := json.Marshal(res); err == nil {
		w.Write(b)
	}
}

func (route *Route) responseError(status int, errMsg string, w http.ResponseWriter) {
	res := make(map[string]interface{})
	res["status"] = "error"
	res["error"] = errMsg
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if b, err := json.Marshal(res); err == nil {
		w.Write(b)
	}
}
