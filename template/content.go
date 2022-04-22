package template

import "embed"

//go:embed init/workflows/my-workflow/ods/ods.sql
//go:embed init/workflows/my-workflow/workflow.yaml
//go:embed init/workflows/my-workflow/dm/dm.sql
//go:embed init/workflows/my-workflow/dw/dim.sql
//go:embed init/workflows/my-workflow/dw/fact.sql
//go:embed init/workflows/my-workflow/raw/ingestion.py
//go:embed init/requirements-glue.txt
//go:embed init/.gitignore
//go:embed init/vendor/__init__.py
//go:embed init/vendor/helper.py

var Content embed.FS
