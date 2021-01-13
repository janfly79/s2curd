# struct2curd: auto generate curd from  struct


A Swiss Army Knife helps you generate sql from [gorm](https://github.com/jinzhu/gorm) model struct.


## Installation

```
go get github.com/janfly79/struct2curd
```

## Usage

`blacklist.go`:

```

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
}
```

```
struct2curd curd -f ./testdata/blacklist.go -s Blacklist 
```

Result:

```
appFile to blacklist.go



```


## How it works



