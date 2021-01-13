package sqlorm

import (
	"errors"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/janfly79/s2curd/gencode"
	"github.com/janfly79/s2curd/program"

	"github.com/janfly79/s2curd/util/writefile"

	log "github.com/liudanking/goutil/logutil"
	"github.com/urfave/cli"
)

func SqlCommand() cli.Command {
	return cli.Command{
		Name:  "sql",
		Usage: "generate sql from golang model struct",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "file, f",
				Usage: "source file or dir, default: current dir",
			},
			&cli.StringFlag{
				Name:  "struct, s",
				Usage: "struct name or pattern: https://golang.org/pkg/path/filepath/#Match",
			},
			&cli.StringFlag{
				Name:  "out, o",
				Usage: "output file",
			},
		},
		Action: SqlCommandAction,
	}
}

func CurdCommand() cli.Command {
	return cli.Command{
		Name:  "curd",
		Usage: "generate curd from golang model struct",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "file, f",
				Usage: "source file or dir, default: current dir",
			},
			&cli.StringFlag{
				Name:  "struct, s",
				Usage: "struct name or pattern: https://golang.org/pkg/path/filepath/#Match",
			},
		},
		Action: CurdCommandAction,
	}
}

func SqlCommandAction(c *cli.Context) error {
	file := c.String("file")
	if file == "" {
		file, _ = os.Getwd()
	}
	fi, err := os.Stat(file)
	if err != nil {
		log.Warning("get file info [%s] failed:%v", file, err)
		return err
	}

	pattern := c.String("struct")
	if pattern == "" {
		return errors.New("struct is empty")
	}

	out := c.String("out")
	if out == "" {
		return errors.New("output file is empty")
	}

	matchFunc := func(structName string) bool {
		match, _ := filepath.Match(pattern, structName)
		return match
	}

	var types []*ast.TypeSpec
	if !fi.IsDir() {
		fset := token.NewFileSet()
		data, err := ioutil.ReadFile(file)
		if err != nil {
			log.Warning("read [file:%s] failed:%v", file, err)
			return err
		}
		f, err := parser.ParseFile(fset, file, string(data), parser.ParseComments)
		if err != nil {
			log.Warning("parse [file:%s] failed:%v", file, err)
			return err
		}
		types = program.FindMatchStruct([]*ast.File{f}, matchFunc)
	} else {
		absPath, err := gencode.AbsPath(file)
		if err != nil {
			log.Warning("get [path:%s] absPath failed:%v", file, err)
			return err
		}
		srcPkg, err := build.ImportDir(absPath, build.IgnoreVendor)
		if err != nil {
			log.Warning("get package [%s] info failed:%v", absPath, err)
			return err
		}

		prog, err := program.NewProgram([]string{srcPkg.ImportPath})
		if err != nil {
			log.Warning("new program failed:%v", err)
			return err
		}
		pi, err := prog.GetPkgByName(srcPkg.ImportPath)
		if err != nil {
			log.Warning("get package [%s] failed:%v", srcPkg.ImportPath, err)
			return err
		}
		types = program.FindMatchStruct(pi.Files, matchFunc)
	}

	log.Info("get %d matched struct", len(types))

	sqls := []string{}
	for _, typ := range types {
		ms, err := NewSqlGenerator(typ)
		if err != nil {
			log.Warning("create model struct failed:%v", err)
			return err
		}

		sql, err := ms.GetCreateTableSql()
		if err != nil {
			log.Warning("generate sql failed:%v", err)
			return err
		}

		sqls = append(sqls, sql)
	}

	return ioutil.WriteFile(out, []byte(strings.Join(sqls, "\n\n")), 0666)
}

func CurdCommandAction(c *cli.Context) error {
	file := c.String("file")
	if file == "" {
		file, _ = os.Getwd()
	}
	fi, err := os.Stat(file)
	if err != nil {
		log.Warning("get file info [%s] failed:%v", file, err)
		return err
	}

	pattern := c.String("struct")
	if pattern == "" {
		return errors.New("struct is empty")
	}

	matchFunc := func(structName string) bool {
		match, _ := filepath.Match(pattern, structName)
		return match
	}

	var types []*ast.TypeSpec
	if !fi.IsDir() {
		fset := token.NewFileSet()
		data, err := ioutil.ReadFile(file)
		if err != nil {
			log.Warning("read [file:%s] failed:%v", file, err)
			return err
		}
		f, err := parser.ParseFile(fset, file, string(data), parser.ParseComments)
		if err != nil {
			log.Warning("parse [file:%s] failed:%v", file, err)
			return err
		}
		types = program.FindMatchStruct([]*ast.File{f}, matchFunc)

		//if strings.LastIndex(string(data), "Scan") > 0 || strings.LastIndex(string(data), "db.Row") > 0 {
		//	return errors.New("curd is already create")
		//}

	} else {
		absPath, err := gencode.AbsPath(file)
		log.Info("abspath %s",  absPath)
		if err != nil {
			log.Warning("get [path:%s] absPath failed:%v", file, err)
			return err
		}
		srcPkg, err := build.ImportDir(absPath, build.IgnoreVendor)
		if err != nil {
			log.Warning("get package [%s] info failed:%v", absPath, err)
			return err
		}

		log.Info("import path %s", srcPkg.ImportPath)

		prog, err := program.NewProgram([]string{srcPkg.ImportPath})
		if err != nil {
			log.Warning("new program failed:%v", err)
			return err
		}
		pi, err := prog.GetPkgByName(srcPkg.ImportPath)
		if err != nil {
			log.Warning("get package [%s] failed:%v", srcPkg.ImportPath, err)
			return err
		}
		types = program.FindMatchStruct(pi.Files, matchFunc)
	}

	log.Info("get %d matched struct", len(types))

	//writefile.AddImportModule("context", file)

	for _, typ := range types {
		log.Info("types %+v", typ)
		ms, err := NewSqlGenerator(typ)
		log.Info("NewSqlGenerator types %+v", *ms)
		if err != nil {
			log.Warning("create model struct failed:%v", err)
			return err
		}

		str, err := ms.AddFuncStr()

		if err != nil {
			log.Warning("create curd string failed:%v", err)
			return err
		}

		//log.Info(str)

		err = writefile.WriteAppendFile(file, "\n\n\n"+str)

		if err != nil {
			log.Warning("write curd file failed:%v", err)
			return err
		}

		writefile.Gofmt(file)
	}

	return nil
}
