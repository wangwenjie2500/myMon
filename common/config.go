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
	"github.com/BurntSushi/toml"
	"os"
)

// BaseConf config about dir, log, etc.
type BaseConf struct {
	BaseDir      string `toml:"basedir"`
	SnapshotDir  string `toml:"snapshot_dir"`
	SnapshotDay  int    `toml:"snapshot_day"`
	LogDir       string `toml:"log_dir"`
	LogFile      string `toml:"log_file"`
	LogLevel     int    `toml:"log_level"`
	FalconClient string `toml:"falcon_client"`
	IgnoreFile   string `toml:"ignore_file"`
}

// DatabaseConf config about database
type DatabaseConf struct {
	User     string `toml:"user"`
	Password string `toml:"password"`
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Endpoint string `toml:"endpoint"`
}

// Config for initializing. This can be loaded from TOML file with -c
type Config struct {
	Base  BaseConf       `toml:"default"`
	Mysql []DatabaseConf `toml:"mysql"`
}

// NewConfig the constructor of config
func NewConfig(file string) (*Config, error) {
	conf, err := readConf(file)
	return &conf, err
}

func readConf(file string) (conf Config, err error) {
	_, err = os.Stat(file)
	if err != nil {
		file = fmt.Sprint("etc/", file)
		_, err = os.Stat(file)
		if err != nil {
			panic(err)
		}
	}
	tomlConfig := Config{}
	_, err = toml.DecodeFile(file, &tomlConfig)
	if err != nil {
		panic(err)
	}

	return tomlConfig, nil
	//cfg, err := ini.Load(file)
	//
	//if err != nil {
	//	panic(err)
	//}
	//snapshotDay, err := cfg.Section("default").Key("snapshot_day").Int()
	//if err != nil {
	//	fmt.Println("No Snapshot!")
	//	snapshotDay = -1
	//}
	//logLevel, err := cfg.Section("default").Key("log_level").Int()
	//if err != nil {
	//	fmt.Println("Log level default: 7!")
	//	logLevel = 7
	//}
	//
	////host := strings.Split(cfg.Section("mysql").Key("host").String(), ",")
	////
	//snapshotDir := cfg.Section("default").Key("snapshot_dir").String()
	//if snapshotDir == "" {
	//	fmt.Println("SnapshotDir default current dir ")
	//	snapshotDir = "."
	//}
	////
	////port, err := cfg.Section("mysql").Key("port").Int()
	////if err != nil {
	////	fmt.Println("Port: default 3306!")
	////	port = 3306
	////	err = nil
	////}
	//dbs := DatabaseList{}
	//err = cfg.Section("mysql").MapTo(&dbs)
	//if err != nil {
	//	fmt.Println("Parse dblist error, %v", err)
	//}
	//
	//conf = Config{
	//	BaseConf{
	//		BaseDir:      cfg.Section("default").Key("basedir").String(),
	//		SnapshotDir:  snapshotDir,
	//		SnapshotDay:  snapshotDay,
	//		LogDir:       cfg.Section("default").Key("log_dir").String(),
	//		LogFile:      cfg.Section("default").Key("log_file").String(),
	//		Endpoint:     cfg.Section("default").Key("endpoint").String(),
	//		LogLevel:     logLevel,
	//		FalconClient: cfg.Section("default").Key("falcon_client").String(),
	//		IgnoreFile:   cfg.Section("default").Key("ignore_file").String(),
	//	},
	//	dbs,
	//}
}
