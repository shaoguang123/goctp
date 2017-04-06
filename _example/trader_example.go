package main

import (
	"flag"
	"github.com/shaoguang123/goctp"
	"github.com/shaoguang123/goctp/goctp_api"
	"log"
	"os"
)

var (
	broker_id    = flag.String("BrokerID", "9999", "经纪公司编号,SimNow BrokerID统一为：9999")
	investor_id  = flag.String("InvestorID", "<InvestorID>", "交易用户代码")
	pass_word    = flag.String("Password", "<Password>", "交易用户密码")
	market_front = flag.String("MarketFront", "tcp://180.168.146.187:10031", "行情前置,SimNow的测试环境: tcp://180.168.146.187:10031")
	trade_front  = flag.String("TradeFront", "tcp://180.168.146.187:10030", "交易前置,SimNow的测试环境: tcp://180.168.146.187:10030")
)

var CTP goctp_api.GoCTPClient

func init() {
	log.SetFlags(log.LstdFlags)
	log.SetPrefix("CTP: ")
}

func main() {

	if len(os.Args) < 2 {
		log.Fatal("usage: ./goctp_trader_example -BrokerID 9999 -InvestorID 000000 -Password 000000 -MarketFront tcp://180.168.146.187:10010 -TradeFront tcp://180.168.146.187:10000")
	}

	flag.Parse()

	CTP = goctp_api.GoCTPClient{
		BrokerID:   *broker_id,
		InvestorID: *investor_id,
		Password:   *pass_word,

		MdFront: *market_front,
		MdApi:   goctp.CThostFtdcMdApiCreateFtdcMdApi(),

		TraderFront: *trade_front,
		TraderApi:   goctp.CThostFtdcTraderApiCreateFtdcTraderApi(),

		MdRequestID:     0,
		TraderRequestID: 0,
	}

	pTraderSpi := goctp.NewDirectorCThostFtdcTraderSpi(&goctp_api.GoCThostFtdcTraderSpi{Client: CTP})

	CTP.TraderApi.RegisterSpi(pTraderSpi)                         // 注册事件类
	CTP.TraderApi.SubscribePublicTopic(1 /*THOST_TERT_RESTART*/)  // 注册公有流
	CTP.TraderApi.SubscribePrivateTopic(1 /*THOST_TERT_RESTART*/) // 注册私有流
	CTP.TraderApi.RegisterFront(CTP.TraderFront)

	CTP.TraderApi.Init()

	CTP.TraderApi.Join()
	CTP.TraderApi.Release()
}
