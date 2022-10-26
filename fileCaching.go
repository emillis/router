package veryFastRouter

import (
	"io/fs"
	"path/filepath"
)

//===========[STATIC]====================================================================================================

var defaultFileCaching = FileCaching{
	MinSize:                0,
	MaxSize:                10000,
	DefaultExtensionStatus: false,
	ByExtension:            make(map[string]bool),
	ByPattern:              make(map[string]bool),
}

//===========[STRUCTS]====================================================================================================

//FileCaching allows the user to define the list of conditions to be matched in
//order for a file to be cached in-memory
type FileCaching struct {
	//Defines minimum file size in bytes that's allowed to be cached
	MinSize int

	//Defines maximum file size in bytes that's allowed to be cached
	MaxSize int

	//DefaultExtensionStatus allows to set what to do with all extensions by default.
	//If not set, it will refuse to cache all extensions, then you can exclude some
	//using "byExtension" field.
	DefaultExtensionStatus bool

	//ByExtension allows defining file types to be or not to be cached
	ByExtension map[string]bool

	//TODO: implement!
	ByPattern map[string]bool
}

func (fc *FileCaching) ToCache(f fs.FileInfo) bool {
	if int(f.Size()) < fc.MinSize {
		return false
	}

	if int(f.Size()) > fc.MaxSize {
		return false
	}

	allowExt, exist := fc.ByExtension[filepath.Ext(f.Name())]
	if (!exist && !fc.DefaultExtensionStatus) || (exist && !allowExt) {
		return false
	}

	return true
}

//copyStruct makes a perfect copy of FileCaching struct
func (fc *FileCaching) copyStruct() *FileCaching {
	newFc := FileCaching{
		MinSize:                fc.MinSize,
		MaxSize:                fc.MaxSize,
		DefaultExtensionStatus: fc.DefaultExtensionStatus,
		ByExtension:            make(map[string]bool),
		ByPattern:              make(map[string]bool),
	}

	for k, v := range fc.ByExtension {
		newFc.ByExtension[k] = v
	}

	for k, v := range fc.ByPattern {
		newFc.ByPattern[k] = v
	}

	return &newFc
}

//Copy returns FileCaching as an interface
func (fc *FileCaching) Copy() FileCacher {
	return fc.copyStruct()
}

//===========[FUNCTIONALITY]====================================================================================================

//NewFileCaching returns FileCaching struct with the default values
func NewFileCaching() *FileCaching {
	return defaultFileCaching.copyStruct()
}
