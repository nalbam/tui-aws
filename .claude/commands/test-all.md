# Run Full Test Suite

Execute all tests and static analysis for the tui-aws project.

## Steps

1. Run static analysis:
   ```bash
   go vet ./...
   ```

2. Run all tests with verbose output:
   ```bash
   go test ./... -v
   ```

3. If tests fail:
   - Show the failing test name and error
   - Read the test file to understand expected behavior
   - Suggest a fix

4. Report summary:
   - Total tests run
   - Pass/fail count
   - Any vet warnings
