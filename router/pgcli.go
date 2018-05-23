package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	//"github.com/gin-gonic/gin/binding"
	"encoding/json"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Users struct {
	gorm.Model
	OpenID     string `json:"OpenID"`
	SessionKey string `json:"SessionKey"`
	Name       string `json:"Name"`
	Area       string `json:"Area"`
	Phone      string `json:"Phone"`
}

type CNCProgList struct {
	gorm.Model
	//ID        uint    //`json:"id"`
	Consumer     string `form:"Consumer" json:"Consumer"`
	Workpiece    string `form:"Workpiece" json:"Workpiece"`
	Programer    string `form:"Programer" json:"Programer"`
	CNCPath      string `form:"CNCPath" json:"CNCPath"`
	Snapshot1    string `form:"Snapshot1" json:"Snapshot1"`
	Verticalview string `form:"Verticalview" json:"Verticalview"`
	Sideview     string `form:"Sideview" json:"Sideview"`

	TP []ToolPath //`json:"刀具路径"`
}

type ToolPath struct {
	//ID	 uint `json:"序号"`
	Style       string `json:"类别"`
	Name        string `json:"程序名称"`
	Tool        string `json:"刀具"`
	ToolsType   string `json:"刀具类型"`
	ToolsNumber string `json:"刀号"`
	Allowance   string `json:"余量"`
	Bottom      string `json:"最低点"`
	Convtime    string `json:"时间"`
	Note        string `json:"备注"`
}

func RedisCache(c *gin.Context) {
	//如果在缓存中有
	v, err := RedisGet("Mygo")
	if err == nil && v != "" {
		c.JSON(http.StatusOK, getCNProgListFromStream(v))//
	} else {
		db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=gorm sslmode=disable password=68957423")
		if err != nil {
			c.Error(err)
		}
		defer db.Close()
		cncProgList := make([]CNCProgList, 0, 20)
		if err := db.Find(&cncProgList).Error; err != nil {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			slice := getCNProList10(cncProgList)
			v, err := json.Marshal(slice)
			if err == nil {
				err := RedisSet("Mygo", v)
				if err != nil {
					c.Error(err)
				}
			}
			c.JSON(http.StatusOK, slice)
		}
	}
}

func getCNProgListFromStream(val string) []CNCProgList {
	cncProgList := make([]CNCProgList, 0, 20)
	json.Unmarshal([]byte(val), &cncProgList)
	return getCNProList10(cncProgList)
}

func getCNProList10(list []CNCProgList) []CNCProgList {
	lenth := len(list)
	if lenth > 10 {
		return list[lenth-10 : lenth]
	} else {
		return list
	}
}

//抓取toop100的数据
func LevelDBCache(c *gin.Context) {
	//levelDB缓存
	v, err := DBGet("MYGO")
	if err == nil && v != "" {
		c.JSON(http.StatusOK, getCNProgListFromStream(v))
	} else {
		db, err := gorm.Open("postgres", "host=127.0.0.1 user=postgres dbname=gorm sslmode=disable password=68957423")
		if err != nil {
			c.Error(err)
		}
		defer db.Close()
		cncProgList := make([]CNCProgList, 0, 20)
		if err := db.Find(&cncProgList).Error; err != nil {
			c.AbortWithStatus(404)
		} else {
			slice := getCNProList10(cncProgList)
			v, err := json.Marshal(slice)
			if err == nil {
				err := DBSet("MYGO", v)
				if err != nil {
					c.Error(err)
				}
			}
			c.JSON(http.StatusOK, slice)
		}
	}
}
