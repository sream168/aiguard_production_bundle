# Delivery Notes (2026-03-26)

## Completed updates

1. Converted all Chinese text in `scripts/` to English to avoid console encoding issues on Windows.
2. Updated `github.com/wailsapp/wails/v2` in `go.mod` to `v2.11.0`.
3. Added OpenAI proxy configuration:
   - `openai.proxy.enabled`
   - `openai.proxy.url`
   - `openai.proxy.http`
   - `openai.proxy.https`
   - `openai.proxy.no_proxy`
4. Added runtime log support:
   - default log file path: `workspace/logs/aiguard.log`
   - UI button to view real-time logs
5. Added cache cleanup support:
   - clears repos, worktrees, cache, reports and logs
6. Reworked repository parsing and clone fallback strategy:
   - supports GitLab `/-/merge_requests/{id}`
   - supports nested groups/subgroups
   - prefers SSH clone URL, falls back to HTTPS automatically
   - SSH / HTTPS host and port are configurable independently
7. Split workflow into two actions:
   - `PullCode`
   - `StartReview`
   `StartReview` now validates both branches and LLM connectivity before review starts.
8. Added Windows process execution option to hide flashing console windows during git operations.
9. Updated desktop layout to support the new buttons and log panel cleanly.
10. Preserved the original review pipeline logic as much as possible while extending configuration and UI behavior.

## Config examples

See:
- `examples/config.yaml`
- `examples/config.production.yaml`

## Recommended GitLab config pattern

```yaml
git:
  preferred_protocol: "ssh"
  gitlab:
    ssh:
      host: "ssh.gitlab.example.com"
      port: "2222"
      user: "git"
    https:
      scheme: "https"
      host: "gitlab.example.com"
      port: "8443"
```

## Offline verification completed in this environment

The execution environment used for this delivery does not have access to Go module or npm package registries, so full Wails build and end-to-end packaging could not be executed here.

Completed checks:
- `gofmt` on updated Go files
- `go test ./internal/gitops`
- `go test ./internal/findings ./internal/history ./internal/scanner ./internal/projectctx ./internal/report ./internal/gitops`
- script scan confirms `scripts/` contains no Chinese characters
- config example scan confirms proxy and git host/port blocks are present

## Remaining local build step on a networked machine

On a machine that can access Go and npm registries, run:

```bash
cd frontend
npm install
npm run build
cd ..
go mod tidy
wails build
```

Windows release packaging:

```powershell
./scripts/build_windows.ps1
```
