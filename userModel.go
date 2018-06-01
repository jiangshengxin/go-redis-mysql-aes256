package main

import (
	//"database/sql" //这包一定要引用，是底层的sql驱动
	_ "github.com/go-sql-driver/mysql"
	"encoding/json"
	//"fmt"
	//"github.com/garyburd/redigo/redis"
)

func getUserByDingId(dingid string) (string) {

	if err := dbService.Ping(); err != nil {

	}

	query, err := dbService.Query("select * from ding_users where dingid=?", dingid)

	if err != nil {

	}

	//读出查询出的列字段名
	cols, _ := query.Columns()
	//values是每个列的值，这里获取到byte里
	values := make([][]byte, len(cols))
	//query.Scan的参数，因为每次查询出来的列是不定长的，用len(cols)定住当次查询的长度
	scans := make([]interface{}, len(cols))
	//让每一行数据都填充到[][]byte里面
	for i := range values {
		scans[i] = &values[i]
	}

	//最后得到的map
	results := make(map[int]map[string]string)
	i := 0
	for query.Next() {
		//循环，让游标往下推
		if err := query.Scan(scans...); err != nil {
			//query.Scan查询出来的不定长值放到scans[i] = &values[i],也就是每行都放在values里

			return ""
		}

		row := make(map[string]string) //每行数据

		for k, v := range values {
			//每行数据是放在values里面，现在把它挪到row里
			key := cols[k]
			row[key] = string(v)
		}
		row["errcode"] = string("0")
		row["errmsg"] = string("success")
		results[i] = row //装入结果集中
		i++
	}

	//查询出来的数组
	// for k, v := range results {
	// 	fmt.Println(k, v)
	// }
	b, err := json.Marshal(results[0])

	//fmt.Println(results)
	//fmt.Println(string(b))
	//cache_key:="ding_user_"+results[0]["dingid"]
	//defer dbService.Close()
	return string(b)
}
