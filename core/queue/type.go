package queue

type Pusher interface {
	Name() string
	Push(string) error
}

// MessageQueue接口代表一个消息队列
type MessageQueue interface {
	Start()
	Stop()
}

type (

	// 一个消费者接口代表一个可以消费字符串信息的消费者
	Consumer interface {
		Consume(string) error // * 表示该消费者可以消费字符串类型的消息，并返回错误
		OnEvent(event any)    // * 表示该消费者可以处理任何类型的事件。
	}

	// 生成消费者的工厂方法
	ConsumerFactory func() (Consumer, error)

	// 定义生产者接口
	Producer interface {
		AddListener(listener ProduceListener)
		Produce() (string, bool)
	}

	// 定义生产者监听器接口
	ProduceListener interface {
		OnProducerPause()
		OnProducerResume()
	}

	// 生成生产者的工厂方法
	ProducerFactory func() (Producer, error)

	// 定义一个监听器接口，用于监听队列的暂停和恢复事件
	Listener interface {
		OnPause()
		OnResume()
	}

	// 投票
	// Poller interface {
	// 	Name() string
	// 	Poll() string
	// }
)
