# IAM Tab
## Role
IAM Users — list users with groups, policies, creation date, last password used.
## Key Files
- `model.go` — IAMModel implementing TabModel, 4-state viewState (table/search/actionMenu/detail), status bar shows account ID
- `table.go` — Columns: UserName, UserID, ARN, Created, LastUsed
- `detail.go` — Full user detail: groups list, attached policies list, password last used
## Rules
- Uses `iam.Client` (ListUsers, ListGroupsForUser, ListAttachedUserPolicies)
- Account ID fetched via `STS GetCallerIdentity` and displayed in status bar
- IAM is a global service — data is region-independent
- Groups and policies loaded per-user in the detail fetch
- PasswordLastUsed defaults to "Never" when empty
- Groups/policies show "(none)" when empty
- Action menu offers only "User Details"
- Search filters on username and ARN
- Overlay closes on Esc; standard 4-state viewState flow
