package environment

import (
	"testing"
)

func TestReadUsernameWorkingPath(t *testing.T) {
	path := "../../test/rabbitmq"
	user, err := ReadUsername(path)
	if err != nil {
		t.Errorf("Username could not be read! %v", err)
	}

	if len(user) <= 0 {
		t.Errorf("The length of username is %d", len(user))
	}
}

func TestReadUsernameNoPath(t *testing.T) {
	path := ""
	user, err := ReadUsername(path)
	if err == nil {
		t.Error("Expected error not to be nil")
	}

	if len(user) > 0 {
		t.Errorf("The length of username is %d", len(user))
	}
}

func TestReadUsernameFaultyPath(t *testing.T) {
	path := "../../test\n/rabbitmq"
	user, err := ReadUsername(path)
	if err == nil {
		t.Error("Expected error not to be nil")
	}

	if len(user) > 0 {
		t.Errorf("The length of username is %d", len(user))
	}
}