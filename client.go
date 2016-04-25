package configurationservice

import (
	log "github.com/Sirupsen/logrus"
	pb "github.com/russellchadwick/configurationservice/proto"
	"github.com/russellchadwick/rpc"
	"golang.org/x/net/context"
	"time"
)

// Client is used to speak with the configuration service via rpc
type Client struct{}

// GetConfiguration gets a configuration value by given a key
func (c *Client) GetConfiguration(name string) (*string, error) {
	log.Debug("-> client.GetConfiguration")
	start := time.Now()

	client := rpc.Client{}
	clientConn, err := client.Dial("configuration")
	if err != nil {
		log.WithField("error", err).Error("error during dial")
		return nil, err
	}
	defer func() {
		closeErr := clientConn.Close()
		if closeErr != nil {
			log.WithField("error", closeErr).Error("error during close")
		}
	}()

	grpcClient := pb.NewConfigurationClient(clientConn)

	response, err := grpcClient.GetConfiguration(context.Background(), &pb.ConfigurationRequest{Name: name})
	if err != nil {
		log.WithField("error", err).Error("error from rpc")
		return nil, err
	}

	elapsed := float64(time.Since(start)) / float64(time.Microsecond)
	log.WithField("elapsed", elapsed).Debug("<- client.GetConfiguration")

	return &response.Value, nil
}
