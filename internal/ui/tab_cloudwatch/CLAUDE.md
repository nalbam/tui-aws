# CloudWatch Tab
## Role
CloudWatch Alarms — list alarms with state, metric, namespace, threshold.
## Key Files
- `model.go` — CWModel implementing TabModel, 4-state viewState (table/search/actionMenu/detail)
- `table.go` — State column: OK=green, ALARM=red, INSUFFICIENT_DATA=yellow
- `detail.go` — Full alarm detail: metric, namespace, dimensions, threshold, comparison, actions (OK/Alarm/InsufficientData)
## Rules
- Uses `cloudwatch.Client` (DescribeAlarms)
- State colors match AWS Console convention (green/red/yellow)
- Dimensions rendered as Key=Value pairs
- Actions grouped by trigger type (OK actions, Alarm actions, InsufficientData actions)
- Action menu offers only "Alarm Details"
- Search filters on alarm name and metric name
- Overlay closes on Esc; standard 4-state viewState flow
