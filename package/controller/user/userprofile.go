package controller

import (
	"net/http"
	"project1/package/initializer"
	"project1/package/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func UserProfile(c *gin.Context) {
	var userProfile models.Users
	var userAddress []models.Address
	userId := c.GetUint("userid")

	if err := initializer.DB.First(&userProfile, userId).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "failed to find user"})
	} else {
		c.JSON(200, gin.H{
			"user name":  userProfile.Name,
			"user email": userProfile.Email,
			"user phone": userProfile.Phone,
			"user id":    userProfile.ID,
		})
	}
	if err := initializer.DB.Find(&userAddress, "user_id=?", userId).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "failed to find address"})
	} else {
		for _, val := range userAddress {
			c.JSON(200, gin.H{
				"user address":  val.Address,
				"user city":     val.City,
				"user pin code": val.Pincode,
				"user id":       val.ID,
				"user phone":    val.Phone,
			})
		}
	}
}
func AddressStore(c *gin.Context) {
	var userCheck models.Users
	var addAddress models.Address
	userId := c.GetUint("userid")
	if err := c.ShouldBindJSON(&addAddress); err != nil {
		c.JSON(500, gin.H{"error": err})
	}
	if err := initializer.DB.First(&userCheck, userId).Error; err != nil {
		c.JSON(404, gin.H{
			"error": "no user found",
		})
	} else {
		addAddress.UserId = int(userId)
		if result := initializer.DB.Create(&addAddress); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add address"})
		} else {
			c.JSON(201, gin.H{"message": "new address added successfully"})
		}
	}
}
func AddressEdit(c *gin.Context) {
	var addressEdit models.Address
	id := c.Param("ID")
	err := initializer.DB.First(&addressEdit, id)
	if err.Error != nil {
		c.JSON(404, gin.H{"error": "can't find address"})
	} else {
		err := c.ShouldBindJSON(&addressEdit)
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to binding data"})
		} else {
			if err := initializer.DB.Save(&addressEdit).Error; err != nil {
				c.JSON(500, gin.H{"error": "failed to update details"})
			} else {
				c.JSON(201, gin.H{"message": "address updated successfully"})
			}
		}
	}
}
func AddressDelete(c *gin.Context) {
	var deleteAddress models.Address
	session := sessions.Default(c)
	id := session.Get("userid")
	err := initializer.DB.First(&deleteAddress, id)
	if err.Error != nil {
		c.JSON(404, gin.H{"error": "can't find address"})
	} else {
		err := initializer.DB.Delete(&deleteAddress).Error
		if err != nil {
			c.JSON(500, gin.H{
				"error": "failed to delete address",
			})
		} else {
			c.JSON(200, gin.H{
				"message": "address deleted successfully",
			})
		}
	}
}
func EditUserProfile(c *gin.Context) {
	var editProfile models.Users
	userId := c.GetUint("userid")
	if err := initializer.DB.First(&editProfile, userId).Error; err != nil {
		c.JSON(404, gin.H{
			"error": "user not found",
		})
	} else {
		err := c.ShouldBindJSON(&editProfile)
		if err != nil {
			c.JSON(500, gin.H{
				"error": "failed to bind data",
			})
		} else {
			if err := initializer.DB.Save(&editProfile).Error; err != nil {
				c.JSON(500, gin.H{
					"error": "failed to update data",
				})
			} else {
				c.JSON(500, gin.H{
					"message": "updated data",
				})
			}
		}
	}
}
