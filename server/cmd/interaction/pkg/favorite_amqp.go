package pkg

import (
	"context"
	"fmt"
	"github.com/CyanAsterisk/TikGok/server/cmd/interaction/dao"
	"github.com/CyanAsterisk/TikGok/server/cmd/interaction/model"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/streadway/amqp"
)

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

// Publish publishes a favorite model.
func (p *FavoritePublisher) Publish(_ context.Context, fav *model.Favorite) error {
	body, err := sonic.Marshal(fav)
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

// Subscribe subscribes and returns a channel with Favorite model.
func (s *FavoriteSubscriber) Subscribe(c context.Context) (chan *model.Favorite, func(), error) {
	msgCh, cleanUp, err := s.SubscribeRaw(c)
	if err != nil {
		return nil, cleanUp, err
	}

	ch := make(chan *model.Favorite)
	go func() {
		for msg := range msgCh {
			var fav model.Favorite
			err := sonic.Unmarshal(msg.Body, &fav)
			if err != nil {
				klog.Errorf("cannot unmarshal %s", err.Error())
			}
			ch <- &fav
		}
		close(ch)
	}()
	return ch, cleanUp, nil
}

func FavoriteSubscribeRoutine(subscriber *FavoriteSubscriber) error {
	favCh, cleanUp, err := subscriber.Subscribe(context.Background())
	defer cleanUp()
	if err != nil {
		klog.Error("cannot subscribe", err)
		return err
	}
	for fav := range favCh {
		fr, err := dao.GetFavoriteInfo(fav.UserId, fav.VideoId)
		if err == nil && fr == nil {
			err = dao.CreateFavorite(fav)
			if err != nil {
				klog.Error("favorite action error", err)
			}
			continue
		}
		if err != nil {
			klog.Error("get favorite info error", err)
		}
		err = dao.UpdateFavorite(fav.UserId, fav.VideoId, fav.ActionType)
		if err != nil {
			klog.Error("update favorite error", err)
		}
	}
	return nil
}
