package main

import (
	"flag"

	"github.com/datacenter-client/business"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

func main() {
	var configPath string
	var orgName string
	var peerName string
	flag.StringVar(&configPath, "c", "", "config file path")
	flag.StringVar(&orgName, "o", "", "org name")
	flag.StringVar(&peerName, "p", "", "peer name")
	flag.Parse()
	sdk, err := fabsdk.New(config.FromFile(configPath))
	if err != nil {
		business.Log.WithError(err).Error("cannot load config file")
		return
	}
	ccp := sdk.ChannelContext(business.ChannelID, fabsdk.WithOrg(orgName), fabsdk.WithUser(business.UserName))
	channelClient, err := channel.New(ccp)
	if err != nil {
		business.Log.WithError(err).Error("cannot get channel client")
		return
	}
	ec, err := event.New(ccp, event.WithBlockEvents())
	if err != nil {
		business.Log.WithError(err).Error("cannot get event client")
		return
	}
	dataReg, dataNotify, err := ec.RegisterChaincodeEvent(business.DataChainCodeID, business.DataEventName)
	if err != nil {
		business.Log.WithError(err).Error("cannot register data event")
		return
	}
	mpcReg, mpcNotify, err := ec.RegisterChaincodeEvent(business.MpcChainCodeID, business.MpcEventName)
	if err != nil {
		business.Log.WithError(err).Error("cannot register mpc event")
		return
	}
	business.Log.Info("register success")
	defer ec.Unregister(dataReg)
	defer ec.Unregister(mpcReg)

	go business.DataEventHandler(channelClient, dataNotify, peerName)
	go business.MpcEventHandler(channelClient, mpcNotify, peerName)
	select {}
}
