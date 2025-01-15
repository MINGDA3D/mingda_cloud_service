package rabbitmq

import (
	"fmt"

	"github.com/streadway/amqp"
	"mingda_cloud_service/internal/pkg/config"
)

var Conn *amqp.Connection
var Channel *amqp.Channel

// Init 初始化RabbitMQ连接
func Init(cfg config.RabbitMQConfig) error {
	var err error
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.VHost,
	)

	// 建立连接
	Conn, err = amqp.Dial(url)
	if err != nil {
		return fmt.Errorf("connect rabbitmq failed: %v", err)
	}

	// 创建Channel
	Channel, err = Conn.Channel()
	if err != nil {
		return fmt.Errorf("create channel failed: %v", err)
	}

	return nil
}

// Close 关闭连接
func Close() {
	if Channel != nil {
		Channel.Close()
	}
	if Conn != nil {
		Conn.Close()
	}
}

// IsConnected 检查连接状态
func IsConnected() bool {
	if Conn == nil || Channel == nil {
		return false
	}

	// 检查连接是否关闭
	if Conn.IsClosed() {
		return false
	}

	// 尝试声明一个临时队列来测试channel是否正常
	_, err := Channel.QueueDeclare(
		"",    // name
		false, // durable
		true,  // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)

	return err == nil
} 