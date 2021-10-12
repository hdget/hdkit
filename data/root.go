package data

import (
	"embed"
)

type EmbedFs struct {
	Fs  embed.FS
	Dir string
}

var (
	//go:embed script/*
	scriptFs embed.FS
	Script   = EmbedFs{
		Fs:  scriptFs,
		Dir: "script",
	}

	//go:embed proto/*
	protoFs embed.FS
	Proto   = EmbedFs{
		Fs:  protoFs,
		Dir: "proto",
	}

	//go:embed setting/*
	settingFs embed.FS
	Setting   = EmbedFs{
		Fs:  settingFs,
		Dir: "setting",
	}

	//go:embed message/install_compiler.txt
	MsgInstallProtoc string
	//go:embed message/setup_windows.txt
	MsgWinSetup string
	//go:embed message/setup_linux.txt
	MsgLinuxSetup string
)
