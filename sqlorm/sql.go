package sqlorm

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"html/template"
	"strconv"
	"strings"

	"github.com/janfly79/s2curd/util"

	"github.com/janfly79/s2curd/bindata"

	"github.com/pinzolo/casee"

	log "github.com/liudanking/goutil/logutil"
)

type SqlGenerator struct {
	structName string
	modelType  *ast.StructType
}

// 生成select,update,insert,delete所需信息
type SqlInfo struct {
	TableName           string              // 表名
	PrimaryKey          string              // 主键字段
	PrimaryType         string              // 主键类型
	StructTableName     string              // 结构表名称
	NullStructTableName string              // 判断为空的表名
	PkgEntity           string              // 实体空间名称
	PkgTable            string              // 表的空间名称
	UpperTableName      string              // 大写的表名
	AllFieldList        string              // 所有字段列表,如: id,name
	InsertFieldList     string              // 插入字段列表,如:id,name
	InsertMark          string              // 插入字段使用多少个?,如: ?,?,?
	UpdateFieldList     string              // 更新字段列表
	SecondField         string              // 存放第二个字段
	UpdateListField     []string            // 更新字段列表
	FieldsInfo          []*SqlFieldInfo     // 字段信息
	NullFieldsInfo      []*NullSqlFieldInfo // 判断为空时
	InsertInfo          []*SqlFieldInfo
	OriginTableName string // 真正的表名
	DBName string // 数据库名
}

// 查询使用的字段结构信息
type SqlFieldInfo struct {
	HumpName string // 驼峰字段名称
	Comment  string // 字段注释
}
type NullSqlFieldInfo struct {
	GoType       string // golang类型
	HumpName     string // 驼峰字段名称
	OriFieldType string // 原数据库类型
	Comment      string // 字段注释
}

// 表名与表注释
type TableNameAndComment struct {
	Index   int
	Name    string
	Comment string
}

func NewSqlGenerator(typeSpec *ast.TypeSpec) (*SqlGenerator, error) {
	structType, ok := typeSpec.Type.(*ast.StructType)
	if !ok {
		return nil, errors.New("typeSpec is not struct type")
	}

	return &SqlGenerator{
		structName: typeSpec.Name.Name,
		modelType:  structType,
	}, nil
}

func (ms *SqlGenerator) AddFuncStr() (string, error) {

	var (
		columnList []string
		//primaryKey string
		InsertInfo    = make([]*SqlFieldInfo, 0)
		InsertMark    string
		insertFields  = make([]string, 0)
		allFields     = make([]string, 0)
		nullFieldList = make([]*NullSqlFieldInfo, 0)
	)

	for _, field := range ms.getStructFieds(ms.modelType) {
		columnName := getColumnName(field)

		nullFieldList = append(nullFieldList, &NullSqlFieldInfo{
			HumpName: field.Names[0].Name,
			Comment:  "",
		})

		allFields = append(allFields, columnName)

		if isPrimaryKey(field) || isTimeKey(field) {
			continue
		}

		insertFields = append(insertFields, columnName)
		InsertInfo = append(InsertInfo, &SqlFieldInfo{
			HumpName: field.Names[0].Name,
			Comment:  "",
		})
		// 拼出SQL所需要结构数据
		InsertMark = strings.Repeat("?,", len(insertFields))
		columnList = append(columnList, columnName)
	}

	sqlInfo := &SqlInfo{
		TableName: ms.tableName(),
		//PrimaryKey:          AddQuote(PrimaryKey),
		//PrimaryType:         primaryType,
		StructTableName: ms.structName,
		PkgEntity:       ".",
		PkgTable:        ".",
		AllFieldList:    strings.Join(allFields, ","),
		InsertFieldList: strings.Join(columnList, ","),
		InsertMark:      strings.TrimRight(InsertMark, ","),
		//UpdateFieldList:     strings.Join(updateList, ","),
		//UpdateListField:     updateListField,
		//FieldsInfo:          fieldsList,
		NullFieldsInfo: nullFieldList,
		InsertInfo:     InsertInfo,
		//SecondField:         AddQuote(secondField),
		OriginTableName:"xxx",// 真正的表名
		DBName:"comic",//数据库名
	}

	// 解析模板
	tplByte, err := bindata.Asset("assets/tpl/curd.tpl")
	if err != nil {
		return "", err
	}
	tpl, err := template.New("CURD").Parse(string(tplByte))
	if err != nil {
		return "", err
	}
	// 解析
	content := bytes.NewBuffer([]byte{})
	err = tpl.Execute(content, sqlInfo)
	if err != nil {
		return "", err
	}
	return content.String(), nil
}


func (ms *SqlGenerator) AddCurdFuncStr(originDBName, originTableName string) (string, error) {

	var (
		columnList []string
		//primaryKey string
		InsertInfo    = make([]*SqlFieldInfo, 0)
		InsertMark    string
		insertFields  = make([]string, 0)
		allFields     = make([]string, 0)
		nullFieldList = make([]*NullSqlFieldInfo, 0)
	)

	for _, field := range ms.getStructFieds(ms.modelType) {
		columnName := getColumnName(field)

		nullFieldList = append(nullFieldList, &NullSqlFieldInfo{
			HumpName: field.Names[0].Name,
			Comment:  "",
		})

		allFields = append(allFields, columnName)

		if isPrimaryKey(field) || isTimeKey(field) {
			continue
		}

		insertFields = append(insertFields, columnName)
		InsertInfo = append(InsertInfo, &SqlFieldInfo{
			HumpName: field.Names[0].Name,
			Comment:  "",
		})
		// 拼出SQL所需要结构数据
		InsertMark = strings.Repeat("?,", len(insertFields))
		columnList = append(columnList, columnName)
	}

	sqlInfo := &SqlInfo{
		TableName: ms.tableName(),
		//PrimaryKey:          AddQuote(PrimaryKey),
		//PrimaryType:         primaryType,
		StructTableName: ms.structName,
		PkgEntity:       ".",
		PkgTable:        ".",
		AllFieldList:    strings.Join(allFields, ","),
		InsertFieldList: strings.Join(columnList, ","),
		InsertMark:      strings.TrimRight(InsertMark, ","),
		//UpdateFieldList:     strings.Join(updateList, ","),
		//UpdateListField:     updateListField,
		//FieldsInfo:          fieldsList,
		NullFieldsInfo: nullFieldList,
		InsertInfo:     InsertInfo,
		//SecondField:         AddQuote(secondField),
		OriginTableName:originTableName,// 真正的表名
		DBName:originDBName,//数据库名
	}

	// 解析模板
	tplByte, err := bindata.Asset("assets/tpl/curd.tpl")
	if err != nil {
		return "", err
	}
	tpl, err := template.New("CURD").Parse(string(tplByte))
	if err != nil {
		return "", err
	}
	// 解析
	content := bytes.NewBuffer([]byte{})
	err = tpl.Execute(content, sqlInfo)
	if err != nil {
		return "", err
	}
	return content.String(), nil
}


func (ms *SqlGenerator) AddCacheFuncStr() (string, error) {

	var (
		columnList []string
		//primaryKey string
		InsertInfo    = make([]*SqlFieldInfo, 0)
		InsertMark    string
		insertFields  = make([]string, 0)
		allFields     = make([]string, 0)
		nullFieldList = make([]*NullSqlFieldInfo, 0)
	)

	for _, field := range ms.getStructFieds(ms.modelType) {
		columnName := getColumnName(field)

		nullFieldList = append(nullFieldList, &NullSqlFieldInfo{
			HumpName: field.Names[0].Name,
			Comment:  "",
		})

		allFields = append(allFields, columnName)

		if isPrimaryKey(field) || isTimeKey(field) {
			continue
		}

		insertFields = append(insertFields, columnName)
		InsertInfo = append(InsertInfo, &SqlFieldInfo{
			HumpName: field.Names[0].Name,
			Comment:  "",
		})
		// 拼出SQL所需要结构数据
		InsertMark = strings.Repeat("?,", len(insertFields))
		columnList = append(columnList, columnName)
	}

	sqlInfo := &SqlInfo{
		TableName: ms.tableName(),
		//PrimaryKey:          AddQuote(PrimaryKey),
		//PrimaryType:         primaryType,
		StructTableName: ms.structName,
		PkgEntity:       ".",
		PkgTable:        ".",
		AllFieldList:    strings.Join(allFields, ","),
		InsertFieldList: strings.Join(columnList, ","),
		InsertMark:      strings.TrimRight(InsertMark, ","),
		//UpdateFieldList:     strings.Join(updateList, ","),
		//UpdateListField:     updateListField,
		//FieldsInfo:          fieldsList,
		NullFieldsInfo: nullFieldList,
		InsertInfo:     InsertInfo,
		//SecondField:         AddQuote(secondField),
		OriginTableName:"xxx",// 真正的表名
		DBName:"comic",//数据库名
	}

	// 解析模板
	tplByte, err := bindata.Asset("assets/tpl/mc.tpl")
	if err != nil {
		return "", err
	}
	tpl, err := template.New("mc").Parse(string(tplByte))
	if err != nil {
		return "", err
	}
	// 解析
	content := bytes.NewBuffer([]byte{})
	err = tpl.Execute(content, sqlInfo)
	if err != nil {
		return "", err
	}
	return content.String(), nil
}

func (ms *SqlGenerator) GetCreateTableSql() (string, error) {
	var tags []string
	var primaryKeys []string
	indices := map[string][]string{}
	uniqIndces := map[string][]string{}
	log.Info("mode types %+v", ms.modelType)
	for _, field := range ms.getStructFieds(ms.modelType) {
		log.Info("===<<<<<<=========")
		log.Info("field unknow string %+v", util.GetFieldName(field))
		log.Info("========>>>>>>====")
		switch t := field.Type.(type) {
		case *ast.Ident:
			tag, err := generateSqlTag(field)
			log.Info("mode tag string %+v", tag)
			if err != nil {
				log.Warning("generateSqlTag [%s] failed:%v", t.Name, err)
			} else {
				tags = append(tags, fmt.Sprintf("%s %s", getColumnName(field), tag))
			}
		case *ast.SelectorExpr:
			log.Info("columnName node ")
			tag, err := generateSqlTag(field)
			if err != nil {
				log.Warning("generateSqlTag [%s] failed:%v", t.Sel.Name, err)
			} else {
				tags = append(tags, fmt.Sprintf("%s %s", getColumnName(field), tag))
			}
		default:
			log.Warning("field %s not supported, ignore", util.GetFieldName(field))
		}

		columnName := getColumnName(field)
		log.Info("columnName this is my column %+v", columnName)
		if isPrimaryKey(field) {
			primaryKeys = append(primaryKeys, columnName)
		}
		log.Info("primary key name %+v", primaryKeys)
		sqlSettings := ParseTagSetting(util.GetFieldTag(field, "sql").Name)
		if idxName, ok := sqlSettings["INDEX"]; ok {
			keys := indices[idxName]
			keys = append(keys, columnName)
			indices[idxName] = keys
		}
		if idxName, ok := sqlSettings["UNIQUE_INDEX"]; ok {
			keys := uniqIndces[idxName]
			keys = append(keys, columnName)
			uniqIndces[idxName] = keys
		}

	}

	var primaryKeyStr string
	if len(primaryKeys) > 0 {
		primaryKeyStr = fmt.Sprintf("PRIMARY KEY (%v)", strings.Join(primaryKeys, ", "))
	}

	indicesStrs := []string{}
	for idxName, keys := range indices {
		indicesStrs = append(indicesStrs, fmt.Sprintf("INDEX %s (%s)", idxName, strings.Join(keys, ", ")))
	}

	uniqIndicesStrs := []string{}
	for idxName, keys := range uniqIndces {
		uniqIndicesStrs = append(uniqIndicesStrs, fmt.Sprintf("UNIQUE INDEX %s (%s)", idxName, strings.Join(keys, ", ")))
	}

	options := []string{
		"engine=innodb",
		"DEFAULT charset=utf8mb4",
	}

	return fmt.Sprintf(`CREATE TABLE %v 
(
  %v,
  %v
) %v;`,
		"`"+ms.tableName()+"`",
		strings.Join(append(tags, append(indicesStrs, uniqIndicesStrs...)...), ",\n  "),
		primaryKeyStr,
		strings.Join(options, " ")), nil
}

func structSelection(node ast.Node) (int, int, error) {

	encStruct, ok := node.(*ast.StructType)

	if !ok {
		return 0, 0, errors.New("struct name does not exist")
	}

	fset := token.NewFileSet()


	start := fset.Position(encStruct.Pos()).Line
	end := fset.Position(encStruct.End()).Line

	return start, end, nil
}

func StructFileLine(node ast.Node) (string, error){
	var buf bytes.Buffer
	fset := token.NewFileSet()
	err := format.Node(&buf, fset, node)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}



func (ms *SqlGenerator) getStructFieds(node ast.Node) []*ast.Field {
	var fields []*ast.Field
	nodeType, ok := node.(*ast.StructType)
	if !ok {
		return nil
	}

	//structStr, err := StructFileLine(node)
	//
	//writefile.WriteFile("cc.log", structStr)
	//
	//log.Info("", structStr, err)
	//
	//start,end, err := structSelection(node)
	//
	//log.Info("", start, end, err)



	for _, field := range nodeType.Fields.List {
		if util.GetFieldTag(field, "sql").Name == "-" {
			continue
		}

		switch t := field.Type.(type) {
		case *ast.Ident:
			if t.Obj != nil && t.Obj.Kind == ast.Typ {
				if typeSpec, ok := t.Obj.Decl.(*ast.TypeSpec); ok {
					fields = append(fields, ms.getStructFieds(typeSpec.Type)...)
				}
			} else {
				fields = append(fields, field)
			}
		case *ast.SelectorExpr:
			fields = append(fields, field)
		default:
			fields = append(fields, field)
			//log.Warning("filed %s not supported, ignore", util.GetFieldName(field))
		}
	}

	return fields
}

func (ms *SqlGenerator) tableName() string {
	return casee.ToSnakeCase(ms.structName)
}

func generateSqlTag(field *ast.Field) (string, error) {
	var sqlType string
	var err error

	tagStr := util.GetFieldTag(field, "sql").Name
	log.Info("sql tag string %+v", tagStr)
	sqlSettings := ParseTagSetting(tagStr)

	if value, ok := sqlSettings["TYPE"]; ok {
		sqlType = value
	}

	if _, found := sqlSettings["NOT NULL"]; !found { // default: not null
		sqlSettings["NOT NULL"] = "NOT NULL"
	}

	additionalType := sqlSettings["NOT NULL"] + " " + sqlSettings["UNIQUE"]
	if value, ok := sqlSettings["DEFAULT"]; ok {
		additionalType = additionalType + " DEFAULT " + value
	}

	if sqlType == "" {
		var size = 128

		if value, ok := sqlSettings["SIZE"]; ok {
			size, _ = strconv.Atoi(value)
		}

		_, autoIncrease := sqlSettings["AUTO_INCREMENT"]
		if isPrimaryKey(field) {
			autoIncrease = true
		}

		sqlType, err = mysqlTag(field, size, autoIncrease)
		if err != nil {
			log.Warning("get mysql field tag failed:%v", err)
			return "", err
		}
	}

	if strings.TrimSpace(additionalType) == "" {
		return sqlType, nil
	} else {
		return fmt.Sprintf("%v %v", sqlType, additionalType), nil
	}

}

func getColumnName(field *ast.Field) string {
	tagStr := util.GetFieldTag(field, "gorm").Name
	gormSettings := ParseTagSetting(tagStr)
	if columnName, ok := gormSettings["COLUMN"]; ok {
		return columnName
	}

	if len(field.Names) > 0 {
		return fmt.Sprintf("`%s`", casee.ToSnakeCase(field.Names[0].Name))
	}

	return ""
}

func isPrimaryKey(field *ast.Field) bool {
	tagStr := util.GetFieldTag(field, "gorm").Name
	gormSettings := ParseTagSetting(tagStr)
	if _, ok := gormSettings["PRIMARY_KEY"]; ok {
		return true
	}

	if len(field.Names) > 0 && strings.ToUpper(field.Names[0].Name) == "ID" {
		return true
	}

	return false
}

func isTimeKey(field *ast.Field) bool {
	if len(field.Names) > 0 && (strings.ToUpper(field.Names[0].Name) == "CTIME"  || strings.ToUpper(field.Names[0].Name) == "MTIME"){
		return true
	}
	return false
}


func mysqlTag(field *ast.Field, size int, autoIncrease bool) (string, error) {
	typeName := ""
	switch t := field.Type.(type) {
	case *ast.Ident:
		typeName = t.Name
	case *ast.SelectorExpr:
		typeName = t.Sel.Name
	default:
		return "", errors.New(fmt.Sprintf("field %s not supported", util.GetFieldName(field)))
	}

	switch typeName {
	case "bool":
		return "boolean", nil
	case "int", "int8", "int16", "int32", "uint", "uint8", "uint16", "uint32", "uintptr":
		if autoIncrease {
			return "int AUTO_INCREMENT", nil
		}
		return "int", nil
	case "int64", "uint64":
		if autoIncrease {
			return "bigint AUTO_INCREMENT", nil
		}
		return "bigint", nil
	case "float32", "float64":
		return "double", nil
	case "string", "NullString":
		if size > 0 && size < 65532 {
			return fmt.Sprintf("varchar(%d)", size), nil
		}
		return "longtext", nil
	case "Time":
		return "datetime", nil
	default:
		return "", errors.New(fmt.Sprintf("type %s not supported", typeName))

	}
}

func ParseTagSetting(str string) map[string]string {
	tags := strings.Split(str, ";")
	setting := map[string]string{}
	for _, value := range tags {
		v := strings.Split(value, ":")
		k := strings.TrimSpace(strings.ToUpper(v[0]))
		if len(v) == 2 {
			setting[k] = v[1]
		} else {
			setting[k] = k
		}
	}
	return setting
}
