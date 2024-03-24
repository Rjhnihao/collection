package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
	"math/big"
	"time"
)

type Collection struct {
	Type                string                   `json:"Type"`
	Name                string                   `json:"Name" `            					//藏品名称
	Owner               string                   `json:"Owner" `           					//藏品拥有者id
	Introduce           string                   `json:"Introduce" `       					//藏品介绍
	Id                  string                   `json:"Id" `             					//藏品id
	CurrentPrice        string                   `json:"CurrentPrice" `    					//藏品当前价格
	RemainingNumber     string                   `json:"RemainingNumber" ` 					//藏品剩余数量
	CollectionHash      string                   `json:"CollectionHash"`  					//藏品哈希值
	TimeStamp           int64                    											//unix时间戳
	TransactionHistorys []TransactionHistoryItem 											//交易历史记录
}

type TransactionHistoryItem struct {
	TransactionPrice  string `json:"TransactionPrice"`  									//交易价格
	TransactionNumber string `json:"TransactionNumber"` 									//交易数量
	Seller            string `json:"Seller" `            									//卖家
	Buyer             string `json:"Buyer"`             									//买家
	TransactionId     string `json:"TransactionId"`											//每笔交易的id
	TimeStamp           int64                    											//unix时间戳
	Collection        Collection
}

type NFTChaincode struct {
}

const Type = "数字藏品"

// 初始化链码
func (t *NFTChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println("初始化")

	return shim.Success(nil)
}

// 函数的调用
func (t *NFTChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// 获取用户意图
	fun, args := stub.GetFunctionAndParameters()

	if fun == "saveCollection" {
		return t.saveCollection(stub, args) //  根据id添加藏品信息
	} else if fun == "delCollection" {
		return t.delCollection(stub, args) // 根据id删除信息
	} else if fun == "updateCollection" {
		return t.updateCollection(stub, args) // 根据id更新信息
	} else if fun == "queryCollectionById" {
		return t.queryCollectionById(stub, args) // 根据id查询
	} else if fun == "queryCollectionByHash" {
		return t.queryCollectionByHash(stub, args) // 根据哈希值查询
	} else if fun == "queryCollectionByOwner" {
		return t.queryCollectionByOwner(stub, args) // 根据用户查询
	}else if fun == "addTransaction" {
		return t.addTransaction(stub, args) // 根据用户查询
	}else if fun == "queryTransaction" {
		return t.queryTransaction(stub, args) // 根据用户查询
	}
	return shim.Error("指定的函数名称错误")

}


/*
对交易的处理
*/
// 保存藏品
func PutCollection(stub shim.ChaincodeStubInterface, collection Collection) ([]byte, bool) {

	collection.Type = Type

	b, err := json.Marshal(collection)
	if err != nil {
		return nil, false
	}

	err = stub.PutState(collection.Id, b)
	if err != nil {
		return nil, false
	}

	return b, true
}

// 根据id查询藏品当前信息
func getCollection(stub shim.ChaincodeStubInterface, id string) (Collection, bool) {
	var Collection Collection
	// 根据身份证号码查询信息状态
	a, err := stub.GetState(id)
	if err != nil {
		return Collection, false
	}

	if a == nil {
		return Collection, false
	}

	// 对查询到的状态进行反序列化
	err = json.Unmarshal(a, &Collection)
	if err != nil {
		return Collection, false
	}

	// 返回结果
	return Collection, true
}

// 根据指定的查询字符串实现富查询
func getCollectionByQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer 是一个包含 QueryRecords 的 JSON 数组
	var buffer bytes.Buffer

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}

	fmt.Printf("查询结果:\n%s\n", buffer.String())

	return buffer.Bytes(), nil

}

// (以下是通过key值进行操作)
// 根据id添加藏品信息
func (t *NFTChaincode) saveCollection(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 2 {
		return shim.Error("给定的参数个数不符合要求")
	}


	var Collection Collection
	err := json.Unmarshal([]byte(args[0]), &Collection)
	if err != nil {
		return shim.Error("反序列化失败")
	}

	_, exist := getCollection(stub, Collection.Id)
	if exist {
		return shim.Error("要添加的藏品编号已存在")
	}


	Collection.TimeStamp = time.Now().Unix()

	n, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	strN := fmt.Sprintf("%d", n)
	hash := sha256.Sum256([]byte(strN))

	dataToHash := fmt.Sprintf("%s%s%s%s%s%s%s%s",
		Collection.Name,
		Collection.Owner,
		Collection.Introduce,
		Collection.Id,
		Collection.CurrentPrice,
		Collection.RemainingNumber,
		Collection.TimeStamp,
		hash,
	)

	hash = sha256.Sum256([]byte(dataToHash))
	Collection.CollectionHash = hex.EncodeToString(hash[:])
	Collection.TransactionHistorys=nil

	_, bl := PutCollection(stub, Collection)
	if !bl {
		return shim.Error("保存信息时发生错误")
	}

	//注册事件
	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("信息添加成功"))


}

// 根据id删除信息
func (t *NFTChaincode) delCollection(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("给定的参数个数不符合要求")
	}

	var Collection Collection
	err := json.Unmarshal([]byte(args[0]), &Collection)
	if err != nil {
		return shim.Error("反序列化失败")
	}

	result, exist := getCollection(stub, Collection.Id)
	if  !exist {
		return shim.Error("藏品编号不存在")
	}

	if result.Owner!=Collection.Owner ||result.Name!=Collection.Name ||result.Introduce!=Collection.Introduce ||result.CurrentPrice!=Collection.CurrentPrice || result.RemainingNumber!=Collection.RemainingNumber{
		shim.Error("要删除的藏品信息错误")
	}


	err = stub.DelState(Collection.Id)
	if err != nil {
		return shim.Error("删除信息时发生错误")
	}

	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("信息删除成功"))
}

// 根据id更新信息
func (t *NFTChaincode) updateCollection(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("给定的参数个数不符合要求")
	}

	var Collection Collection
	err := json.Unmarshal([]byte(args[0]), &Collection)
	if err != nil {
		return shim.Error("反序列化信息失败")
	}

	result, bl := getCollection(stub, Collection.Id)
	if !bl {
		return shim.Error("根据藏品id查询信息时发生错误")
	}



	result.Name = Collection.Name
	result.Owner = Collection.Owner
	result.Introduce = Collection.Introduce
	result.Id = Collection.Id
	result.CurrentPrice = Collection.CurrentPrice
	result.RemainingNumber = Collection.RemainingNumber
	result.TimeStamp = time.Now().Unix()

	if	result.Owner != Collection.Owner{
		dataToHash := fmt.Sprintf("%s%s%s%s%s%s%s%s",
			result.Name,
			result.Owner,
			result.Introduce,
			result.Id,
			result.CurrentPrice,
			result.RemainingNumber,
			result.TimeStamp,
			Collection.CollectionHash,
			)
		hash := sha256.Sum256([]byte(dataToHash))
		result.CollectionHash = hex.EncodeToString(hash[:])
	}
	result.TransactionHistorys=nil
	_, bl = PutCollection(stub, result)
	if !bl {
		return shim.Error("保存信息信息时发生错误")
	}

	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("信息更新成功"))
}

// 根据id查询藏品历史信息
func (t *NFTChaincode) queryCollectionById(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("给定的参数个数不符合要求")
	}

	a, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("根据藏品id没有查询到相关的信息")
	}

	if a == nil {
		return shim.Error("根据藏品id没有查询到相关的信息")
	}

	// 对查询到的状态进行反序列化
	var b Collection
	err = json.Unmarshal(a, &b)
	if err != nil {
		return shim.Error("反序列化edu信息失败")
	}

	// 获取历史变更数据
	iterator, err := stub.GetHistoryForKey(b.Id)
	if err != nil {
		return shim.Error("根据指定的藏品id查询对应的历史变更数据失败")
	}
	defer iterator.Close()

	// 迭代处理
	var historys []TransactionHistoryItem
	var hisEdu Collection
	for iterator.HasNext() {
		hisData, err := iterator.Next()
		if err != nil {
			return shim.Error("获取Collection的历史变更数据失败")
		}

		var TransactionHistoryItem TransactionHistoryItem
		json.Unmarshal(hisData.Value, &hisEdu)

		if hisData.Value == nil {
			var empty Collection
			TransactionHistoryItem.Collection = empty
		} else {
			TransactionHistoryItem.Collection = hisEdu
		}

		historys = append(historys, TransactionHistoryItem)

	}

	b.TransactionHistorys = historys

	result, err := json.Marshal(b)
	if err != nil {
		return shim.Error("序列化信息时发生错误")
	}
	return shim.Success(result)
}

// (以下是通过富查询)
// 根据哈希值
func (t *NFTChaincode) queryCollectionByHash(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 1 {
		return shim.Error("给定的参数个数不符合要求")
	}
	CollectionHash := args[0]

	// 拼装CouchDB所需要的查询字符串(是标准的一个JSON串),具体的可以看couchD的文档
	queryString := fmt.Sprintf("{\"selector\":{\"Type\":\"%s\", \"CollectionHash\":\"%s\"}}", Type, CollectionHash)

	// 查询数据
	result, err := getCollectionByQueryString(stub, queryString)
	if err != nil {
		return shim.Error("根据藏品Hash值查询信息时发生错误")
	}
	if result == nil {
		return shim.Error("根据藏品藏品Hash值查询信息失败")
	}
	return shim.Success(result)
}

// 根据用户
func (t *NFTChaincode) queryCollectionByOwner(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 1 {
		return shim.Error("给定的参数个数不符合要求")
	}
	Owner := args[0]
	// 拼装CouchDB所需要的查询字符串(是标准的一个JSON串),具体的可以看couchD的文档
	queryString := fmt.Sprintf("{\"selector\":{\"Type\":\"%s\", \"Owner\":\"%s\"}}", Type, Owner)

	// 查询数据
	result, err := getCollectionByQueryString(stub, queryString)
	if err != nil {
		return shim.Error("根据用户id查询信息时发生错误")
	}
	if result == nil {
		return shim.Error("根据用户id没有查询到相关的信息")
	}
	return shim.Success(result)
}


/*
对交易的处理
 */
// 保存交易
func PutTransaction(stub shim.ChaincodeStubInterface, transactionHistoryItem TransactionHistoryItem) ([]byte, bool) {

	b, err := json.Marshal(transactionHistoryItem)
	if err != nil {
		return nil, false
	}

	err = stub.PutState(transactionHistoryItem.TransactionId, b)
	if err != nil {
		return nil, false
	}

	return b, true
}

// 根据id查询藏品当前信息
func getTransaction(stub shim.ChaincodeStubInterface, id string) (TransactionHistoryItem, bool) {

	var transactionHistoryItem TransactionHistoryItem
	// 根据身份证号码查询信息状态
	a, err := stub.GetState(id)
	if err != nil {
		return transactionHistoryItem, false
	}

	if a == nil {
		return transactionHistoryItem, false
	}

	// 对查询到的状态进行反序列化
	err = json.Unmarshal(a, &transactionHistoryItem)
	if err != nil {
		return transactionHistoryItem, false
	}

	// 返回结果
	return transactionHistoryItem, true
}

//上传交易记录
func (t *NFTChaincode) addTransaction(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("给定的参数个数不符合要求")
	}

	var TransactionHistoryItem TransactionHistoryItem
	err := json.Unmarshal([]byte(args[0]), &TransactionHistoryItem)
	if err != nil {
		return shim.Error("反序列化信息失败")
	}
	//交易时间戳
	TransactionHistoryItem.TimeStamp=time.Now().Unix()
	//生成交易id
	str:=fmt.Sprintf("%s%s%s%s%s%s",
		TransactionHistoryItem.TransactionPrice,
		TransactionHistoryItem.TransactionNumber,
		TransactionHistoryItem.Collection.Id,
		TransactionHistoryItem.TimeStamp,
		TransactionHistoryItem.Seller,
		TransactionHistoryItem.Buyer,
		)
	hash1 := sha256.Sum256([]byte(str))

	TransactionHistoryItem.TransactionId = hex.EncodeToString(hash1[:])
	//在每个藏品后面添加交易

	result, bl := getCollection(stub, TransactionHistoryItem.Collection.Id)
	if !bl {
		return shim.Error("根据藏品id查询信息时发生错误")
	}

	result.TransactionHistorys=append(result.TransactionHistorys,TransactionHistoryItem)

	_, bl = PutCollection(stub, result)
	if !bl {
		return shim.Error("保存信息信息时发生错误")
	}

	//检查交易是否唯一
	_, exist := getTransaction(stub, TransactionHistoryItem.TransactionId)
	if exist {
		return shim.Error("要添加的交易已存在")
	}

	_, bl = PutTransaction(stub, TransactionHistoryItem)
	if !bl {
		errstr:=fmt.Sprintf("%s",bl)
		return shim.Error(errstr+"保存交易时发生错误")
	}

	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	b := map[string]string{"TransactionId": TransactionHistoryItem.TransactionId}
	jsonData, err := json.Marshal(b)
	if err != nil {
		return shim.Error("序列化信息失败")
	}
	return shim.Success(jsonData)
}

//查询交易记录
func (t *NFTChaincode) queryTransaction(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("给定的参数个数不符合要求")
	}

	b, bl := getTransaction(stub, args[0])
	if !bl {
		return shim.Error("根据交易id查询信息时发生错误")
	}

	result, err := json.Marshal(b)
	if err != nil {
		return shim.Error("序列化信息时发生错误")
	}
	return shim.Success(result)
}

func main() {
	err := shim.Start(new(NFTChaincode))
	if err != nil {
		fmt.Printf("启动NFTChaincode时发生错误: %s", err)
	}
}
