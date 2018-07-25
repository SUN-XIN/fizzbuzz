package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	locip "github.com/SUN-XIN/iplocation"
)

// handler to process client's request
func handlerServerRun(s *server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// get obj from body
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Failed ReadAll body: %+v", err)
			quickResponse(w, http.StatusInternalServerError, []byte("Failed Read Body"))
			return
		}

		var cr clientRequest
		err = decodeObj(b, &cr)
		if err != nil {
			log.Printf("Failed decodeObj: %+v", err)
			quickResponse(w, http.StatusInternalServerError, []byte("Failed decode Body"))
			return
		}

		err = cr.validate()
		if err != nil {
			log.Printf("Failed validate: %+v", err)
			quickResponse(w, http.StatusBadRequest, []byte(err.Error()))
			return
		}

		// client info
		clientIP := r.RemoteAddr
		fb := NewFizzBuzz(&cr, clientIP)
		ir, err := locip.IPLocationFromIPStack(clientIP)
		if err != nil {
			log.Printf("Failed IPLocationFromIPStack for IP (%s): %+v", clientIP, err)
		} else {
			fb.locip = ir
		}

		// add new call into HistoryCall
		s.AddNewCall(fb)

		// process
		res := processRequest(&cr)
		b, err = json.Marshal(clientResponse{Result: res})
		if err != nil {
			log.Printf("Failed Marshal response: %+v", err)
			quickResponse(w, http.StatusInternalServerError, []byte("Failed encode response"))
			return
		}

		// set stat "ok" for this call
		fb.SetOK()
		switch {
		case cr.ResponseGzip,
			s.CheckForceGzip(cr.Limit):
			responseGzip(w, http.StatusOK, b)
		default:
			quickResponse(w, http.StatusOK, b)
		}
	})
}

// handler to update server's configuration
func handlerUpdateConf(s *server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// another methode to read body (different with handlerServerRun)
		dec := json.NewDecoder(r.Body)

		var conf configuration
		err := dec.Decode(&conf)
		if err != nil {
			log.Printf("Failed Decode Body: %+v", err)
			quickResponse(w, http.StatusInternalServerError, []byte("Failed Decode Body"))
			return
		}

		err = conf.validate()
		if err != nil {
			quickResponse(w, http.StatusBadRequest, []byte(err.Error()))
		}

		s.UpdateConf(&conf)
	})
}

// handler to show server's stat, ONLY ADMIN CAN ACCESS
func handlerServerStat(s *server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.FormValue("key")
		if key == "" {
			quickResponse(w, http.StatusUnauthorized, []byte("Only admin can access"))
		}

		if s.CheckAdmin(key) {
			res := s.ShowInfo()
			fmt.Fprintf(w, "%s", res)
		} else {
			quickResponse(w, http.StatusForbidden, []byte("Bad Password"))
		}
	})
}

func handlerHeartbeats(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}

func decodeObj(b []byte, obj interface{}) error {
	return json.Unmarshal(b, obj)
}

func quickResponse(w http.ResponseWriter, httpCode int, msg []byte) {
	w.WriteHeader(httpCode)
	w.Write(msg)
}

func responseGzip(w http.ResponseWriter, httpCode int, msg []byte) {
	w.Header().Set("Content-Encoding", "gzip")
	w.WriteHeader(httpCode)

	zw := gzip.NewWriter(w)
	zw.Write(msg)
}
