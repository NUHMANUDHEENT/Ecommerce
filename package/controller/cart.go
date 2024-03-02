package controller

import (
	"project1/package/initializer"
	"project1/package/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CartStore(c *gin.Context) {
	var cartStore models.Cart
	id := c.Param("ID")
	err := c.ShouldBind(&cartStore)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "failed to bind data",
		})
	} else {
		prodcutId, _ := strconv.Atoi(id)
		cartStore.ProductId = prodcutId
		if err := initializer.DB.Create(&cartStore).Error; err != nil {
			c.JSON(500, gin.H{
				"error": "can't find product",
			})
		} else {
			c.JSON(500, gin.H{
				"message": "product added to cart",
			})
		}
	}
}
func CartView(c *gin.Context) {
	var cartView []models.Cart
	var cartBind models.Cart
	var totalAmount int
	err := c.ShouldBindJSON(&cartBind)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "failed to bind data",
		})
	} else {
		if err := initializer.DB.Joins("Products").Find(&cartView).Where("UserId=?", cartBind.UserId).Error; err != nil {
			c.JSON(500, gin.H{
				"error": "failed to fetch data",
			})
		} else {
			for _, val := range cartView {
				c.JSON(200, gin.H{
					"product name":  val.Product.Name,
					"product image": val.Product.ImagePath1,
					"product price": val.Product.Price,
				})
				totalAmount += int(val.Product.Price)
			}
			c.JSON(200, gin.H{
				"Total Amount": totalAmount,
			})
		}
	}
}