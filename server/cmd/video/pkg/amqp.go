package pkg

import (
	"context"
	"fmt"

	"github.com/CyanAsterisk/TikGok/server/cmd/video/dao"
	"github.com/CyanAsterisk/TikGok/server/cmd/video/model"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/streadway/amqp"
)

// Publisher implements an amqp publisher.
type Publisher struct {
	ch       *amqp.Channel
	exchange string
}

func NewPublisher(conn *amqp.Connection, exchange string) (*Publisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("cannot allocate channel: %v", err)
	}

	if err = declareExchange(ch, exchange); err != nil {
		return nil, fmt.Errorf("cannot declare exchange: %v", err)
	}
	return &Publisher{
		ch:       ch,
		exchange: exchange,
	}, nil
}

// Publish publishes a video record.
func (p *Publisher) Publish(_ context.Context, video *model.Video) error {
	body, err := sonic.Marshal(video)
	if err != nil {
		return fmt.Errorf("cannot marshal: %v", err)
	}

	return p.ch.Publish(p.exchange, "", false, false, amqp.Publishing{
		Body: body,
	})
}

// Subscriber implements an amqp subscriber.
type Subscriber struct {
	conn     *amqp.Connection
	exchange string
}

// NewSubscriber creates an amqp subscriber.
func NewSubscriber(conn *amqp.Connection, exchange string) (*Subscriber, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("cannot allocate channel: %v", err)
	}
	defer ch.Close()

	if err = declareExchange(ch, exchange); err != nil {
		return nil, fmt.Errorf("cannot declare exchange: %v", err)
	}

	return &Subscriber{
		conn:     conn,
		exchange: exchange,
	}, nil
}

// SubscribeRaw subscribes and returns a channel with raw amqp delivery.
func (s *Subscriber) SubscribeRaw(_ context.Context) (<-chan amqp.Delivery, func(), error) {
	ch, err := s.conn.Channel()
	if err != nil {
		return nil, func() {}, fmt.Errorf("cannot allocate channel: %v", err)
	}
	closeCh := func() {
		err := ch.Close()
		if err != nil {
			klog.Errorf("cannot close channel %s", err.Error())
		}
	}

	q, err := ch.QueueDeclare("", false, true, false, false, nil)
	if err != nil {
		return nil, closeCh, fmt.Errorf("cannot declare queue: %v", err)
	}

	cleanUp := func() {
		_, err := ch.QueueDelete(q.Name, false, false, false)
		if err != nil {
			klog.Errorf("cannot delete queue %s : %s", q.Name, err.Error())
		}
		closeCh()
	}

	err = ch.QueueBind(q.Name, "", s.exchange, false, nil)
	if err != nil {
		return nil, cleanUp, fmt.Errorf("cannot bind: %v", err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return nil, cleanUp, fmt.Errorf("cannot consume queue: %v", err)
	}
	return msgs, cleanUp, nil
}

// Subscribe subscribes and returns a channel with video publish request.
func (s *Subscriber) Subscribe(c context.Context) (videoChan chan *model.Video, cleanUpFunc func(), err error) {
	msgCh, cleanUp, err := s.SubscribeRaw(c)
	if err != nil {
		return nil, cleanUp, err
	}

	ch := make(chan *model.Video)
	go func() {
		for msg := range msgCh {
			var v model.Video
			err := sonic.Unmarshal(msg.Body, &v)
			if err != nil {
				klog.Errorf("cannot unmarshal %s", err.Error())
			}
			ch <- &v
		}
		close(ch)
	}()
	return ch, cleanUp, nil
}

func declareExchange(ch *amqp.Channel, exchange string) error {
	return ch.ExchangeDeclare(exchange, "fanout", true, false, false, false, nil)
}

func SubscribeRoutine(subscriber *Subscriber) error {
	videoChan, cleanUp, err := subscriber.Subscribe(context.Background())
	defer cleanUp()
	if err != nil {
		klog.Error("cannot subscribe", err)
	}
	for video := range videoChan {
		err = dao.CreateVideo(video)
		if err != nil {
			klog.Error("create video err", err)
		}
	}
	return nil
}
