package veryFastRouter

import "io/fs"

//===========[STATIC]====================================================================================================

//defaultFileServingOptions define default values that FileServingOptions will come with
var defaultFileServingOptions = FileServingOptions{
	IgnorePatterns:      []string{"_*"},
	AllowDynamicServing: false,
	Filter:              nil,
}

//===========[INTERFACES]====================================================================================================

//FileCacher defines interface for deciding whether a file should be cached
type FileCacher interface {
	ToCache(fs.FileInfo) bool
}

//===========[STRUCTS]====================================================================================================

//FileServingOptions is the default way of defining advanced rules for file serving
type FileServingOptions struct {
	//FileCaching makes in-memory file caching possible for defined file types, sizes, etc..
	FileCaching FileCacher

	//IgnorePatterns allows you to ignore certain paths, individual files or patterns within
	//the defined root file system folder. You can use "*" (wildcard) to generalize the selection
	//of paths to ignore. By default, "_*" is present, to ignore everything that starts with underscore.
	IgnorePatterns []string

	//AllowDynamicServing allows files being served that have been added after the
	//initial indexing of the root file system location. If this is set to false,
	//the newly added files will be ignored.
	AllowDynamicServing bool

	//Filter function will be invoked for each file. This is used for custom file processing logic.
	//If used, it will ultimately decide what to do with the file
	Filter func(file fs.File) bool
}

//Copy will return a perfect copy of FileServingOptions
func (o *FileServingOptions) Copy() FileServingOptions {
	newOpt := FileServingOptions{
		IgnorePatterns:      make([]string, 0, len(o.IgnorePatterns)),
		AllowDynamicServing: o.AllowDynamicServing,
		Filter:              o.Filter,
	}

	for i := 0; i < len(o.IgnorePatterns); i++ {
		newOpt.IgnorePatterns[i] = o.IgnorePatterns[i]
	}

	return newOpt
}

//===========[Functionality]====================================================================================================

//makeFileServingOptionsReasonable will check FileServingOptions for nonsensical values
func makeFileServingOptionsReasonable(o *FileServingOptions) error {
	return nil
}

//NewFileServingOptions will return a newly initiated FileServingOptions struct with default values
func NewFileServingOptions() FileServingOptions {
	return defaultFileServingOptions.Copy()
}
