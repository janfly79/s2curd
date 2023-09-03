package sqlmodel

import (
	"time"
)

type Blacklist struct {
	ID        int32 // 主键
	Reason    string
	UID       int64     // 用户 uid
	CouponID  int32     // 卡券id
	StartTime time.Time // 开始时间
	EndTime   time.Time // 结束时间
	Cuser     string    // 创建者
	Ctime     time.Time
	Mtime     time.Time
	Ext       BlackExt
}

type BlackExt struct {
	Name string
	Age  int32
}
