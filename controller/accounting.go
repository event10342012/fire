package controller

import (
	"fire/model"
	"fire/server"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"net/http"
)

func GetAccountingTransactions(c *gin.Context) {
	db := server.GetDB()

	if c.Query("skip") != "" {
		skip := cast.ToInt(c.Query("skip"))
		db = db.Offset(skip)
	}
	if c.Query("limit") != "" {
		limit := cast.ToInt(c.Query("limit"))
		db = db.Limit(limit)
	}

	var accountingTransactions []model.AccountingTransaction
	db.Find(&accountingTransactions)

	c.JSON(200, accountingTransactions)
}

func GetAccountingTransaction(c *gin.Context) {
	id := c.Param("id")

	var accountTransaction model.AccountingTransaction
	db := server.GetDB()
	db.First(&accountTransaction, id)
	c.JSON(200, accountTransaction)
}

func CreateAccountingTransaction(c *gin.Context) {
	var accountTransaction model.AccountingTransaction
	if err := c.BindJSON(&accountTransaction); err != nil {
		c.JSON(http.StatusBadRequest, "invalid payload")
		return
	}

	db := server.GetDB()
	db.Create(&accountTransaction)
	c.JSON(200, accountTransaction)
}

func UpdateAccountingTransaction(c *gin.Context) {
	id := c.Param("id")

	var accountTransaction model.AccountingTransaction
	db := server.GetDB()
	db.First(&accountTransaction, id)

	if err := c.BindJSON(&accountTransaction); err != nil {
		c.JSON(http.StatusBadRequest, "invalid payload")
		return
	}

	db.Save(&accountTransaction)
	c.JSON(200, accountTransaction)
}

func DeleteAccountingTransaction(c *gin.Context) {
	id := c.Param("id")

	var accountTransaction model.AccountingTransaction
	db := server.GetDB()
	db.First(&accountTransaction, id)

	db.Delete(&accountTransaction)
	c.JSON(200, "deleted")
}
