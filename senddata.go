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
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/open-falcon/mymon/common"
)

// SendData Post the json of all result to falcon-agent
func SendData(pdata []*MetaData, config *common.Config) ([]byte, error) {

	data := filterIgnoreData(pdata)
	js, err := json.Marshal(data)

	if err != nil {
		Log.Debug("parse json data error: %+v", err)
		return nil, err
	}
	Log.Info("Send to %s, size: %d", config.Base.FalconClient, len(data))
	for _, m := range data {
		Log.Debug("%v", m)
	}

	res, err := http.Post(config.Base.FalconClient, "Content-Type: application/json", bytes.NewBuffer(js))
	if err != nil {
		Log.Debug("send data to falcon-agent error: %+v", err)
		return nil, err
	}

	defer func() { _ = res.Body.Close() }()
	return ioutil.ReadAll(res.Body)
}

func parseLine(line string) map[string]string {
	var parseRes map[string]string
	// remove space and \n
	line = strings.Replace(strings.TrimSpace(line), "\n", "", -1)

	// match metric, tag and value
	reMetricTagValue, _ := regexp.Compile("^([0-9A-Za-z_,]+)" + TagSplitChar + "?([0-9A-Za-z_,=]*)" + ValueSplitChar + "?([0-9A-Za-z_]*)$")
	matchMetricTagValue := reMetricTagValue.FindSubmatch([]byte(line))
	if len(matchMetricTagValue) > 0 {
		parseRes = map[string]string{
			"metric": string(matchMetricTagValue[1]),
			"tag":    string(matchMetricTagValue[2]),
			"value":  string(matchMetricTagValue[3]),
		}
	} else {
		Log.Info("Error format of ignorefile: %s", line)
	}
	return parseRes
}

func filterIgnoreData(data []*MetaData) []*MetaData {
	pdata := []*MetaData{}
	for _, metric := range data {
		if IngoreMetric(metric.Metric) {
			metric.Metric = AliasType(metric.Metric)
			pdata = append(pdata, metric)
		}
	}
	return pdata
}

func tagSame(tag1, tag2 string) bool {
	x, y := strings.Split(tag1, ","), strings.Split(tag2, ",")
	sort.Strings(x)
	sort.Strings(y)
	return reflect.DeepEqual(x, y)
}

// Snapshot make a record of note, some metric should be noted before sending
func Snapshot(conf *common.Config, note string, fileNameDay string, fileNameOldDay string) error {
	if conf.Base.SnapshotDay < 0 {
		// Just remind but do not stop
		Log.Info("snapshot_day setted error!")
	}
	f, err := os.OpenFile(fileNameDay, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		Log.Debug("open snapshot file %s error: %+v", fileNameDay, err)
		return err
	}
	defer f.Close()
	_, err = f.WriteString(note)
	if err != nil {
		Log.Debug("write info to snapshot file error: %+v", err)
		return err
	}
	e := os.Remove(fileNameOldDay)
	if e != nil {
		// Just remind but do not stop
		Log.Info("Error remove %s, %s", fileNameOldDay, e.Error())
	}
	return err
}
