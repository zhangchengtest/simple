package params

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zhangchengtest/simple/sqls"
	"strconv"
	"strings"
	"time"

	"github.com/iris-contrib/schema"
	"github.com/zhangchengtest/simple/common/dates"
	"github.com/zhangchengtest/simple/common/strs"
)

var (
	decoder = schema.NewDecoder() // form, url, schema.
)

func init() {
	decoder.AddAliasTag("form", "json")
	decoder.ZeroEmpty(true)
}

// param error
func paramError(name string) error {
	return errors.New(fmt.Sprintf("unable to find param value '%s'", name))
}

// ReadForm read object from FormData
func ReadForm(ctx *gin.Context, obj interface{}) error {
	values := ctx.Request.PostForm
	if len(values) == 0 {
		return nil
	}
	return decoder.Decode(obj, values)
}

func PostForm(ctx *gin.Context, name string) string {
	return ctx.PostForm(name)
}

func FormValueRequired(ctx *gin.Context, name string) (string, error) {
	str := PostForm(ctx, name)
	if len(str) == 0 {
		return "", errors.New("参数：" + name + "不能为空")
	}
	return str, nil
}

//func FormValueDefault(ctx *gin.Context, name, def string) string {
//	return ctx.FormValueDefault(name, def)
//}

func FormValueInt(ctx *gin.Context, name string) (int, error) {
	str := ctx.PostForm(name)
	if str == "" {
		return 0, paramError(name)
	}
	return strconv.Atoi(str)
}

func FormValueIntDefault(ctx *gin.Context, name string, def int) int {
	if v, err := FormValueInt(ctx, name); err == nil {
		return v
	}
	return def
}

func FormValueInt64(ctx *gin.Context, name string) (int64, error) {
	str := ctx.PostForm(name)
	if str == "" {
		return 0, paramError(name)
	}
	return strconv.ParseInt(str, 10, 64)
}

func FormValueInt64Default(ctx *gin.Context, name string, def int64) int64 {
	if v, err := FormValueInt64(ctx, name); err == nil {
		return v
	}
	return def
}

func FormValueInt64Array(ctx *gin.Context, name string) []int64 {
	str := ctx.PostForm(name)
	if str == "" {
		return nil
	}
	ss := strings.Split(str, ",")
	if len(ss) == 0 {
		return nil
	}
	var ret []int64
	for _, v := range ss {
		item, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			continue
		}
		ret = append(ret, item)
	}
	return ret
}

func FormValueStringArray(ctx *gin.Context, name string) []string {
	str := ctx.PostForm(name)
	if len(str) == 0 {
		return nil
	}
	ss := strings.Split(str, ",")
	if len(ss) == 0 {
		return nil
	}
	var ret []string
	for _, s := range ss {
		s = strings.TrimSpace(s)
		if len(s) == 0 {
			continue
		}
		ret = append(ret, s)
	}
	return ret
}

func FormValueBool(ctx *gin.Context, name string) (bool, error) {
	str := ctx.PostForm(name)
	if str == "" {
		return false, paramError(name)
	}
	return strconv.ParseBool(str)
}

func FormValueBoolDefault(ctx *gin.Context, name string, def bool) bool {
	str := ctx.PostForm(name)
	if str == "" {
		return def
	}
	value, err := strconv.ParseBool(str)
	if err != nil {
		return def
	}
	return value
}

// 从请求中获取日期
func FormDate(ctx *gin.Context, name string) *time.Time {
	value := PostForm(ctx, name)
	if strs.IsBlank(value) {
		return nil
	}
	layouts := []string{dates.FmtDateTime, dates.FmtDate, dates.FmtDateTimeNoSeconds}
	for _, layout := range layouts {
		if ret, err := dates.Parse(value, layout); err == nil {
			return &ret
		}
	}
	return nil
}

func GetPaging(ctx *gin.Context) *sqls.Paging {
	page := FormValueIntDefault(ctx, "page", 1)
	limit := FormValueIntDefault(ctx, "limit", 20)
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	return &sqls.Paging{Page: page, Limit: limit}
}
