package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/yiqinguo/armyant/pkg/models"

	"github.com/influxdata/influxdb/client/v2"
)

var _ Storage = &influxdbStorage{}

const (
	precision string = "ns"
)

type influxdbStorage struct {
	client client.Client
	db     string
}

func NewInfluxdbStorage(addr, db, username, password string) (*influxdbStorage, error) {
	is := &influxdbStorage{db: db}
	client, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     addr,
		Username: username,
		Password: password,
	})
	if err != nil {
		return is, err
	}
	is.client = client

	return is, nil
}

func (i *influxdbStorage) AddJob(job *models.Job) error {
	return nil
}

func (i *influxdbStorage) GetJob(id string) (*models.Job, error) {
	return &models.Job{}, nil
}

func (i *influxdbStorage) GetJobResult(jobId, instanceId string) (*models.JobResultResponse, error) {
	result := &models.JobResultResponse{
		Result: map[string]float32{},
		Code:   1,
	}
	query := `select sum(spec_numberOfConnections) as spec_numberOfConnections, sum(spec_numberOfRequests) as spec_numberOfRequests, sum(result_bytesRead) as result_bytesRead, sum(result_bytesWritten) as result_bytesWritten, sum(result_timeTakenSeconds) as result_timeTakenSeconds, sum(result_others) as result_others, sum(result_req1xx) as result_req1xx, sum(result_req2xx) as result_req2xx, sum(result_req3xx) as result_req3xx, sum(result_req4xx) as result_req4xx, sum(result_req5xx) as result_req5xx, mean(result_latency_max) as result_latency_max, mean(result_latency_mean) as result_latency_mean, mean(result_latency_stddev) as result_latency_stddev, mean(result_latency_percentiles_50) as result_latency_percentiles_50, mean(result_latency_percentiles_75) as result_latency_percentiles_75, mean(result_latency_percentiles_90) as result_latency_percentiles_90, mean(result_latency_percentiles_95) as result_latency_percentiles_95, mean(result_latency_percentiles_99) as result_latency_percentiles_99, mean(rps_max) as rps_max, mean(rps_mean) as rps_mean, mean(rps_stddev) as rps_stddev from benchmark_job where jobId = '` + jobId + "'"
	if instanceId != "" {
		query += " AND instancdId = '" + instanceId + "'"
	}

	q := client.NewQuery(query, i.db, precision)
	resp, err := i.client.Query(q)
	if err != nil {
		return result, err
	}
	result.Code = 0
	results := resp.Results
	if len(results) == 0 {
		return result, nil
	}

	rows := results[0].Series
	if len(rows) == 0 {
		return result, nil
	}

	keys := rows[0].Columns
	values := rows[0].Values[0]
	for i, key := range keys {
		value, err := values[i].(json.Number).Float64()
		if err != nil {
			log.Println(err)
			result.Result[key] = float32(0)
		} else {
			result.Result[key] = float32(value)
		}
	}

	return result, nil

}

func (i *influxdbStorage) ListJobs() ([]models.Job, error) {
	return []models.Job{}, nil
}

func (i *influxdbStorage) DeleteJob(id string) error {
	return nil
}

func (i *influxdbStorage) UpdateJob(job *models.Job) error {
	return nil
}

func (i *influxdbStorage) AddResult(stats *models.Stats) error {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  i.db,
		Precision: precision,
	})
	if err != nil {
		return err
	}

	tags := map[string]string{
		"jobId":     stats.JobId,
		"historyId": stats.HistoryId,
	}

	result, err := statsConverToMap(stats)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(result)
	log.Println(string(body))

	pt, err := client.NewPoint("benchmark_job", tags, result, time.Now())
	if err != nil {
		return err
	}
	bp.AddPoint(pt)

	err = i.client.Write(bp)
	return err
}

func (i *influxdbStorage) TestWrite() error {
	str := `{"spec":{"numberOfConnections":125,"testType":"number-of-requests","numberOfRequests":1000,"method":"GET","url":"http://10.6.24.107","timeoutSeconds":2,"client":"fasthttp"},"result":{"bytesRead":4226000,"bytesWritten":59000,"timeTakenSeconds":0.39235067,"req2xx":1000,"latency":{"mean":44930.805,"stddev":110565.13,"max":389353},"rps":{"mean":2787.829,"stddev":4354.456,"max":12698.79,"percentiles":{"50":0,"75":0,"90":11964.633,"95":12395.243,"99":12698.79}}}}`

	var stats models.Stats
	err := json.Unmarshal([]byte(str), &stats)
	fmt.Println("============: 222", err)
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  i.db,
		Precision: precision,
	})
	if err != nil {
		return err
	}

	tags := map[string]string{
		"jobId":     stats.JobId,
		"historyId": stats.HistoryId,
	}

	var result map[string]interface{}
	//body, err := proto.Marshal(&stats)
	body, err := json.Marshal(stats)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &result)
	fmt.Println("============: 11", err, i.db)
	if err != nil {
		return err
	}

	/*result = map[string]interface{}{
		"idle":   10.1,
		"system": 53.3,
		"user":   46.6,
	}*/

	pt, err := client.NewPoint("benchmark_job", tags, result, time.Now())
	if err != nil {
		return err
	}
	bp.AddPoint(pt)

	/*bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "armyant",
		Precision: "ns",
	})

	tags := map[string]string{"cpu": "cpu-total"}
	fields := map[string]interface{}{
		"idle":   10.1,
		"system": 53.3,
		"user":   46.6,
	}
	pt, err := client.NewPoint("cpu_usage", tags, fields, time.Now())
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return err
	}
	bp.AddPoint(pt)

	// Write the batch
	return i.client.Write(bp)*/
	//return nil
	return i.client.Write(bp)
}

func statsConverToMap(stats *models.Stats) (map[string]interface{}, error) {
	var root map[string]interface{}
	rootBody, err := json.Marshal(stats)
	if err != nil {
		return root, err
	}

	err = json.Unmarshal(rootBody, &root)
	if err != nil {
		return root, err
	}
	delete(root, "jobId")
	delete(root, "historyId")

	delete(root, "spec")
	delete(root, "result")

	spec := stats.Spec
	result := stats.Result
	var latency *models.Latency
	var rps *models.Rps
	var percentiles map[string]float32
	if result.Latency != nil {
		latency = result.Latency
	}
	if result.Rps != nil {
		rps = result.Rps
	}

	specBody, err := json.Marshal(spec)
	if err != nil {
		return root, err
	}

	var specMap map[string]interface{}
	err = json.Unmarshal(specBody, &specMap)
	if err != nil {
		return root, err
	}

	for k, v := range specMap {
		root["spec_"+k] = v
	}

	resultBody, err := json.Marshal(result)
	if err != nil {
		return root, err
	}

	var resultMap map[string]interface{}
	err = json.Unmarshal(resultBody, &resultMap)
	if err != nil {
		return root, err
	}
	delete(resultMap, "latency")
	delete(resultMap, "rps")

	for k, v := range resultMap {
		root["result_"+k] = v
	}

	var latencyMap map[string]interface{}
	latencyBody, err := json.Marshal(latency)
	if err != nil {
		return root, err
	}
	err = json.Unmarshal(latencyBody, &latencyMap)
	if err != nil {
		return root, err
	}
	delete(latencyMap, "percentiles")
	for k, v := range latencyMap {
		root["result_latency_"+k] = v
	}

	if latency != nil {
		percentiles = latency.Percentiles
	}
	for k, v := range percentiles {
		root["result_latency_percentiles_"+k] = v
	}

	var rpsMap map[string]interface{}
	rpsBody, err := json.Marshal(rps)
	if err != nil {
		return root, err
	}

	err = json.Unmarshal(rpsBody, &rpsMap)
	if err != nil {
		return root, err
	}
	delete(rpsMap, "percentiles")
	for k, v := range rpsMap {
		root["rps_"+k] = v
	}

	if rps != nil {
		percentiles = rps.Percentiles
	}
	for k, v := range percentiles {
		root["rps_percentiles_"+k] = v
	}

	if v, ok := root["result_req1xx"].(float64); !ok {
		root["result_req1xx"] = float64(0)
	} else {
		root["result_req1xx"] = v
	}

	if v, ok := root["result_req2xx"].(float64); !ok {
		root["result_req2xx"] = float64(0)
	} else {
		root["result_req2xx"] = v
	}

	if v, ok := root["result_req3xx"].(float64); !ok {
		root["result_req3xx"] = float64(0)
	} else {
		root["result_req3xx"] = v
	}

	if v, ok := root["result_req4xx"].(float64); !ok {
		root["result_req4xx"] = float64(0)
	} else {
		root["result_req4xx"] = v
	}

	if v, ok := root["result_req5xx"].(float64); !ok {
		root["result_req5xx"] = float64(0)
	} else {
		root["result_req5xx"] = v
	}

	if v, ok := root["result_others"].(float64); !ok {
		root["result_others"] = float64(0)
	} else {
		root["result_others"] = v
	}

	return root, nil
}
