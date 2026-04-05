# Lambda Tab
## Role
Lambda Functions — list functions with runtime, memory, timeout, VPC config, layers.
## Key Files
- `model.go` — LambdaModel implementing TabModel, 4-state viewState (table/search/actionMenu/detail)
- `table.go` — State: Active=green, Pending=yellow, Inactive=gray, Failed=red
- `detail.go` — Full function detail: handler, code size (formatted bytes), VPC subnets/SGs, layers, env vars
## Rules
- Uses `lambda.Client` (ListFunctions)
- State defaults to "Active" when empty (Lambda returns empty string for active functions)
- VPC config section only rendered when VpcID is non-empty
- Layers rendered with ARN list
- formatBytes helper converts code size to human-readable format
- Action menu offers only "Function Details"
- Search filters on function name and runtime
- Overlay closes on Esc; standard 4-state viewState flow
