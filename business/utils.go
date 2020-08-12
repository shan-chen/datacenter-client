package business

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	log "github.com/sirupsen/logrus"
)

func Query(client *channel.Client, method string, args [][]byte) (channel.Response, error) {
	log.WithField("method", method).Info("query begin")
	req := channel.Request{
		ChaincodeID: ChainCodeID,
		Fcn:         method,
		Args:        args,
	}
	return client.Query(req, channel.WithTargetEndpoints(Peer))
}

func Execute(client *channel.Client, method string, args [][]byte) (channel.Response, error) {
	log.WithField("method", method).Info("execute begin")
	req := channel.Request{
		ChaincodeID: ChainCodeID,
		Fcn:         method,
		Args:        args,
	}
	return client.Execute(req)
}
