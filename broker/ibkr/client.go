package ibkr

import (
	"time"

	"github.com/scmhub/ibapi"
	"github.com/scmhub/ibsync"
)

// Client 封装 IBKR 客户端
type Client struct {
	ib *ibsync.IB
}

// Config 连接配置
type Config struct {
	Host     string
	Port     int
	ClientID int
	Timeout  time.Duration
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Host:     "127.0.0.1",
		Port:     7497,
		ClientID: 0, // 0 表示接收手动订单
		Timeout:  30 * time.Second,
	}
}

// NewClient 创建新的 IBKR 客户端
func NewClient() *Client {
	ibsync.SetConsoleWriter() // 设置美化的控制台日志
	return &Client{
		ib: ibsync.NewIB(),
	}
}

// Connect 连接到 IBKR
func (c *Client) Connect(config *Config) error {
	if config == nil {
		config = DefaultConfig()
	}

	err := c.ib.Connect(ibsync.NewConfig(
		ibsync.WithHost(config.Host),
		ibsync.WithPort(config.Port),
		ibsync.WithClientID(int64(config.ClientID)),
		ibsync.WithTimeout(config.Timeout),
		WithReadOnly(true),
	))
	if err != nil {
		return err
	}

	log := ibapi.Logger()
	accounts := c.ib.ManagedAccounts()
	log.Info().Strs("accounts", accounts).Msg("连接成功，管理的账户列表")

	return nil
}

func WithReadOnly(readonly bool) func(c *ibsync.Config) {
	return func(c *ibsync.Config) {
		c.ReadOnly = readonly
	}
}

// Disconnect 断开连接
func (c *Client) Disconnect() {
	if c.ib != nil {
		c.ib.Disconnect()
	}
}

// ManagedAccounts 获取管理的账户列表
func (c *Client) ManagedAccounts() []string {
	return c.ib.ManagedAccounts()
}

// IB 获取底层的 IB 客户端（用于高级操作）
func (c *Client) IB() *ibsync.IB {
	return c.ib
}
