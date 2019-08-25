package baserabbitmq

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"register-go/infra"
)

var conn *amqp.Connection

func GetConn() *amqp.Connection {
	return conn
}

type RabbitMQStarter struct {
	infra.BaseStarter
}

func (r *RabbitMQStarter) Init(ctx infra.StarterContext) {
	var (
		err error
	)
	config := ctx.Yaml().RabbitMQConfig
	// amqp://username:password@192.168.56.130:5672/vhost
	if conn, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/%s", config.UserName, config.Password, config.IpAddr, config.Port, config.Vhost)); err != nil {
		logrus.Error(err)
	}
}
