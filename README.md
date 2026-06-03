# e2b-go-sdk

Go SDK for E2B API, including:
- Python SDK exceptions mapping (`v2.14.1`)
- Sandboxes APIs
- Templates APIs
- Code Interpreter SDK surface (`code_interpreter` Python SDK `v2.5.0`)

Module:
- `github.com/sinhang/e2b-go-sdk`

## Quick start

```go
client := e2b.NewClient()
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
- Compatibility mode is enabled by default (`WithCompatMode(true)`).
- In compatibility mode, selected APIs will fallback on `404`:
  - `StartProcess`: `/process/start` -> `/sandbox/exec` -> `:8089/cube/sandbox/exec`
  - `CreateTemplateV2` / `CreateTemplateV3`: fallback to `POST /templates`

Sandbox command execution:
- `client.Commands(sandboxID).RunSimple(ctx, "echo hello cube")`
- `client.RunCommand(ctx, e2b.RunCommandRequest{SandboxID: sandboxID, Cmd: "echo hello cube"})`
- Returns `Stdout`/`Stderr` directly from the sandbox command response when the deployment exposes them.

## Code Interpreter

The SDK includes a Go surface corresponding to the E2B Python `code_interpreter` SDK (`v2.5.0`).

Implemented Go APIs:
- `e2b.NewCodeInterpreter(client, sandboxID...)`
- `(*CodeInterpreter).Run(ctx, RunCodeRequest)`
- `(*CodeInterpreter).RunSimple(ctx, code, sandboxID)`
- `e2b.ExecuteCode(ctx, client, code, sandboxID)`
- `(*CodeInterpreter).CreateCodeContext(ctx, CreateCodeContextRequest)`
- `(*CodeInterpreter).ListCodeContexts(ctx, sandboxID)`
- `(*CodeInterpreter).RemoveCodeContext(ctx, RemoveCodeContextRequest)`
- `(*CodeInterpreter).RestartCodeContext(ctx, RestartCodeContextRequest)`

Core types:
- `RunCodeRequest`
- `Execution`
- `Result`
- `Logs`
- `OutputMessage`
- `ExecutionError`
- `CodeContext`

Current implementation notes:
- The SDK shape is implemented and unit-tested.
- The default HTTP routes assumed by the SDK are:
  - `POST /sandboxes/{sandboxID}/code-interpreter/run`
  - `POST /sandboxes/{sandboxID}/code-interpreter/contexts`
  - `GET /sandboxes/{sandboxID}/code-interpreter/contexts`
  - `DELETE /sandboxes/{sandboxID}/code-interpreter/contexts/{contextID}`
  - `POST /sandboxes/{sandboxID}/code-interpreter/contexts/{contextID}/restart`
- Your current Cube deployment still lacks a real code execution backend for the `code_interpreter` routes, so these APIs are SDK-complete but cannot yet succeed end-to-end until the server-side routes exist.

Example:

```go
client := e2b.NewClient()
interpreter := e2b.NewCodeInterpreter(client, sandboxID)

execution, err := interpreter.Run(ctx, e2b.RunCodeRequest{
    SandboxID: sandboxID,
    Code:      "print('hello')",
    Language:  "python",
    Timeout:   30,
})
if err != nil {
    panic(err)
}

fmt.Println(execution.Text())
fmt.Println(execution.Stdout)
```

Command execution example:

```go
client := e2b.NewClient()
runner := client.Commands(sandboxID)

resp, err := runner.RunSimple(ctx, "echo hello cube")
if err != nil {
    panic(err)
}

fmt.Println(resp.StdoutText())
```

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

## Cube Compatibility Matrix (2026-05-13)

Probe target: `http://192.168.1.28:3002`  
Rule: direct or SDK compat-fallback available => `cube=是`; otherwise => `cube=否`

| method | api | e2b | cube |
|---|---|---|---|
| GET | /sandboxes | 是 | 是 |
| POST | /sandboxes | 是 | 是 |
| GET | /v2/sandboxes | 是 | 是 |
| GET | /sandboxes/metrics | 是 | 否 |
| GET | /sandboxes/{sandboxID}/logs | 是 | 是 |
| GET | /v2/sandboxes/{sandboxID}/logs | 是 | 是 |
| GET | /sandboxes/{sandboxID} | 是 | 是 |
| DELETE | /sandboxes/{sandboxID} | 是 | 是 |
| GET | /sandboxes/{sandboxID}/metrics | 是 | 否 |
| POST | /sandboxes/{sandboxID}/pause | 是 | 是 |
| POST | /sandboxes/{sandboxID}/resume | 是 | 是 |
| POST | /sandboxes/{sandboxID}/connect | 是 | 是 |
| POST | /sandboxes/{sandboxID}/timeout | 是 | 否 |
| PUT | /sandboxes/{sandboxID}/network | 是 | 否 |
| POST | /sandboxes/{sandboxID}/refresh | 是 | 否 |
| POST | /sandboxes/{sandboxID}/snapshots | 是 | 是 |
| GET | /snapshots | 是 | 否 |
| POST | /v3/templates | 是 | 是 |
| POST | /v2/templates | 是 | 是 |
| GET | /templates/upload-link | 是 | 否 |
| GET | /templates | 是 | 是 |
| POST | /templates | 是 | 是 |
| GET | /templates/{templateID} | 是 | 是 |
| POST | /templates/{templateID} | 是 | 是 |
| DELETE | /templates/{templateID} | 是 | 是 |
| PATCH | /templates/{templateID} | 是 | 否 |
| POST | /templates/{templateID}/build | 是 | 否 |
| POST | /v2/templates/{templateID}/build | 是 | 否 |
| PATCH | /v2/templates/{templateID} | 是 | 否 |
| GET | /templates/{templateID}/builds/{buildID} | 是 | 否 |
| GET | /templates/{templateID}/builds/{buildID}/logs | 是 | 是 |
| GET | /templates/aliases/{alias} | 是 | 否 |
| POST | /templates/{templateID}/tags | 是 | 否 |
| DELETE | /templates/{templateID}/tags | 是 | 否 |
| GET | /templates/{templateID}/tags | 是 | 否 |
| GET | /volumes | 是 | 否 |
| POST | /volumes | 是 | 否 |
| GET | /volumes/{volumeID} | 是 | 否 |
| DELETE | /volumes/{volumeID} | 是 | 否 |
| GET | /envd/health | 是 | 否 |
| GET | /envd/stats | 是 | 否 |
| GET | /envd/envs | 是 | 否 |
| GET | /filesystem/download | 是 | 否 |
| POST | /filesystem/upload | 是 | 否 |
| POST | /filesystem/compose | 是 | 否 |
| POST | /filesystem/createwatcher | 是 | 否 |
| POST | /filesystem/getwatcherevents | 是 | 否 |
| POST | /filesystem/listdir | 是 | 否 |
| POST | /filesystem/makedir | 是 | 否 |
| POST | /filesystem/move | 是 | 否 |
| POST | /filesystem/remove | 是 | 否 |
| POST | /filesystem/removewatcher | 是 | 否 |
| POST | /filesystem/stat | 是 | 否 |
| POST | /filesystem/watchdir | 是 | 否 |
| POST | /process/closestdin | 是 | 否 |
| POST | /process/connect | 是 | 否 |
| POST | /process/list | 是 | 否 |
| POST | /process/sendinput | 是 | 否 |
| POST | /process/sendsignal | 是 | 否 |
| POST | /process/start | 是 | 是 |
| POST | /process/streaminput | 是 | 否 |
| POST | /process/update | 是 | 否 |
| GET | /teams | 是 | 否 |
| GET | /teams/metrics | 是 | 否 |
| GET | /teams/metrics/max | 是 | 否 |

## Verified Runtime Parameters

The following APIs were verified against the local CubeSandbox deployment.

| method | api | verified params |
|---|---|---|
| GET | /sandboxes | `metadata` optional |
| POST | /sandboxes | `templateID` required; `volumeMounts` works when each item uses `name` + `path` |
| GET | /v2/sandboxes | `metadata` optional |
| GET | /sandboxes/{sandboxID} | `sandboxID` path param |
| GET | /sandboxes/{sandboxID}/logs | `sandboxID` path param |
| GET | /v2/sandboxes/{sandboxID}/logs | `sandboxID` path param |
| DELETE | /sandboxes/{sandboxID} | `sandboxID` path param |
| POST | /sandboxes/{sandboxID}/pause | `sandboxID` path param |
| POST | /sandboxes/{sandboxID}/resume | `sandboxID` path param |
| POST | /sandboxes/{sandboxID}/connect | `sandboxID` path param; body needs `timeout` |
| POST | /sandboxes/{sandboxID}/snapshots | `sandboxID` path param; body needs `name` |
| POST | /templates | `templateID`, `image`, `source_image_ref`, `writable_layer_size` required |
| GET | /templates | no required body params |
| GET | /templates/{templateID} | `templateID` path param |
| POST | /templates/{templateID} | `templateID` path param; empty body is accepted for rebuild |
| GET | /templates/{templateID}/builds/{buildID}/logs | `templateID` and `buildID` path params |
| DELETE | /templates/{templateID} | `templateID` path param |
| POST | /process/start | `SandboxID` or `env.E2B_SANDBOX_ID`; `Cmd` or `Args`; `ContainerID` optional and defaults to `SandboxID` |

Verified tests:
- `TestCreateSandbox`
- `TestCreateSandboxWithMountedExec`

Notes:
- `CreateTemplateV2` and `CreateTemplateV3` use the same verified create payload via compat fallback.
- For mounted file execution, the working payload used `volumeMounts: [{"name":"tmp","path":"/workspace"}]` and `StartProcessRequest` with `SandboxID` plus `Args`.


### cube-sandbox
```shell
sudo mount -o loop /xfs.img /data/cubelet
sudo /usr/local/services/cubetoolbox/scripts/one-click/down-local.sh
sudo /usr/local/services/cubetoolbox/scripts/one-click/up.sh
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


cubemastercli -a 127.0.0.1 -p 8089 tpl create-from-image \
    --image cube-sandbox-cn.tencentcloudcr.com/cube-sandbox/sandbox-code:latest \
    --writable-layer-size 1G \
    --expose-port 49996 \
    --expose-port 49980 \
    --template-id code-template

# tpl-b574464e57cf457fa42caa2f

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
[gpt](https://chatgpt.com/c/6a053383-f298-8321-8fed-44cfa3f01138)  
[deepseek](https://chat.deepseek.com/a/chat/s/ad583503-fe28-4708-949e-b3fd7a57219b)

### test
```shell
go test -v -run TestCreateSandbox ./test/
go test -v -run TestListSandbox ./test/
go test -v -run TestCreateTemplateV2 ./test/
go test -v -run TestCreateSandbox ./test/
go test -v -run TestCreateSandboxWithMountedExec ./test/

go test -v -run TestRunValidation ./e2b/
go test -v -run TestCodeInterpreterRunSimple ./e2b/
go test -v -run TestRunCode1 ./test/
```

### tag
```shell
git push -u origin main
git tag v0.0.3
git push origin v0.0.3
```

### python example
阅读：docs/e2b-dev-sidecar.md
dir path: `/mnt/nvme2/develope/develope/code/CubeSandbox/examples/e2b-dev-sidecar`
```shell
source /mnt/nvme2/develope/develope/code/CubeSandbox/examples/.venv/bin/activate
source /mnt/nvme2/develope/develope/code/CubeSandbox/examples/e2b-dev-sidecar/.env
python /mnt/nvme2/develope/develope/code/CubeSandbox/examples/e2b-dev-sidecar/demo.py
```
### 输出
```
python demo.py 
SandboxInfo(sandbox_id='243a708bb9c048bf8ba4bf41fa54a8ae', sandbox_domain='cube.app', template_id='tpl-3a05aafec23c4d928cfa1850', name=None, metadata={'cube.product': 'cubebox', 'cube.numa_node': '0', 'cube.master.runtime.snapshot.attached_at': '2026-06-03T08:10:25.512852494Z', 'cube.master.runtime.restore.snapshot.attached_at': '2026-06-03T08:10:25.512852494Z', 'cube.master.instance.type': 'cubebox', 'cube.master.runtime.snapshot.id': 'tpl-3a05aafec23c4d928cfa1850', 'cube.master.runtime.restore.snapshot.id': 'tpl-3a05aafec23c4d928cfa1850', 'cube.master.appsnapshot.template.id': 'tpl-3a05aafec23c4d928cfa1850', 'X-Caller': 'X-Caller'}, started_at=datetime.datetime(2026, 6, 3, 8, 10, 25, 512907, tzinfo=tzutc()), end_at=datetime.datetime(2026, 6, 3, 8, 10, 25, 512907, tzinfo=tzutc()), state=<SandboxState.RUNNING: 'running'>, cpu_count=2, memory_mb=2000, envd_version='0.2.0', _envd_access_token=None, allow_internet_access=None, network=None, lifecycle=None, volume_mounts=[]) http://127.0.0.1:12580/sandboxes/router/243a708bb9c048bf8ba4bf41fa54a8ae/49983
Hello world Cube！
```