package business

import "github.com/sirupsen/logrus"

const (
	UserName        = "Admin"
	OrgName         = "Org1"
	ChannelID       = "mych"
	DataChainCodeID = "ss"
	MpcChainCodeID  = "mpc"
	DataEventName   = "data"
	MpcEventName    = "mpc"
	Peer            = "peer1.org1.example.com"
	URLPrefix       = "http://localhost:8888"
)

var Log *logrus.Logger

func init() {
  Log = logrus.New()
  Log.SetFormatter(&logrus.TextFormatter{
    FullTimestamp: true,
  })
}
