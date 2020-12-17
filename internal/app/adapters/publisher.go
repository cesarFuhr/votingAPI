package adapters

import (
	"encoding/json"
	"fmt"

	"github.com/cesarFuhr/votingAPI/internal/app/domain/session"
	"github.com/cesarFuhr/votingAPI/internal/pkg/logger"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

// Publisher used to publish messages to the broker
type Publisher struct {
	c MQTT.Client
	l logger.Logger
}

// NewMQTTPublisher created a new mqtt publisher
func NewMQTTPublisher(connStr string, l logger.Logger) Publisher {
	p := Publisher{}
	p.l = l
	p.connect(connStr)
	return p
}

func (p *Publisher) connect(brokerURL string) {
	opts := MQTT.NewClientOptions().AddBroker(brokerURL)
	id, _ := uuid.NewUUID()
	clientName := fmt.Sprintf("Pub-%s", id)
	opts.SetClientID(clientName)

	p.c = MQTT.NewClient(opts)
	if token := p.c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	p.l.Info("Broker connected")
}

// PublishResult publishes a result message
func (p *Publisher) PublishResult(r session.Result) error {
	m, err := json.Marshal(&r)
	if err != nil {
		return err
	}

	p.l.Info("Published -> ", r.ID, " ", string(m))
	token := p.c.Publish("result", 0, false, m)
	token.Wait()
	return nil
}
