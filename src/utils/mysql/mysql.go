package mysql

import (
	"conf"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"strings"
	"sync"
	"time"
)

type Mysql struct {
	sqldb *sql.DB
}

var mysql *Mysql
var once sync.Once

// 获取MySQL实例(单例)
func GetMysql() *Mysql{
	once.Do(func() {
		mysql = &Mysql{}
		mysql.Conn()
	})
	return mysql
}

// 连接mysql
func (this *Mysql) Conn () error {
	var err error
	this.sqldb, err = sql.Open("mysql", conf.GetConfig().Mysql)
	if err != nil {
		glog.Fatalln(err)
	}
	// 设置连接池
	this.sqldb.SetConnMaxLifetime(2 * time.Hour)
	this.sqldb.SetMaxOpenConns(10)
	this.sqldb.SetMaxIdleConns(50)
	err = this.sqldb.Ping()
	if err != nil {
		glog.Fatalln(err)
	}
	return err
}

// 执行存储过程、查询数据库
func (this *Mysql) Call (procName string, args ...interface{}) (string, error) {
	// 1. 组装参数
	argv := make([]string, 0)
	for i := 0; i < len(args); i++ {
		argv = append(argv, "?")
	}
	query := fmt.Sprintf(" CALL `%s`(%s) ", procName, strings.Join(argv, ","))
	glog.Infoln("Query: " , query, args)

	// 2. 预执行语句处理
	stmt, err := this.sqldb.Prepare(query)
	if err != nil {
		glog.Infof("Prepare Error: %v\n" , err)
		return "", err
	}
	defer stmt.Close()

	// 3. 执行存储过程
	rows, err := stmt.Query(args...)
	if err != nil {
		glog.Infof("Rows Error: %v\n" , err)
		return "", err
	}
	defer rows.Close()

	// 4. 处理返回结果集
	iResult := make(map[string]interface{}, 0)
	for {
		columns, err := rows.Columns()
		if err != nil {
			glog.Infof("Columns Error: %v\n" , err)
			return "", err
		}
		colCount := len(columns)
		values := make([]interface{}, colCount)
		valuePtrs := make([]interface{}, colCount)
		list := make([]map[string]interface{}, 0)
		for rows.Next() {
			for i := 0; i < len(valuePtrs); i++ {
				valuePtrs[i] = &values[i]
			}
			rows.Scan(valuePtrs...)
			info := make(map[string]interface{})
			for i, col := range columns {
				var v interface{}
				val := values[i]
				b, ok := val.([]byte)
				if ok {
					v = string(b)
				} else {
					v = val
				}
				info[col] = v
			}
			list = append(list, info)
		}
		if len(list) > 0 {
			if _, ok := list[0]["total"]; ok {
				iResult["total"] = list[0]["total"]
			} else {
				iResult["list"] = list
			}
		}
		if !rows.NextResultSet() {
			break
		}
	}

	// 4. json encode
	var jsonStr []byte
	if len(iResult) == 1 {
		jsonStr, err = json.Marshal(iResult["list"])
	} else {
		jsonStr, err = json.Marshal(iResult)
	}
	if err != nil {
		glog.Infof("Json Error: %v\n", err)
		return "", err
	}
	return string(jsonStr), nil
}