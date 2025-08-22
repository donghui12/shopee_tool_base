package pool

import (
    "fmt"
    "sync"

    "github.com/panjf2000/ants/v2"
    "go.uber.org/zap"
    "github.com/shopee_tool_base/pkg/logger"
    "github.com/shopee_tool_base/pkg/constant"
)

var (
	workerPool *WorkerPool
	once sync.Once
)

type Task struct {
    Topic   string
    Execute    func() error
}

func InitWorkerPool() {
	once.Do(func() {
		workerPool = NewWorkerPool(constant.WorkerPoolSize)
		for _, topic := range constant.Topics {
			workerPool.RegisterTopic(topic)
		}
	})
}

func GetWorkerPool() *WorkerPool {
	return workerPool
}

type WorkerPool struct {
    pools       map[string]*ants.Pool
    isRunning   bool
    mu          sync.RWMutex
    poolSize    int
    poolOptions *ants.Options
}

// 创建新的 WorkerPool
func NewWorkerPool(size int) *WorkerPool {
    options := &ants.Options{
        PreAlloc:       true,
        MaxBlockingTasks: 1000,
        Nonblocking:    false,
        PanicHandler: func(err interface{}) {
            logger.Error("Worker pool panic",
                zap.Any("error", err),
            )
        },
    }

    return &WorkerPool{
        pools:       make(map[string]*ants.Pool),
        poolSize:    size,
        poolOptions: options,
    }
}

// 注册新的 topic
func (p *WorkerPool) RegisterTopic(topic string) error {
    p.mu.Lock()
    defer p.mu.Unlock()

    if _, exists := p.pools[topic]; exists {
        return fmt.Errorf("topic %s already registered", topic)
    }

    pool, err := ants.NewPool(p.poolSize, ants.WithOptions(*p.poolOptions))
    if err != nil {
        return fmt.Errorf("failed to create pool for topic %s: %v", topic, err)
    }

    p.pools[topic] = pool
    logger.Info("Registered new topic",
        zap.String("topic", string(topic)),
        zap.Int("pool_size", p.poolSize),
    )

    return nil
}

// 提交任务
func (p *WorkerPool) Submit(task Task) error {
    p.mu.RLock()
    pool, exists := p.pools[task.Topic]
    p.mu.RUnlock()

    if !exists {
        return fmt.Errorf("unknown topic: %s", task.Topic)
    }

    err := pool.Submit(func() {
        if err := task.Execute(); err != nil {
            logger.Error("Task execution failed",
                zap.String("topic", string(task.Topic)),
                zap.Error(err),
            )
        }
    })

    if err != nil {
        logger.Error("Failed to submit task",
            zap.String("topic", string(task.Topic)),
            zap.Error(err),
        )
        return fmt.Errorf("failed to submit task: %v", err)
    }

    logger.Debug("Task submitted successfully",
        zap.String("topic", string(task.Topic)),
    )
    return nil
}

// 关闭指定 topic 的池
func (p *WorkerPool) ClosePool(topic string) error {
    p.mu.Lock()
    defer p.mu.Unlock()

    pool, exists := p.pools[topic]
    if !exists {
        return fmt.Errorf("unknown topic: %s", topic)
    }

    pool.Release()
    delete(p.pools, topic)
    logger.Info("Closed pool for topic",
        zap.String("topic", string(topic)),
    )

    return nil
}

// 关闭所有池
func (p *WorkerPool) Release() {
    p.mu.Lock()
    defer p.mu.Unlock()

    for topic, pool := range p.pools {
        pool.Release()
        logger.Info("Released pool",
            zap.String("topic", string(topic)),
        )
    }
    p.pools = make(map[string]*ants.Pool)
}
