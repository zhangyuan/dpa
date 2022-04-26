package template

import "embed"

//go:embed lakeformation/workflows/my-workflow/ods/ods.sql
//go:embed lakeformation/workflows/my-workflow/workflow.yaml
//go:embed lakeformation/workflows/my-workflow/dm/dm.sql
//go:embed lakeformation/workflows/my-workflow/dw/dim.sql
//go:embed lakeformation/workflows/my-workflow/dw/fact.sql
//go:embed lakeformation/workflows/my-workflow/raw/ingestion.py
//go:embed lakeformation/requirements-glue.txt
//go:embed lakeformation/.gitignore
//go:embed lakeformation/vendor/__init__.py
//go:embed lakeformation/vendor/helper.py

var Content embed.FS
