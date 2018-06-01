package main

import (
	"database/sql" //这包一定要引用，是底层的sql驱动
	"github.com/gin-gonic/gin"
	_"github.com/go-sql-driver/mysql"
	"net/http"
	"github.com/Unknwon/goconfig"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"time"
)

//配置
var CONFIG map[string]string
//缓存
var Cache redis.Conn
//数据库
var dbService *sql.DB
//返回参数
var RETURN map[string]interface{}
//当前时间
var timeNow string

func main() {

	//------------配置初始化 开始
	cfg, err := goconfig.LoadConfigFile("config.ini")
	if err != nil {

	}
	CONFIG = make(map[string]string)
	//全局设置环境
	CONFIG["APP_DEBUG"], _ = cfg.GetValue("APP", "APP_DEBUG")
	//密钥配置
	CONFIG["APP_KEY"], _ = cfg.GetValue("APP", "APP_KEY")
	//数据库配置
	CONFIG["DB_HOST"], _ = cfg.GetValue("DB", "DB_HOST")
	CONFIG["DB_PORT"], _ = cfg.GetValue("DB", "DB_PORT")
	CONFIG["DB_USERNAME"], _ = cfg.GetValue("DB", "DB_USERNAME")
	CONFIG["DB_PASSWORD"], _ = cfg.GetValue("DB", "DB_PASSWORD")
	CONFIG["DB_DATABASE"], _ = cfg.GetValue("DB", "DB_DATABASE")
	//缓存配置
	CONFIG["CACHE_PREFIX"], _ = cfg.GetValue("CACHE", "CACHE_PREFIX")
	CONFIG["CACHE_REDIS_CONNECTION"], _ = cfg.GetValue("CACHE", "CACHE_REDIS_CONNECTION")
	CONFIG["CACHE_PORT"], _ = cfg.GetValue("CACHE", "CACHE_PORT")
	//--------------配置初始化 结束

	//全局设置运行环境
	if CONFIG["APP_DEBUG"] == "true" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	//--------------------缓存初始化开始
	Cache, err = redis.Dial("tcp", CONFIG["CACHE_REDIS_CONNECTION"] + ":" + CONFIG["CACHE_PORT"])
	//-------------------缓存初始化结束

	//---------------数据库配置初始化开始
	dbService, err = sql.Open("mysql", CONFIG["DB_USERNAME"] + ":" + CONFIG["DB_PASSWORD"] + "@tcp(" + CONFIG["DB_HOST"] + ":" + CONFIG["DB_PORT"] + ")/" + CONFIG["DB_DATABASE"] + "?parseTime=true")
	//defer dbService.Close()
	if err != nil {

	}
	dbService.SetMaxIdleConns(20)
	dbService.SetMaxOpenConns(20)
	//-------------数据库初始化结束

	//-----------------路由初始化
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	//----------------路由初始化结束

	//---------------------php_BasicController 开始
	basic := router.Group("/basic")
	{
		//获取接口ticket
		basic.GET("/get_ticket", func(c *gin.Context) {
			//参数初始化
			initialize()
			appid := c.Query("appid")
			token := c.Query("appkey")
			get_ticket(appid,token)
			//响应请求
			c.JSON(http.StatusOK, RETURN)
		})

		//获取app配置接口get_js_config

	}
	//---------------------php_BasicController 结束


	//------------------测试开始
	//token:="eyJpdiI6ImFqeFVTRTY0aWVcL0FPZFZ5YUlOT0NBPT0iLCJ2YWx1ZSI6IlNcL2lQbys4ZWx0TU9yUzdNQlJnTytvSVFyVHBPRERkNmxOekRcL0MyK3dBMD0iLCJtYWMiOiJlZTliNjY0ZDM0NDkxMmYwMWVhODAzMzk3ZDgxYjI1NzUzMGU4NWRhOTNlYzM3MzhlZWNhODRjZDRjNmYwYTE2In0="
	//fmt.Println(getUserIdByToken(token));
	api := router.Group("/api")
	{
		api.GET("/sso/verify", func(c *gin.Context) {
			token := c.Query("token")
			//fmt.Println(token)
			dingid := getUserIdByToken(token)
			//判断token是否正确解析
			if dingid == "" {

				row := make(map[string]string)
				row["errcode"] = string("3")
				row["errmsg"] = string("wrong ticket!")
				c.JSON(http.StatusOK, row)

			} else {
				//测试
				//dingid := c.Query("dingid")
				//c.String(http.StatusOK,dingid)

				var userData string
				//测试缓存查询
				userData = CacheGet("dingid_" + dingid)

				if userData == "null" {

					userData = getUserByDingId(dingid)

					//测试缓存添加 true false
					CacheSet("dingid_" + dingid, userData,"10")

				}
				if userData == "null" {
					userReg := make(map[string]string)
					userReg["errcode"] = string("1")
					userReg["errmsg"] = string("The user does not exist")
					c.JSON(http.StatusOK, userReg)
				} else {
					//加密示例
					encodeStr, err := encode("1111|2222")
					if err != nil {

					}
					c.String(http.StatusOK, encodeStr)

					//解密示例
					decodeStr, err := decode(encodeStr)
					if err != nil {

					}

					c.String(http.StatusOK, decodeStr)

					c.String(http.StatusOK, userData)
					//c.JSON(http.StatusOK, userData)
				}

			}
		})
	}
	//-----------------------------测试结束

	router.Run(":8800")
}

//参数初始化
func initialize()  {

	//初始化当前时间
	timeNow = strconv.FormatInt(time.Now().Unix(),10)
	//返回参数初始化
	RETURN = make(map[string]interface{})
	RETURN["errcode"] = 1

}