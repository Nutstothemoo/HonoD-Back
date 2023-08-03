package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("GetOrders")
	}
}
func GetOrderByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("GetOrderByID")
	}
}
func PostOrder() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
func UpdateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		
	}
}
func DeleteOrder() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
