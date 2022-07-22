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