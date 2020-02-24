package infrastructure

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/influxdata/influxdb1-client/models"
	influxdb "github.com/influxdata/influxdb1-client/v2"
	"github.com/oliviermichaelis/home-sensor/pkg/domain"
	"html/template"
	"net"
	"time"
)

const EnvInfluxDatabase = "INFLUX_DATABASE"

type influxdbHandler struct {
	client influxdb.Client
	logger Logger
}

// TODO needs to be refactored!!
type Climate struct {
	Timestamp   time.Time `json:"timestamp"`
	Station     string    `json:"station"`
	Temperature float64   `json:"temperature"`
	Humidity    float64   `json:"humidity"`
	Pressure    float64   `json:"pressure"`
}

func NewInfluxdbHandler(host string, port string, username string, password string) (*influxdbHandler, error) {
	// Input validation
	if len(host) == 0 {
		return nil, errors.New("influxdb: host is empty")
	}
	if len(port) == 0 {
		return nil, errors.New("influxdb: port is empty")
	}
	if len(username) == 0 {
		return nil, errors.New("influxdb: username is empty")
	}
	if len(password) == 0 {
		return nil, errors.New("influxdb: password is empty")
	}

	// TODO implement TLS
	config := influxdb.HTTPConfig{
		Addr:               influxURL(host, port),
		Username:           username,
		Password:           password,
		UserAgent:          "",
		Timeout:            0,
		InsecureSkipVerify: false,
		TLSConfig:          nil,
		Proxy:              nil,
	}

	influx := influxdbHandler{
		logger: Logger{},
	}

	client, err := influxdb.NewHTTPClient(config)
	if err != nil {
		influx.logger.Fatal(err.Error())
	}
	influx.client = client
	return &influx, nil
}

func (handler *influxdbHandler) Insert(measurement domain.Measurement) {
	// transform measurement to influxdb point
	tags := map[string]string{
		"station": measurement.Station,
	}
	fields := map[string]interface{}{
		"temperature": measurement.Temperature,
		"humidity":    measurement.Humidity,
		"pressure":    measurement.Pressure,
	}
	timestamp, err := time.Parse("20060102150405", measurement.Timestamp)
	if err != nil {
		handler.logger.Log("influxdb: could not parse timestamp: ", err)
		return
	}

	database, err := GetConfig(EnvInfluxDatabase)
	if err != nil {
		handler.logger.Fatal("influxdb: could not retrieve database name")
	}

	var measurementBatchPoint = influxdb.BatchPointsConfig{
		Precision:        "",
		Database:         database,
		RetentionPolicy:  "",
		WriteConsistency: "",
	}

	points, err := influxdb.NewBatchPoints(measurementBatchPoint)
	if err != nil {
		message := "influxdb: could not create batchpoints: "
		handler.logger.Log(message, err)
		return
	}

	point, err := influxdb.NewPoint("sensor", tags, fields, timestamp)
	if err != nil {
		message := "influxdb: could not create batchpoint: "
		handler.logger.Log(message, err)
	}
	points.AddPoint(point)

	// Each goroutine will try to reconnect every 10s in case the data could'nt be persisted to influxdb
	for {
		err = nil
		if err = handler.client.Write(points); err != nil {
			handler.logger.Log("influxdb:", err)
			if _, ok := err.(net.Error); !ok {
				// In case the returned error is not a net.Error, the error isn't a networking malfunction/error and can
				// therefore not be resolved by continuous retries
				handler.logger.Log("influxdb: error is not net.Error and retries are not prohibited")
				return
			}
		}

		// In case no error has occured, break out of infinite retry loop
		if err == nil {
			break
		}
		time.Sleep(time.Second * 10)
	}

	message := fmt.Sprintf("influxdb: successfully inserted: %v", measurement)
	handler.logger.Log(message)
}

// TODO break up into smaller functions. Reuse code
func (handler *influxdbHandler) RetrieveLastWindow(station string, duration time.Duration) (*[]domain.Measurement, error) {
	// Subtract duration from current time in UTC to get last window
	t := time.Now().UTC().Add(-duration)

	database, err := GetConfig(EnvInfluxDatabase)
	if err != nil {
		return nil, err
	}

	// query template populated with values
	data := map[string]interface{}{
		"measurement": database,
		"station":     station,
		"timestamp":   t.Format(time.RFC3339),
	}

	tmpl, err := template.New("window").Parse("SELECT temperature, humidity, pressure FROM {{ .measurement }} WHERE station = '{{ .station }}' and time > '{{ .timestamp }}'")
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	if tmpl.Execute(&b, data) != nil {
		return nil, err
	}

	// execute query
	q := influxdb.Query{
		Command:  b.String(),
		Database: "sensor",
	}
	resp, err := handler.client.Query(q)
	if err != nil {
		return nil, err
	}
	if resp.Error() != nil {
		return nil, resp.Error()
	}

	// parse query
	// TODO use domain.Measurement
	var values []*Climate
	for _, r := range resp.Results {
		for _, s := range r.Series {
			v, err := parseSeries(&s, "")
			if err != nil {
				return nil, err
			}
			values = append(values, v...)
		}
	}

	v := climatesToMeasurements(values)
	return &v, nil
}

// Returns the connection URL for the influxdb client
// TODO add flag for https
func influxURL(host string, port string) string {
	return fmt.Sprintf("http://%s:%s", host, port)
}

func parseSeries(row *models.Row, station string) ([]*Climate, error) {
	// index the order of columns
	index := make(map[string]int)
	for i, v := range row.Columns {
		index[v] = i
	}

	// convert values to climate object
	var values []*Climate
	for _, v := range row.Values {
		//t := index["timestamp"]
		t, err := time.Parse(time.RFC3339, v[index["time"]].(string))
		if err != nil {
			return nil, err
		}

		temp, err := v[index["temperature"]].(json.Number).Float64()
		humid, err := v[index["humidity"]].(json.Number).Float64()
		press, err := v[index["pressure"]].(json.Number).Float64()

		values = append(values, &Climate{
			Timestamp:   t,
			Station:     station,
			Temperature: temp,
			Humidity:    humid,
			Pressure:    press,
		})
	}
	return values, nil
}

// TODO delete when refactor
func climateToMeasurement(climate *Climate) domain.Measurement {
	return domain.Measurement{
		Timestamp:   climate.Timestamp.Format("20060102150405"),
		Station:     climate.Station,
		Temperature: climate.Temperature,
		Humidity:    climate.Humidity,
		Pressure:    climate.Pressure,
	}
}

func climatesToMeasurements(c []*Climate) []domain.Measurement {
	var m []domain.Measurement
	for _, climate := range c {
		m = append(m, climateToMeasurement(climate))
	}
	return m
}

/*
func addMissingStation(c []*Climate, s string) []*Climate {
	for _, v := range c {
		if v.Station != "" {
			continue
		}
		v.Station = s
	}
	return c
}
 */
