



func get{{.StructTableName}}TableName() string {
    return "{{.TableName}}"
}

const insert{{.StructTableName}}Fields = ""

const list{{.StructTableName}}Fields = ""

const db{{.StructTableName}} = "user"

func(item *{{.StructTableName}}) Scan(row db.Row) (err error){
    if err = row.Scan(
        {{range .NullFieldsInfo}}
        &row.{{.HumpName}},
        {{end}}
    ); err != nil {
        return
    }
    return
}

// 添加
func add(ctx context.Context, value {{.StructTableName}}(lastId int64, err error){
    conn := db.Get(ctx, "user")
    sql := "INSERT INTO " + get{{.StructTableName}}TableName() + " ({{.InsertFieldList}}) " +
            "VALUES ({{.InsertMark}})"
    q := db.SQLInsert(get{{.StructTableName}}TableName(), sql)
    res, err := conn.ExecContext(ctx, q,
        {{range .InsertInfo}}value.{{.HumpName}},
        {{end}})
    if err != nil {
        return
    }
    lastId, _ = res.LastInsertId()
    return
}

func del(ctx context.Context, id interface{}) (err error){
    conn := db.Get(ctx, "user")
    sql := "delete from " + get{{.StructTableName}}TableName() + " where id = ?"
    q := db.SQLDelete(get{{.StructTableName}}TableName(), sql)
    _, err = conn.ExecContext(ctx, q, id)
    return
}





func Add{{.StructTableName}}(ctx context.Context, value {{.StructTableName}})(lastId int64, err error) {
    conn := db.Get(ctx, {{.StructTableName}}_DB)
    sql := "INSERT INTO " + get{{.StructTableName}}TableName() + " ({{.InsertFieldList}}) " +
        "VALUES ({{.InsertMark}})"
    q := db.SQLInsert(get{{.StructTableName}}TableName(), sql)
    res, err := conn.ExecContext(ctx, q,
        {{range .InsertInfo}}value.{{.HumpName}},
        {{end}})
    if err != nil {
        return
    }
    lastId, _ = res.LastInsertId()
    return
}

// 删除单条记录
func Del{{.StructTableName}}(ctx context.Context, where string, args []interface{}) (err error) {
    conn := db.Get(ctx, {{.StructTableName}}_DB)
	sql := "delete from {{.TableName}} " + where
	q := db.SQLDelete("{{.TableName}}", sql)

	_, err = conn.ExecContext(ctx, q, args...)
    return
}

// 获取单条记录
func Get{{.StructTableName}}(ctx context.Context, fields string, where string, args []interface{})(row {{.StructTableName}}, err error){
    conn := db.Get(ctx, {{.StructTableName}}_DB)
    sqlText := "select " + fields + " from {{.TableName}} " + where
    q := db.SQLSelect("{{.TableName}}", sqlText)
    err = conn.QueryRowContext(ctx, q, args...).Scan(
            		{{range .NullFieldsInfo}}&row.{{.HumpName}},// {{.Comment}}
            		{{end}})
    return
}

// 更新
func Update{{.StructTableName}}(ctx context.Context, updateStr string, where string, args []interface{})(err error) {
    conn := db.Get(ctx, {{.StructTableName}}_DB)
    sqlText := "update {{.TableName}} set " + updateStr + " " + where
    q := db.SQLUpdate("{{.TableName}}", sqlText)
    _, err = conn.ExecContext(ctx, q, args...)
    return
}

// 列表
func list(ctx context.Context, where string, args []interface{}) (rowsResult []{{.StructTableName}}, err error) {
    conn := db.Get(ctx, {{.StructTableName}}_DB)
    sqlText := "select " + list{{.StructTableName}}Fields + " from {{.TableName}} " + where
    q := db.SQLSelect("{{.TableName}}", sqlText)
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


// 列表
func List{{.StructTableName}}(ctx context.Context, fields string, where string, args []interface{})(rowsResult []{{.StructTableName}}, err error) {
    conn := db.Get(ctx, {{.StructTableName}}_DB)
    sqlText := "select " + fields + " from {{.TableName}} " + where
    q := db.SQLSelect("{{.TableName}}", sqlText)
    rows, err := conn.QueryContext(ctx, q, args...)
    if err != nil {
    		return
    	}
    defer rows.Close()

    for rows.Next() {
        row := {{.StructTableName}}{}
        if err = rows.Scan(
            {{range .NullFieldsInfo}}&row.{{.HumpName}},// {{.Comment}}
                        		{{end}}
        ); err != nil {
            return
        }
        rowsResult = append(rowsResult, row)
    }

    err = rows.Err()

    return
}


func count(ctx context.Context, where string, args []interface{})(total int32, err error){
    conn := db.Get(ctx, db{{.StructTableName}})
    sqlText := "select count(*) from {{.TableName}} " + where
    q := db.SQLSelect("{{.TableName}}", sqlText)
    err = conn.QueryRowContext(ctx, q, args...).Scan(&total)
    return
}




