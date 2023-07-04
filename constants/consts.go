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

const (
	LoginTokenLength = 32
)

var (
	WanZBlogConfigFilePath = path.Join(
		WanZBlogConfigPath,
		fmt.Sprintf("%s.%s", WanZBlogConfigName, WanZBlogConfigType),
	)
)
