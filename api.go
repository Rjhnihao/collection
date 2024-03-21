package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)


func saveCollectionhander(c *gin.Context) {
	if err := c.ShouldBindJSON(&collection); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	txid, err := serviceSetup.Save(collection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg":"藏品信息保存成功！","txid": txid})
}

func delCollectionhander(c *gin.Context) {

	if err := c.ShouldBindJSON(&collection); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	txid, err := serviceSetup.Del(collection.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg":"藏品信息删除成功！","txid": txid})
}

func updateCollectionhander(c *gin.Context) {
	if err := c.ShouldBindJSON(&collection); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	txid, err := serviceSetup.Update(collection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg":"藏品信息更新成功","txid": txid})
}

func queryCollectionByIdhander(c *gin.Context) {

	if err := c.ShouldBindJSON(&collection); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := serviceSetup.QueryById(collection.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var msg interface{}
	err =json.Unmarshal([]byte(result),&msg)
	if err!=nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}


	c.JSON(http.StatusOK, gin.H{"msg":"查询成功","result": msg})
}

func queryCollectionByHashhander(c *gin.Context) {
	if err := c.ShouldBindJSON(&collection); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := serviceSetup.QueryByHash(collection.CollectionHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var msg interface{}
	err =json.Unmarshal([]byte(result),&msg)
	if err!=nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg":"查询成功","result": msg})
}

func queryCollectionByOwnerhander(c *gin.Context) {

	if err := c.ShouldBindJSON(&collection); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := serviceSetup.QueryByOwner(collection.Owner)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var msg interface{}
	err =json.Unmarshal([]byte(result),&msg)
	if err!=nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg":"查询成功","result": msg})
}
