/*
* Open-Falcon
*
* Copyright (c) 2014-2018 Xiaomi, Inc. All Rights Reserved.
*
* This product is licensed to you under the Apache License, Version 2.0 (the "License").
* You may not use this product except in compliance with the License.
*
* This product may include a number of subcomponents with separate copyright notices
* and license terms. Your use of these subcomponents is subject to the terms and
* conditions of the subcomponent's license, as noted in the LICENSE file.
 */

package common

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native" //with mysql
)

type DbConn struct {
	Conn   mysql.Conn
	Health bool
}

// NewMySQLConnection the constructor of mysql connecting
func NewMySQLConnection(conf DatabaseConf) (*DbConn, error) {
	return initMySQLConnection(conf)
}

// QueryResult the result of query
func initMySQLConnection(conf DatabaseConf) (conn *DbConn, err error) {
	db := mysql.New("tcp", "", fmt.Sprintf(
		"%s:%d", conf.Host, conf.Port),
		conf.User, conf.Password)
	db.SetTimeout(5000 * time.Millisecond) //2S超时
	if err = db.Connect(); err != nil {
		err = errors.Wrap(err, "Building mysql connection failed!")
	}
	conn = &DbConn{
		Conn:   db,
		Health: true,
	}
	if err := conn.Conn.Ping(); err != nil {
		conn.Health = false //失败
	}
	//go conn.HealthCheck() //开启健康检查
	return
}

func (dbConn DbConn) HealthCheck() {
	t := time.NewTicker(10000 * time.Millisecond) //每10S中检查一次
	defer func() {
		t.Stop()
	}()
	for {
		if _, _, err := dbConn.Conn.Query("select @@version;"); err != nil { //连接异常
			dbConn.Health = false //首先将这个连接健康度设置为False
			if err = dbConn.Conn.Reconnect(); err == nil {
				dbConn.Health = true
			}
		}
		<-t.C
	}

}
