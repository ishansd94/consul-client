package consul_sds

import (
	"testing"
	"time"

	"github.com/hashicorp/consul/api"
)


func TestNewClientWrapper(t *testing.T) {
	_, _ = NewTestClient(t, "test")
}


func TestWrapper_RegisterService(t *testing.T) {
	client, _  := NewTestClient(t, "test1")

	err := client.RegisterService(time.Second*5)
	if err != nil {
		t.Errorf("error registering service - %s", err)
		t.Fail()
	}
}


func TestWrapper_AddServiceCheck(t *testing.T) {
	client, _  := NewTestClient(t, "test2")

	err := client.RegisterService(time.Second*5)
	if err != nil {
		t.Errorf("error registering service - %s", err)
		t.Fail()
	}

	err = client.AddServiceCheck("test2", 5*time.Second, "", func() (b bool, err error) {
		return true, nil
	})
	if err != nil {
		t.Errorf("error adding service check - %s", err)
		t.Fail()
	}
}


func TestWrapper_DeregisterService(t *testing.T) {
	client, _  := NewTestClient(t , "test3")

	err := client.RegisterService(time.Second*5)
	if err != nil {
		t.Errorf("error registering service - %s", err)
		t.Fail()
	}

	err = client.DeregisterService()
	if err != nil {
		t.Errorf("error de-registering service - %s", err)
		t.Fail()
	}
}


func NewTestClient(t *testing.T, name string) (*Wrapper,error) {
	cl , err := NewClientWrapper(name, &api.Config{
		Address: "localhost:8500",
	})

	if err != nil {
		t.Errorf("error creating test client - %s" , err)
		t.Fail()

		return nil, err
	}

	return cl, nil
}

