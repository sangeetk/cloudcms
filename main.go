package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "git.urantiatech.com/cloudcms/cloudcms/client"
	s "git.urantiatech.com/cloudcms/cloudcms/service"
	"github.com/gorilla/mux"
	h "github.com/urantiatech/kit/transport/http"
)

func main() {
	// Parse command line parameters
	var port int
	var dbFile string
	flag.IntVar(&port, "port", 8080, "Port")
	flag.StringVar(&dbFile, "dbFile", "cloudcms.db", "The database filename")
	flag.Parse()

	log.SetFlags(log.Lshortfile)

	if os.Getenv("PORT") != "" {
		p, err := strconv.ParseInt(os.Getenv("PORT"), 10, 32)
		if err != nil {
			port = int(p)
		}
	}

	if err := s.Initialize(dbFile); err != nil {
		log.Fatal(err)
	}

	var svc s.Service
	svc = s.Service{}

	r := mux.NewRouter()

	r.Handle("/create", h.NewServer(s.CreateEndpoint(svc), s.DecodeCreateReq, s.Encode))
	r.Handle("/read", h.NewServer(s.ReadEndpoint(svc), s.DecodeReadReq, s.Encode))
	r.Handle("/update", h.NewServer(s.UpdateEndpoint(svc), s.DecodeUpdateReq, s.Encode))
	r.Handle("/delete", h.NewServer(s.DeleteEndpoint(svc), s.DecodeDeleteReq, s.Encode))
	r.Handle("/search", h.NewServer(s.SearchEndpoint(svc), s.DecodeSearchReq, s.Encode))
	r.Handle("/ping", h.NewServer(s.PingEndpoint(svc), s.DecodePingReq, s.Encode))

	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}
