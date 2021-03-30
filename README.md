# s2curd: auto generate curd from  struct


A Swiss Army Knife helps you generate sql from [gorm](https://github.com/jinzhu/gorm) model struct.


## Installation

```
go get github.com/janfly79/s2curd

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
s2curd  curd -f testdata/sqlmodel/blacklist.go -s Blacklist -d user -t t_black_list
```

Result:

```
appFile to blacklist.go

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

const insertBlacklistFields = "`reason`,`uid`,`coupon_id`,`start_time`,`end_time`,`cuser`"

const listBlacklistFields = "`id`,`reason`,`uid`,`coupon_id`,`start_time`,`end_time`,`cuser`,`ctime`,`mtime`"

func (item *Blacklist) Scan(row db.Row) (err error) {
	if err = row.Scan(
		&item.ID,
		&item.Reason,
		&item.UID,
		&item.CouponID,
		&item.StartTime,
		&item.EndTime,
		&item.Cuser,
		&item.Ctime,
		&item.Mtime,
	); err != nil {
		return
	}
	return
}

// 添加
func addBlacklist(ctx context.Context, item Blacklist) (lastId int64, err error) {
	conn := db.Get(ctx, "user")
	sql := "insert into t_black_list (" + insertBlacklistFields + ") " +
		"VALUES (?,?,?,?,?,?)"
	q := db.SQLInsert("t_black_list", sql)
	res, err := conn.ExecContext(ctx, q,
		item.Reason,
		item.UID,
		item.CouponID,
		item.StartTime,
		item.EndTime,
		item.Cuser,
	)
	if err != nil {
		return
	}
	lastId, _ = res.LastInsertId()
	return
}

func delBlacklist(ctx context.Context, id interface{}) (err error) {
	conn := db.Get(ctx, "user")
	sql := "delete from t_black_list where id = ?"
	q := db.SQLDelete("t_black_list", sql)
	_, err = conn.ExecContext(ctx, q, id)
	return
}

// 获取单条记录
func getBlacklist(ctx context.Context, where string, args []interface{}) (row Blacklist, err error) {
	conn := db.Get(ctx, "user")
	sqlText := "select " + listBlacklistFields + " from t_black_list " + where
	q := db.SQLSelect("t_black_list", sqlText)
	sqlRow := conn.QueryRowContext(ctx, q, args...)
	if err = row.Scan(sqlRow); db.IsNoRowsErr(err) {
		err = nil
	}
	return
}

// 更新
func updateBlacklist(ctx context.Context, updateStr string, where string, args []interface{}) (err error) {
	conn := db.Get(ctx, "user")
	sqlText := "update t_black_list set " + updateStr + " " + where
	q := db.SQLUpdate("t_black_list", sqlText)
	_, err = conn.ExecContext(ctx, q, args...)
	return
}

// 列表
func listBlacklist(ctx context.Context, where string, args []interface{}) (rowsResult []Blacklist, err error) {
	conn := db.Get(ctx, "user")
	sqlText := "select " + listBlacklistFields + " from t_black_list " + where
	q := db.SQLSelect("t_black_list", sqlText)
	rows, err := conn.QueryContext(ctx, q, args...)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		row := Blacklist{}
		if err = row.Scan(rows); err != nil {
			return
		}
		rowsResult = append(rowsResult, row)
	}
	err = rows.Err()
	return
}

func countBlacklist(ctx context.Context, where string, args []interface{}) (total int32, err error) {
	conn := db.Get(ctx, "user")
	sqlText := "select count(*) from t_black_list " + where
	q := db.SQLSelect("t_black_list", sqlText)
	err = conn.QueryRowContext(ctx, q, args...).Scan(&total)
	return
}

```


## How it works



