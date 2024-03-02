package controller

import (
	"net/http"
	"project1/package/initializer"
	"project1/package/models"

	"github.com/gin-gonic/gin"
)

func AddressStore(c *gin.Context) {
	var userCheck models.Users
	var addAddress models.Address
	if err := c.ShouldBindJSON(&addAddress); err != nil {
		c.JSON(500, gin.H{"error": err})
	}
	if err := initializer.DB.First(&userCheck, addAddress.UserId).Error; err != nil {
		c.JSON(500, "no user found")
	} else {
		if result := initializer.DB.Create(&addAddress); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add address"})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "new address added successfully"})
		}
	}
}
func AddressEdit(c *gin.Context) {
	var addressEdit models.Address
	id := c.Param("ID")
	err := initializer.DB.First(&addressEdit, id)
	if err.Error != nil {
		c.JSON(500, gin.H{"error": "can't find address"})
	} else {
		err := c.ShouldBindJSON(&addressEdit)
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to binding data"})
		} else {
			if err := initializer.DB.Save(&addressEdit).Error; err != nil {
				c.JSON(500, gin.H{"error": "failed to update details"})
			} else {
				c.JSON(200, gin.H{"message": "address updated successfully"})
			}
		}
	}
}
func AddressDelete(c *gin.Context) {
	var deleteAddress models.Address
	id := c.Param("ID")
	err := initializer.DB.First(&deleteAddress, id)
	if err.Error != nil {
		c.JSON(500, gin.H{"error": "can't find address"})
	} else {
		err := initializer.DB.Delete(&deleteAddress).Error
		if err != nil {
			c.JSON(500, "failed to delete address")
		} else {
			c.JSON(200, "address deleted successfully")
		}
	}
}
func UserProfile(c *gin.Context) {
	var userProfile models.Users
	var userAddress []models.Address
	id := c.Param("ID")
	if err := initializer.DB.First(&userProfile, id).Error; err != nil {
		c.JSON(500, "failed to find user")
	} else {
		c.JSON(200, gin.H{
			"user name":  userProfile.Name,
			"user email": userProfile.Email,
			"user phone": userProfile.Phone,
		})
	}
	if err := initializer.DB.Find(&userAddress, "user_id=?", id).Error; err != nil {
		c.JSON(500, "failed to find address")
	} else {
		for _, val := range userAddress {
			c.JSON(200, gin.H{
				"user address":  val.Address,
				"user city":     val.City,
				"user pin code": val.Pincode,
			})
		}
	}
}