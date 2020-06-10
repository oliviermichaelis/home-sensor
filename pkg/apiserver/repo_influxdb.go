package apiserver

import (
	"context"
	"fmt"
	"github.com/influxdata/influxdb-client-go"
	"time"
)

type InfluxDbRepository struct {
	client influxdb2.Client
	debug  bool
	url    string
	token  string
}

func (i *InfluxDbRepository) getUrl() string {
	return i.url
}

// returns empty struct if database is reachable
func (i *InfluxDbRepository) attemptConnect() RepoError {

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	isReady, err := i.client.Ready(ctx)
	if isReady && err == nil {
		return RepoError{Url: i.url}
	}

	return RepoError{
		Err: fmt.Errorf("influxdb not ready: %v", err),
		Url: i.url,
	}
}

func (i *InfluxDbRepository) insert(measurement measurement) error {
	panic("implement me")
}

func (i *InfluxDbRepository) retrieveWindow(station string, start time.Time, end time.Time) (*[]measurement, error) {
	panic("implement me")
}

func NewInfluxDbRepository(debug bool, url string, token string) Repository {
	return &InfluxDbRepository{
		client: influxdb2.NewClient(url, token),
		debug:  debug,
		url:    url,
		token:  token,
	}
}
