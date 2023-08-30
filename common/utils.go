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
	"os"
	"strconv"
	"strings"
	"time"
)

// GetFileNameDayAndOldDay get name of new file and old file
func GetFileNameDayAndOldDay(conf *DatabaseConf, prefix string) (string, string) {
	t := time.Now()
	ld := t.Add(-time.Duration(7) * 24 * time.Hour).Day()
	fileName := fmt.Sprintf(
		"%s/%s_%s:%d", "./", prefix,
		conf.Endpoint, conf.Port)
	fileNameDay := fmt.Sprintf("%s_%d", fileName, t.Day())
	fileNameOldDay := fmt.Sprintf("%s_%d", fileName, ld)
	return fileNameDay, fileNameOldDay
}

// GetLastNum get file name number suffix
func GetLastNum(str string, split string) int {
	parts := strings.Split(str, split)
	if len(parts) < 2 {
		return -1
	}
	ans, err := strconv.ParseInt(parts[1], 10, 60)
	if err != nil {
		return -2
	}
	return int(ans)
}

// Hostname get current hostname
func Hostname(conf *DatabaseConf) string {
	host, err := os.Hostname()
	if err != nil {
		host = conf.Endpoint
	}
	return host
}

//CompatibleLog for making log_file and log_dir compatible
func CompatibleLog(conf *Config) string {
	logDir := conf.Base.LogDir
	logFile := conf.Base.LogFile
	initLogFile := "myMon.log"
	if logFile != "" {
		return logFile
	}
	if logDir != "" {
		return fmt.Sprintf("%s/%s", logDir, initLogFile)
	}
	return initLogFile
}
