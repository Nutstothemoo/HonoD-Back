package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("GetUsers")
	}
}
func GetUserByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("GetUserByID")
	}
}
func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("SignUp")
	}
}
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("Login")
	}
}
func UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("UpdateUser")
	}
}
func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("DeleteUser")
	}
}
func HashPassword( password string) string {
}
func VerifyPassword( Userpassword string, providedPassword string) bool {	
}