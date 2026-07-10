# go-app-shared

The versioned cross-service contracts shared by [`auth`](https://github.com/guille1988/auth), [`email`](https://github.com/guille1988/email), and [`broadcasting`](https://github.com/guille1988/broadcasting) — the Kafka DTOs/routing keys and the gRPC proto definitions that flow between them.

This module has exactly one job: make sure that if a service changes the shape of a contract (a Kafka event or a gRPC message), every other service compiled against it **fails to compile** instead of silently breaking at runtime.

---

## What's in here

```text
messaging/kafka/
├── dtos/
│   ├── welcome_email.go     # published by auth on user.created, consumed by email
│   ├── user_logged_in.go    # published by auth on user.logged_in, consumed by broadcasting
│   └── stress_email.go      # synthetic load payload, published by auth/email's /api/stress
└── constants/
    └── routing_key.go       # the Kafka topic names, defined once
rpc/
└── auth/v1/
    ├── auth.proto           # AuthService gRPC contract (served by auth, called by broadcasting)
    ├── auth.pb.go           # generated — do not edit; regenerate with `make proto`
    └── auth_grpc.pb.go      # generated — do not edit; regenerate with `make proto`
```

The layout under `rpc/` is `rpc/<owning-service>/<version>/`: each service that exposes RPCs owns its own proto package (e.g. a future `rpc/email/v1/`), and the generated Go code is committed so consumers build without needing protoc.

---

## How it's consumed

Each of the three services checks this repository out as a **git submodule**, nested at `<service>/internal/shared`. All three point at the same commit of this repository, which the root `go-app` repo's Makefile enforces:

```bash
make check-shared-drift   # fails if the 3 services aren't on the same commit of this repo
make sync-shared FROM=auth  # propagate a change made in one service to the other two
```

Each service's `go.mod` currently resolves this module via a local `replace` directive against that submodule checkout, rather than a tagged, remotely-fetched version:

```
require github.com/guille1988/go-app-shared v0.0.0
replace github.com/guille1988/go-app-shared => ./internal/shared
```

This keeps the workflow simple for a 3-service system with a single maintainer, at the cost of relying on the Makefile (rather than Go's module resolution) to guarantee the three checkouts never drift. Migrating to a real tagged dependency (`go get github.com/guille1988/go-app-shared@v0.1.0`, no `replace`) is a natural next step if this system grows to more services or more contributors.

---

## Adding a new event

1. Add the DTO here, in `messaging/kafka/dtos/`.
2. Add its routing key to `messaging/kafka/constants/routing_key.go`.
3. Push this repo, then run `make sync-shared FROM=<service>` from `go-app` to propagate the new commit to the other two services' submodule checkouts.
4. Register the DTO on the publishing side and the handler on the consuming side (see the "Messaging" section in each service's own README).

## Adding or changing a gRPC contract

1. Edit (or add) the `.proto` under `rpc/<owning-service>/<version>/` in **auth's** checkout of this repo (`microservices/auth/internal/shared`).
2. From the `go-app` root, run `make proto` — it runs protoc in docker with pinned plugin versions and regenerates the `*.pb.go` files in place.
3. Commit the `.proto` together with the regenerated files, then propagate with `make sync-shared FROM=auth`.
4. Enum style note: enums are nested inside the message that uses them so their values don't need an enum-name prefix (top-level proto3 enum values share the package scope and would collide).
