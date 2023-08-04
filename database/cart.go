package database

import (
	"errors"
)

var {
	ErrorCantFindProdruct = errors.New("Can't find product")
	ErrorCantFindUser = errors.New("Can't find user, this userId is not valid")
	ErrorCantUpdateUser = errors.New("Can't update user")
	ErrorCantGetItemFromCart = errors.New("Can't get item from cart")
	ErrorCantAddItemToCart = errors.New("Can't add item to cart")
	ErrorCantRemoveItemFromCart = errors.New("Can't remove item from cart")
}

func AddProductToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}
func RemoveItemFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}
func BuyItemFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}
func InstantBuy() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}
