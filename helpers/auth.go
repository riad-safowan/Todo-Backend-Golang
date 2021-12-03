package helpers

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func MatchUserTypeToUid(c *gin.Context, userId string) (err error) {
	userType := c.GetString("user_type")
	uid := c.GetString("user_id")
	err = nil
	if userType == "USER" && uid != userId {
		err = errors.New("unauthorized to access this resources")
		return err
	}
	return err
}

func CheckUserType(c *gin.Context, role string) (err error) {
	userType := c.GetString("user_type")
	err = nil
	if userType != role{
		err = errors.New("unauthorized to access this resources")
	}
	return err
}
