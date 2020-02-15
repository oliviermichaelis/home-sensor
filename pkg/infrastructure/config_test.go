package infrastructure

import (
	"testing"
)

func TestRegister(t *testing.T) {
	c := make(config)

	// assert that fallback value is given back
	fallback := "test123"
	s, err := c.register("searchKey", fallback)
	if err != nil {
		t.Errorf("test: unexpected error: %v", err)
	}
	if s != fallback {
		t.Errorf("test: unexpected fallback value: %s, expected: %s", s, fallback)
	}

	// assert that error is raised if key wasn't found and no fallback value was given
	s, err = c.register("searchKey", "")
	if err == nil || len(s) > 0 {
		t.Errorf("test: key was found")
	}

	// assert that error is raised if key is of null value
	s, err = c.register("", fallback)
	if err == nil || len(s) > 0 {
		t.Errorf("test: key was found")
	}
}

func TestGet(t *testing.T) {
	c := make(config)
	c["testKey"] = "testValue"

	// input validation
	s, err := c.get("")
	if err == nil || len(s) > 0 {
		t.Errorf("test: returned value should be empty: %v", err)
	}

	// case if config doesn't contain key
	s, err = c.get("test")
	if err == nil || len(s) > 0 {
		t.Errorf("test: returned value should be empty: %v", err)
	}

	// valid case
	s, err = c.get("testKey")
	if err != nil || len(s) <= 0 {
		t.Errorf("test: expected nil, is: %v", err)
	}
}

func TestRegisterConfig(t *testing.T) {
	f := "fallback"
	s, err := RegisterConfig("test", f)
	if err != nil {
		t.Errorf("test: %v", err)
	}

	if s != f {
		t.Errorf("test: expected: %s but was: %s", f, s)
	}
}

func TestGetConfig(t *testing.T) {
	s, err := GetConfig("test2")
	if len(s) > 0 {
		t.Errorf("test: exptected null value but was: %d", len(s))
	}
	if err == nil {
		t.Errorf("test: expected error but was nil")
	}
}
