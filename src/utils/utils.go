package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"qiniupkg.com/x/errors.v7"
	"strconv"
	"strings"
	"time"
)

//md5加密
func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

//获取post json 文本参数
func GetRawData (c *gin.Context, keys []string) (map[string]string, error) {
	buf, _ := c.GetRawData()
	requstMap := make(map[string]string)
	err := json.Unmarshal(buf, &requstMap)
	if keys != nil && err == nil {
		for _, v := range keys {
			if _, ok := requstMap[v]; !ok {
				return requstMap, errors.New("参数错误")
			}
		}
	}
	return requstMap, err
}

//根据生日获取年龄
func GetAgeByBirthday(birthday string) int {
	birthdayArr := strings.Split(birthday, "-")
	birYear, _ := strconv.Atoi(birthdayArr[0])
	birMonth, _ := strconv.Atoi(birthdayArr[1])
	age := time.Now().Year() - birYear
	if int(time.Now().Month()) < birMonth {
		age--
	}
	return age
}

//获取本周一日期
func GetThisWeekMonday() string {
	return time.Now().AddDate(0, 0, 1+(-1*int(time.Now().Weekday()))).Format("20060102")
}