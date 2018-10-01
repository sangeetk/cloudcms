package cloudcms

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	s "git.urantiatech.com/cloudcms/cloudcms/service"
	"git.urantiatech.com/cloudcms/cloudcms/worker"
	"github.com/gorilla/mux"
	h "github.com/urantiatech/kit/transport/http"
)

// Run method should be called from main function
func Run() {
	// Parse command line parameters
	var host, dbFile, syncFile string
	var port int

	log.SetFlags(log.Lshortfile)

	flag.StringVar(&host, "host", "localhost", "The local hostname/IP address")
	flag.IntVar(&port, "port", 8080, "The local port number")
	flag.StringVar(&dbFile, "dbFile", "cloudcms.db", "The database filename")
	flag.StringVar(&syncFile, "syncFile", "cloudcms.sync", "The workers database filename")

	flag.Parse()

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

	localworker := worker.Worker{Host: host, Port: port}

	if err := s.Initialize(dbFile, syncFile, &localworker); err != nil {
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