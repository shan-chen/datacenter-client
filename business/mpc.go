package business

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/datacenter-client/model"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	log "github.com/sirupsen/logrus"
)

func MpcEventHandler(client *channel.Client, notify <-chan *fab.CCEvent, peerName string) {
	taskMap := make(map[string][]model.MpcTask)
	for {
		select {
		case event, ok := <-notify:
			if !ok {
				return
			}
			log.Info("receive mpc event")
			var data model.MpcTask
			if err := json.Unmarshal(event.Payload, &data); err != nil {
				log.WithError(err).Error("unmarshal failed")
				continue
			}
			if data.Sponsor != peerName {
				continue
			}
			taskMap[data.Nonce] = append(taskMap[data.Nonce], data)
			if len(taskMap[data.Nonce]) == 3 {
				var args [][]byte
				for _, item := range taskMap[data.Nonce] {
					args = append(args, []byte(item.Data))
				}
				success := true
				resp, err := Execute(MpcChainCodeID, client, "executeMpcTask", args)
				if err != nil {
					log.WithError(err).Error("executeMpcTask failed")
					success = false
				}
				delete(taskMap, data.Nonce)
				req := make(map[string]interface{})
				req["success"] = success
				req["nonce"] = data.Nonce
				req["result"] = string(resp.Responses[0].Response.Payload)
				reqBytes, _ := json.Marshal(req)
				_, err = http.Post(URLPrefix+"/mpc/callback", "application/json", bytes.NewReader(reqBytes))
				if err != nil {
					log.WithError(err).Error("http post failed")
					continue
				}
			}
		}
	}
}
