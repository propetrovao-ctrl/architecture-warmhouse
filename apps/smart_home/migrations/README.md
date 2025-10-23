# Database Migrations

This directory contains database migrations for the Smart Home API.

## Migration Files

- `001_create_sensors_table.sql` - Creates the main sensors table with indexes
- `002_create_migrations_table.sql` - Creates the migrations tracking table
- `003_add_sensor_constraints.sql` - Adds constraints and validation rules

## Running Migrations

### Using Docker Compose (Recommended)

Migrations are automatically run when starting the application:

```bash
docker-compose up
```

### Manual Migration

#### Using the migration tool:

```bash
# Run migrations
./migrate.sh

# Check migration status
./migrate.sh -c status

# Use custom database URL
./migrate.sh -d "postgres://user:pass@host:port/db"
```

#### Using PowerShell (Windows):

```powershell
# Run migrations
.\migrate.ps1

# Check migration status
.\migrate.ps1 -Command status

# Use custom database URL
.\migrate.ps1 -DatabaseUrl "postgres://user:pass@host:port/db"
```

#### Direct Go execution:

```bash
cd smart_home
go run cmd/migrate/main.go -command=up
go run cmd/migrate/main.go -command=status
```

## Migration System

The migration system:

1. **Tracks applied migrations** in a `migrations` table
2. **Runs migrations in order** based on version numbers
3. **Prevents duplicate execution** of the same migration
4. **Calculates checksums** to detect file changes
5. **Uses transactions** to ensure atomicity

## Adding New Migrations

1. Create a new SQL file with format: `XXX_description.sql`
2. Use sequential numbers for versioning (e.g., `004_add_new_feature.sql`)
3. Include descriptive comments in the SQL file
4. Test the migration locally before committing

## Migration Best Practices

- Always test migrations on a copy of production data
- Never modify existing migration files after they've been applied
- Use `IF NOT EXISTS` clauses where appropriate
- Include rollback instructions in comments
- Keep migrations small and focused on single changes
