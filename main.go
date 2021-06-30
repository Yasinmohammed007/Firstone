package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"runtime"
	"strings"

	src "github.com/SophosNSG/paloma-config-service/src"

	"github.com/gorilla/mux"
	opentracing "github.com/opentracing/opentracing-go"
	oLog "github.com/opentracing/opentracing-go/log"

	ext "github.com/opentracing/opentracing-go/ext"
	"github.com/prnvkv/my-nats/pkg/util"
	log "github.com/sirupsen/logrus"
	jaeger "github.com/uber/jaeger-client-go"
	config "github.com/uber/jaeger-client-go/config"
)

func init() {
	formatter := new(log.JSONFormatter)
	formatter.CallerPrettyfier = func(f *runtime.Frame) (string, string) {
		s := strings.Split(f.Function, ".")
		funcname := s[len(s)-1]
		_, filename := path.Split(f.File)
		return funcname, filename
	}
	log.SetFormatter(formatter)
	log.SetReportCaller(true)
	log.SetLevel(log.DebugLevel) // Debug for Development Purpose
}

// initJaeger returns an instance of Jaeger Tracer that samples 100% of traces and logs all spans to stdout.
func initJaeger(service string) (opentracing.Tracer, io.Closer) {
	log.Debug("Initializing Jaeger")
	cfg, err := config.FromEnv()
	if err != nil {
		// parsing errors might happen here, such as when we get a string where we expect a number
		log.Printf("Could not parse Jaeger env vars: %s", err.Error())
		panic(fmt.Sprintf("Could not parse Jaeger env vars: %s", err.Error()))
	}

	cfg.ServiceName = service

	log.WithFields(log.Fields{
		"cfg": cfg,
	}).Debug("Jaeger config from ENV: ")

	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	return tracer, closer
}

var dnsSidecarHost string = util.GetEnv("SIDECAR_HOST", "listener.networked")
var dnsSidecarPort string = util.GetEnv("SIDECAR_PORT", "8082")
var DnsSidecarUrl string = fmt.Sprintf("http://%s:%s", dnsSidecarHost, dnsSidecarPort)

func testGetAPI(write http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(write, "Welcome to the API server!!")
}

func getDNSConfigAPI(write http.ResponseWriter, request *http.Request) {
	// Tracing demo
	var traceID string

	tracer, closer := initJaeger("paloma-config-service")
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(request.Header))
	span := tracer.StartSpan("API call to DNS Service", ext.RPCServerOption(spanCtx))
	if sc, ok := span.Context().(jaeger.SpanContext); ok {
		traceID = sc.TraceID().String()
		log.WithFields(log.Fields{
			"traceID": traceID,
		}).Debug("Fetching Trace ID from span")
	}
	defer span.Finish()

	getConfigURI := "/v1/config/dns"

	url := DnsSidecarUrl + getConfigURI

	client := &http.Client{}

	log.WithFields(log.Fields{
		"traceID": traceID,
	}).Debug("Requesting config from URL: ", url)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")

	// Set some tags on the span to annotate that it's the client span. The additional HTTP tags are useful for debugging purposes.
	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, url)
	ext.HTTPMethod.Set(span, "GET")

	//Inject the span context into headers
	tracer.Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header))
	span.LogFields(
		oLog.String("Event", "API call to get DNS entries"),
		oLog.String("URI", url))

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Errored when sending request to the server")
		ext.LogError(span, err)
		log.WithFields(log.Fields{
			"traceID": traceID,
			"error":   err,
		}).Error("Error sending request to the server")
		return
	}

	defer resp.Body.Close()
	resp_body, _ := ioutil.ReadAll(resp.Body)

	log.WithFields(log.Fields{
		"traceID":     traceID,
		"Status Code": resp.StatusCode,
		"Status":      resp.Status,
		"Response":    string(resp_body),
	}).Debug("Received response for DNS entries")

	if resp.StatusCode != 200 {
		msg := "ERROR"
		switch resp.StatusCode {
		case 400:
			msg = "CLIENT_ERROR"
		case 500:
			msg = "SERVER_ERROR"
		}
		ext.Error.Set(span, true)
		span.LogFields(
			oLog.String("Event", "API call to get DNS entries failed"),
			oLog.String("Message", msg))

		log.WithFields(log.Fields{
			"traceID":     traceID,
			"Status Code": resp.StatusCode,
			"Status":      resp.Status,
			"Error":       string(resp_body),
		}).Error("Error getting DNS entries.")

		write.WriteHeader(resp.StatusCode)
		fmt.Fprintf(write, msg)
		return
	}

	fmt.Fprintf(write, string(resp_body))
}

func postDNSConfigAPI(write http.ResponseWriter, request *http.Request) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		panic(err)
	}

	log.Debug("Recieved body as", string(body))
	result, err := src.ParseAndPublishDNSConfig(body)
	if err != nil {
		write.WriteHeader(400)
		fmt.Fprintf(write, err.Error())
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Error occurred while parsing the body")

	} else {
		fmt.Fprintf(write, result)
	}
}

func deleteDNSRecordAPI(write http.ResponseWriter, request *http.Request) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		panic(err)
	}
	log.Println("Recieved body as", string(body))
	result, err := src.ParseAndDeleteDNSRecodConfig(body)
	containsError := strings.Contains(result, "ERROR")
	log.Println("Containes Error", containsError)
	if err != nil || containsError {
		if containsError {
			err = fmt.Errorf(result)
		}
		write.WriteHeader(400)
		fmt.Fprintf(write, err.Error())
		log.Println("Error occurred while requesting")

	} else {
		fmt.Fprintf(write, result)
	}
}

func postDNSRecordAPI(write http.ResponseWriter, request *http.Request) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		panic(err)
	}
	log.Debug("Recieved body as", string(body))
	result, err := src.ParseAndWriteDNSRecodConfig(body)
	containsError := strings.Contains(result, "ERROR")
	log.Println("Containes Error", containsError)
	if err != nil || containsError {
		if containsError {
			err = fmt.Errorf(result)
		}
		write.WriteHeader(400)
		fmt.Fprintf(write, err.Error())
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Error occurred while parsing the body")

	} else {
		fmt.Fprintf(write, result)
	}
}

func handler() {
	router := mux.NewRouter()
	router.HandleFunc("/test", testGetAPI).Methods("GET")

	router.HandleFunc("/v1/config/dns/addHost", postDNSRecordAPI).Methods("POST")
	router.HandleFunc("/v1/config/dns/deleteHost", deleteDNSRecordAPI).Methods("DELETE")
	router.HandleFunc("/v1/config/dns/entry", getDNSConfigAPI).Methods("GET")
	router.HandleFunc("/v1/config/dns/entry", postDNSConfigAPI).Methods("POST")

	log.Println("Starting the HTTP Server")
	log.Fatal(http.ListenAndServe(":8081", router))
}

func main() {
	log.Info("Registering APIs")
	handler()
}
