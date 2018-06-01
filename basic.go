package main


//获取接口ticket
func get_ticket(appid, appkey string) (bool) {
	if len(appid)<3 || len(appkey)<3 {
		RETURN["errmsg"] = "wrong appid or wrong appkey!"
		return true
	}

	//检查缓存是否存在
	access_info := CacheGetArray("access_"+appid)
	if(access_info["null"] == "null") {
		sql := "select * from api_keys where appid=?"
		//检查数据状态
		dbExamine()
		//执行sql
		query, err := dbService.Query(sql, appid)

		if err != nil {
			RETURN["errmsg"] = "filename:basic,message:sql error"
			return true
		}

		access_reg, err := sqlArray(query)
		if err != nil {
			RETURN["errmsg"] = "wrong appid or wrong appkey!"
			return true
		} else {
			access_info = access_reg[0]
			CacheSetArray("access_"+appid,access_info,"3600")
		}
	}

	//验证appid+appkey是否正确
	if access_info["appkey"] == appkey {
		//加密tocket
		ticket, err := encode(access_info["appid"] + "|" + timeNow)
		if err != nil {
			RETURN["errmsg"] = "filename:basic,message:Encryption errors"
		} else {
			//ticket加入缓存 7300秒...
			//CacheSet("ticket_"+ticket,"true","7300")

			RETURN["errcode"] = 0
			RETURN["errmsg"] = "ok"
			RETURN["ticket"] = ticket
		}
	} else {
		RETURN["errmsg"] = "wrong appid or wrong appkey!"
	}
	return true

}