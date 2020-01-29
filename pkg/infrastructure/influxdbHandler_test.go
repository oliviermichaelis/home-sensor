package infrastructure

import (
	"testing"
)

func TestNewInfluxdbHandler(t *testing.T) {
	// test input validation
	if c, err := NewInfluxdbHandler("", "8086", "user", "pass"); c != nil || err == nil {
		t.Errorf("test: input not validated properly")
	}

	if c, err := NewInfluxdbHandler("localhost", "", "user", "pass"); c != nil || err == nil {
		t.Errorf("test: input not validated properly")
	}

	if c, err := NewInfluxdbHandler("localhost", "8086", "", "pass"); c != nil || err == nil {
		t.Errorf("test: input not validated properly")
	}

	if c, err := NewInfluxdbHandler("localhost", "8086", "user", ""); c != nil || err == nil {
		t.Errorf("test: input not validated properly")
	}

	if c, err := NewInfluxdbHandler("localhost", "8086", "user", "pass"); c == nil || err != nil {
		t.Errorf("test: %v", err)
	}
}

//func TestInfluxdbHandler_Insert(t *testing.T) {
//	measurement := domain.Measurement{}
//	measurement.PopulateRandomValues()
//
//	c, err := NewInfluxdbHandler("localhost", "8086", "user", "pass")
//	if c == nil || err != nil {
//		t.Errorf("test: %v", err)
//	}
//
//	c.Insert(measurement)
//
//}
