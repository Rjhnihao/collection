package main

import (
	"collection/sdkInit"
	"collection/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
)

const (
	cc_name = "simplecc"
	cc_version = "1.0.0"
)

var serviceSetup *service.ServiceSetup

func main() {
	// init orgs information
	orgs := []*sdkInit.OrgInfo{
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org1",
			OrgMspId:      "Org1MSP",
			OrgUser:       "User1",
			OrgPeerNum:    1,
			OrgAnchorFile: os.Getenv("GOPATH") + "/src/collection/fixtures/channel-artifacts/Org1MSPanchors.tx",
		},

	}

	// init sdk env info
	info := sdkInit.SdkEnvInfo{
		ChannelID:        "mychannel",
		ChannelConfig:    os.Getenv("GOPATH") + "/src/collection/fixtures/channel-artifacts/channel.tx",
		Orgs:             orgs,
		OrdererAdminUser: "Admin",
		OrdererOrgName:   "OrdererOrg",
		OrdererEndpoint:  "orderer.example.com",
		ChaincodeID:      cc_name,
		ChaincodePath:    os.Getenv("GOPATH")+"/src/collection/chaincode/",
		ChaincodeVersion: cc_version,
	}

	// sdk setup
	sdk, err := sdkInit.Setup("config.yaml", &info)
	if err != nil {
		fmt.Println(">> SDK setup error:", err)
		os.Exit(-1)
	}

	// create channel and join
	if err := sdkInit.CreateAndJoinChannel(&info); err != nil {
		fmt.Println(">> Create channel and join error:", err)
		os.Exit(-1)
	}

	// create chaincode lifecycle
	if err := sdkInit.CreateCCLifecycle(&info, 1, false, sdk); err != nil {
		fmt.Println(">> create chaincode lifecycle error: %v", err)
		os.Exit(-1)
	}

	// invoke chaincode set status
	fmt.Println(">> 通过链码外部服务设置链码状态......")

	serviceSetup, err = service.InitService(info.ChaincodeID, info.ChannelID, info.Orgs[0], sdk)
	if err!=nil{
		fmt.Println("初始化错误!")
		os.Exit(-1)
	}

	router := gin.Default()

	// 设置路由
	router.POST("/saveCollection", saveCollectionhander)
	router.DELETE("/delCollection", delCollectionhander)
	router.POST("/updateCollection", updateCollectionhander)
	router.GET("/queryCollectionInfoById", queryCollectionByIdhander)
	router.GET("/queryCollectionByHash", queryCollectionByHashhander)
	router.GET("/queryCollectionInfoByOwner", queryCollectionByOwnerhander)
	router.POST("/addCollectionTransaction", addCollectionTransactionhander)
	router.GET("/queryCollectionTransaction",queryCollectionTransactionhander)
	router.Run(":8080")
}
