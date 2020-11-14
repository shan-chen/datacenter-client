package business

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
)

func DataEventHandler(client *channel.Client, notify <-chan *fab.CCEvent, peerName string) {
	for {
		select {
		case event, ok := <-notify:
			if ok {
				fmt.Println(time.Now().UnixNano())
				go handler(client, peerName, string(event.Payload))
			} else {
				return
			}
		}
	}
}

func handler(client *channel.Client, peerName string, keyword string) {
	resp, err := Query(DataChainCodeID, client, "queryIDs", [][]byte{[]byte(keyword)}, peerName)
	Log.Infoln("query finished")
	if err != nil {
		Log.WithError(err).Error("query failed")
		callback([]byte{}, peerName, keyword)
		return
	}
	if len(resp.Responses) <= 0 || resp.Responses[0] == nil || resp.Responses[0].Response == nil {
		Log.Error("query result is nil")
		callback([]byte{}, peerName, keyword)
		return
	}

	bytesData := callback(resp.Responses[0].Response.Payload, peerName, keyword)
	args := [][]byte{[]byte(peerName), bytesData, []byte(strconv.FormatInt(time.Now().Unix(), 10))}
	_, err = Execute(DataChainCodeID, client, "logQuery", args)
	if err != nil {
		Log.WithError(err).Error("log query failed")
	}
}

func callback(data []byte, peerName string, keyword string) []byte {
	req := make(map[string]interface{})
	req["payload"] = data
	req["owner"] = peerName
	req["keyword"] = keyword
	bytesData, _ := json.Marshal(req)
	_, err := http.Post(URLPrefix+"/data/callback", "application/json", bytes.NewReader(bytesData))
	if err != nil {
		Log.WithError(err).Error("http post failed")
		return nil
	}
	Log.Infoln("callback finished")
	return bytesData
}
