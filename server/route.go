package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func newAPIRouter() http.Handler {
	router := httprouter.New()
	router.GET("/api/token/:id", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		id := params.ByName("id")
		body := make(map[string]interface{})
		if id == "" {
			responseError(http.StatusBadRequest, "need id in url query", w)
			return
		}
		body["token"] = house.GetToken(id)
		responseOK(body, w)
	})

	router.POST("/api/room/add", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		room := r.URL.Query().Get("room")
		members := strings.Split(r.URL.Query().Get("member"), "-")
		house.RoomAdd(room, members)
		responseOK(nil, w)
		return
	})
	return router
}

func newWSRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws/connect", func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "" {
			responseError(http.StatusBadRequest, "need a token", w)
			return
		}
		member := house.GetMemberFromToken(token)
		if member != nil {
			member.GenerateConnection(w, r, house.ConnPool)
		}
	})
	return mux
}

func responseOK(body interface{}, w http.ResponseWriter) {
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

func responseError(status int, errMsg string, w http.ResponseWriter) {
	res := make(map[string]interface{})
	res["status"] = "error"
	res["error"] = errMsg
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if b, err := json.Marshal(res); err == nil {
		w.Write(b)
	}
}
