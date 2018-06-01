package main

import (
	"database/sql"
)


//检查数据库状态 正常返回true
func dbExamine() (bool) {

	if err := dbService.Ping(); err != nil {
		RETURN["errmsg"] = "filename:db,message:Cannot link to the database"
		return false
	} else {
		return true
	}
}


//数据库查询结果集转数组
func sqlArray(query *sql.Rows) (map[int]map[string]string,error) {

	//定义最后得到的map
	results := make(map[int]map[string]string)

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

	i := 0
	for query.Next() {
		//循环，让游标往下推
		if err := query.Scan(scans...); err != nil {
			//query.Scan查询出来的不定长值放到scans[i] = &values[i],也就是每行都放在values里
			return results,err
		}

		row := make(map[string]string) //每行数据

		for k, v := range values {
			//每行数据是放在values里面，现在把它挪到row里
			key := cols[k]
			row[key] = string(v)
		}
		results[i] = row //装入结果集中
		i++
	}

	return results,nil

}