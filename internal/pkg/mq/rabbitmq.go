package mq

import (
    "fmt"
    "mingda_cloud_service/internal/pkg/config"
    "github.com/streadway/amqp"
)

func NewRabbitMQConnection(cfg *config.RabbitMQConfig) (*amqp.Connection, error) {
    url := fmt.Sprintf("amqp://%s:%s@%s:%d/",
        cfg.Username,
        cfg.Password,
        cfg.Host,
        cfg.Port,
    )
    
    return amqp.Dial(url)
}
