package controller

import (
	"fmt"
	"net/http"
	"project1/package/handler"
	"project1/package/initializer"
	"project1/package/middleware"
	"project1/package/models"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var LogJs models.Users
var otp string
var RoleUser = "user"

func UserSignUp(c *gin.Context) {
	LogJs = models.Users{}
	var otpStore models.OtpMail
	err := c.ShouldBindJSON(&LogJs)
	if err != nil {
		c.JSON(501, gin.H{"error": "json binding error"})
		return
	}

	if err := initializer.DB.First(&LogJs, "email=?", LogJs.Email).Error; err == nil {
		c.JSON(501, gin.H{"error": "Email address already exist"})
		return
	}
	otp = handler.GenerateOtp()
	fmt.Println("----------------", otp, "-----------------")
	err = handler.SendOtp(LogJs.Email, otp)
	if err != nil {
		c.JSON(500, "failed to send otp")
	} else {
		c.JSON(202, "otp send to mail  "+otp)
		result := initializer.DB.First(&otpStore, "email=?", LogJs.Email)
		if result.Error != nil {
			otpStore = models.OtpMail{
				Otp:       otp,
				Email:     LogJs.Email,
				CreatedAt: time.Now(),
				ExpireAt:  time.Now().Add(30 * time.Second),
			}
			err := initializer.DB.Create(&otpStore)
			if err.Error != nil {
				c.JSON(500, gin.H{"error": "failed to save otp details"})
				return
			}
		} else {
			err := initializer.DB.Model(&otpStore).Where("email=?", LogJs.Email).Updates(models.OtpMail{
				Otp:      otp,
				ExpireAt: time.Now().Add(15 * time.Second),
			})
			if err.Error != nil {
				c.JSON(500, "failed too update data")
			}
		}
	}
}

func OtpCheck(c *gin.Context) {
	var otpcheck models.OtpMail
	var otpExistTable models.OtpMail
	initializer.DB.First(&otpExistTable, "email=?", LogJs.Email)

	err := c.ShouldBindJSON(&otpcheck)
	if err != nil {
		c.JSON(500, "failed to bind otp details")
	}
	var existingOTP models.OtpMail
	if err := initializer.DB.Where("otp = ? AND expire_at > ?", otpcheck.Otp, time.Now()).First(&existingOTP).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired OTP"})
		return
	} else {
		fmt.Println("currect otp")
		HashPass, err := bcrypt.GenerateFromPassword([]byte(LogJs.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(501, gin.H{"error": "hashing error"})
		}

		LogJs.Password = string(HashPass)
		LogJs.Blocking = true
		erro := initializer.DB.Create(&LogJs)
		if erro.Error != nil {
			c.JSON(500, "failed to signup")
		} else {
			if err := initializer.DB.Delete(&otpExistTable).Error; err != nil {
				c.JSON(500, "failed to delete otp data")
			}
			if err := initializer.DB.First(&LogJs).Error; err != nil {
				c.JSON(501, gin.H{
					"error": "failed to fetch user details for wallet",
				})
				return
			}
			initializer.DB.Create(&models.Wallet{
				User_id: int(LogJs.ID),
			})
			c.JSON(201, gin.H{"message": "user created successfully"})
		}
	}
}
func ResendOtp(c *gin.Context) {
	var otpStore models.OtpMail
	otp = handler.GenerateOtp()
	err := handler.SendOtp(LogJs.Email, otp)
	if err != nil {
		c.JSON(500, "failed to send otp")
	} else {
		c.JSON(202, "otp resend to mail  "+otp)
		result := initializer.DB.First(&otpStore, "email=?", LogJs.Email)
		if result.Error != nil {
			otpStore = models.OtpMail{
				Otp:       otp,
				Email:     LogJs.Email,
				CreatedAt: time.Now(),
				ExpireAt:  time.Now().Add(15 * time.Second),
			}
			err := initializer.DB.Create(&otpStore)
			if err.Error != nil {
				c.JSON(500, gin.H{"error": "failed to save otp details"})
			}
		} else {
			err := initializer.DB.Model(&otpStore).Where("email=?", LogJs.Email).Updates(models.OtpMail{
				Otp:      otp,
				ExpireAt: time.Now().Add(15 * time.Second),
			})
			if err.Error != nil {
				c.JSON(500, "failed to update data")
			}
		}
	}
}

func UserLogin(c *gin.Context) {
	LogJs = models.Users{}
	var userPass models.Users
	err := c.ShouldBindJSON(&LogJs)
	if err != nil {
		c.JSON(501, gin.H{"error": "error binding data"})
	}
	fmt.Println(LogJs)
	initializer.DB.First(&userPass, "email=?", LogJs.Email)
	err = bcrypt.CompareHashAndPassword([]byte(userPass.Password), []byte(LogJs.Password))
	if err != nil {
		c.JSON(501, gin.H{"error": "invalid username or password"})
	} else {
		if !userPass.Blocking {
			c.JSON(300, gin.H{
				"message": "User blocked"})
		} else {
			token := middleware.JwtTokenStart(c, userPass.ID, userPass.Email, RoleUser)
			c.SetCookie("jwtToken",token,int((time.Hour* 1).Seconds()),"/","localhost",false,true)
			c.JSON(200, gin.H{
				"message": "login successfully",
				"token":   token,
			})
		}
	}
}
func UserLogout(c *gin.Context) {

	c.SetCookie("jwtToken", "", -1, "", "", false, false)
	c.JSON(200, gin.H{
		"message": "logout Successfull",
	})

}
