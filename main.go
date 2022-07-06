package main

import (
	"net/http"

	"github.com/AntoineAugusti/moduluschecking/models"
	"github.com/AntoineAugusti/moduluschecking/parsers"
	"github.com/AntoineAugusti/moduluschecking/resolvers"
	"github.com/gin-gonic/gin"
)

type bankAccountRequest struct {
	SortCode      string `json:"sort_code"`
	AccountNumber string `json:"account_number"`
}

type ValidityResponse struct {
	SortCode      string `json:"sort_code"`
	AccountNumber string `json:"account_number"`
	Valid         bool   `json:"is_valid"`
}

func handleGetHealth(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"message": "healthcheck is good"})

}

func handlePostVerify(c *gin.Context) {
	var verifyBank bankAccountRequest

	if err := c.BindJSON(&verifyBank); err != nil {
		return
	}

	sortCodeLength := len(verifyBank.SortCode)
	accountNumberLength := len(verifyBank.AccountNumber)

	cksum := sortCodeLength == 6 && (accountNumberLength >= 6 && accountNumberLength <= 10)
	if !cksum {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Expected a 6 digits sort code and an account number between 6 and 10 digits."})
		return
	}

	parser := parsers.CreateFileParser()
	resolver := resolvers.NewResolver(parser)
	bankAccount := models.CreateBankAccount(verifyBank.SortCode, verifyBank.AccountNumber)
	isValid := resolver.IsValid(bankAccount)

	response := ValidityResponse{
		Valid:         isValid,
		SortCode:      verifyBank.SortCode,
		AccountNumber: verifyBank.AccountNumber,
	}

	c.IndentedJSON(http.StatusOK, response)
}

func main() {
	router := gin.Default()
	router.POST("/verify", handlePostVerify)
	router.GET("/health", handleGetHealth)

	router.Run()
}
