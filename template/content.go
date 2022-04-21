package template

import "embed"

//go:embed init/workflows/orders/ods/ods.sql
//go:embed init/workflows/orders/workflow.yaml
//go:embed init/workflows/orders/dm/dm.sql
//go:embed init/workflows/orders/dw/dim.sql
//go:embed init/workflows/orders/dw/fact.sql
//go:embed init/workflows/orders/raw/ingestion.py
//go:embed init/requirements-glue.txt
//go:embed init/.gitignore
//go:embed init/vendor/__init__.py
//go:embed init/vendor/helper.py

var Content embed.FS
