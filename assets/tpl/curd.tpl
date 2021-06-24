




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
    conn := db.Get(ctx, "{{.DBName}}")
    sql := "insert into {{.OriginTableName}} (" +
            {{range $i,$el := .InsertFieldList}}{{if $i}},"+
            {{else}}{{end}}"{{$el}}{{end}}"+
            ") values ({{.InsertMark}})"
    q := db.SQLInsert("{{.OriginTableName}}", sql)

    res, err := conn.ExecContext(ctx, q,
        {{range .InsertInfo}}item.{{.HumpName}},
        {{end}})
    if err != nil {
        return
    }

    lastId, err = res.LastInsertId()

    return
}

func del{{.StructTableName}} (ctx context.Context, id interface{}) (err error){
    conn := db.Get(ctx, "{{.DBName}}")
    sql := "delete from {{.OriginTableName}} where id = ?"
    q := db.SQLDelete("{{.OriginTableName}}", sql)

    _, err = conn.ExecContext(ctx, q, id)

    return
}

func update{{.StructTableName}}(ctx context.Context, where string, args []interface{})(err error) {
    conn := db.Get(ctx, "{{.DBName}}")
    sql := "update {{.OriginTableName}} set " +
                {{range $i,$el := .InsertFieldList}}{{if $i}},"+
                {{else}}{{end}}"{{$el}}=?{{end}} " + where
    q := db.SQLUpdate("{{.OriginTableName}}", sql)

    _, err = conn.ExecContext(ctx, q, args...)

    return
}


func list{{.StructTableName}} (ctx context.Context, where string, args []interface{}) (rowsResult []{{.StructTableName}}, err error) {
    conn := db.Get(ctx, "{{.DBName}}")
    sql := "select "+
                    {{range $i,$el := .AllFieldList}}{{if $i}},"+
                    {{else}}{{end}}"{{$el}}{{end}} "+
                    "from {{.OriginTableName}} " + where
    q := db.SQLSelect("{{.OriginTableName}}", sql)

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
    conn := db.Get(ctx, "{{.DBName}}")
    sql := "select count(*) from {{.OriginTableName}} " + where
    q := db.SQLSelect("{{.OriginTableName}}", sql)

    err = conn.QueryRowContext(ctx, q, args...).Scan(&total)

    return
}




