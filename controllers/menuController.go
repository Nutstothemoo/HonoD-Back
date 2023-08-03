package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetMenus() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("GetMenu")
	}
}
func GetMenuByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("GetMenuByID")
	}
}
func PostMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("PostMenu")
	}
}
func UpdateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("UpdateMenu")
	}
}
func DeleteMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("DeleteMenu")
	}
}