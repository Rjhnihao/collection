package service

import (
	"education/sdkInit"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"time"
)

type Collection struct {
	Type                string                   `json:"Type"`
	Name                string                   `json:"Name" `            					//藏品名称
	Owner               string                   `json:"Owner"`           					//藏品拥有者id
	Introduce           string                   `json:"Introduce"`       					//藏品介绍
	Id                  string                   `json:"Id" binding:"required"`             //藏品id
	CurrentPrice        string                   `json:"CurrentPrice"`    					//藏品当前价格
	RemainingNumber     string                   `json:"RemainingNumber"` 					//藏品剩余数量
	CollectionHash      string                   `json:"CollectionHash"`  					//藏品哈希值
	TimeStamp           int64                    											//unix时间戳
	TransactionHistorys []TransactionHistoryItem 											//交易历史记录
}

type TransactionHistoryItem struct {
	TransactionPrice  string `json:"TransactionPrice"`  									//交易价格
	TransactionNumber string `json:"TransactionNumber"` 									//交易数量
	Seller            string `json:"Seller"`            									//卖家
	Buyer             string `json:"Buyer"`             									//买家
	Id            	  string `json:"Id"`													//每笔交易的id
	TimeStamp           int64                    											//unix时间戳
	Collection        Collection
}


type ServiceSetup struct {
	ChaincodeID string
	Client      *channel.Client
	evClient    *event.Client
}

func regitserEvent(client *channel.Client, chaincodeID, eventID string) (fab.Registration, <-chan *fab.CCEvent) {

	reg, notifier, err := client.RegisterChaincodeEvent(chaincodeID, eventID)
	if err != nil {
		fmt.Println("注册链码事件失败: %s", err)
	}
	return reg, notifier
}

func eventResult(notifier <-chan *fab.CCEvent, eventID string) error {
	select {
	case ccEvent := <-notifier:
		fmt.Printf("接收到链码事件: %v\n", ccEvent)
	case <-time.After(time.Second * 20):
		return fmt.Errorf("不能根据指定的事件ID接收到相应的链码事件(%s)", eventID)
	}
	return nil
}

func InitService(chaincodeID, channelID string, org *sdkInit.OrgInfo, sdk *fabsdk.FabricSDK) (*ServiceSetup, error) {
	handler := &ServiceSetup{
		ChaincodeID: chaincodeID,
	}
	//prepare channel client context using client context
	clientChannelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(org.OrgUser), fabsdk.WithOrg(org.OrgName))
	// Channel client is used to query and execute transactions (Org1 is default org)
	client, err := channel.New(clientChannelContext)
	if err != nil {
		return nil, fmt.Errorf("Failed to create new channel client: %s", err)
	}
	handler.Client = client
	return handler, nil
}
