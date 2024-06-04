package cache

import (
	"fs/src/config"
	"github.com/zyedidia/generic/cache"
	"io/fs"
)

type FileData struct {
	FileBytes  *[]byte
	DirEntries *[]fs.DirEntry
	BaseDir    *string
	IsDir      bool
}

type Cache struct {
	*cache.Cache[string, *FileData]
}

func New(cfg *config.Cache) *Cache {
	return &Cache{
		Cache: cache.New[string, *FileData](cfg.MaxCount),
	}
}

func (c *Cache) Clear() {
	c.Cache = cache.New[string, *FileData](c.Capacity())
}
