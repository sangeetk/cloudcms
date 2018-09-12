package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	s "git.urantiatech.com/cloudcms/cloudcms/service"
	h "github.com/urantiatech/kit/transport/http"
)

func main() {
	// Parse command line parameters
	var port int
	var dbFile, workDir string
	flag.IntVar(&port, "port", 8080, "Port")
	flag.StringVar(&dbFile, "dbFile", "cloudcms.db", "The database file path")
	flag.StringVar(&workDir, "workDir", ".", "The directory for storing db and index files")
	flag.Parse()

	log.SetFlags(log.Lshortfile)

	if os.Getenv("PORT") != "" {
		p, err := strconv.ParseInt(os.Getenv("PORT"), 10, 32)
		if err != nil {
			port = int(p)
		}
	}

	if err := s.Initialize(dbFile, workDir); err != nil {
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
