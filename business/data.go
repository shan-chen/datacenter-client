package business

import (
	"bytes"
	"encoding/json"
	"github.com/datacenter-client/model"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	//log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
  "fmt"
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
  Log.Info("query finished")
	if err != nil {
	  Log.WithError(err).Error("query failed")
    callback([]model.Article{}, peerName, keyword)
		return
	}
	if len(resp.Responses) <= 0 || resp.Responses[0] == nil || resp.Responses[0].Response == nil {
		Log.Error("query result is nil")
    callback([]model.Article{}, peerName, keyword)
		return
	}
	var data model.QueryIDsRes
	err = json.Unmarshal(resp.Responses[0].Response.Payload, &data)
	if err != nil {
		Log.WithError(err).Error("unmarshal failed")
		callback([]model.Article{}, peerName, keyword)
    return
	}
  //if len(data.A) == 0 {
	//log.Info("query result empty")
  //continue
	//}
  bytesData := callback(data.A, peerName, keyword)
	args := [][]byte{[]byte(peerName), bytesData, []byte(strconv.FormatInt(time.Now().Unix(), 10))}
	_, err = Execute(DataChainCodeID, client, "logQuery", args)
	if err != nil {
		Log.WithError(err).Error("log query failed")
  }
}

func callback(articles []model.Article, peerName string, keyword string) []byte {
	req := make(map[string]interface{})
	req["ids"] = articles
	req["owner"] = peerName
	req["keyword"] = keyword
	bytesData, _ := json.Marshal(req)
  _, err := http.Post(URLPrefix+"/data/callback", "application/json", bytes.NewReader(bytesData))
	if err != nil {
	  Log.WithError(err).Error("http post failed")
    return nil
	}
  Log.Info("callback finished")
  return bytesData
}
