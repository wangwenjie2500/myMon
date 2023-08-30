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

package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/open-falcon/mymon/common"
	"github.com/robfig/cron"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"syscall"
	"time"
)

// Global tag var
var (
	IsSlave    int
	IsReadOnly int
	Tag        string
)

//Log logger of project
var Log *logs.BeeLogger

func main() {
	// parse config file
	var confFile string
	flag.StringVar(&confFile, "c", "myMon.toml", "myMon configure file")
	version := flag.Bool("v", false, "show version")
	flag.Parse()
	if *version {
		fmt.Println(fmt.Sprintf("%10s: %s", "Version", Version))
		os.Exit(0)
	}
	conf, err := common.NewConfig(confFile)
	if err != nil {
		fmt.Printf("NewConfig Error: %s\n", err.Error())
		return
	}
	if conf.Base.LogDir != "" {
		err = os.MkdirAll(conf.Base.LogDir, 0755)
		if err != nil {
			fmt.Printf("MkdirAll Error: %s\n", err.Error())
			return
		}
	}

	// init log and other necessary
	Log = common.MyNewLogger(conf, common.CompatibleLog(conf))
	c := cron.New()
	_, err = c.AddFunc("*/1 * * * *", func() {
		Log.Info("Start MySql Monitor...")
		// init cmdb mysql
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30) //如果30S内这个节点没有拿到数据就强制中断
		defer cancel()
		for _, dbConfig := range conf.Mysql {
			Log.Debug(fmt.Sprintf("Create New Instance Conn Host: %s user: %s Port: %d", dbConfig.Host, dbConfig.User, dbConfig.Port))
			go fetchData(conf, dbConfig, ctx)
		}
		Log.Info("End MySQL Monitor")
	})
	c.Start()
	stop()
	c.Stop()

}

func timeout() {
	time.AfterFunc(TimeOut*time.Second, func() {
		Log.Error("Execute timeout")
		os.Exit(1)
	})
}

func stop() {
	sigle := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigle, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigle
		Log.Debug("quit .....")
		done <- true
	}()
	<-done
}

func fetchData(conf *common.Config, dbConfg common.DatabaseConf, ctx context.Context) {
	db, err := common.NewMySQLConnection(dbConfg)
	if err != nil {
		Log.Error("NewMySQLConnection Error: %s Host: %s\n", err.Error(), dbConfg.Host)
		return
	}
	defer func() {
		MySQLAlive(dbConfg, err == nil)
		if err := recover(); err != nil {
			Log.Error("fetch data error panic: ", err)
		} else {
			Log.Debug("fetch data success")
		}
	}()

	funcs := []func(common.DatabaseConf, mysql.Conn) ([]*MetaData, error){
		ShowSlaveStatus,
		ShowGlobalStatus,
		ShowGlobalVariables,
		ShowInnodbStatus,
		ShowBinaryLogs,
	}

	go func(fs []func(common.DatabaseConf, mysql.Conn) ([]*MetaData, error)) {

		// SHOW XXX Metric
		var data []*MetaData

		for _, fn := range fs {
			value, err := fn(dbConfg, db.Conn)
			if err != nil {
				fname := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
				Log.Error("Get fnanem: %s error: %+v", fname, err)
			}
			data = append(data, value...)
		}
		//Set Tag
		for _, metric := range data {
			metric.SetTag(fmt.Sprintf("port=%d", dbConfg.Port))
		}
		// Send Data to falcon-agent
		msg, err := SendData(data, conf)
		if err != nil {
			Log.Error("Send response Error %s:%d - %s", dbConfg.Host, dbConfg.Port, string(msg))
		} else {
			Log.Debug("Send response Success %s:%d - %s", dbConfg.Host, dbConfg.Port, string(msg))
		}

	}(funcs)
	select {
	case <-ctx.Done():
		t, b := ctx.Deadline()
		if b {
			Log.Debug(fmt.Sprintf("instance fetdata success: %s %s daedline: %d", dbConfg.Endpoint, dbConfg.Host, t.Second()))
		} else {
			Log.Error(fmt.Sprintf("instance fetdata timeout: %s %s daedline: %d", dbConfg.Endpoint, dbConfg.Host, t.Second()))
		}
		return
	}
}
