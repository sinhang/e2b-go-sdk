# e2b-go-sdk

Go SDK for E2B API, including:
- Python SDK exceptions mapping (`v2.14.1`)
- Sandboxes APIs
- Templates APIs

Module:
- `github.com/sinhang/e2b-go-sdk`

## Quick start

```go
client := e2b.NewClient("<E2B_API_KEY>")
ctx := context.Background()

sandboxes, err := client.ListSandboxes(ctx, "")
if err != nil {
    panic(err)
}
_ = sandboxes
```

## Implemented APIs

Sandboxes:
- `GET /sandboxes`
- `POST /sandboxes`
- `GET /v2/sandboxes`
- `GET /sandboxes/metrics`
- `GET /sandboxes/{sandboxID}/logs` (deprecated)
- `GET /v2/sandboxes/{sandboxID}/logs`
- `GET /sandboxes/{sandboxID}`
- `DELETE /sandboxes/{sandboxID}`
- `GET /sandboxes/{sandboxID}/metrics`
- `POST /sandboxes/{sandboxID}/pause`
- `POST /sandboxes/{sandboxID}/resume`
- `POST /sandboxes/{sandboxID}/connect` (deprecated)
- `POST /sandboxes/{sandboxID}/timeout`
- `PUT /sandboxes/{sandboxID}/network`
- `POST /sandboxes/{sandboxID}/refresh`
- `POST /sandboxes/{sandboxID}/snapshots`
- `GET /snapshots`

Templates:
- `POST /v3/templates`
- `POST /v2/templates` (deprecated)
- `GET /templates/upload-link`
- `GET /templates`
- `POST /templates` (deprecated)
- `GET /templates/{templateID}`
- `POST /templates/{templateID}` (rebuild, deprecated)
- `DELETE /templates/{templateID}`
- `PATCH /templates/{templateID}` (deprecated)
- `POST /templates/{templateID}/build` (deprecated)
- `POST /v2/templates/{templateID}/build`
- `PATCH /v2/templates/{templateID}`
- `GET /templates/{templateID}/builds/{buildID}`
- `GET /templates/{templateID}/builds/{buildID}/logs`
- `GET /templates/aliases/{alias}`

## Notes

- The SDK uses `X-API-Key` header for auth.
- Methods accept `url.Values` for query parameters where applicable.
- Response types are either typed structs (`Sandbox`, `Template`) or `map[string]any` for flexible schema compatibility.

Additional APIs implemented:

Tags:
- `POST /templates/{templateID}/tags`
- `DELETE /templates/{templateID}/tags`
- `GET /templates/{templateID}/tags`

Volumes:
- `GET /volumes`
- `POST /volumes`
- `GET /volumes/{volumeID}`
- `DELETE /volumes/{volumeID}`

Envd:
- `GET /envd/health`
- `GET /envd/stats`
- `GET /envd/envs`

Filesystem:
- `GET /filesystem/download`
- `POST /filesystem/upload`
- `POST /filesystem/compose`
- `POST /filesystem/createwatcher`
- `POST /filesystem/getwatcherevents`
- `POST /filesystem/listdir`
- `POST /filesystem/makedir`
- `POST /filesystem/move`
- `POST /filesystem/remove`
- `POST /filesystem/removewatcher`
- `POST /filesystem/stat`
- `POST /filesystem/watchdir`

Process:
- `POST /process/closestdin`
- `POST /process/connect`
- `POST /process/list`
- `POST /process/sendinput`
- `POST /process/sendsignal`
- `POST /process/start`
- `POST /process/streaminput`
- `POST /process/update`

Teams:
- `GET /teams`
- `GET /teams/metrics`
- `GET /teams/metrics/max`


### cube-sandbox
```shell
sublime-text.subl /home/mercury/cube-sandbox/cube-sandbox-one-click-9c16021/.env
sublime-text.subl /usr/local/services/cubetoolbox/CubeMaster/conf.yaml
sudo /home/mercury/cube-sandbox/cube-sandbox-one-click-9c16021/install.sh 
sudo /home/mercury/cube-sandbox/cube-sandbox-one-click-9c16021/down.sh 
sudo /home/mercury/cube-sandbox/cube-sandbox-one-click-9c16021/smoke.sh
sudo /usr/local/services/cubetoolbox/scripts/one-click/up.sh
sudo /usr/local/services/cubetoolbox/scripts/one-click/down-local.sh
# log
/data/log/CubeMaster/


cubemastercli -a 127.0.0.1 -p 8089 tpl create-from-image \
    --image cube-sandbox-cn.tencentcloudcr.com/cube-sandbox/sandbox-browser:latest \
    --writable-layer-size 1G \
    --expose-port 49999 \
    --expose-port 49983
2026/05/13 15:21:56 job_id: 3f494f91-930e-47f4-8243-d29432510364
2026/05/13 15:21:56 template_id: tpl-3a864cb982224e97ac2168b5
2026/05/13 15:21:56 attempt_no: 1
2026/05/13 15:21:56 operation: CREATE
2026/05/13 15:21:56 artifact_id: 
2026/05/13 15:21:56 status: PENDING
2026/05/13 15:21:56 phase: PULLING
2026/05/13 15:21:56 progress: 0%
2026/05/13 15:21:56 distribution: 0/0 ready, 0 failed


```

```shell
cubemastercli -a 127.0.0.1 -p 8089 tpl create-from-imag
- --image
- --template-id
- --writable-layer-size
- --expose-port（可重复）
- --instance-type（默认 cubebox）
- --network-type（默认 tap）
- --node（可重复）
- --allow-internet-access
- --allow-out-cidr（可重复）
- --deny-out-cidr（可重复）
- --registry-username
- --registry-password
- --cmd（可重复）
- --arg（可重复）
- --env（KEY=VALUE，可重复）
- --dns（可重复）
- --probe
- --probe-path（默认 /health）
- --cpu（默认 2000）
- --memory（默认 2000）
- --json
```

### docs
[e2b.dev](https://www.e2b.dev/docs/api-reference/sandboxes/list-sandboxes)

### test
```shell
go test -v -run TestCreateSandbox ./test/
```