package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

const (
	ADMIN_PASSWORD         = "SUNXIN" // TODO: must be in ouside configuration, or fetch from DB
	MAX_HIST               = 100
	DEFAULT_FORCE_GZIP_NUM = 10000
)

func main() {
	var httpPort = flag.String("port", "8080", "Server http posrt, defqult is 8080, ex: 8081")
	var forceGzip = flag.Bool("force_gzip", false, "When limit is a huge number, return response by encoding gzip, default value is false")
	var forceGzipNum = flag.Int64("force_gzip_num", DEFAULT_FORCE_GZIP_NUM, "When (limit >= force_gzip_num)+(force_gzip=true), return response by encoding gzip, default value is 10000")
	flag.Parse()

	s := NewServer(*forceGzip, *forceGzipNum)
	http.HandleFunc("/heartbeats", handlerHeartbeats)
	http.Handle("/stat", handlerServerStat(s))
	http.Handle("/update_conf", handlerUpdateConf(s))
	http.Handle("/run", handlerServerRun(s))

	err := http.ListenAndServe(fmt.Sprintf(":%s", *httpPort), nil)

	log.Printf("PROG FAILED: %+v", err)
}
