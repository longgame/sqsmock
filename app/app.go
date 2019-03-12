package app

import (
	"github.com/gorilla/mux"
	"github.com/greenac/artemis/logger"
	"github.com/greenac/sqsmock/sqs"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type serverArguments struct {
	ServerUrl *string
	WorkerUrls *[]string
	Delay     int64
}

const (
	delay        = "delay"
	serverUrlArg = "url"
	workerUrlsArg = "workerUrls"
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

func (sa *serverArguments) formatWorkerUrls() {
	for i:=0; i < len(*sa.WorkerUrls); i += 1 {
		url := sa.formatWorkerUrl((*sa.WorkerUrls)[i])
		if url != (*sa.WorkerUrls)[i] {
			(*sa.WorkerUrls)[i] = url
		}
	}
}

func (sa *serverArguments) formatWorkerUrl(url string) string {
	wUrl := url
	if !strings.Contains(url, "http://") && !strings.Contains(url, "https://") {
		wUrl = "http://" + url
	}

	return wUrl
}

func (sa *serverArguments) formatUrls() {
	sa.formatServerUrl()
	sa.formatWorkerUrls()
}

func (sa *serverArguments) getUrls() *[]string {
	urls := make([]string, 0)
	if sa.ServerUrl != nil {
		urls = append(urls, *(sa.ServerUrl))
	}

	if sa.WorkerUrls != nil {
		urls = append(urls, *(sa.WorkerUrls)...)
	}

	return &urls
}

func Start() {
	args := parseArgs()

	if args.WorkerUrls == nil {
		logger.Warn("No worker urls provided. If a worker is needed to forward messages, use --workerUrls flag")
	}

	args.formatUrls()
	logger.Log("Running sqs mock server on urls:", *(args.getUrls()))

	handler := sqs.RequestHandler{WorkerUrls: args.WorkerUrls, Delay: args.Delay}
	router := mux.NewRouter()
	router.HandleFunc(sqs.AddMessageEndpoint, handler.Add).Methods("POST")
	router.HandleFunc(sqs.RetrieveMessageEndpoint, handler.Retrieve).Methods("POST")
	router.HandleFunc(sqs.DeleteMessageEndpoint, handler.Delete).Methods("POST")
	router.HandleFunc(sqs.PrintQueueEndpoint, handler.Print).Methods("GET")

	logger.Error(http.ListenAndServe(*(args.ServerUrl), router))
}

func parseWorkerUrls(urls string) []string {
	parts := strings.Split(urls, ",")
	pUrls := make([]string, len(parts))
	for i, p := range parts {
		pUrls[i] = strings.Replace(p, " ", "", -1)
	}

	return pUrls
}

func parseArgs() *serverArguments {
	args := os.Args
	sa := serverArguments{}
	for i := 0; i < len(args); i += 1 {
		a := args[i]
		if strings.Contains(a, "--") {
			arg := strings.Replace(a, "-", "", -1)
			if i+1 < len(args) {
				val := args[i+1]
				switch arg {
				case delay:
					logger.Log("Delay:", val)
					d, err := strconv.Atoi(val)
					if err != nil {
						logger.Error("Could not convert arg:", val, "to integer for delay")
						panic(err)
					}
					sa.Delay = int64(d)
				case serverUrlArg:
					sa.ServerUrl = &val
					break
				case workerUrlsArg:
					wUrls := parseWorkerUrls(val)
					sa.WorkerUrls = &wUrls
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
