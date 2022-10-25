package veryFastRouter

import (
	"io/fs"
	"path/filepath"
)

//===========[STRUCTS]====================================================================================================

//FileCaching allows the user to define the list of conditions to be matched in
//order for a file to be cached in-memory
type FileCaching struct {
	MinSize int
	MaxSize int
	//ByExtension allows defining file types to be or not to be cached
	ByExtension map[string]bool
	ByPattern   map[string]bool
}

func (fc *FileCaching) ToCache(f fs.FileInfo) bool {
	if !fc.ByExtension[filepath.Ext(f.Name())] {
		return false
	}

	return true
}
