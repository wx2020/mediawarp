package static

import (
	"io/fs"
)

// EmbeddedStaticAssets 内嵌的静态资源文件系统
// 注意：已移除git子模块依赖，当前为空文件系统
// 如需使用Web美化功能，请通过 config 中的 custom 目录自行提供静态资源
var EmbeddedStaticAssets fs.FS = emptyFS{}

// emptyFS 空文件系统实现
type emptyFS struct{}

func (emptyFS) Open(name string) (fs.File, error) {
	return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
}
