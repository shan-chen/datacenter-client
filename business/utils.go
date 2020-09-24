package business

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

func Query(cc string, client *channel.Client, method string, args [][]byte) (channel.Response, error) {
	Log.WithField("method", method).Info("query begin")
	req := channel.Request{
		ChaincodeID: cc,
		Fcn:         method,
		Args:        args,
	}
	return client.Query(req, channel.WithTargetEndpoints(Peer))
}

func Execute(cc string, client *channel.Client, method string, args [][]byte) (channel.Response, error) {
	Log.WithField("method", method).Info("execute begin")
	req := channel.Request{
		ChaincodeID: cc,
		Fcn:         method,
		Args:        args,
	}
	return client.Execute(req, channel.WithTargetEndpoints("peer1.org1.example.com","peer2.org2.example.com","peer3.org3.example.com"))
}
