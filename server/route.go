package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-kit/kit/log/level"
	"github.com/julienschmidt/httprouter"
	"github.com/tomasen/realip"
)

func checkClientIP(h httprouter.Handle, address *Address) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		clientIP := realip.FromRequest(r)
		if address.Contains(clientIP) {
			h(w, r, ps)
		} else {
			w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	}
}

func newAPIRouter(config *ServerConfig) http.Handler {
	router := httprouter.New()
	router.GET("/api/token/:id", checkClientIP(func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		id := params.ByName("id")
		if id == "" {
			responseError(http.StatusBadRequest, "need id in url query", w)
			return
		}
		body := make(map[string]interface{})
		body["token"] = house.NewToken(id)
		responseOK(body, w)
		return
	}, config.APIAddress))

	router.POST("/api/room/add", checkClientIP(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		room := r.URL.Query().Get("room")
		members := r.URL.Query().Get("member")
		if room == "" || members == "" {
			responseError(http.StatusBadRequest, "bad request", w)
		}
		house.RoomAdd(room, strings.Split(members, "-"))
		responseOK(nil, w)
		return
	}, config.APIAddress))

	return router
}

func newWSRouter(config *ServerConfig) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws/connect", func(w http.ResponseWriter, r *http.Request) {
		// check client ip
		clientIP := realip.FromRequest(r)
		if !config.WSAddress.Contains(clientIP) {
			return
		}
		token := r.URL.Query().Get("token")
		mid := r.URL.Query().Get("member")
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
