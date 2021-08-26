package cmdnew

type ProjectFactory interface {
	createProjectDirs() error
	copyScriptFiles() error
	copySettingFiles() error
	copy3rdProtoFiles() error
	createGoModuleFile() error

	Create() error
}
