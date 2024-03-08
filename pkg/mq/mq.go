// Package mq package mq
package mq

import (
	"fmt"
	"io"
	"log/slog"
	"sync"
	"time"

	"tmt/pkg/utils"

	"github.com/google/uuid"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/mochi-mqtt/server/v2/packets"
)

const (
	_defaultPort = "18883"
)

var (
	srv  *MQSrv
	once sync.Once
)

type MQSrv struct {
	server             *mqtt.Server
	subscriptionTopic  map[int]string
	subscriptionID     int
	subscriptionIDLock sync.Mutex
}

func Serve() error {
	errChan := make(chan error, 2)
	newServer := &MQSrv{
		subscriptionTopic: make(map[int]string),
		server: mqtt.New(&mqtt.Options{
			InlineClient:             true,
			Logger:                   slog.New(slog.NewTextHandler(io.Discard, nil)),
			ClientNetWriteBufferSize: 4096,
			ClientNetReadBufferSize:  4096,
			SysTopicResendInterval:   10,
		}),
	}
	once.Do(func() {
		_ = newServer.server.AddHook(&auth.AllowHook{}, nil)
		tcp := listeners.NewTCP(uuid.NewString(), fmt.Sprintf("127.0.0.1:%s", _defaultPort), nil)
		err := newServer.server.AddListener(tcp)
		if err != nil {
			errChan <- err
			return
		}
		go func() {
			err = newServer.server.Serve()
			if err != nil {
				errChan <- err
			}
		}()
	})
	for {
		select {
		case err := <-errChan:
			return err
		case <-time.After(1 * time.Second):
			if utils.GetPortIsUsed("127.0.0.1", _defaultPort) {
				srv = newServer
				return nil
			}
		}
	}
}

func Get() *MQSrv {
	if srv == nil {
		panic("MQSrv is not initialized")
	}
	return srv
}

func (m *MQSrv) getSubscriptionID(topic string) int {
	m.subscriptionIDLock.Lock()
	defer m.subscriptionIDLock.Unlock()
	m.subscriptionID++
	m.subscriptionTopic[m.subscriptionID] = topic
	return m.subscriptionID
}

func (m *MQSrv) getTopic(id int) string {
	if id < 0 {
		return ""
	}

	m.subscriptionIDLock.Lock()
	defer m.subscriptionIDLock.Unlock()
	topic, ok := m.subscriptionTopic[id]
	if !ok {
		return ""
	}
	delete(m.subscriptionTopic, id)
	return topic
}

func (m *MQSrv) Subscribe(topic string, callbackFn func(cl *mqtt.Client, sub packets.Subscription, pk packets.Packet)) int {
	id := m.getSubscriptionID(topic)
	err := m.server.Subscribe(topic, id, callbackFn)
	if err != nil {
		return -1
	}
	return id
}

func (m *MQSrv) Unsubscribe(id int) {
	topic := m.getTopic(id)
	if topic == "" {
		return
	}
	_ = m.server.Unsubscribe(topic, id)
}

func (m *MQSrv) Publish(topic string, payload []byte) error {
	return m.server.Publish("direct/publish", []byte("packet scheduled message"), false, 0)
}
