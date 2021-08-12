package g

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Songmu/prompter"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

// KitFs wraps an afero.Fs
type KitFs struct {
	Fs afero.Fs
}

// global kitfs
var gKitFs *KitFs

// NewKitFs creates a KitFs with `dir` as root.
func NewKitFs(rootDir string) *KitFs {
	// get project base filesystem
	var inFs afero.Fs
	if viper.GetBool("testing") {
		inFs = afero.NewMemMapFs()
	} else {
		cliFolder := viper.GetString("folder")
		if cliFolder != "" {
			inFs = afero.NewBasePathFs(afero.NewOsFs(), cliFolder)
		} else {
			inFs = afero.NewOsFs()
		}
	}

	fs := inFs
	if rootDir != "" {
		fs = afero.NewBasePathFs(inFs, rootDir)
	}

	return &KitFs{Fs: fs}
}

// GetFs returns a new KitFs if it was not initiated before or
// it returns the existing gKitFs if it is initiated.
func GetFs() *KitFs {
	if gKitFs == nil {
		gKitFs = NewKitFs("")
	}
	return gKitFs
}

// ReadFile reads the file from `path` and returns the content in string format
// or returns an error if it occurs.
func (kf *KitFs) ReadFile(path string) ([]byte, error) {
	return afero.ReadFile(kf.Fs, path)
}

func (kf *KitFs) IsDir(path string) (bool, error) {
	return afero.IsDir(kf.Fs, path)
}

// ReadFile reads the file from `path` and returns the content in string format
// or returns an error if it occurs.
//func (f *KitFs) ReadFile(path string) ([]byte, error) {
//	return afero.ReadFile(f.Fs, path)
//}

// WriteFile writs a file to the `path` with `data` as content, if `force` is set
// to true it will override the file if it already exists.
func (kf *KitFs) WriteFile(path string, data []byte, force bool) error {
	if exists, _ := kf.Exists(path); exists && !force {
		bs, _ := kf.ReadFile(path)
		if bytes.Equal(bs, data) {
			return nil
		}

		yes := prompter.YN(fmt.Sprintf("`%s` already exists do you want to override it ?", path), false)
		if !yes {
			return nil
		}
	}
	return afero.WriteFile(kf.Fs, path, data, os.ModePerm)
}

// Exists returns true,nil if the dir/file exists or false,nil if
// the dir/file does not exist, it will return an error if something
// went wrong.
func (kf *KitFs) Exists(path string) (bool, error) {
	return afero.Exists(kf.Fs, path)
}

// MakeDir create dif if not exists
func (kf *KitFs) MakeDir(path string) error {
	exists, err := kf.Exists(path)
	if err != nil {
		return err
	}

	if !exists {
		return kf.Fs.MkdirAll(path, os.ModePerm)
	}
	return nil
}

func (kf *KitFs) Copy(sourcePath, destPath string) error {
	data, err := ioutil.ReadFile(sourcePath)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(destPath, data, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
