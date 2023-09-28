# branching-gorm-example

> **Warning:** This repo is only for testing. Please don't use it in production.

This repo is a gorm example of CI/CD workflow powered by the TiDB Serverless branching. From this repo, you can learn:

- How to connect to TiDB Serverless in gorm.
- How to use branching GitHub integration.

## About this repo

This repo is based on [gorm playground](https://github.com/go-gorm/playground), with some changes:

- A tidb dialect is added to the repo to test the TiDB Cloud.
- The [gormigrate](https://github.com/go-gormigrate/gormigrate) is used in `RunMigrations` to help migration.
- Delete some useless files like GitHub actions, docker-compose, etc.

## Connect to TiDB Serverless in gorm

> Make sure you have installed the go environment.

1. clone the code

```
git clone git@github.com:tidbcloud/branching-gorm-example.git
cd branching-gorm-example
```

2. Fill in the following environment variable. You can find the information in the TiDB Serverless console.

```
export GORM_ENABLE_CACHE=true
export GORM_DIALECT=tidb
export GORM_DSN="${{ username }}:${{ password }}@tcp(${{ host }}:${{ port }})/test?parseTime=true&tls=tidb"
```

3. Connect to the TiDB Serverless (migration will be executed automatically)

```
./test.sh
```

## Use branching GitHub integration

This repo has been connected to a TiDB Serverless using the [Branching GitHub integration](https://docs.pingcap.com/tidbcloud/branch-github-integration). This brings database branches to your GitHub workflows, and a TiDB Cloud App will automatically manage database branches for you in the pull request.

**CI workflow**

The repo has a [Test GitHub Action](./.github/workflows/tests.yml) to run the test on the created TiDB Serverless branch. This action uses the [wait-for-tidbcloud-branch](https://github.com/tidbcloud/wait-for-tidbcloud-branch) to get branch connection information and pass it by environment variables. We can do it because the repo accepts the `GORM_DSN` environment variable as connection information. See the [code](https://github.com/tidbcloud/branching-gorm-example/blob/9ca6e4037edd25abafc35e1a378fc29ad62b5f68/db.go#L49) for more details.

Check the [pull request](https://github.com/tidbcloud/branching-gorm-example/pulls) to see how we use the CI workflow!

**CD workflow**

The CD workflow works well with native frameworks.

Take DDL as an example, you can use the [gormigrate](https://github.com/go-gormigrate/gormigrate) to manage your database migrations. Any DDL changes can be applied to the production cluster when the PR is merged. Don't worry about the influence of production business, the TiDB Serverless cluster supports online DDL without blocking your business.

