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
	var host, externalHost, upstreamHost, dbFile, syncFile, key string
	var port, externalPort, upstreamPort int

	flag.StringVar(&host, "host", "localhost", "The local hostname/IP address")
	flag.IntVar(&port, "port", 8080, "The local port number")
	flag.StringVar(&externalHost, "externalHost", "", "The external hostname/IP address")
	flag.IntVar(&externalPort, "externalPort", 8080, "The external port number")
	flag.StringVar(&upstreamHost, "upstreamHost", "", "The uppstream hostname/IP address")
	flag.IntVar(&upstreamPort, "upstreamPort", 8080, "The upstream port number")
	flag.StringVar(&dbFile, "dbFile", "cloudcms.db", "The database filename")
	flag.StringVar(&syncFile, "syncFile", "cloudcms.sync", "The workers database filename")
	flag.StringVar(&key, "key", "", "Key that is used by other servers/clusters to join master")

	flag.Parse()

	log.SetFlags(log.Lshortfile)

	// Internal Host/Port for this service
	if os.Getenv("LOCAL_HOST") != "" {
		host = os.Getenv("LOCAL_HOST")
	}
	if os.Getenv("LOCAL_PORT") != "" {
		p, err := strconv.ParseInt(os.Getenv("LOCAL_PORT"), 10, 32)
		if err != nil {
			port = int(p)
		}
	}

	// External Host/Port for this service
	if os.Getenv("EXTERNAL_HOST") != "" {
		externalHost = os.Getenv("EXTERNAL_HOST")
	}
	if os.Getenv("EXTERNAL_PORT") != "" {
		p, err := strconv.ParseInt(os.Getenv("EXTERNAL_PORT"), 10, 32)
		if err != nil {
			externalPort = int(p)
		}
	}

	// Upstream Host/Port
	if os.Getenv("UPSTREAM_HOST") != "" {
		upstreamHost = os.Getenv("UPSTREAM_HOST")
	}
	if os.Getenv("UPSTREAM_PORT") != "" {
		p, err := strconv.ParseInt(os.Getenv("UPSTREAM_PORT"), 10, 32)
		if err != nil {
			upstreamPort = int(p)
		}
	}

	local := worker.Worker{Host: host, Port: port}
	upstream := worker.Worker{Host: upstreamHost, Port: upstreamPort}

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
	if upstreamHost == "" && key != "" {
		// This is master server/cluster, allow other servers/clusters to join this
		sync.Handle("/join", h.NewServer(s.JoinEndpoint(svc), s.DecodeJoinReq, s.Encode))
	}
	go http.ListenAndServe(fmt.Sprintf(":%d", port+1), sync)

	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}
