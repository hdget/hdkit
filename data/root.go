package data

import "embed"

var (
	//go:embed script/*
	ScriptFs embed.FS
	//go:embed proto/*
	ProtoFs embed.FS
	//go:embed setting/*
	SettingFs embed.FS
	//go:embed message/install_compiler.txt
	MsgInstallProtoc string
	//go:embed message/setup_windows.txt
	MsgWinSetup string
	//go:embed message/setup_linux.txt
	MsgLinuxSetup string
)
