package service

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

//(以下是通过key值进行操作)
//调用链码增加藏品信息
func (t *ServiceSetup) Save(collection Collection) ( string, error) {

	eventID := "SaveCollection"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)

	b, err := json.Marshal(collection)
	if err != nil {
		return  "", fmt.Errorf("指定的collection对象序列化时发生错误")
	}

	fmt.Println("开始执行链码(调用链码增加藏品信息)")

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "saveCollection", Args: [][]byte{b, []byte(eventID)}}
	respone, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}

	err = eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}

	txid :=respone.TransactionID
	return	string(txid),nil
}

//调用链码删除藏品信息
func (t *ServiceSetup) Del(Id string) (string, error) {

	eventID := "DelCollection"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)
	fmt.Println("开始执行链码(调用链码删除藏品信息)")

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "delCollection", Args: [][]byte{[]byte(Id), []byte(eventID)}}
	respone, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}

	err = eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}

	txid :=respone.TransactionID
	return	string(txid),nil
}

//调用链码更新藏品信息
func (t *ServiceSetup) Update(collection Collection) (string, error) {

	eventID := "UpdateCollection"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)


	fmt.Println("开始执行链码(调用链码更新藏品信息)")
	b, err := json.Marshal(collection)
	if err != nil {
		return "", fmt.Errorf("指定的collection对象序列化时发生错误")
	}

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "updateCollection", Args: [][]byte{b, []byte(eventID)}}
	respone, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}

	err = eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}

	txid :=respone.TransactionID

	return	string(txid),nil
}

//调用链码通过藏品id查询数据
func (t *ServiceSetup) QueryById(Id string) (string, error){
	fmt.Println("开始执行链码(调用链码通过藏品id查询数据)")
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "queryCollectionById", Args: [][]byte{[]byte(Id)}}
	respone, err := t.Client.Query(req)
	if err != nil {
		return "[]byte{0x00}", err
	}

	Payload :=respone.Payload

	return	string(Payload),nil
}


//(以下是通过富查询)
//调用链码通过Hash查询数据
func (t *ServiceSetup) QueryByHash(CollectionHash string) (string, error){

	fmt.Println("开始执行链码(调用链码通过Hash查询数据)")
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "queryCollectionByHash", Args: [][]byte{[]byte(CollectionHash)}}
	respone, err := t.Client.Query(req)
	if err != nil {
		return "[]byte{0x00}", err
	}

	Payload :=respone.Payload

	return	string(Payload),nil
}

//调用链码通过拥有者查信息
func (t *ServiceSetup) QueryByOwner(Owner string) (string, error){

	fmt.Println("开始执行链码(调用链码通过拥有者查询数据)")
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "queryCollectionByOwner", Args: [][]byte{[]byte(Owner)}}
	respone, err := t.Client.Query(req)
	if err != nil {
		return "[]byte{0x00}", err
	}

	Payload :=respone.Payload

	return	string(Payload),nil
}