package controller

import (
	"project1/package/initializer"
	"project1/package/middleware"
	"project1/package/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AdminPage(c *gin.Context) {
	c.JSON(200, gin.H{"message": "welcome admin page"})
}

func AdminLogin(c *gin.Context) {
	var AdminCheck models.Users
	session := sessions.Default(c)
	check := session.Get("admin")
	if check != nil {
		c.JSON(200, "Already logged")
	} else {
		err := c.ShouldBindJSON(&AdminCheck)
		if err != nil {
			c.JSON(501, gin.H{"error": "error binding data"})
		}
		if AdminCheck.Email == "nuhman@gmail.com" && AdminCheck.Password == "nuhman@1234" {
			middleware.SessionCreate("nuhman@gmail.com", "admin", c)
			c.JSON(202, gin.H{"message": "successfully logged"})
		} else {
			c.JSON(501, gin.H{"error": "invalid username or password"})
		}
	}
}
func AdminLogout(c *gin.Context) {
	session := sessions.Default(c)
	check := session.Get("admin")
	if check == nil {
		c.JSON(200, "Unauthorized")
	} else {
		session.Delete("admin")
		session.Save()
		c.JSON(200, gin.H{
			"message": "logout successfull",
		})
	}
}
func UserList(c *gin.Context) {
	var user_managment []models.Users
	initializer.DB.Order("ID").Find(&user_managment)
	for _, val := range user_managment {
		c.JSON(200, gin.H{
			"ID":         val.ID,
			"name":       val.Name,
			"Email":      val.Email,
			"gender":     val.Gender,
			"created At": val.CreatedAt,
			"status":     val.Blocking,
		})
	}
}
func EditUserDetails(c *gin.Context) {
	var userEdit models.Users
	id := c.Param("ID")
	err := initializer.DB.First(&userEdit, id)
	if err.Error != nil {
		c.JSON(500, gin.H{"error": "can't find user"})
	} else {
		err := c.ShouldBindJSON(&userEdit)
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to bindinng data"})
		} else {
			if err := initializer.DB.Save(&userEdit).Error; err != nil {
				c.JSON(500, gin.H{"error": "failed to update details"})
			} else {
				c.JSON(200, gin.H{"message": "User updated successfully"})
			}
		}
	}
}
func BlockUser(c *gin.Context) {
	var blockUser models.Users
	id := c.Param("ID")
	err := initializer.DB.First(&blockUser, id)
	if err.Error != nil {
		c.JSON(500, gin.H{"error": "can't find user"})
	} else {
		if blockUser.Blocking {
			blockUser.Blocking = false
			c.JSON(200, "user blocked")
		} else {
			blockUser.Blocking = true
			c.JSON(200, "user unblocked")
		}
		if err := initializer.DB.Save(&blockUser).Error; err != nil {
			c.JSON(500, "failed to block/unblock user")
		}
	}
}
func DeleteUser(c *gin.Context) {
	var deleteUser models.Users
	id := c.Param("ID")
	err := initializer.DB.First(&deleteUser, id)
	if err.Error != nil {
		c.JSON(500, gin.H{"error": "can't find user"})
	} else {
		err := initializer.DB.Delete(&deleteUser).Error
		if err != nil {
			c.JSON(500, "failed to delete user")
		} else {
			c.JSON(200, "user deleted successfully")
		}
	}
}