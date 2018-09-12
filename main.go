package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	s "git.urantiatech.com/cloudcms/cloudcms/service"
	"github.com/boltdb/bolt"
	h "github.com/urantiatech/kit/transport/http"
)

func main() {
	// Parse command line parameters
	var port int
	var masterPath, workerPath, cloudsyncDNS string
	flag.IntVar(&port, "port", 8080, "Port")
	flag.StringVar(&masterPath, "masterPath", "master", "The path for Master process")
	flag.StringVar(&workerPath, "workerPath", "worker", "The path for Worker process")
	flag.StringVar(&cloudsyncDNS, "cloudsyncDNS", "cloudsync.default.svc.cluster.local",
		"The dns for CloudSync service")
	flag.Parse()

	log.SetFlags(log.Lshortfile)

	if os.Getenv("PORT") != "" {
		p, err := strconv.ParseInt(os.Getenv("PORT"), 10, 32)
		if err != nil {
			port = int(p)
		}
	}

	// Open the data file in your current directory.
	// It will be created if it doesn't exist.
	var err error
	s.DB, err = bolt.Open(masterPath+"/cloudcms.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer s.DB.Close()

	if err := s.InitIndexMap(masterPath); err != nil {
		log.Fatal(err)
	}

	var svc s.Service
	svc = s.Service{}

	http.Handle("/create", h.NewServer(s.CreateEndpoint(svc), s.DecodeCreateReq, s.Encode))
	http.Handle("/read", h.NewServer(s.ReadEndpoint(svc), s.DecodeReadReq, s.Encode))
	http.Handle("/update", h.NewServer(s.UpdateEndpoint(svc), s.DecodeUpdateReq, s.Encode))
	http.Handle("/delete", h.NewServer(s.DeleteEndpoint(svc), s.DecodeDeleteReq, s.Encode))
	http.Handle("/search", h.NewServer(s.SearchEndpoint(svc), s.DecodeSearchReq, s.Encode))
	http.Handle("/ping", h.NewServer(s.PingEndpoint(svc), s.DecodePingReq, s.Encode))

	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
