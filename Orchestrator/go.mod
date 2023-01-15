module orchestrator

go 1.18

replace arrowhead.eu/common/datamodels => ../common/datamodels

replace arrowhead.eu/common/auth => ../common/auth

replace arrowhead.eu/common/database => ../common/database

replace arrowhead.eu/common/util => ../common/util

require (
	arrowhead.eu/common/util v0.0.0-00010101000000-000000000000
	arrowhead.eu/common/auth v0.0.0-00010101000000-000000000000
	arrowhead.eu/common/database v0.0.0-00010101000000-000000000000
	arrowhead.eu/common/datamodels v0.0.0-00010101000000-000000000000
	github.com/BurntSushi/toml v1.1.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gorilla/mux v1.8.0
)

