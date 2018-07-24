package app

import (
	"os"
	"github.com/greenac/artemis/logger"
	"github.com/greenac/sqsmock/sqs"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

type serverArguments struct {
	ServerUrl *string
	WorkerUrl *string
}

const(
	serverUrlArg = "url"
	workerUrlArg = "workerurl"
)

func (sa *serverArguments) formatServerUrl() {
	if sa.ServerUrl == nil {
		url := "localhost:4242"
		sa.ServerUrl = &url
	} else if strings.Contains(*(sa.ServerUrl), "http://") {
		url := strings.Replace(*(sa.ServerUrl), "http://", "", -1)
		sa.ServerUrl = &url
	}

	if strings.Contains(*(sa.ServerUrl), "https://") {
		url := strings.Replace(*(sa.ServerUrl), "https://", "", -1)
		sa.ServerUrl = &url
	}
}

func (sa *serverArguments) formatWorkerUrl() {
	if !strings.Contains(*(sa.WorkerUrl), "http://") && !strings.Contains(*(sa.WorkerUrl), "https://"){
		url := "http://" + *(sa.WorkerUrl)
		sa.WorkerUrl = &url
	}
}

func (sa *serverArguments) formatUrls() {
	sa.formatServerUrl()
	sa.formatWorkerUrl()
}

func (sa *serverArguments) getUrls() *[]string {
	urls := make([]string, 0)
	if sa.ServerUrl != nil {
		urls = append(urls, *(sa.ServerUrl))
	}

	if sa.WorkerUrl != nil {
		urls = append(urls, *(sa.WorkerUrl))
	}

	return &urls
}

func Start() {
	args := parseArgs()

	if args.WorkerUrl == nil {
		logger.Warn("No worker url provided. If a worker is needed to forward messages, use --workerUrl flag")
	}

	args.formatUrls()
	logger.Log("Running sqs mock server on urls:", *(args.getUrls()))

	handler := sqs.RequestHandler{WorkerUrl: args.WorkerUrl}
	router := mux.NewRouter()
	router.HandleFunc(sqs.AddMessageEndpoint, handler.Add).Methods("POST")
	router.HandleFunc(sqs.RetrieveMessageEndpoint, handler.Retrieve).Methods("POST")
	router.HandleFunc(sqs.DeleteMessageEndpoint, handler.Delete).Methods("POST")
	router.HandleFunc(sqs.PrintQueueEndpoint, handler.Print).Methods("GET")

	logger.Error(http.ListenAndServe(*(args.ServerUrl), router))
}

func parseArgs() *serverArguments {
	args := os.Args
	sa := serverArguments{}
	for i := 0; i < len(args); i += 1 {
		a := args[i]
		if strings.Contains(a, "--") {
			arg := strings.ToLower(strings.Replace(a, "-", "", -1))
			if i + 1 < len(args) {
				val := args[i+1]
				switch arg {
				case serverUrlArg:
					sa.ServerUrl = &val
					break
				case workerUrlArg:
					sa.WorkerUrl = &val
					break
				default:
					logger.Warn("Cannot get argument for:", arg, "unknown argument type")
					break
				}
			}

			i += 1
		}
	}

	return &sa
}
