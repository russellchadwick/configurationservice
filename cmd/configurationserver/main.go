package main

import (
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	consulapi "github.com/hashicorp/consul/api"
	pb "github.com/russellchadwick/configurationservice/proto"
	"github.com/russellchadwick/rpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type server struct{}

func (s *server) GetConfiguration(ctx context.Context, in *pb.ConfigurationRequest) (*pb.ConfigurationResponse, error) {
	log.WithField("name", in.Name).Info("GetConfiguration")

	client, err := newConsulAPIClient()
	if err != nil {
		return nil, err
	}

	kvPair, _, err := client.KV().Get("configuration/"+in.Name, nil)
	if err != nil {
		return nil, err
	}

	if kvPair == nil {
		log.WithField("name", in.Name).Warn("Not found")
		return nil, errors.New("Not found")
	}

	log.WithField("name", in.Name).Info("Found")
	return &pb.ConfigurationResponse{Value: string(kvPair.Value)}, nil
}

func main() {
	rpcServer := rpc.Server{}
	go serve(&rpcServer)
	defer func() {
		err := rpcServer.Stop()
		if err != nil {
			log.WithField("error", err).Error("error during stop")
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	signal.Notify(signalChan, syscall.SIGTERM)
	<-signalChan
	log.Info("Received shutdown signal")
}

func serve(rpcServer *rpc.Server) {
	err := rpcServer.Serve("configuration", func(grpcServer *grpc.Server) {
		pb.RegisterConfigurationServer(grpcServer, &server{})
	})
	if err != nil {
		log.WithField("error", err).Error("error from rpc serve")
	}
}

func newConsulAPIClient() (*consulapi.Client, error) {
	config := consulapi.DefaultConfig()
	config.HttpClient.Timeout = 2 * time.Second
	client, err := consulapi.NewClient(config)
	if err != nil {
		return nil, err
	}
	return client, nil
}
