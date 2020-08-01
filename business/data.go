package business

import (
	"bytes"
	"encoding/json"
	"github.com/datacenter-client/model"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

func EventHandler(client *channel.Client, notify <-chan *fab.CCEvent, peerName string) {
	for {
		select {
		case event, ok := <-notify:
			if ok {
				resp, err := Query(client, "queryIDs", [][]byte{event.Payload})
				if err != nil {
					log.WithError(err).Error("query failed")
					continue
				}
				if len(resp.Responses) <= 0 || resp.Responses[0] == nil || resp.Responses[0].Response == nil {
					log.Error("query result is nil")
					continue
				}
				var data model.QueryIDsRes
				err = json.Unmarshal(resp.Responses[0].Response.Payload, &data)
				if err != nil {
					log.WithError(err).Error("unmarshal failed")
					continue
				}
				if len(data.A) == 0 {
					continue
				}
				req := make(map[string]interface{})
				req["ids"] = data.A
				req["owner"] = peerName
				req["keyword"] = string(event.Payload)
				bytesData, _ := json.Marshal(req)
				_, err = http.Post(URLPrefix+"/data/callback", "application/json", bytes.NewReader(bytesData))
				if err != nil {
					log.WithError(err).Error("http post failed")
					continue
				}
				args := [][]byte{[]byte(peerName), bytesData, []byte(strconv.FormatInt(time.Now().Unix(), 10))}
				_, err = Execute(client, "logQuery", args)
				if err != nil {
					log.WithError(err).Error("log query failed")
				}
			} else {
				return
			}
		}
	}
}

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
