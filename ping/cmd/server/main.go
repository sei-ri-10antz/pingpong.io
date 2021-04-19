package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/nats-io/nats.go"
	"github.com/sei-ri/pingpong.io/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	var srv Server
	if err := envconfig.Process("", &srv); err != nil {
		log.Fatal(err)
	}
	srv.Serve(context.Background())
}

const Topic = "ping"

type Server struct {
	Name string `envconfig:"NAME" default:"service"`
	Host string `envconfig:"HOST" default:"0.0.0.0"`
	Port int    `envconfig:"PORT" default:"10010"`

	// 3rd party service
	ServiceName string `envconfig:"SERVICE_NAME" default:":service"`
	ServiceAddr string `envconfig:"SERVICE_ADDR" default:":10010"`

	NatsURL string `envconfig:"NATS_URL" default:"nats://127.0.0.1:4222"`

	log      *log.Logger
	eventbus *nats.Conn
}

func (s *Server) Ping(ctx context.Context, req *protos.Request) (*protos.Response, error) {
	var resp *protos.Response

	switch req.Path {
	case protos.Request_none:
		resp = &protos.Response{
			Result: fmt.Sprintf("[%s] pong!", s.Name),
		}
	case protos.Request_nats:
		subj := fmt.Sprintf("%s.%s", s.ServiceName, Topic)
		payload := []byte("ping")
		msg, err := s.eventbus.Request(subj, payload, 2*time.Second)
		if err != nil {
			return nil, err
		}
		// s.log.Printf("Published [%s]: %s", subj, payload)
		// s.log.Printf("Received [%s]: %s", msg.Subject, string(msg.Data))
		resp = &protos.Response{
			Result: fmt.Sprintf("[%s] %s!", s.Name, string(msg.Data)),
		}
	case protos.Request_client:
		conn, err := grpc.Dial(s.ServiceAddr, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
		defer conn.Close()

		client := protos.NewServiceClient(conn)
		ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()

		resp, err = client.Ping(ctx, &protos.Request{
			Path: protos.Request_Path(protos.Request_none),
		})
		if err != nil {
			return nil, err
		}
	}

	return resp, nil
}

func (s *Server) Serve(ctx context.Context) {
	s.log = log.New(os.Stdout, fmt.Sprintf("[%s] ", s.Name), log.LstdFlags)

	if v, err := nats.Connect(s.NatsURL); err != nil {
		s.log.Fatal(err)
	} else {
		s.eventbus = v
	}

	subj := fmt.Sprintf("%s.%s", s.ServiceName, Topic)
	s.eventbus.QueueSubscribe(subj, s.Name, func(msg *nats.Msg) {
		msg.Respond([]byte("pong"))
	})
	s.eventbus.Flush()
	if err := s.eventbus.LastError(); err != nil {
		s.log.Println("eventbus subscribe err: ", err)
	}

	srv := grpc.NewServer()

	protos.RegisterServiceServer(srv, s)
	reflection.Register(srv)

	lis, err := net.Listen("tcp", net.JoinHostPort(s.Host, strconv.Itoa(s.Port)))
	if err != nil {
		s.log.Fatal(err)
	}
	defer lis.Close()

	s.log.Println("gRPC server listening at", lis.Addr())

	if err := srv.Serve(lis); err != nil {
		s.log.Fatal(err)
	}
}
