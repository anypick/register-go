package rabbitmqutil

/**
定义RabbitMQ的生产者与消费者的基本操作，没有定义队列以及交换机声明相关操作的原因是：队列或者交换机声明还是希望用户在
使用之前定义好，统一维护，不要和业务代码掺杂在一起。
*/

import (
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"register-go/infra/base/rabbitmq"
)

type RabbitTemplate interface {
	// 发送消息
	Send(data []byte) error

	// 消息监听,
	// autoAck表示是否自动签收，true自动签收，无论消息是否被成功消费，消息都会被丢失
	// false, 如果消息消费失败，消息会被重新丢回队列，进行重新消费。
	// 一般情况下选择true, 提高系统性能
	MessageListener(consumer func([]byte) error, autoAck bool) error
}

type RabbitOperator struct {
	Exchange   string // 交换机名称
	Queue      string // 队列名称
	RoutingKey string // 路由键
}

func (r *RabbitOperator) Send(publishing amqp.Publishing) error {
	var (
		conn    = baserabbitmq.GetConn()
		channel *amqp.Channel
		err     error
	)
	if channel, err = conn.Channel(); err != nil {
		logrus.Error("get channel error", err)
		return err
	}
	defer channel.Close()
	if err = channel.Publish(r.Exchange, r.RoutingKey, false, false, publishing); err != nil {
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
