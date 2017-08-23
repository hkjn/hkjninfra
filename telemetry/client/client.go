// client.go implements a GRPC Report client.
package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	googletime "github.com/golang/protobuf/ptypes/timestamp"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "hkjn.me/hkjninfra/telemetry/report"
)

const (
	defaultAddr      = "localhost:50051"
	defaultName      = "world"
	defaultFactsPath = "facts.json"
)

var debugging = os.Getenv("REPORT_DEBUGGING") == "true"

func debug(format string, a ...interface{}) {
	if !debugging {
		return
	}
	log.Printf(format, a...)
}

// getAddr returns the address to report in to, given a default.
func getAddr(d string) string {
	addrEnv := os.Getenv("REPORT_ADDR")
	if addrEnv != "" {
		return addrEnv
	}
	return d
}

// getInfo returns the extra info to use when reporting in.
func getInfo(d string) (*pb.ClientInfo, error) {
	factsPath := os.Getenv("REPORT_FACTS_PATH")
	if factsPath == "" {
		factsPath = d
	}
	debug("Reading facts.json from %q..\n", factsPath)
	info := &pb.ClientInfo{}
	f, err := os.Open(factsPath)
	if err != nil {
		return nil, err
	}
	if err := json.NewDecoder(f).Decode(info); err != nil {
		return nil, err
	}
	return info, nil
}

func getClient(addr string) (pb.ReportClient, func() error) {
	// TODO: security.
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return pb.NewReportClient(conn), conn.Close
}

// send reports to the server.
func send(c pb.ReportClient) error {
	info, err := getInfo(defaultFactsPath)
	if err != nil {
		return err
	}
	req := &pb.ReportRequest{
		Ts: &googletime.Timestamp{
			Seconds: time.Now().Unix(),
			Nanos:   int32(time.Now().Nanosecond()),
		},
		Info: info,
	}
	log.Printf("Sending request: %v\n", req)
	r, err := c.Send(context.Background(), req)
	if err != nil {
		return err
	}
	log.Printf("Got message from server: %q", r.Message)
	return nil
}

func main() {
	log.Printf("report_client %s starting..\n", Version)
	addr := getAddr(defaultAddr)

	log.Printf("Contacting server at tcp addr %q..\n", addr)
	c, close := getClient(addr)
	defer close()

	if err := send(c); err != nil {
		log.Fatalf("Could not report: %v", err)
	}
}
