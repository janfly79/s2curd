



func get{{.StructTableName}}Table() string {
    return "{{.TableName}}"
}

const insert{{.StructTableName}}Fields = "{{.InsertFieldList}}"

const list{{.StructTableName}}Fields = "{{.AllFieldList}}"

const db{{.StructTableName}} = "user"

func(item *{{.StructTableName}}) Scan(row db.Row) (err error){
    if err = row.Scan(
        {{range .NullFieldsInfo}}&item.{{.HumpName}},
    {{end}}); err != nil {
        return
    }
    return
}

// 添加
func add{{.StructTableName}} (ctx context.Context, item {{.StructTableName}})(lastId int64, err error){
    conn := db.Get(ctx, "user")
    sql := "INSERT INTO " + get{{.StructTableName}}Table() + " (" + insert{{.StructTableName}}Fields + ") " +
            "VALUES ({{.InsertMark}})"
    q := db.SQLInsert(get{{.StructTableName}}Table(), sql)
    res, err := conn.ExecContext(ctx, q,
        {{range .InsertInfo}}item.{{.HumpName}},
        {{end}})
    if err != nil {
        return
    }
    lastId, _ = res.LastInsertId()
    return
}

func del{{.StructTableName}} (ctx context.Context, id interface{}) (err error){
    conn := db.Get(ctx, "user")
    sql := "delete from " + get{{.StructTableName}}Table() + " where id = ?"
    q := db.SQLDelete(get{{.StructTableName}}Table(), sql)
    _, err = conn.ExecContext(ctx, q, id)
    return
}


// 获取单条记录
func get{{.StructTableName}} (ctx context.Context, where string, args []interface{})(row {{.StructTableName}}, err error){
    conn := db.Get(ctx, db{{.StructTableName}})
    sqlText := "select " + list{{.StructTableName}}Fields + " from " + get{{.StructTableName}}Table() + " " + where
    q := db.SQLSelect(get{{.StructTableName}}Table(), sqlText)
    sqlRow := conn.QueryRowContext(ctx, q, args...)
    if err = row.Scan(sqlRow); db.IsNoRowsErr(err) {
    	err = nil
    }
    return
}

// 更新
func update{{.StructTableName}}(ctx context.Context, updateStr string, where string, args []interface{})(err error) {
    conn := db.Get(ctx, db{{.StructTableName}})
    sqlText := "update " + get{{.StructTableName}}Table() + " set " + updateStr + " " + where
    q := db.SQLUpdate(get{{.StructTableName}}Table(), sqlText)
    _, err = conn.ExecContext(ctx, q, args...)
    return
}

// 列表
func list{{.StructTableName}} (ctx context.Context, where string, args []interface{}) (rowsResult []{{.StructTableName}}, err error) {
    conn := db.Get(ctx, db{{.StructTableName}})
    sqlText := "select " + list{{.StructTableName}}Fields + " from "+get{{.StructTableName}}Table()+" " + where
    q := db.SQLSelect(get{{.StructTableName}}Table(), sqlText)
    rows, err := conn.QueryContext(ctx, q, args...)
    if err != nil {
            return
    }
    defer rows.Close()
    for rows.Next() {
            row := {{.StructTableName}}{}
            if err = row.Scan(rows); err != nil {
                return
            }
            rowsResult = append(rowsResult, row)
    }
    err = rows.Err()
    return
}



func count{{.StructTableName}} (ctx context.Context, where string, args []interface{})(total int32, err error){
    conn := db.Get(ctx, db{{.StructTableName}})
    sqlText := "select count(*) from "+ get{{.StructTableName}}Table()+" " + where
    q := db.SQLSelect(get{{.StructTableName}}Table(), sqlText)
    err = conn.QueryRowContext(ctx, q, args...).Scan(&total)
    return
}




