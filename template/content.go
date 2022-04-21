package template

import "embed"

//go:embed assets/workflows/orders/ods/ods.sql
//go:embed assets/workflows/orders/workflow.yaml
//go:embed assets/workflows/orders/dm/dm.sql
//go:embed assets/workflows/orders/dw/dim.sql
//go:embed assets/workflows/orders/dw/fact.sql
//go:embed assets/workflows/orders/raw/ingestion.py
//go:embed assets/requirements-glue.txt
//go:embed assets/.gitignore
//go:embed assets/vendor/__init__.py
//go:embed assets/vendor/helper.py

var Content embed.FS
