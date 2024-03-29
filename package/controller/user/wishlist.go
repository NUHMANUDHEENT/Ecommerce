package controller

import (
	"project1/package/initializer"
	"project1/package/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func WishlistProducts(c *gin.Context) {
	var wishlist []models.Wishlist
	userid := c.GetUint("userid")
	if err := initializer.DB.Where("user_id=?", userid).Preload("Product").Find(&wishlist).Error; err != nil {
		c.JSON(502, gin.H{
			"error ": "failed to fetch wishlist items",
		})
		return
	}
	if len(wishlist) == 0 {
		c.JSON(502, gin.H{
			"message ": "No item found in wishlist",
		})
		return
	}
	for _, val := range wishlist {
		c.JSON(502, gin.H{
			"product id":    val.ProductId,
			"product name":  val.Product.Name,
			"product Image": val.Product.ImagePath1,
			"Product Price": val.Product.Price,
		})
	}
}
func WishlistAdd(c *gin.Context) {
	var wishAdd models.Wishlist
	userId := c.GetUint("userid")
	id := c.Param("ID")
	err := initializer.DB.Where("user_id=? AND product_id=?", userId, id).First(&wishAdd)
	if err.Error != nil {
		wishAdd.UserId = int(userId)
		wishAdd.ProductId, _ = strconv.Atoi(id)
		if err := initializer.DB.Create(&wishAdd).Error; err != nil {
			c.JSON(500, gin.H{
				"error": "Failed to add to wishlist",
			})
			return
		}
		c.JSON(500, gin.H{
			"message": "Item added to wishlist",
		})
	} else {
		c.JSON(500, gin.H{
			"error": "This item already added",
		})
	}
}

// ============== delete wishlist item ==============
func WishlistDelete(c *gin.Context) {
	var wishlistDelete models.Wishlist
	userId := c.GetUint("userid")
	id := c.Param("ID")
	if err := initializer.DB.Where("product_id=? AND user_id=?", id, userId).Delete(&wishlistDelete).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "failed to remove Item",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "Item remove successfully",
	})
}
