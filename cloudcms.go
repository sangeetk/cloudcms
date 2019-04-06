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
	"golang.org/x/text/language"
)

// Languages supported
func Languages(languages []language.Tag) {
	s.Languages = languages
}

// Run method should be called from main function
func Run(port int) {
	// Parse command line parameters
	var host, dbFile, syncFile string

	flag.StringVar(&host, "host", "localhost", "The local hostname/IP address")
	flag.StringVar(&dbFile, "dbFile", "db/cloudcms.db", "The database filename")
	flag.StringVar(&syncFile, "syncFile", "db/cloudcms.sync", "The workers database filename")
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
	if err := s.RebuildIndex(); err != nil {
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
	r.Handle("/facets", h.NewServer(s.FacetsSearchEndpoint(svc), s.DecodeFacetsSearchReq, s.Encode))
	r.Handle("/list", h.NewServer(s.ListEndpoint(svc), s.DecodeListReq, s.Encode))
	r.Handle("/schema", h.NewServer(s.SchemaEndpoint(svc), s.DecodeSchemaReq, s.Encode))
	r.Handle("/pull", h.NewServer(s.PullEndpoint(svc), s.DecodePullReq, s.Encode))
	r.Handle("/push", h.NewServer(s.PushEndpoint(svc), s.DecodePushReq, s.Encode))

	r.PathPrefix("/drive/").Handler(http.StripPrefix("/drive/", http.FileServer(http.Dir("drive"))))

	sync := mux.NewRouter()
	sync.Handle("/sync", h.NewServer(s.SyncEndpoint(svc), s.DecodeSyncReq, s.Encode))
	sync.Handle("/ping", h.NewServer(s.PingEndpoint(svc), s.DecodePingReq, s.Encode))
	go http.ListenAndServe(fmt.Sprintf(":%d", port+1), sync)

	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}
