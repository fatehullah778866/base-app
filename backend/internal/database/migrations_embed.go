package database

import "embed"

//go:embed ../../migrations/*.up.sql
var embeddedMigrations embed.FS
