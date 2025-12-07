package messaging

import (
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
)

// NATSClient wraps NATS connection
type NATSClient struct {
	conn *nats.Conn
}

// NewNATSClient creates a new NATS client
func NewNATSClient(url string) (*NATSClient, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	return &NATSClient{conn: conn}, nil
}

// Publish publishes a message to a subject
func (c *NATSClient) Publish(subject string, data interface{}) error {
	if c == nil || c.conn == nil {
		return nil
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return c.conn.Publish(subject, jsonData)
}

// Subscribe subscribes to a subject
func (c *NATSClient) Subscribe(subject string, handler func([]byte)) (*nats.Subscription, error) {
	if c == nil || c.conn == nil {
		return nil, fmt.Errorf("NATS client not connected")
	}

	return c.conn.Subscribe(subject, func(msg *nats.Msg) {
		handler(msg.Data)
	})
}

// QueueSubscribe subscribes to a subject with queue group
func (c *NATSClient) QueueSubscribe(subject, queue string, handler func([]byte)) (*nats.Subscription, error) {
	if c == nil || c.conn == nil {
		return nil, fmt.Errorf("NATS client not connected")
	}

	return c.conn.QueueSubscribe(subject, queue, func(msg *nats.Msg) {
		handler(msg.Data)
	})
}

// Close closes the NATS connection
func (c *NATSClient) Close() {
	if c != nil && c.conn != nil {
		c.conn.Close()
	}
}

// Event subjects
const (
	SubjectPatientCreated    = "patient.created"
	SubjectPatientUpdated    = "patient.updated"
	SubjectEncounterCreated  = "encounter.created"
	SubjectEncounterUpdated  = "encounter.updated"
	SubjectOrderCreated      = "order.created"
	SubjectOrderCompleted    = "order.completed"
	SubjectResultsAvailable  = "results.available"
	SubjectAppointmentBooked = "appointment.booked"
	SubjectAppointmentCancelled = "appointment.cancelled"
	SubjectNotificationSend  = "notification.send"
	SubjectERPSync           = "erp.sync"
)
