package apiserver

import "time"

type Repository interface {
	attemptConnect() RepoError
	insert(measurement measurement) error
	retrieveWindow(station string, start time.Time, end time.Time) (*[]measurement, error)
}

type RepoError struct {
	Err error
	Url string
}

func (r *RepoError) Unwrap() error {
	return r.Err
}
