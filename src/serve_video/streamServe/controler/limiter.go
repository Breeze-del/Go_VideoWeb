package controler

import "log"

// 流控
// 利用缓存channel来控制连接个数, 一个request便添加连接直到response便释放连接
type ConnLimiter struct {
	concurrentConn int // 并发连接数
	bucket         chan int
}

// 创建一个流控制器
func NewConnLimiter(cc int) *ConnLimiter {
	return &ConnLimiter{
		concurrentConn: cc,
		bucket:         make(chan int, cc),
	}
}

// 添加一个连接
func (c *ConnLimiter) GetConn() bool {
	if len(c.bucket) >= c.concurrentConn {
		log.Print("Reached the rate limitation")
		return false
	}
	c.bucket <- 1
	return true
}

// 释放一个连接
func (c *ConnLimiter) RealseConn() {
	num := <-c.bucket
	log.Printf("New connection coming: %d", num)
}
