package gateway

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

type ServiceServerHandle func(context.Context, *gwruntime.ServeMux, *grpc.ClientConn) error

type Service struct {
	Name    string
	Address string
	Handler ServiceServerHandle
}

func NewGateway(ctx context.Context, opts []gwruntime.ServeMuxOption, services []Service) (http.Handler, error) {
	mux := gwruntime.NewServeMux(opts...)

	for _, item := range services {
		// new dial client
		conn, err := dial(ctx, "tcp", item.Address)
		if err != nil {
			log.Printf("[%s] did dial connect", item.Name)
		}
		go func() {
			<-ctx.Done()
			if err = conn.Close(); err != nil {
				log.Printf("[%s] gRPC client connect failed: %v", item.Name, err)
			}
		}()

		// register microservices server
		item.Handler(ctx, mux, conn)
	}

	return mux, nil
}

func dial(ctx context.Context, network, addr string) (*grpc.ClientConn, error) {
	switch network {
	case "tcp":
		return dialTCP(ctx, addr)
	case "unix":
		return dialUnix(ctx, addr)
	default:
		return nil, fmt.Errorf("unsupported network type %q", network)
	}
}

func dialTCP(ctx context.Context, addr string) (*grpc.ClientConn, error) {
	return grpc.DialContext(ctx, addr, grpc.WithInsecure())
}

func dialUnix(ctx context.Context, addr string) (*grpc.ClientConn, error) {
	d := func(ctx context.Context, addr string) (net.Conn, error) {
		return net.Dial("unix", addr)
	}
	return grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithContextDialer(d))
}
