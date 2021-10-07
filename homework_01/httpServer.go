package main

import(
	"flag"
	"fmt"
	"github.com/golang/glog"
	"io"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"strings"
	"time"
)

func main() {
	err := flag.Set("v", "4")
	if err != nil {
		return 
	}
	glog.V(2).Info("Http Server starting~")
	listenPort := os.Getenv("LISTEN_ADDR")
	if listenPort == "" {
		listenPort = ":80"
	}
	mux:= http.NewServeMux()
	startHTTP(mux)
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/healthz", healthzHandler)
	httpErr := http.ListenAndServe(listenPort, mux)
	if httpErr != nil {
		log.Fatal(httpErr)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	for k, v := range r.Header {
		w.Header().Add(k, strings.Join(v, ","))
		fmt.Println(k, v)
	}
	vEnv := os.Getenv("VERSION")
	if vEnv == "" {
		vEnv = "0.0.0"
	}
	w.Header().Add("VERSION", vEnv)
	w.WriteHeader(http.StatusOK)
	serverLog(http.StatusOK, r)
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, "ok")
}

func startHTTP(mux *http.ServeMux) {
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
}

func serverLog(statusCode int, r *http.Request) {
	t := time.Now().Format("2006-01-02 15-04-05")
	fmt.Printf("time: %s, remote: %s, status: %d\n", t, r.RemoteAddr, statusCode)
}