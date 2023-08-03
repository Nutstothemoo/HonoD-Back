package controller

import "github.com/gin-gonic/gin"

func GetOrderItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		
	}
}
func GetOrderItemByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		
	}
}
func GetOrderItemByOrderID() gin.HandlerFunc {
	return func(c *gin.Context) {
		
	}
}

func PostOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
func UpdateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		
	}
}
func DeleteOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
func ItemsByOrderID( id string ) (OrderItems []primitive.M, err error) {
	
}
