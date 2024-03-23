package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Respone struct {
	Msg   string `json:"msg"`
	Result	Result `json:"result"`
}

type Result struct {
	Payload interface{} `json:"payload"`
	Txid	string `json:"txid"`
	Error	string `json:"error"`
}

var respone Respone

func saveCollectionhander(c *gin.Context) {

	respone.Msg="藏品保存失败"
	respone.Result.Payload=""
	if err := c.ShouldBindJSON(&collection); err != nil {

		respone.Result.Error=err.Error()
		c.JSON(http.StatusBadRequest, gin.H{"respone":respone})
		return
	}
	payload,txid, err := serviceSetup.Save(collection)
	if err != nil {
		respone.Result.Error=err.Error()
		c.JSON(http.StatusInternalServerError, gin.H{"respone":respone})
		return
	}

	respone.Msg="藏品保存成功！"
	respone.Result.Payload=payload
	respone.Result.Txid=txid
	respone.Result.Error=""

	c.JSON(http.StatusOK, gin.H{"respone":respone})
}

func delCollectionhander(c *gin.Context) {

	respone.Msg="藏品删除失败"
	respone.Result.Payload=""
	if err := c.ShouldBindJSON(&collection); err != nil {
		respone.Result.Error=err.Error()
		c.JSON(http.StatusBadRequest, gin.H{"respone":respone})
		return
	}

	payload,txid, err := serviceSetup.Del(collection.Id)
	if err != nil {
		respone.Result.Error=err.Error()
		c.JSON(http.StatusInternalServerError, gin.H{"respone":respone})
		return
	}

	respone.Msg="藏品删除成功！"
	respone.Result.Payload=payload
	respone.Result.Txid=txid
	respone.Result.Error=""

	c.JSON(http.StatusOK, gin.H{"respone":respone})
}

func updateCollectionhander(c *gin.Context) {

	respone.Msg="藏品更新失败"
	respone.Result.Payload=""
	if err := c.ShouldBindJSON(&collection); err != nil {
		respone.Result.Error=err.Error()
		c.JSON(http.StatusBadRequest,  gin.H{"respone":respone})
		return
	}

	payload,txid, err := serviceSetup.Update(collection)
	if err != nil {
		respone.Result.Error=err.Error()
		c.JSON(http.StatusInternalServerError,  gin.H{"respone":respone})
		return
	}

	respone.Msg="藏品更新成功！"
	respone.Result.Payload=payload
	respone.Result.Txid=txid
	respone.Result.Error=""

	c.JSON(http.StatusOK, gin.H{"respone":respone})
}

func queryCollectionByIdhander(c *gin.Context) {

	respone.Msg="藏品查询失败"
	respone.Result.Payload=""
	if err := c.ShouldBindJSON(&collection); err != nil {
		respone.Result.Error=err.Error()
		c.JSON(http.StatusBadRequest, gin.H{"respone":respone})
		return
	}

	payload,txid, err := serviceSetup.QueryById(collection.Id)
	if err != nil {
		respone.Result.Error=err.Error()
		c.JSON(http.StatusInternalServerError, gin.H{"respone":respone})
		return
	}

	var msg interface{}
	err =json.Unmarshal([]byte(payload),&msg)
	if err!=nil {
		respone.Result.Error=err.Error()
		c.JSON(http.StatusInternalServerError, gin.H{"respone":respone})
		return
	}

	respone.Msg="藏品查询成功！"
	respone.Result.Payload=msg
	respone.Result.Txid=txid
	respone.Result.Error=""

	c.JSON(http.StatusOK, gin.H{"respone":respone})
}

func queryCollectionByHashhander(c *gin.Context) {

	respone.Msg="藏品查询失败"
	respone.Result.Payload=""

	if err := c.ShouldBindJSON(&collection); err != nil {
		respone.Result.Error=err.Error()
		c.JSON(http.StatusBadRequest, gin.H{"respone":respone})
		return
	}

	payload,txid, err := serviceSetup.QueryByHash(collection.CollectionHash)
	if err != nil {
		respone.Result.Error=err.Error()
		c.JSON(http.StatusInternalServerError, gin.H{"respone":respone})
		return
	}

	var msg interface{}
	err =json.Unmarshal([]byte(payload),&msg)
	if err!=nil {
		respone.Result.Error=err.Error()
		c.JSON(http.StatusInternalServerError, gin.H{"respone":respone})
		return
	}

	respone.Msg="藏品查询成功！"
	respone.Result.Payload=msg
	respone.Result.Txid=txid
	respone.Result.Error=""
	c.JSON(http.StatusOK, gin.H{"respone":respone})
}

func queryCollectionByOwnerhander(c *gin.Context) {

	respone.Msg="藏品查询失败"
	respone.Result.Payload=""

	if err := c.ShouldBindJSON(&collection); err != nil {
		respone.Result.Error=err.Error()
		c.JSON(http.StatusBadRequest, gin.H{"respone":respone})
		return
	}
	payload,txid, err := serviceSetup.QueryByOwner(collection.Owner)
	if err != nil {
		respone.Result.Error=err.Error()
		c.JSON(http.StatusInternalServerError, gin.H{"respone":respone})
		return
	}
	var msg interface{}
	err =json.Unmarshal([]byte(payload),&msg)
	if err!=nil {
		respone.Result.Error=err.Error()
		c.JSON(http.StatusInternalServerError, gin.H{"respone":respone})
		return
	}

	respone.Msg="藏品查询成功！"
	respone.Result.Payload=msg
	respone.Result.Txid=txid
	respone.Result.Error=""

	c.JSON(http.StatusOK, gin.H{"respone":respone})
}

func addCollectionTransactionhander(c *gin.Context) {

	respone.Msg="交易上传失败"
	respone.Result.Payload=""
	if err := c.ShouldBindJSON(&transactionHistoryItem); err != nil {
		respone.Result.Error=err.Error()
		c.JSON(http.StatusBadRequest,  gin.H{"respone":respone})
		return
	}

	payload,txid, err := serviceSetup.AddT(transactionHistoryItem)
	if err != nil {
		respone.Result.Error=err.Error()
		c.JSON(http.StatusInternalServerError,  gin.H{"respone":respone})
		return
	}

	respone.Msg="交易上传成功！"
	respone.Result.Payload=payload
	respone.Result.Txid=txid
	respone.Result.Error=""

	c.JSON(http.StatusOK, gin.H{"respone":respone})
}

func queryCollectionTransactionhander(c *gin.Context) {

	respone.Msg="交易查询失败"
	respone.Result.Payload=""
	if err := c.ShouldBindJSON(&transactionHistoryItem); err != nil {
		respone.Result.Error=err.Error()
		c.JSON(http.StatusBadRequest, gin.H{"respone":respone})
		return
	}
	payload,txid, err := serviceSetup.QueryT(transactionHistoryItem.Id)
	if err != nil {
		respone.Result.Error=err.Error()
		c.JSON(http.StatusInternalServerError, gin.H{"respone":respone})
		return
	}
	var msg interface{}
	err =json.Unmarshal([]byte(payload),&msg)
	if err!=nil {
		respone.Result.Error=err.Error()
		c.JSON(http.StatusInternalServerError, gin.H{"respone":respone})
		return
	}
	respone.Msg="交易查询成功！"
	respone.Result.Payload=msg
	respone.Result.Txid=txid
	respone.Result.Error=""
	fmt.Println("5")
	c.JSON(http.StatusOK, gin.H{"respone":respone})
}