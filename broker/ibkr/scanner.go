package ibkr

import (
	"github.com/scmhub/ibsync"
)

// ReqScannerParameters 请求扫描器参数
func (c *Client) ReqScannerParameters() (string, error) {
	return c.ib.ReqScannerParameters()
}

// ReqScannerSubscription 请求扫描器订阅
func (c *Client) ReqScannerSubscription(
	subscription *ibsync.ScannerSubscription,
	opts ...ibsync.ScannerSubscriptionOptions,
) ([]ibsync.ScanData, error) {
	return c.ib.ReqScannerSubscription(subscription, opts...)
}

// NewScannerSubscription 创建扫描器订阅
func NewScannerSubscription() *ibsync.ScannerSubscription {
	return ibsync.NewScannerSubscription()
}

// ScannerSubscriptionOptions 扫描器订阅选项类型
type ScannerSubscriptionOptions = ibsync.ScannerSubscriptionOptions

// TagValue 标签值类型
type TagValue = ibsync.TagValue
