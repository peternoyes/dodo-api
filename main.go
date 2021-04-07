package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var httpFlag = flag.String("http", ":3000", "Listen for HTTP connections on this address.")

var svc *dynamodb.DynamoDB

func main() {
	flag.Parse()

	printServingAt(*httpFlag)

	svc = dynamodb.New(session.New(&aws.Config{Region: aws.String("us-west-2")}))

	router := NewRouter()
	//router.PathPrefix("/assets/").Handler(httpgzip.FileServer(assets, httpgzip.FileServerOptions{ServeError: httpgzip.Detailed}))
	//router.NotFoundHandler = http.NotFoundHandler()
	log.Fatal(http.ListenAndServe(":3000", router))
}

func printServingAt(addr string) {
	hostPort := addr
	if strings.HasPrefix(hostPort, ":") {
		hostPort = "localhost" + hostPort
	}
	fmt.Printf("serving at http://%s/\n", hostPort)
}
