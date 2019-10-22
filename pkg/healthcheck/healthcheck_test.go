package healthcheck

import (
	"testing"
)

var serviceUpConsumer = &health{
	Status:  statusFail,
	Details: []service{{
		Name:   ServiceRabbitMQ,
		Status: true,
	}, {
		Name:   ServiceInfluxDB,
		Status: true,
	}},
}

func (h health) deepCopy() health {
	return h
}

func TestHealth_CheckStatusFailedToPassProducer(t *testing.T) {
	isConsumer = false
	h := health{
		Status:  statusFail,
		Details: []service{{
			Name:   ServiceRabbitMQ,
			Status: true,
		}},
	}

	h.checkStatus()
	if h.Status != statusPass {
		t.Error("Status is not \"pass\" after running checkStatus()!")
	}
}

func TestHealth_CheckStatusFailedToPassConsumer(t *testing.T) {
	isConsumer = true
	h := serviceUpConsumer.deepCopy()

	h.checkStatus()
	if h.Status != statusPass {
		t.Error("Status is not \"pass\" after running checkStatus()!")
	}
}

func TestHealth_CheckStatusFailedToFailedConsumer(t *testing.T) {
	isConsumer = true
	h := health{
		Status:  statusFail,
		Details: []service{{
			Name:   ServiceRabbitMQ,
			Status: false,
		}, {
			Name:   ServiceInfluxDB,
			Status: true,
		}},
	}

	h.checkStatus()
	if h.Status != statusFail {
		t.Errorf("Status is not \"%s\" after running checkStatus()!", statusFail)
	}
}

func TestHealth_SetStatus(t *testing.T) {
	isConsumer = true
	h := serviceUpConsumer.deepCopy()
	h.checkStatus()

	h.SetStatus(ServiceInfluxDB, false)
	if h.Status != statusFail {
		t.Errorf("Status is not \"%s\" after running checkStatus()!", statusFail)
	}
}
