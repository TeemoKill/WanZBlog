package constants

import (
	"fmt"
	"path"
)

const (
	WanZBlogConfigName = "wanzblog"
	WanZBlogConfigType = "toml"
	WanZBlogConfigPath = "./"
)

var (
	WanZBlogConfigFilePath = path.Join(
		WanZBlogConfigPath,
		fmt.Sprintf("%s.%s", WanZBlogConfigName, WanZBlogConfigType),
	)
)
