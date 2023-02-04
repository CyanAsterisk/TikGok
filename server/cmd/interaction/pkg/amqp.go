package pkg

import (
	"context"
	"fmt"

	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/interaction"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/streadway/amqp"
)

// CommentPublisher implements an amqp publisher.
type CommentPublisher struct {
	ch       *amqp.Channel
	exchange string
}

func NewCommentPublisher(conn *amqp.Connection, exchange string) (*CommentPublisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("cannot allocate channel: %v", err)
	}

	if err = declareExchange(ch, exchange); err != nil {
		return nil, fmt.Errorf("cannot declare exchange: %v", err)
	}
	return &CommentPublisher{
		ch:       ch,
		exchange: exchange,
	}, nil
}

// Publish publishes a comment action request.
func (p *CommentPublisher) Publish(_ context.Context, car *interaction.DouyinCommentActionRequest) error {
	body, err := sonic.Marshal(car)
	if err != nil {
		return fmt.Errorf("cannot marshal: %v", err)
	}

	return p.ch.Publish(p.exchange, "", false, false, amqp.Publishing{
		Body: body,
	})
}

// CommentSubscriber implements an amqp subscriber.
type CommentSubscriber struct {
	conn     *amqp.Connection
	exchange string
}

// NewCommentSubscriber creates an amqp subscriber.
func NewCommentSubscriber(conn *amqp.Connection, exchange string) (*CommentSubscriber, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("cannot allocate channel: %v", err)
	}
	defer ch.Close()

	if err = declareExchange(ch, exchange); err != nil {
		return nil, fmt.Errorf("cannot declare exchange: %v", err)
	}

	return &CommentSubscriber{
		conn:     conn,
		exchange: exchange,
	}, nil
}

// SubscribeRaw subscribes and returns a channel with raw amqp delivery.
func (s *CommentSubscriber) SubscribeRaw(_ context.Context) (<-chan amqp.Delivery, func(), error) {
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

// FavoritePublisher implements an amqp publisher.
type FavoritePublisher struct {
	ch       *amqp.Channel
	exchange string
}

func NewFavoritePublisher(conn *amqp.Connection, exchange string) (*FavoritePublisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("cannot allocate channel: %v", err)
	}

	if err = declareExchange(ch, exchange); err != nil {
		return nil, fmt.Errorf("cannot declare exchange: %v", err)
	}
	return &FavoritePublisher{
		ch:       ch,
		exchange: exchange,
	}, nil
}

// Publish publishes a favorite action request.
func (p *FavoritePublisher) Publish(_ context.Context, car *interaction.DouyinFavoriteActionRequest) error {
	body, err := sonic.Marshal(car)
	if err != nil {
		return fmt.Errorf("cannot marshal: %v", err)
	}

	return p.ch.Publish(p.exchange, "", false, false, amqp.Publishing{
		Body: body,
	})
}

// FavoriteSubscriber implements an amqp subscriber.
type FavoriteSubscriber struct {
	conn     *amqp.Connection
	exchange string
}

// NewFavoriteSubscriber creates an amqp subscriber.
func NewFavoriteSubscriber(conn *amqp.Connection, exchange string) (*FavoriteSubscriber, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("cannot allocate channel: %v", err)
	}
	defer ch.Close()

	if err = declareExchange(ch, exchange); err != nil {
		return nil, fmt.Errorf("cannot declare exchange: %v", err)
	}

	return &FavoriteSubscriber{
		conn:     conn,
		exchange: exchange,
	}, nil
}

// SubscribeRaw subscribes and returns a channel with raw amqp delivery.
func (s *FavoriteSubscriber) SubscribeRaw(_ context.Context) (<-chan amqp.Delivery, func(), error) {
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

// Subscribe subscribes and returns a channel with Favorite action request.
func (s *FavoriteSubscriber) Subscribe(c context.Context) (chan *interaction.DouyinFavoriteActionRequest, func(), error) {
	msgCh, cleanUp, err := s.SubscribeRaw(c)
	if err != nil {
		return nil, cleanUp, err
	}

	reqCh := make(chan *interaction.DouyinFavoriteActionRequest)
	go func() {
		for msg := range msgCh {
			var req interaction.DouyinFavoriteActionRequest
			err := sonic.Unmarshal(msg.Body, &req)
			if err != nil {
				klog.Errorf("cannot unmarshal %s", err.Error())
			}
			reqCh <- &req
		}
		close(reqCh)
	}()
	return reqCh, cleanUp, nil
}

func declareExchange(ch *amqp.Channel, exchange string) error {
	return ch.ExchangeDeclare(exchange, "fanout", true, false, false, false, nil)
}
