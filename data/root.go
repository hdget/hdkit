package data

import "embed"

var (
	//go:embed scripts/*
	Scripts embed.FS
	//go:embed install_compiler.txt
	MsgInstallProtoc string
	//go:embed setup_windows.txt
	MsgWinSetup string
	//go:embed setup_linux.txt
	MsgLinuxSetup string
)
