package server

import (
	"github.com/go-kit/kit/log/level"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func newAPIRouter() http.Handler {
	router := httprouter.New()
	router.GET("/api/token/:id", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		id := params.ByName("id")
		if id == "" {
			responseError(http.StatusBadRequest, "need id in url query", w)
			return
		}
		body := make(map[string]interface{})
		body["token"] = house.NewToken(id)
		responseOK(body, w)
		return
	})

	router.POST("/api/room/add", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		room := r.URL.Query().Get("room")
		members := r.URL.Query().Get("member")
		if room == "" || members == "" {
			responseError(http.StatusBadRequest, "bad request", w)
		}
		house.RoomAdd(room, strings.Split(members, "-"))
		responseOK(nil, w)
		return
	})

	return router
}

func newWSRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws/connect", func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		mid := r.URL.Query().Get("mid")
		// count, err := strconv.Atoi(r.URL.Query().Get("count"))
		// if err != nil {
		// 	responseError(http.StatusBadRequest, "bad request", w)
		// 	return
		// }
		if token == "" || mid == "" {
			responseError(http.StatusBadRequest, "bad request", w)
			return
		}
		err := house.MemberConnect(w, r, mid, token)
		if err != nil {
			level.Warn(logger).Log("msg", err)
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
