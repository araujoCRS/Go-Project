package service

import (
	"client/configs"
	"encoding/json"

	"errors"
	"log"

	"github.com/streadway/amqp"
)

type RabbitMQChannel struct {
	channel   *amqp.Channel
	queueName string
}

type RabbitMQMessage struct {
	Body        []byte
	deliveryTag uint64
}

func (s *RabbitMQChannel) validateChannel() error {
	if s == nil || s.channel == nil {
		return errors.New("canal RabbitMQ não inicializado")
	}

	return nil
}

func NewRabbitMQService(conf configs.RabbitMQ, confQueue configs.RabbitMQQueue) (*RabbitMQChannel, *amqp.Connection, error) {
	conn, err := amqp.Dial(conf.GetConnectionString())
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, err
	}

	_, err = ch.QueueDeclare(
		confQueue.Name,
		confQueue.Durable,
		confQueue.AutoDelete,
		confQueue.Exclusive,
		confQueue.Passive,
		confQueue.Args,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, nil, errors.New("fila não existe ou conexão falhou: " + err.Error())
	}

	channel := &RabbitMQChannel{
		channel:   ch,
		queueName: confQueue.Name,
	}

	return channel, conn, nil
}

func (s *RabbitMQChannel) ListenClose(conn *amqp.Connection, onClose func(error)) {
	errChan := make(chan *amqp.Error)

	// Notifica fechamento da conexão
	s.channel.NotifyClose(errChan)
	//conn.NotifyClose(errChan)

	go func() {
		err := <-errChan
		if err != nil {
			onClose(err)
		}
	}()
}

func (s *RabbitMQChannel) Publish(data any) error {
	if err := s.validateChannel(); err != nil {
		return err
	}

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = s.channel.Publish(
		"",          // Exchange
		s.queueName, // Nome da fila
		false,       // Mandatório
		false,       // Imediato
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("Erro ao publicar mensagem na fila: %v", err)
		return err
	}

	log.Printf("Cliente enviado para a fila %s: %s", s.queueName, string(body))
	return nil
}

func (s *RabbitMQChannel) Dequeue(ack bool) (*RabbitMQMessage, error) {
	if err := s.validateChannel(); err != nil {
		return nil, err
	}

	msg, ok, err := s.channel.Get(
		s.queueName,
		ack,
	)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("nenhuma mensagem disponível na fila")
	}
	return &RabbitMQMessage{
		Body:        msg.Body,
		deliveryTag: msg.DeliveryTag,
	}, nil
}

func (s *RabbitMQChannel) Consume() (<-chan amqp.Delivery, error) {
	if err := s.validateChannel(); err != nil {
		return nil, err
	}

	delivery, err := s.channel.Consume(
		s.queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return delivery, nil
}

func (s *RabbitMQChannel) SendAck(msg RabbitMQMessage) error {
	if err := s.validateChannel(); err != nil {
		return err
	}
	return s.channel.Ack(msg.deliveryTag, false)
}

func (s *RabbitMQChannel) Close() {
	if s != nil && s.channel != nil {
		s.channel.Close()
	}
}
