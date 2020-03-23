package util

import (
	"encoding/json"
	"fmt"

	"github.com/tidwall/pretty"
)

// JSON 屏蔽错误返回json字节
func JSON(data interface{}) []byte {
	bts, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return bts
}

// JSONWithColor 屏蔽错误返回json格式化带颜色字节
func JSONWithColor(data interface{}) []byte {
	return jsonColor(JSON(data))
}

// PrintJSONWithColor .
func PrintJSONWithColor(data interface{}) {
	fmt.Printf("%s\n", JSONWithColor(data))
}

// jsonColor json字节数组添加格式化及颜色
func jsonColor(bts []byte) []byte {
	return pretty.Color(pretty.Pretty(bts), pretty.TerminalStyle)
}
