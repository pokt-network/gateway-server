package os_gateway

import "embed"

//go:embed db_migrations
var Migrations embed.FS
