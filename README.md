# Makefile Commands

Run Tests:
```
make test_all
```
This will run tests with the minimum required dependencies then delete all containers on success. On failure, containers will be shut down but not deleted.

Returns exit code of 0 on success, 1 on failure (to be used in CI)

**Note: Need to manually delete after resolving error for a clean DB**

Run Tests in Debug Mode:
```
make test_all_debug
```
This will run tests with minmum required dependencies + a Mongo Express server to inspect DB. Containers will also remain up regardless if tests pass or fail

# Troubleshooting

## Updating `init_db_test_data.js`
Because this file is included as a volume, when updated, need to do the following steps:

1. Kill all containers
2. Run `docker volume prune`
3. Restart containers to refresh the init script