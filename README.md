# go-app-shared

The versioned Kafka contract shared by [`auth`](https://github.com/guille1988/auth), [`email`](https://github.com/guille1988/email), and [`broadcasting`](https://github.com/guille1988/broadcasting) — the DTOs and routing keys that flow between them.

This module has exactly one job: make sure that if `auth` changes the shape of an event, every consumer of that event **fails to compile** instead of silently breaking at runtime on a JSON field mismatch.

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
```

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
