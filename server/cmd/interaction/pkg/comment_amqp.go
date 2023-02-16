package pkg

import (
	"context"
	"fmt"

	"github.com/CyanAsterisk/TikGok/server/cmd/interaction/dao"
	"github.com/CyanAsterisk/TikGok/server/cmd/interaction/model"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
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

// Publish publishes a comment model.
func (p *CommentPublisher) Publish(_ context.Context, comment *model.Comment) error {
	body, err := sonic.Marshal(comment)
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

// Subscribe subscribes and returns a channel with Favorite action request.
func (s *CommentSubscriber) Subscribe(c context.Context) (chan *model.Comment, func(), error) {
	msgCh, cleanUp, err := s.SubscribeRaw(c)
	if err != nil {
		return nil, cleanUp, err
	}

	ch := make(chan *model.Comment)
	go func() {
		for msg := range msgCh {
			var comment model.Comment
			err := sonic.Unmarshal(msg.Body, &comment)
			if err != nil {
				klog.Errorf("cannot unmarshal %s", err.Error())
			}
			ch <- &comment
		}
		close(ch)
	}()
	return ch, cleanUp, nil
}

func CommentSubscribeRoutine(subscriber *CommentSubscriber, dao *dao.Comment) error {
	commentCh, cleanUp, err := subscriber.Subscribe(context.Background())
	defer cleanUp()
	if err != nil {
		klog.Error("cannot subscribe", err)
		return err
	}
	for comment := range commentCh {
		if comment.ActionType == consts.ValidComment {
			if err = dao.CreateComment(comment); err != nil {
				klog.Error("create comment err")
			}
		} else if comment.ActionType == consts.InvalidComment {
			if err = dao.DeleteComment(comment.ID); err != nil {
				klog.Error("delete comment err")
			}
		} else {
			klog.Error("invalid comment action type")
		}
	}
	return nil
}

func declareExchange(ch *amqp.Channel, exchange string) error {
	return ch.ExchangeDeclare(exchange, "fanout", true, false, false, false, nil)
}
