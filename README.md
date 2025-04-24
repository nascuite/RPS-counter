# RPS-counter

Counts requests‑per‑second (RPS) with sync/atomic
with:
- Graceful shutdown on SIGINT / SIGTERM.
- Built‑in pprof endpoints

## Quick start
```
go run main.go
```
## get profile

```bash
go tool pprof -http=:6060 http://localhost:8080/debug/pprof/profile\?seconds\=30
```

## web UI

```bash
go tool pprof -http=:6060 profiles/cpu_profile_30s.pb.gz
```