package rabbitmq

import (
	"encoding/json"
	"testing"
)

const host = "amqp://user:bitnami@localhost:32672"

func TestPublish(t *testing.T) {
	q := New(host)
	defer q.Close()
	q.Bind("test")

	expect := "test"
	q.Publish("test", expect)

	c := q.Consume()
	msg := <-c
	var actual interface{}
	err := json.Unmarshal(msg.Body, &actual)
	if err != nil {
		t.Error(err)
	}
	if actual != expect {
		t.Errorf("expected %s, actual %s", expect, actual)
	}
	if msg.ReplyTo != q.Name {
		t.Error(msg)
	}
}

func TestSend(t *testing.T) {
	q := New(host)
	defer q.Close()

	expect := "test"
	expect2 := "test2"
	q.Send(q.Name, expect)
	q.Send(q.Name, expect2)

	c := q.Consume()
	msg := <-c
	var actual interface{}
	err := json.Unmarshal(msg.Body, &actual)
	if err != nil {
		t.Error(err)
	}
	if actual != expect {
		t.Errorf("expected %s, actual %s", expect, actual)
	}

	msg = <-c
	err = json.Unmarshal(msg.Body, &actual)
	if err != nil {
		t.Error(err)
	}
	if actual != expect2 {
		t.Errorf("expected %s, actual %s", expect2, actual)
	}
}
