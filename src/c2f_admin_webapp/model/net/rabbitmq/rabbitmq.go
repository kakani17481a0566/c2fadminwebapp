package rabbitmq

import (
	"errors"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
	done    chan error
}

type RabbitMQ struct {
	Ip             string
	Port           int
	Ctag           string
	Username       string
	Password       string
	Vhost          string
	Exchange       string
	ExchangeType   string
	Queue          string
	Key            string
	Durable        bool
	NoWait         bool
	AutoAck        bool
	ConnectionPtr  *consumer
	QueuePtr       amqp.Queue
	DeliveriesChan <-chan amqp.Delivery
}

func (obj *RabbitMQ) Connect() (err error) {
	obj.ConnectionPtr = &consumer{
		conn:    nil,
		channel: nil,
		tag:     obj.Ctag,
		done:    make(chan error),
	}

	amqpURI := fmt.Sprint("amqp://", obj.Username, ":", obj.Password, "@", obj.Ip,
		":", obj.Port, "/", obj.Vhost)

	obj.ConnectionPtr.conn, err = amqp.Dial(amqpURI)
	if err != nil {
		log.Println("Error", "Error", err)
		err = errors.New(amqpURI + " : dial fail")
		return
	}

	obj.ConnectionPtr.channel, err = obj.ConnectionPtr.conn.Channel()
	if err != nil {
		log.Println("Error", "Error", err)
		err = errors.New(amqpURI + " : channel fail")
		return
	}
	return
}

func (obj *RabbitMQ) Send(message []byte, header map[string]interface{}) (err error) {
	exchange := header["DestinationExchange"]
	key := header["DestinationKey"]
	if err = obj.ConnectionPtr.channel.Publish(
		exchange.(string), // publish to an exchange
		key.(string),      // routing to 0 or more queues
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			Headers:         header,
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            message,
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,
		},
	); err != nil {
		log.Println(err)
		err = errors.New(exchange.(string) + " : publish fail")
		return
	}
	return
}

func (obj *RabbitMQ) DeclareExchange() (err error) {
	err = obj.ConnectionPtr.channel.ExchangeDeclare(
		obj.Exchange,     // name of the exchange
		obj.ExchangeType, // type
		obj.Durable,      // durable
		false,            // delete when complete
		false,            // internal
		obj.NoWait,       // noWait
		nil,              // arguments
	)
	if err != nil {
		log.Println("Error", "Error", err)
		err = errors.New(obj.Exchange + " : declare fail")
		return
	}
	return
}

func (obj *RabbitMQ) DeclareQueue() (err error) {
	obj.QueuePtr, err = obj.ConnectionPtr.channel.QueueDeclare(
		obj.Queue,   // name of the queue
		obj.Durable, // durable
		false,       // delete when usused
		false,       // exclusive
		obj.NoWait,  // noWait
		nil,         // arguments
	)
	if err != nil {
		log.Println("Error", "Error", err)
		err = errors.New(obj.Queue + " : declare fail")
		return
	}
	return
}

func (obj *RabbitMQ) BindQueue() (err error) {
	err = obj.ConnectionPtr.channel.QueueBind(
		obj.QueuePtr.Name, // name of the queue
		obj.Key,           // bindingKey
		obj.Exchange,      // sourceExchange
		obj.NoWait,        // noWait
		nil,               // arguments
	)
	if err != nil {
		log.Println("Error", "Error", err)
		err = errors.New(obj.QueuePtr.Name + " : queue bind fail")
		return
	}
	return
}

func (obj *RabbitMQ) Consume() (err error) {
	obj.DeliveriesChan, err = obj.ConnectionPtr.channel.Consume(
		obj.QueuePtr.Name, //queue
		obj.ConnectionPtr.tag,
		obj.AutoAck, //Auto Ack
		false,       //exclusive
		false,       //noLocal
		obj.NoWait,  //noWait
		nil)
	if err != nil {
		log.Println("Error", "Error", err)
		err = errors.New(obj.QueuePtr.Name + " : consume fail")
		return
	}
	return
}

func (obj *RabbitMQ) Shutdown() (err error) {
	if err := obj.ConnectionPtr.channel.Cancel(obj.ConnectionPtr.tag, true); err != nil {
		log.Println("Error", "Error", err)
		err = errors.New(obj.Queue + " : cancel fail")
		return err
	}

	if err = obj.ConnectionPtr.conn.Close(); err != nil {
		log.Println("Error", "Error", err)
		err = errors.New(obj.Queue + " : close fail")
		return
	}
	return
}

func (obj *RabbitMQ) CheckAliveness() (err error) {
	err = errors.New(fmt.Sprintf("%s", <-obj.ConnectionPtr.conn.NotifyClose(make(chan *amqp.Error))))
	return
}
