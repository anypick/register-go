package rabbitmqutil

import (
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"register-go/infra/base/rabbitmq"
)

type RabbitTemplate interface {
	// 发送消息, 交换机类型direct
	Send(data []byte) error

	// 消息监听
	MessageListener(consumer func([]byte) error, autoAck bool) error
}

type RabbitOperator struct {
	Exchange   string // 交换机名称
	Queue      string // 队列名称
	RoutingKey string // 路由键
}

func (r *RabbitOperator) Send(data []byte) error {
	var (
		conn    = baserabbitmq.GetConn()
		channel *amqp.Channel
		err     error
	)
	if channel, err = conn.Channel(); err != nil {
		logrus.Error("channel get error", err)
		return err
	}
	defer channel.Close()
	if err = channel.Publish(r.Exchange, r.RoutingKey, false, false, amqp.Publishing{Body: data}); err != nil {
		logrus.Error("publish message error,", err)
		return err
	}
	return err
}

func (r *RabbitOperator) MessageListener(consumer func([]byte) error, autoAck bool) error {
	var (
		channel    *amqp.Channel
		conn       = baserabbitmq.GetConn()
		deliveries <-chan amqp.Delivery
		err        error
	)
	if channel, err = conn.Channel(); err != nil {
		logrus.Error("channel get error", err)
		return err
	}
	if deliveries, err = channel.Consume(r.Queue, "", autoAck, false, false, false, nil); err != nil {
		logrus.Error(err)
		return err
	}
	go func(d <-chan amqp.Delivery) {
		for d := range d {
			e := consumer(d.Body)
			if !autoAck {
				if e == nil {
					e = d.Ack(true)
					if e != nil {
						logrus.Error("ack error, ", e)
					}
				} else {
					logrus.Error(e)
					e = d.Nack(true, false)
					if e != nil {
						logrus.Error("nack error, ", e)
					}
				}
			}

		}
	}(deliveries)
	return nil
}
