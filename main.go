package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "git.urantiatech.com/cloudcms/cloudcms/client"
	_ "git.urantiatech.com/cloudcms/cloudcms/content"
	s "git.urantiatech.com/cloudcms/cloudcms/service"
	"git.urantiatech.com/cloudcms/cloudcms/worker"
	"github.com/gorilla/mux"
	h "github.com/urantiatech/kit/transport/http"
)

func main() {
	// Parse command line parameters
	var port, upstreamPort int
	var ip, upstreamIP, dbFile, syncFile string
	flag.StringVar(&ip, "ip", "127.0.0.1", "The local worker IP address")
	flag.IntVar(&port, "port", 8080, "Port number of local worker")
	flag.StringVar(&dbFile, "dbFile", "cloudcms.db", "The database filename")
	flag.StringVar(&syncFile, "syncFile", "cloudcms.sync", "The workers database filename")
	flag.StringVar(&upstreamIP, "upstreamIP", "", "Upstream server hostname/IP")
	flag.IntVar(&upstreamPort, "upstreamPort", 8081, "The Port number of upstream server")
	flag.Parse()

	log.SetFlags(log.Lshortfile)

	if os.Getenv("LOCAL_WORKER_IP") != "" {
		ip = os.Getenv("LOCAL_WORKER_IP")
	}
	if os.Getenv("LOCAL_WORKER_PORT") != "" {
		p, err := strconv.ParseInt(os.Getenv("LOCAL_WORKER_PORT"), 10, 32)
		if err != nil {
			port = int(p)
		}
	}
	if os.Getenv("UPSTREAM_IP") != "" {
		upstreamIP = os.Getenv("UPSTREAM_IP")
	}
	if os.Getenv("UPSTREAM_PORT") != "" {
		p, err := strconv.ParseInt(os.Getenv("UPSTREAM_PORT"), 10, 32)
		if err != nil {
			upstreamPort = int(p)
		}
	}

	local := worker.Worker{Host: ip, Port: port}
	upstream := worker.Worker{Host: upstreamIP, Port: upstreamPort}

	if err := s.Initialize(dbFile, syncFile, &local, &upstream); err != nil {
		log.Fatal(err.Error())
	}
	var svc s.Service
	svc = s.Service{}

	r := mux.NewRouter()
	r.Handle("/create", h.NewServer(s.CreateEndpoint(svc), s.DecodeCreateReq, s.Encode))
	r.Handle("/read", h.NewServer(s.ReadEndpoint(svc), s.DecodeReadReq, s.Encode))
	r.Handle("/update", h.NewServer(s.UpdateEndpoint(svc), s.DecodeUpdateReq, s.Encode))
	r.Handle("/delete", h.NewServer(s.DeleteEndpoint(svc), s.DecodeDeleteReq, s.Encode))
	r.Handle("/search", h.NewServer(s.SearchEndpoint(svc), s.DecodeSearchReq, s.Encode))

	sync := mux.NewRouter()
	sync.Handle("/sync", h.NewServer(s.SyncEndpoint(svc), s.DecodeSyncReq, s.Encode))
	sync.Handle("/ping", h.NewServer(s.PingEndpoint(svc), s.DecodePingReq, s.Encode))

	go http.ListenAndServe(fmt.Sprintf(":%d", port+1), sync)
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}
