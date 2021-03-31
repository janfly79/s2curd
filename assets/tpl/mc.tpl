package xxx

import (
    "context"
    "encoding/json"
    "fmt"

    "sniper/util/cachekey"
    "sniper/util/mc"
)

const (
    // Fix Me
    {{.StructTableName}}Key = "xxxxx"
)

func init(){
    cachekey.Register(cachekey.KeyInfo{
    		// Fix Me
    		Name: {{.StructTableName}}Key,
    		Doc:  "",
    	})
}

func Set{{.StructTableName}}Cache(ctx context.Context, id interface{})(err error){
    c := mc.Get(ctx, "default")
    key := fmt.Sprint({{.StructTableName}}Key)

    b, err := json.Marshal(id)
    if err != nil {
        return
    }

    err = c.Set(ctx, &mc.Item{
        Key:        key,
        Value:      b,
        Expiration: 300, // 缓存 5 分钟
    })

    return
}

func Get{{.StructTableName}}Cache(ctx context.Context, id interface{})(err error){
    // Fix Me
    c := mc.Get(ctx, "default")
	key := "xxx"

	_, err = c.Get(ctx, key)
	if err != nil {
		return
	}
    return
}


func Del{{.StructTableName}}Cache(ctx context.Context, id interface{})(err error){
    // Fix Me
    c := mc.Get(ctx, "default")
    key := "xxxx"
    err = c.Delete(ctx, key)
    return
}
