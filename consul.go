package consul_client

import (
	"errors"
	"log"
	"time"

	"github.com/hashicorp/consul/api"
)

type Wrapper struct {
	Client      *api.Client
	ServiceName string
}

type CheckFunc func() (bool, error)

func NewClientWrapper(name string, config *api.Config) (*Wrapper, error) {

	consulClient, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	client := Wrapper{
		Client:      consulClient,
		ServiceName: name,
	}

	return &client, nil
}

func (w *Wrapper) RegisterService(ttl time.Duration) error {

	agent := w.agent()

	defaultCheckID := "default"
	defaultCheckName := "liveness"

	srvDef := api.AgentServiceRegistration{
		Name: w.ServiceName,
		Check: &api.AgentServiceCheck{
			TTL:     ttl.String(),
			CheckID: defaultCheckID,
			Name:    defaultCheckName,
		},
		Checks: api.AgentServiceChecks{},
	}

	err := agent.ServiceRegister(&srvDef)
	if err != nil {
		return err
	}

	defaultCheckFunc := func() (bool, error) {
		return true, nil
	}

	go func(w *Wrapper) {
		err := w.updateCheck(defaultCheckID, ttl, "service up", defaultCheckFunc)
		if err != nil {
			log.Println(err)
		}
	}(w)

	return nil
}

func (w *Wrapper) DeregisterService() error {

	if err := w.agent().ServiceDeregister(w.ServiceName); err != nil {
		return err
	}

	return nil
}

func (w *Wrapper) AddServiceCheck(name string, ttl time.Duration, notes string, checkFunc CheckFunc) error {

	check := &api.AgentCheckRegistration{
		Name:      name,
		ServiceID: w.ServiceName,
		AgentServiceCheck: api.AgentServiceCheck{
			Name:  name,
			Notes: notes,
			TTL:   ttl.String(),
		},
	}

	if err := w.agent().CheckRegister(check); err != nil {
		return err
	}

	go func(w *Wrapper) {
		err := w.updateCheck(name, ttl, "liveness check running", checkFunc)
		if err != nil {
			log.Println(err)
		}
	}(w)

	return nil
}

func (w *Wrapper) updateCheck(checkId string, ttl time.Duration, note string, checkFunc CheckFunc) error {
	tick := time.NewTicker(ttl)
	agent := w.agent()

	for range tick.C {
		ok, err := checkFunc()
		if err != nil {
			return err
		}

		if ok {
			err := agent.PassTTL(checkId, note)
			if err != nil {
				return err
			}
		} else {
			err := agent.FailTTL(checkId, note)
			if err != nil {
				return err
			}
		}
	}

	return errors.New("updateCheck exited")
}

func (w *Wrapper) agent() *api.Agent {
	return w.Client.Agent()
}

func (w *Wrapper) client() *api.Client {
	return w.Client
}

