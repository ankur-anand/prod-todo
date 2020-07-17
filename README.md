## What and how of a writing modern production ready web application in Go(Golang).

**IMP: This is currently an work in progress.**

This repository outlines an example todo web application written in Go(Golang), which can be useful for people learning
Go(Golang) to get started with writing web application in Go employing best practices in modern cloud landscape.

> Disclaimer: This repository is not an empirical definition of "good structure" or "good code". What this repository contains is
my opinions about:

1. How to structure the web application in Go that is easier to maintain as well as test while the codebase grows.
2. How to use Interface to advantage to make testing easier and avoid any tight coupling with dependencies.
3. How to programmatically write integration testing for a database with Docker API.
4. How to use Distributed tracing and metrics collection using opentelemetry for modern cloud landscape.
5. Some best practice to follow for logging in application.
6. An example Docker file for docker image creation with go module system.

### Structuring and Interface usage.

```bash
├── cmd
│   └── rest-server
└── pkg
    ├── authstrategy
    ├── observability
    ├── resthandler
    └── storage
```
`cmd` will contain binaries of my application. For example, `cmd/rest-server` `main.go`, that will build binary for
rest server. It will tie all the dependency together.

`pkg` will have all the code to perform all logical operation for my example todo application.

Top level contains code, that just are specific to the domain of the web application for our case.

1. User
2. Todo

```bash
.
├── authstrategy
├── observability
├── resthandler
├── storage
├── todo.go
├── todo_test.go
├── user.go
└── user_test.go
```
`user.go` contains all the business logic that could be essential for user operation in our system.

For example `UserModel` represent a user in the system.
```go
// UserModel represents individual user registered in the system
type UserModel struct {
	ID        uuid.UUID
	Email     string
	Password  string
	FirstName string
	LastName  string
	Username  string
}
```
It also needs to validate logic for user, like duplicate registration, valid email and password. So it will also need
to interact with some storage. example, postgres, mongodb etc.

As this storage's implementation can be database dependent and domain being the consumer of it,
we will define an interface for all storage implementation.

```go
// UserStorage define a contract for storage, to interact
// with the UserModel.
type UserStorage interface {
	Find(ctx context.Context, id uuid.UUID) (UserModel, error)
	FindByEmail(ctx context.Context, email string) (UserModel, error)
	Update(ctx context.Context, user UserModel) error
	Store(ctx context.Context, user UserModel) (uuid.UUID, error)
}
```

Similar approach we will take for other domains of our business. In our example web application we will do it for `todo`.

### Subpackages inside `pkg`

It's grouped by dependency.

```bash
.
├── authstrategy
├── observability
├── resthandler
└── storage
```
1. `authstrategy` currently, contains logic to generate jwt token and validation. You can put other authentication
strategy like OAuth, SAML anything that your application would like to use.

2. `observability` provide functionality for metric, trace and logging.

3. `resthandler` contains all `http.Handler` for providing rest-api capabilities for the web-application.

4. `storage` contains actual implementation against a database.
As there can be different implementation against each particular database, we have split the test into `testsuite` that
contain test case against an interface instead of actual implementation.
`serror` provides errors independent of database used for the consumer.

```bash
├── postgres
│   ├── migrations
│   │   ├── 01_create_user_table.down.sql
│   │   ├── 01_create_user_table.up.sql
│   │   ├── 02_create_todo_table.down.sql
│   │   └── 02_create_todo_table.up.sql
│   ├── sqlquery.go
│   ├── todo.go
│   └── user.go
├── pqsql.go
├── pqsql_integration_test.go
├── serror
│   └── queryerror.go
└── testsuite
    └── user_storage_psql.go
```

### Programmatically writing integration testing for a database with Docker API.

Also instead of mocking the `sql` to test database interactions,
we can take advantage of `docker` api to spin up a database container and run tests against it directly.
Check the `psql_integration_test.go` for implementation details.

### Distributed tracing and Metrics
Distributed tracing is sort of correlated logging. We will use it to gain the visibility into the operation of request,
and database call, for the use cases such as performance profiling, debugging and RCA.
