# Academy Seeder

Seeder uses the application's `src/.env`, `src/config/config.yaml`, Mongo driver,
and domain models. Run commands from `src/`.

## Commands

```bash
go run ./cmd/seeder validate --all
go run ./cmd/seeder migrate up
go run ./cmd/seeder seed --all
go run ./cmd/seeder seed --only go-programming
go run ./cmd/seeder clear --only go-programming
go run ./cmd/seeder clear --all
go run ./cmd/seeder migrate down
```

Available topic keys:

- `go-programming`
- `backend-go`
- `modern-js-ts`
- `react-engineering`
- `python-automation-data`
- `devops-cicd`
- `database-engineering`
- `system-design`
- `cloud-native-microservices`
- `web-security`

## Safety

- `seed` validates data and upserts stable ObjectIDs, so it can be run repeatedly.
- `clear` reads `seeder_runs` manifests and deletes only IDs owned by the selected seed dataset.
- `migrate up` refuses to create a unique index when duplicate data already exists.
- `migrate down` removes only seeder-managed indexes and validators; it does not delete data.
- Authors are selected deterministically from a seed author pool, so topics have varied authors without changing between runs.
- Every lesson has lesson-specific concepts, production scenario, pitfalls, exercise, code, media, and quiz facts.
- Dataset validation rejects short lesson content and duplicated content or quiz questions.

The opt-in integration test creates and drops its own uniquely named database:

```bash
MONGO_TEST_URI=mongodb://localhost:27017 go test ./seeders/runner -run Integration -v
```
