package goctp_api

import (
	"github.com/shaoguang123/goctp"
	"log"
	"reflect"
	"time"
)

type GoCTPClient struct {
	BrokerID   string
	InvestorID string
	Password   string

	MdFront string
	MdApi   goctp.CThostFtdcMdApi

	TraderFront string
	TraderApi   goctp.CThostFtdcTraderApi

	MdRequestID     int
	TraderRequestID int
}

func (g *GoCTPClient) GetMdRequestID() int {
	g.MdRequestID += 1
	return g.MdRequestID
}

func (g *GoCTPClient) GetTraderRequestID() int {
	g.TraderRequestID += 1
	return g.TraderRequestID
}

type GoCThostFtdcMdSpi struct {
	Client GoCTPClient
}

type GoCThostFtdcTraderSpi struct {
	Client GoCTPClient

	tradingDate string
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////
//##########################################################################################################//
//////////////////////////////////////////////////////////////////////////////////////////////////////////////
///行情所需接口

///错误应答
func (p *GoCThostFtdcMdSpi) OnRspError(pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	log.Println("GoCThostFtdcMdSpi.OnRspError.")
	p.IsErrorRspInfo(pRspInfo)
}

///当客户端与交易后台通信连接断开时，该方法被调用。当发生这个情况后，API会自动重新连接，客户端可不做处理。
///@param nReason 错误原因
///        0x1001 网络读失败
///        0x1002 网络写失败
///        0x2001 接收心跳超时
///        0x2002 发送心跳失败
///        0x2003 收到错误报文
func (p *GoCThostFtdcMdSpi) OnFrontDisconnected(nReason int) {
	log.Printf("GoCThostFtdcMdSpi.OnFrontDisconnected: %#v\n", nReason)
}

///心跳超时警告。当长时间未收到报文时，该方法被调用。
///@param nTimeLapse 距离上次接收报文的时间
func (p *GoCThostFtdcMdSpi) OnHeartBeatWarning(nTimeLapse int) {
	log.Printf("GoCThostFtdcMdSpi.OnHeartBeatWarning: %v", nTimeLapse)
}

///当客户端与交易后台建立起通信连接时（还未登录前），该方法被调用。
func (p *GoCThostFtdcMdSpi) OnFrontConnected() {
	log.Println("GoCThostFtdcMdSpi.OnFrontConnected.")
	p.ReqUserLogin()
}

///用户登录请求
func (p *GoCThostFtdcMdSpi) ReqUserLogin() {
	log.Println("GoCThostFtdcMdSpi.ReqUserLogin.")

	req := goctp.NewCThostFtdcReqUserLoginField()
	req.SetBrokerID(p.Client.BrokerID)
	req.SetUserID(p.Client.InvestorID)
	req.SetPassword(p.Client.Password)

	iResult := p.Client.MdApi.ReqUserLogin(req, p.Client.GetMdRequestID())

	if iResult != 0 {
		log.Println("发送用户登录请求: 失败.")
	} else {
		log.Println("发送用户登录请求: 成功.")
	}
}

///错误判断
func (p *GoCThostFtdcMdSpi) IsErrorRspInfo(pRspInfo goctp.CThostFtdcRspInfoField) bool {
	// 如果ErrorID != 0, 说明收到了错误的响应
	bResult := (pRspInfo.GetErrorID() != 0)
	if bResult {
		log.Printf("ErrorID=%v ErrorMsg=%v\n", pRspInfo.GetErrorID(), pRspInfo.GetErrorMsg())
	}
	return bResult
}

///登录请求响应
func (p *GoCThostFtdcMdSpi) OnRspUserLogin(pRspUserLogin goctp.CThostFtdcRspUserLoginField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {

	if bIsLast && !p.IsErrorRspInfo(pRspInfo) {

		log.Printf("获取当前版本信息: %#v\n", goctp.CThostFtdcTraderApiGetApiVersion())
		log.Printf("获取当前交易日期: %#v\n", p.Client.MdApi.GetTradingDay())
		log.Printf("获取用户登录信息: %#v %#v %#v\n", pRspUserLogin.GetLoginTime(), pRspUserLogin.GetSystemName(), pRspUserLogin.GetSessionID())

		ppInstrumentID := []string{"IF1706", "IF1709"}

		p.SubscribeMarketData(ppInstrumentID)
		//p.SubscribeForQuoteRsp(ppInstrumentID)
	}
}

///订阅行情
func (p *GoCThostFtdcMdSpi) SubscribeMarketData(name []string) {

	iResult := p.Client.MdApi.SubscribeMarketData(name)

	if iResult != 0 {
		log.Println("发送行情订阅请求: 失败.")
	} else {
		log.Println("发送行情订阅请求: 成功.")
	}
}

///订阅询价
func (p *GoCThostFtdcMdSpi) SubscribeForQuoteRsp(name []string) {

	iResult := p.Client.MdApi.SubscribeForQuoteRsp(name)

	if iResult != 0 {
		log.Println("发送询价订阅请求: 失败.")
	} else {
		log.Println("发送询价订阅请求: 成功.")
	}
}

///订阅行情应答
func (p *GoCThostFtdcMdSpi) OnRspSubMarketData(pSpecificInstrument goctp.CThostFtdcSpecificInstrumentField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	log.Printf("GoCThostFtdcMdSpi.OnRspSubMarketData: %#v %#v %#v\n", pSpecificInstrument.GetInstrumentID(), nRequestID, bIsLast)
	p.IsErrorRspInfo(pRspInfo)
}

///订阅询价应答
func (p *GoCThostFtdcMdSpi) OnRspSubForQuoteRsp(pSpecificInstrument goctp.CThostFtdcSpecificInstrumentField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	log.Printf("GoCThostFtdcMdSpi.OnRspSubForQuoteRsp: %#v %#v %#v\n", pSpecificInstrument.GetInstrumentID(), nRequestID, bIsLast)
	p.IsErrorRspInfo(pRspInfo)
}

///取消订阅行情应答
func (p *GoCThostFtdcMdSpi) OnRspUnSubMarketData(pSpecificInstrument goctp.CThostFtdcSpecificInstrumentField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	log.Printf("GoCThostFtdcMdSpi.OnRspUnSubMarketData: %#v %#v %#v\n", pSpecificInstrument.GetInstrumentID(), nRequestID, bIsLast)
	p.IsErrorRspInfo(pRspInfo)
}

///取消订阅询价应答
func (p *GoCThostFtdcMdSpi) OnRspUnSubForQuoteRsp(pSpecificInstrument goctp.CThostFtdcSpecificInstrumentField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	log.Printf("GoCThostFtdcMdSpi.OnRspUnSubForQuoteRsp: %#v %#v %#v\n", pSpecificInstrument.GetInstrumentID(), nRequestID, bIsLast)
	p.IsErrorRspInfo(pRspInfo)
}

///深度行情通知
func (p *GoCThostFtdcMdSpi) OnRtnDepthMarketData(pDepthMarketData goctp.CThostFtdcDepthMarketDataField) {

	log.Println("GoCThostFtdcMdSpi.OnRtnDepthMarketData: ", pDepthMarketData.GetTradingDay(), "\tInstrumentID: ",
		pDepthMarketData.GetInstrumentID(), "\tExchangeID: ",
		pDepthMarketData.GetExchangeID(), "\tExchangeInstID: ",
		pDepthMarketData.GetExchangeInstID(), "\tLastPrice: ",
		pDepthMarketData.GetLastPrice(), "\tPreSettlementPrice: ",
		pDepthMarketData.GetPreSettlementPrice(), "\tPreClosePrice: ",
		pDepthMarketData.GetPreClosePrice(), "\tPreOpenInterest: ",
		pDepthMarketData.GetPreOpenInterest(), "\tOpenPrice: ",
		pDepthMarketData.GetOpenPrice(), "\tHighestPrice: ",
		pDepthMarketData.GetHighestPrice(), "\tLowestPrice: ",
		pDepthMarketData.GetLowestPrice(), "\tVolume: ",
		pDepthMarketData.GetVolume(), "\tTurnover: ",
		pDepthMarketData.GetTurnover(), "\tOpenInterest: ",
		pDepthMarketData.GetOpenInterest())
}

///询价通知
func (p *GoCThostFtdcMdSpi) OnRtnForQuoteRsp(pForQuoteRsp goctp.CThostFtdcForQuoteRspField) {
	log.Printf("GoCThostFtdcMdSpi.OnRtnForQuoteRsp: %#v\n", pForQuoteRsp)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////
//##########################################################################################################//
//////////////////////////////////////////////////////////////////////////////////////////////////////////////
///以下均为交易所需接口
func (p *GoCThostFtdcTraderSpi) IsErrorRspInfo(pRspInfo goctp.CThostFtdcRspInfoField) bool {

	iResult := (pRspInfo.GetErrorID() != 0)

	if iResult && pRspInfo.GetErrorID() != 0 {
		log.Printf("ErrorID=%v ErrorMsg=%v\n", pRspInfo.GetErrorID(), pRspInfo.GetErrorMsg())
	}
	return iResult
}

///判断接口内容为空
func (p *GoCThostFtdcTraderSpi) isEmpty(a interface{}) bool {
	v := reflect.ValueOf(a)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v.Interface() == reflect.Zero(v.Type()).Interface()
}

func (p *GoCThostFtdcTraderSpi) OnRspError(pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	log.Println("GoCThostFtdcTraderSpi.OnRspError.")
	p.IsErrorRspInfo(pRspInfo)
}

func (p *GoCThostFtdcTraderSpi) OnFrontDisconnected(nReason int) {
	log.Printf("GoCThostFtdcTraderSpi.OnFrontDisconnected: %#v", nReason)
}

func (p *GoCThostFtdcTraderSpi) OnHeartBeatWarning(nTimeLapse int) {
	log.Printf("GoCThostFtdcTraderSpi.OnHeartBeatWarning: %#v", nTimeLapse)
}

func (p *GoCThostFtdcTraderSpi) OnFrontConnected() {
	log.Println("GoCThostFtdcTraderSpi.OnFrontConnected.")
	p.ReqAuthenticate()
}

///客户端认证
func (p *GoCThostFtdcTraderSpi) ReqAuthenticate() {
	log.Println("GoCThostFtdcTraderSpi.ReqAuthenticate.")

	req := goctp.NewCThostFtdcReqAuthenticateField()
	req.SetBrokerID(p.Client.BrokerID)
	req.SetUserID(p.Client.InvestorID)
	req.SetUserProductInfo("JY95000165")
	req.SetAuthCode("NUM6DX8QK8DS39N0")

	iResult := p.Client.TraderApi.ReqAuthenticate(req, p.Client.GetTraderRequestID())

	if iResult != 0 {
		log.Println("客户端认证请求: 失败.")
	} else {
		log.Println("客户端认证请求: 成功.")
	}
}

///客户端认证应答
func (p *GoCThostFtdcTraderSpi) OnRspAuthenticate(pRspAuthenticateField goctp.CThostFtdcRspAuthenticateField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {

	log.Println("GoCThostFtdcTraderSpi.OnRspAuthenticate.")
	if bIsLast && (p.isEmpty(pRspInfo) || !p.IsErrorRspInfo(pRspInfo)) {
		log.Println("客户端认证成功")
		p.ReqUserLogin()
	}
}

func (p *GoCThostFtdcTraderSpi) ReqUserLogin() {
	log.Println("GoCThostFtdcTraderSpi.ReqUserLogin.")

	req := goctp.NewCThostFtdcReqUserLoginField()
	req.SetBrokerID(p.Client.BrokerID)
	req.SetUserID(p.Client.InvestorID)
	req.SetPassword(p.Client.Password)

	iResult := p.Client.TraderApi.ReqUserLogin(req, p.Client.GetTraderRequestID())

	if iResult != 0 {
		log.Println("发送用户登录请求: 失败.")
	} else {
		log.Println("发送用户登录请求: 成功.")
	}
}

func (p *GoCThostFtdcTraderSpi) OnRspUserLogin(pRspUserLogin goctp.CThostFtdcRspUserLoginField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {

	log.Println("GoCThostFtdcTraderSpi.OnRspUserLogin.")
	if p.isEmpty(pRspInfo) || !p.IsErrorRspInfo(pRspInfo) {
		p.tradingDate = pRspUserLogin.GetTradingDay()
		log.Printf("获取用户登录信息: %#v %#v %#v\n", pRspUserLogin.GetFrontID(), pRspUserLogin.GetSessionID(), pRspUserLogin.GetMaxOrderRef())

		///投资者结算结果确认
		if bIsLast {
			p.ReqQrySettlementInfoConfirm()
		}

	}
}

///请求查询结算单确认日期
func (p *GoCThostFtdcTraderSpi) ReqQrySettlementInfoConfirm() {
	req := goctp.NewCThostFtdcQrySettlementInfoConfirmField()
	req.SetBrokerID(p.Client.BrokerID)
	req.SetInvestorID(p.Client.InvestorID)
	for {
		iResult := p.Client.TraderApi.ReqQrySettlementInfoConfirm(req, p.Client.GetTraderRequestID())
		if iResult == 0 {
			log.Println("请求查询结算单确认日期: 成功, iResult=", iResult)
			break
		} else {
			log.Println("请求查询结算单确认日期: 受到流控, iResult=", iResult)
			time.Sleep(1 * time.Second)
		}
	}
}

///结算单确认日期查询请求应答（若当日已查询过结算单则可以直接进行相关交易操作）
func (p *GoCThostFtdcTraderSpi) OnRspQrySettlementInfoConfirm(pSettlementInfoConfirm goctp.CThostFtdcSettlementInfoConfirmField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	log.Println("GoCThostFtdcTraderSpi.OnRspQrySettlementInfoConfirm.")
	if bIsLast && (p.isEmpty(pRspInfo) || !p.IsErrorRspInfo(pRspInfo)) {
		if !p.isEmpty(pSettlementInfoConfirm) {
			log.Println(pSettlementInfoConfirm.GetConfirmDate())
			log.Println(pSettlementInfoConfirm.GetConfirmTime())

			lastConfirmData := pSettlementInfoConfirm.GetConfirmDate()
			if lastConfirmData != p.tradingDate {
				p.ReqQrySettlementInfo()
			} else {
				log.Println("添加想要查询或执行的操作")
				p.ReqQryTradingAccount()
				//p.ReqQryInvestorPosition("")
				//p.ReqQryInvestorPositionDetail("")
				//p.ReqQryInvestorPositionCombineDetail("")
				//p.ReqOrderInsert()
				//p.ReqQryOrder()
				//p.ReqParkedOrderInsert()
				//p.ReqQryParkedOrder()
				//p.ReqRemoveParkedOrder()
				//p.ReqQryParkedOrderAction()
				//p.ReqQryInstrument("")
				//p.ReqRemoveParkedOrder("           1")

			}
		} else {
			p.ReqQrySettlementInfo()
		}

	}
}

///请求查询结算单
func (p *GoCThostFtdcTraderSpi) ReqQrySettlementInfo() {
	req := goctp.NewCThostFtdcQrySettlementInfoField()
	req.SetBrokerID(p.Client.BrokerID)
	req.SetInvestorID(p.Client.InvestorID)

	for {
		iResult := p.Client.TraderApi.ReqQrySettlementInfo(req, p.Client.GetTraderRequestID())
		if iResult == 0 {
			log.Println("请求查询结算单: 成功, iResult=", iResult)
			break
		} else {
			log.Println("请求查询结算单: 受到流控, iResult=", iResult)
			time.Sleep(1 * time.Second)
		}
	}
}

///查询结算单请求应答
func (p *GoCThostFtdcTraderSpi) OnRspQrySettlementInfo(pSettlementInfo goctp.CThostFtdcSettlementInfoField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	log.Println("GoCThostFtdcTraderSpi.OnRspQrySettlementInfo.")
	if p.isEmpty(pRspInfo) || !p.IsErrorRspInfo(pRspInfo) {
		if !p.isEmpty(pSettlementInfo) {
			log.Println("查询结算单")
		}
		//确认结算单
		if bIsLast {
			p.ReqSettlementInfoConfirm()
		}
	}
}

///投资者结算结果确认
func (p *GoCThostFtdcTraderSpi) ReqSettlementInfoConfirm() {
	req := goctp.NewCThostFtdcSettlementInfoConfirmField()
	req.SetBrokerID(p.Client.BrokerID)
	req.SetInvestorID(p.Client.InvestorID)

	iResult := p.Client.TraderApi.ReqSettlementInfoConfirm(req, p.Client.GetTraderRequestID())
	if iResult == 0 {
		log.Println("投资者结算结果确认: 成功, iResult=", iResult)
	} else {
		log.Println("投资者结算结果确认: 失败, iResult=", iResult)
	}
}

///投资者结算结果确认应答（之后便可执行想要的操作）
func (p *GoCThostFtdcTraderSpi) OnRspSettlementInfoConfirm(pSettlementInfoConfirm goctp.CThostFtdcSettlementInfoConfirmField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	log.Println("GoCThostFtdcTraderSpi.OnRspSettlementInfoConfirm.")
	if p.isEmpty(pRspInfo) || !p.IsErrorRspInfo(pRspInfo) {
		if !p.isEmpty(pSettlementInfoConfirm) {
			log.Println("ConfirmTime: ", pSettlementInfoConfirm.GetConfirmTime())
		}
		if bIsLast {
			log.Println("仅每天第一次启动时执行")
			log.Println("添加想要查询或执行的操作")
			//p.ReqQryInvestorPosition("")
		}

	}
}

///p.ReqQryInvestorPosition("")空字符串表示查询全部持仓
func (p *GoCThostFtdcTraderSpi) ReqQryInvestorPosition(InstrumentID string) {
	req := goctp.NewCThostFtdcQryInvestorPositionField()
	req.SetBrokerID(p.Client.BrokerID)
	req.SetInvestorID(p.Client.InvestorID)
	req.SetInstrumentID(InstrumentID)

	for {
		iResult := p.Client.TraderApi.ReqQryInvestorPosition(req, p.Client.GetTraderRequestID())
		if iResult == 0 {
			log.Printf("--->>> 请求查询投资者持仓: 成功 %#v\n", iResult)
			break
		} else {
			log.Printf("--->>> 请求查询投资者持仓: 受到流控 %#v\n", iResult)
			time.Sleep(1 * time.Second)
		}
	}
}

func (p *GoCThostFtdcTraderSpi) OnRspQryInvestorPosition(pInvestorPosition goctp.CThostFtdcInvestorPositionField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	log.Println("GoCThostFtdcTraderSpi.OnRspQryInvestorPosition.")

	if p.isEmpty(pRspInfo) || !p.IsErrorRspInfo(pRspInfo) {
		if !p.isEmpty(pInvestorPosition) {
			log.Println("#################################################################")
			log.Println("YdPosition:", pInvestorPosition.GetYdPosition())
			log.Println("Position:", pInvestorPosition.GetPosition())
			log.Println("InstrumentID:", pInvestorPosition.GetInstrumentID())
			log.Println("TodayPosition:", pInvestorPosition.GetTodayPosition())
		} else {
			log.Println("kong")
		}
	}
}

///p.ReqQryInvestorPositionDetail("")空字符串表示查询全部持仓
func (p *GoCThostFtdcTraderSpi) ReqQryInvestorPositionDetail(InstrumentID string) {
	req := goctp.NewCThostFtdcQryInvestorPositionDetailField()
	req.SetBrokerID(p.Client.BrokerID)
	req.SetInvestorID(p.Client.InvestorID)
	req.SetInstrumentID(InstrumentID)

	for {
		iResult := p.Client.TraderApi.ReqQryInvestorPositionDetail(req, p.Client.GetTraderRequestID())

		if iResult == 0 {
			log.Printf("--->>> 请求查询投资者持仓详情: 成功 %#v\n", iResult)
			break
		} else {
			log.Printf("--->>> 请求查询投资者持仓详情: 受到流控 %#v\n", iResult)
			time.Sleep(1 * time.Second)
		}
	}
}

func (p *GoCThostFtdcTraderSpi) OnRspQryInvestorPositionDetail(pInvestorPosition goctp.CThostFtdcInvestorPositionDetailField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	log.Println("GoCThostFtdcTraderSpi.OnRspQryInvestorPositionDetail.")

	if p.isEmpty(pRspInfo) || !p.IsErrorRspInfo(pRspInfo) {
		if !p.isEmpty(pInvestorPosition) {
			log.Println("#################################################################")
			log.Println("InstrumentID:", pInvestorPosition.GetInstrumentID())
			log.Println("Direction:", pInvestorPosition.GetDirection())
			log.Println("Volume:", pInvestorPosition.GetVolume())
		} else {
			log.Println("kong")
		}
	}
}

//p.ReqQryInvestorPositionCombineDetail("")空字符串表示查询全部组合持仓
func (p *GoCThostFtdcTraderSpi) ReqQryInvestorPositionCombineDetail(CombInstrumentID string) {
	req := goctp.NewCThostFtdcQryInvestorPositionCombineDetailField()
	req.SetBrokerID(p.Client.BrokerID)
	req.SetInvestorID(p.Client.InvestorID)
	req.SetCombInstrumentID(CombInstrumentID)

	for {
		iResult := p.Client.TraderApi.ReqQryInvestorPositionCombineDetail(req, p.Client.GetTraderRequestID())
		if iResult == 0 {
			log.Printf("--->>> 请求查询投资者组合持仓详情: 成功 %#v\n", iResult)
			break
		} else {
			log.Printf("--->>> 请求查询投资者组合持仓详情: 失败 %#v\n", iResult)
			time.Sleep(1 * time.Second)
		}
	}
}

func (p *GoCThostFtdcTraderSpi) OnRspQryInvestorPositionCombineDetail(pInvestorPosition goctp.CThostFtdcInvestorPositionCombineDetailField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	log.Println("GoCThostFtdcTraderSpi.OnRspQryInvestorPositionCombineDetail.")

	if p.isEmpty(pRspInfo) || !p.IsErrorRspInfo(pRspInfo) {
		if !p.isEmpty(pInvestorPosition) {
			log.Println("ComTradeID:", pInvestorPosition.GetComTradeID())
			log.Println("TradeID:", pInvestorPosition.GetTradeID())
			log.Println("InstrumentID:", pInvestorPosition.GetInstrumentID())
		} else {
			log.Println("kong")
		}
	}
}

///p.ReqQryInstrument("")空字符串表示查询全部合约
func (p *GoCThostFtdcTraderSpi) ReqQryInstrument(InstrumentID string) {
	req := goctp.NewCThostFtdcQryInstrumentField()
	req.SetInstrumentID(InstrumentID)

	for {
		iResult := p.Client.TraderApi.ReqQryInstrument(req, p.Client.GetTraderRequestID())

		if iResult == 0 {
			log.Printf("--->>> 请求查询合约: 成功 %#v\n", iResult)
			break
		} else {
			log.Printf("--->>> 请求查询合约: 受到流控 %#v\n", iResult)
			time.Sleep(1 * time.Second)
		}
	}
}

func (p *GoCThostFtdcTraderSpi) OnRspQryInstrument(pInstrument goctp.CThostFtdcInstrumentField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	log.Println("GoCThostFtdcTraderSpi.OnRspQryInstrument.")
	if p.isEmpty(pRspInfo) || !p.IsErrorRspInfo(pRspInfo) {
		if !p.isEmpty(pInstrument) {
			log.Println("GoCThostFtdcTraderSpi.OnRspQryInstrument: ", pInstrument.GetInstrumentID(), "#1", pInstrument.GetExchangeID(), "#2",
				pInstrument.GetInstrumentName(), "#3", pInstrument.GetExchangeInstID(), "#4", pInstrument.GetProductID(), "#5", pInstrument.GetProductClass(),
				pInstrument.GetDeliveryYear(), pInstrument.GetDeliveryMonth(), pInstrument.GetMaxMarketOrderVolume(), pInstrument.GetMinMarketOrderVolume(),
				pInstrument.GetMaxLimitOrderVolume(), pInstrument.GetMinLimitOrderVolume(), pInstrument.GetVolumeMultiple(), pInstrument.GetPriceTick(),
				pInstrument.GetCreateDate(), pInstrument.GetOpenDate(), pInstrument.GetExpireDate(), pInstrument.GetStartDelivDate(), pInstrument.GetEndDelivDate())
		} else {
			log.Println("kong")
		}
	}
}

///查询资金账户
func (p *GoCThostFtdcTraderSpi) ReqQryTradingAccount() {
	req := goctp.NewCThostFtdcQryTradingAccountField()
	req.SetBrokerID(p.Client.BrokerID)
	req.SetInvestorID(p.Client.InvestorID)

	for {
		iResult := p.Client.TraderApi.ReqQryTradingAccount(req, p.Client.GetTraderRequestID())
		if iResult == 0 {
			log.Printf("--->>> 请求查询资金账户: 成功 %#v\n", iResult)
			break
		} else {
			log.Printf("--->>> 请求查询资金账户: 受到流控 %#v\n", iResult)
			time.Sleep(1 * time.Second)
		}
	}
}

func (p *GoCThostFtdcTraderSpi) OnRspQryTradingAccount(pTradingAccount goctp.CThostFtdcTradingAccountField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {

	log.Println("GoCThostFtdcTraderSpi.OnRspQryTradingAccount.")

	if p.isEmpty(pRspInfo) || !p.IsErrorRspInfo(pRspInfo) {
		if !p.isEmpty(pTradingAccount) {
			log.Println("Available:", pTradingAccount.GetAvailable())
		} else {
			log.Println("kong")
		}
	}
}

//插入报单
func (p *GoCThostFtdcTraderSpi) ReqOrderInsert() {
	req := goctp.NewCThostFtdcInputOrderField()
	req.SetBrokerID(p.Client.BrokerID)
	req.SetInvestorID(p.Client.InvestorID)
	req.SetInstrumentID("IF1706")
	req.SetDirection(goctp.THOST_FTDC_D_Buy)
	req.SetCombOffsetFlag(string(goctp.THOST_FTDC_OF_Open))
	req.SetCombHedgeFlag(string(goctp.THOST_FTDC_HF_Speculation))
	req.SetVolumeTotalOriginal(1)
	req.SetContingentCondition(goctp.THOST_FTDC_CC_Immediately)
	req.SetVolumeCondition(goctp.THOST_FTDC_VC_AV)
	req.SetMinVolume(0)
	req.SetForceCloseReason(goctp.THOST_FTDC_FCC_NotForceClose)
	req.SetIsAutoSuspend(0)
	req.SetUserForceClose(0)
	req.SetOrderPriceType(goctp.THOST_FTDC_OPT_LimitPrice)
	req.SetLimitPrice(3300.00)
	req.SetTimeCondition(goctp.THOST_FTDC_TC_GFD)

	iResult := p.Client.TraderApi.ReqOrderInsert(req, p.Client.GetTraderRequestID())

	if iResult == 0 {
		log.Println("报单插入: 成功, iResult=", iResult)
	} else {
		log.Println("报单插入: 失败, iResult=", iResult)
	}
}

func (p *GoCThostFtdcTraderSpi) OnRspOrderInsert(pInputOrder goctp.CThostFtdcInputOrderField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	log.Println("GoCThostFtdcTraderSpi.OnRspOrderInsert.")

	if bIsLast && (p.isEmpty(pRspInfo) || !p.IsErrorRspInfo(pRspInfo)) {
		log.Println(1)
	}
}

func (p *GoCThostFtdcTraderSpi) OnErrRtnOrderInsert(pInputOrder goctp.CThostFtdcInputOrderField, pRspInfo goctp.CThostFtdcRspInfoField) {
	log.Println("GoCThostFtdcTraderSpi.OnErrRtnOrderInsert.")

	if !p.isEmpty(pRspInfo) && !p.IsErrorRspInfo(pRspInfo) {
		log.Println(2)
	}
}

func (p *GoCThostFtdcTraderSpi) OnRtnOrder(pOrder goctp.CThostFtdcOrderField) {
	log.Println("GoCThostFtdcTraderSpi.OnRtnOrder.")
	log.Println("CancelTime:", pOrder.GetCancelTime())
	log.Println("交易所编号:", pOrder.GetExchangeID())
	log.Println("合约代码:", pOrder.GetInstrumentID())
	log.Println("FrontID:", pOrder.GetFrontID())
	log.Println("SessionID:", pOrder.GetSessionID())
	log.Println("报单引用:", pOrder.GetOrderRef())
	log.Println("买卖方向:", pOrder.GetDirection())
	log.Println("组合开平标志:", pOrder.GetCombOffsetFlag())
	log.Println("价格:", pOrder.GetLimitPrice())
	log.Println("数量:", pOrder.GetVolumeTotalOriginal())
	log.Println("今成交数量:", pOrder.GetVolumeTraded())
	log.Println("剩余数量:", pOrder.GetVolumeTotal())
	log.Println("报单编号（判断报单是否有效）:", pOrder.GetOrderSysID())
	log.Println("报单提交状态:", pOrder.GetOrderSubmitStatus())
	log.Println("报单状态:", pOrder.GetOrderStatus())
	log.Println("报单日期:", pOrder.GetInsertDate())
	log.Println("序号:", pOrder.GetSequenceNo())
}

func (p *GoCThostFtdcTraderSpi) OnRtnTrade(pTrade goctp.CThostFtdcTradeField) {
	log.Println("GoCThostFtdcTraderSpi.OnRtnTrade.")
}

//撤单
//p.ReqOrderAction("CFFEX","       63288") 注意参数形式，直接p.ReqOrderAction("CFFEX","63288")会找不到订单
//强烈建议直接通过GetExchangeID(),GetOrderSysID()来获取参数，以防止由于字符串不匹配导致的找不到订单问题
func (p *GoCThostFtdcTraderSpi) ReqOrderAction(ExchangeID string, OrderSysID string) {
	log.Println("GoCThostFtdcTraderSpi.ReqOrderAction.")
	req := goctp.NewCThostFtdcInputOrderActionField()
	req.SetBrokerID(p.Client.BrokerID)
	req.SetInvestorID(p.Client.InvestorID)
	req.SetExchangeID(ExchangeID)
	req.SetOrderSysID(OrderSysID)
	req.SetActionFlag(goctp.THOST_FTDC_AF_Delete)

	iResult := p.Client.TraderApi.ReqOrderAction(req, p.Client.GetTraderRequestID())

	if iResult != 0 {
		log.Println("ReqOrderAction: 失败.")
	} else {
		log.Println("ReqOrderAction: 成功.")
	}
}

func (p *GoCThostFtdcTraderSpi) OnRspOrderAction(pInputOrderAction goctp.CThostFtdcInputOrderActionField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	log.Println("GoCThostFtdcTraderSpi.OnRspOrderAction.")

	if p.isEmpty(pRspInfo) || !p.IsErrorRspInfo(pRspInfo) {
		log.Println("1234")
	}
}

func (p *GoCThostFtdcTraderSpi) OnErrRtnOrderAction(pInputOrderAction goctp.CThostFtdcInputOrderActionField, pRspInfo goctp.CThostFtdcRspInfoField) {
	log.Println("GoCThostFtdcTraderSpi.OnErrRtnOrderInsert.")

	if p.isEmpty(pRspInfo) || !p.IsErrorRspInfo(pRspInfo) {
		log.Println("2")
	}
}

///预埋单录入请求
func (p *GoCThostFtdcTraderSpi) ReqParkedOrderInsert() {
	req := goctp.NewCThostFtdcParkedOrderField()
	req.SetBrokerID(p.Client.BrokerID)
	req.SetInvestorID(p.Client.InvestorID)
	req.SetInstrumentID("IF1703")
	req.SetDirection(goctp.THOST_FTDC_D_Buy)
	req.SetCombOffsetFlag(string(goctp.THOST_FTDC_OF_Open))
	req.SetCombHedgeFlag(string(goctp.THOST_FTDC_HF_Speculation))
	req.SetVolumeTotalOriginal(1)
	req.SetContingentCondition(goctp.THOST_FTDC_CC_Immediately)
	req.SetVolumeCondition(goctp.THOST_FTDC_VC_AV)
	req.SetMinVolume(1)
	req.SetForceCloseReason(goctp.THOST_FTDC_FCC_NotForceClose)
	req.SetIsAutoSuspend(0)
	req.SetUserForceClose(0)
	req.SetOrderPriceType(goctp.THOST_FTDC_OPT_LimitPrice)
	req.SetLimitPrice(3412.00)
	req.SetTimeCondition(goctp.THOST_FTDC_TC_GFD)

	iResult := p.Client.TraderApi.ReqParkedOrderInsert(req, p.Client.GetTraderRequestID())

	if iResult != 0 {
		log.Println("reqParkedOrderInsert: 失败.")
	} else {
		log.Println("reqParkedOrderInsert: 成功.")
	}
}

func (p *GoCThostFtdcTraderSpi) OnRspParkedOrderInsert(pParkedOrder goctp.CThostFtdcParkedOrderField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	log.Println("GoCThostFtdcTraderSpi.OnRspParkedOrderInsert.")

	if p.isEmpty(pRspInfo) || !p.IsErrorRspInfo(pRspInfo) {
		log.Println("GoCThostFtdcTraderSpi.OnRtnOrder.")
		log.Println("交易所编号:", pParkedOrder.GetExchangeID())
		log.Println("合约代码:", pParkedOrder.GetInstrumentID())
		log.Println("报单引用:", pParkedOrder.GetOrderRef())
		log.Println("买卖方向:", pParkedOrder.GetDirection())
		log.Println("组合开平标志:", pParkedOrder.GetCombOffsetFlag())
		log.Println("价格:", pParkedOrder.GetLimitPrice())
		log.Println("数量:", pParkedOrder.GetVolumeTotalOriginal())
		log.Println("ParkedOrderID:", pParkedOrder.GetParkedOrderID())
		log.Println("Status:", pParkedOrder.GetStatus())
	}
}

///预埋撤单
//p.ReqParkedOrderAction("CFFEX","       63288","IF1709")注意参数形式，直接p.ReqParkedOrderAction("CFFEX","63288")会找不到订单
//强烈建议直接通过GetExchangeID(),GetOrderSysID(),GetInstrumentID()来获取参数，以防止由于字符串不匹配导致的找不到订单问题
func (p *GoCThostFtdcTraderSpi) ReqParkedOrderAction(ExchangeID string, OrderSysID string, InstrumentID string) {
	req := goctp.NewCThostFtdcParkedOrderActionField()
	req.SetBrokerID(p.Client.BrokerID)
	req.SetInvestorID(p.Client.InvestorID)
	req.SetExchangeID(ExchangeID)
	req.SetOrderSysID(OrderSysID)
	req.SetActionFlag(goctp.THOST_FTDC_AF_Delete)
	req.SetInstrumentID(InstrumentID)

	iResult := p.Client.TraderApi.ReqParkedOrderAction(req, p.Client.GetTraderRequestID())

	if iResult != 0 {
		log.Println("ReqParkedOrderAction: 失败.")
	} else {
		log.Println("ReqParkedOrderAction: 成功.")
	}
}

func (p *GoCThostFtdcTraderSpi) OnRspParkedOrderAction(pParkedOrderAction goctp.CThostFtdcParkedOrderActionField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	log.Println("GoCThostFtdcTraderSpi.OnRspParkedOrderAction.")

	if p.isEmpty(pRspInfo) || !p.IsErrorRspInfo(pRspInfo) {
		log.Println("1")
	} else {
		log.Println("2")
	}

}

///请求查询预埋单
func (p *GoCThostFtdcTraderSpi) ReqQryParkedOrder() {
	req := goctp.NewCThostFtdcQryParkedOrderField()
	req.SetBrokerID(p.Client.BrokerID)
	req.SetInvestorID(p.Client.InvestorID)

	for {
		iResult := p.Client.TraderApi.ReqQryParkedOrder(req, p.Client.GetTraderRequestID())
		if iResult == 0 {
			log.Printf("--->>> ReqQryParkedOrder: 成功 %#v\n", iResult)
			break
		} else {
			log.Printf("--->>> ReqQryParkedOrder: 受到流控 %#v\n", iResult)
			time.Sleep(1 * time.Second)
		}
	}
}

func (p *GoCThostFtdcTraderSpi) OnRspQryParkedOrder(pInvestorPosition goctp.CThostFtdcParkedOrderField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	log.Println("GoCThostFtdcTraderSpi.OnRspQryParkedOrder.")

	if p.isEmpty(pRspInfo) || !p.IsErrorRspInfo(pRspInfo) {
		if !p.isEmpty(pInvestorPosition) {
			log.Printf("InstrumentID:%#v", pInvestorPosition.GetInstrumentID())
			log.Printf("ParkedOrderID:%#v", pInvestorPosition.GetParkedOrderID())
			log.Printf("VolumeTotalOriginal:%#v", pInvestorPosition.GetVolumeTotalOriginal())
			log.Println("Status:", pInvestorPosition.GetStatus())
		} else {
			log.Println("kong")
		}
	}
}

///请求查询预埋撤单
func (p *GoCThostFtdcTraderSpi) ReqQryParkedOrderAction() {
	req := goctp.NewCThostFtdcQryParkedOrderActionField()
	req.SetBrokerID(p.Client.BrokerID)
	req.SetInvestorID(p.Client.InvestorID)

	for {
		iResult := p.Client.TraderApi.ReqQryParkedOrderAction(req, p.Client.GetTraderRequestID())

		if iResult == 0 {
			log.Printf("--->>> ReqQryParkedOrderAction: 成功 %#v\n", iResult)
			break
		} else {
			log.Printf("--->>> ReqQryParkedOrderAction: 受到流控 %#v\n", iResult)
			time.Sleep(1 * time.Second)
		}
	}
}

func (p *GoCThostFtdcTraderSpi) OnRspQryParkedOrderAction(pInvestorPosition goctp.CThostFtdcParkedOrderActionField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	log.Println("GoCThostFtdcTraderSpi.OnRspQryParkedOrderAction.")

	if p.isEmpty(pRspInfo) || !p.IsErrorRspInfo(pRspInfo) {
		if !p.isEmpty(pInvestorPosition) {
			log.Println(pInvestorPosition.GetInstrumentID())
			log.Println(pInvestorPosition.GetOrderRef())
			log.Println(pInvestorPosition.GetStatus())

		} else {
			log.Println("kong")
		}
	}
}

///请求删除预埋单
//强烈建议直接通过GetOParkedOrderID()来获取参数，以防止由于字符串不匹配导致的找不到订单问题
func (p *GoCThostFtdcTraderSpi) ReqRemoveParkedOrder(ParkedOrderID string) {
	req := goctp.NewCThostFtdcRemoveParkedOrderField()
	req.SetBrokerID(p.Client.BrokerID)
	req.SetInvestorID(p.Client.InvestorID)
	req.SetParkedOrderID(ParkedOrderID)

	for {
		iResult := p.Client.TraderApi.ReqRemoveParkedOrder(req, p.Client.GetTraderRequestID())
		if iResult == 0 {
			log.Printf("--->>> ReqRemoveParkedOrder: 成功 %#v\n", iResult)
			break
		} else {
			log.Printf("--->>> ReqRemoveParkedOrder: 受到流控 %#v\n", iResult)
			time.Sleep(1 * time.Second)
		}
	}
}

func (p *GoCThostFtdcTraderSpi) OnRspRemoveParkedOrder(pRemoveParkedOrder goctp.CThostFtdcRemoveParkedOrderField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	log.Println("GoCThostFtdcTraderSpi.OnRspRemoveParkedOrder.")

	if !p.isEmpty(pRspInfo) || !p.IsErrorRspInfo(pRspInfo) {
		if !p.isEmpty(pRemoveParkedOrder) {
			log.Printf(pRemoveParkedOrder.GetParkedOrderID())
		}

	}
}

//请求删除预埋撤单
func (p *GoCThostFtdcTraderSpi) ReqRemoveParkedOrderAction(ParkedOrderActionID string) {
	req := goctp.NewCThostFtdcRemoveParkedOrderActionField()
	req.SetBrokerID(p.Client.BrokerID)
	req.SetInvestorID(p.Client.InvestorID)
	req.SetParkedOrderActionID(ParkedOrderActionID)

	for {
		iResult := p.Client.TraderApi.ReqRemoveParkedOrderAction(req, p.Client.GetTraderRequestID())

		if iResult == 0 {
			log.Printf("--->>> ReqRemoveParkedOrderAction: 成功 %#v\n", iResult)
			break
		} else {
			log.Printf("--->>> ReqRemoveParkedOrderAction: 受到流控 %#v\n", iResult)
			time.Sleep(1 * time.Second)
		}
	}
}

func (p *GoCThostFtdcTraderSpi) OnRspRemoveParkedOrderAction(pRemoveParkedOrderAction goctp.CThostFtdcRemoveParkedOrderActionField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	log.Println("GoCThostFtdcTraderSpi.OnRspRemoveParkedOrderAction.")

	if bIsLast && (p.isEmpty(pRspInfo) || !p.IsErrorRspInfo(pRspInfo)) {
		if !p.isEmpty(pRemoveParkedOrderAction) {
			log.Printf("ok2")
		}

	}
}

func (p *GoCThostFtdcTraderSpi) ReqQryOrder() {
	req := goctp.NewCThostFtdcQryOrderField()
	req.SetBrokerID(p.Client.BrokerID)
	req.SetInvestorID(p.Client.InvestorID)

	for {
		iResult := p.Client.TraderApi.ReqQryOrder(req, p.Client.GetTraderRequestID())

		if iResult == 0 {
			log.Printf("--->>> ReqQryOrder: 成功 %#v\n", iResult)
			break
		} else {
			log.Printf("--->>> ReqQryOrder: 失败 %#v\n", iResult)
			time.Sleep(1 * time.Second)
		}
	}
}

///请求查询报单
func (p *GoCThostFtdcTraderSpi) OnRspQryOrder(pOrder goctp.CThostFtdcOrderField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {

	log.Println("GoCThostFtdcTraderSpi.OnRspQryOrder.")

	if p.isEmpty(pRspInfo) || !p.IsErrorRspInfo(pRspInfo) {
		if !p.isEmpty(pOrder) {
			log.Println("InstrumentID:", pOrder.GetInstrumentID())
			log.Println("OrderStatus:", pOrder.GetOrderStatus())
			log.Println("TraderID:", pOrder.GetTraderID())
			log.Printf("ExchangeID:%#v", pOrder.GetExchangeID())
			log.Printf("OrderSysID:%#v", pOrder.GetOrderSysID())
			log.Printf("OrderRef:%#v", pOrder.GetOrderRef())
			log.Println("Direction:", pOrder.GetDirection())
			log.Println("FrontID:", pOrder.GetFrontID())
			log.Println("SessionID:", pOrder.GetSessionID())
			log.Println("OrderLocalID:", pOrder.GetOrderLocalID())
			log.Println("OrderLimitPrice:", pOrder.GetLimitPrice())
		} else {
			log.Println("kong")
		}
	}
}
