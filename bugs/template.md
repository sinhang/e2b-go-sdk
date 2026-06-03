### 【虚拟环境内部执行】命令行创建
```shell
cubemastercli tpl create-from-image \
--image cube-sandbox-cn.tencentcloudcr.com/cube-sandbox/sandbox-code:latest \
--writable-layer-size 1G \
--expose-port 49999 \
--expose-port 49983 \
--probe 49999
```
创建的模板可以使用。

### 【宿主机】API 请求
```go
result, err := client.CreateTemplateV2(ctx, e2b.JSONMap{
    "image":             "cube-sandbox-cn.tencentcloudcr.com/cube-sandbox/sandbox-code:latest",
    "exposePort":        []int{49999, 49983},
    "probe":             49999,
    "writableLayerSize": "1G",
})
```
**创建的模板不可用。**
`
api 创建的模板
python demo.py
SandboxInfo(sandbox_id='5494701e3b48486a97ff9579570f89f3', sandbox_domain='cube.app', template_id='tpl-698a1ab5376f47e883639fa1', name=None, metadata={'cube.product': 'cubebox', 'X-Caller': 'X-Caller', 'cube.master.appsnapshot.template.id': 'tpl-698a1ab5376f47e883639fa1', 'cube.master.runtime.restore.snapshot.id': 'tpl-698a1ab5376f47e883639fa1', 'cube.numa_node': '0', 'cube.master.runtime.snapshot.attached_at': '2026-06-03T07:11:39.53498368Z', 'cube.master.runtime.restore.snapshot.attached_at': '2026-06-03T07:11:39.53498368Z', 'cube.master.runtime.snapshot.id': 'tpl-698a1ab5376f47e883639fa1', 'cube.master.instance.type': 'cubebox'}, started_at=datetime.datetime(2026, 6, 3, 7, 11, 39, 535019, tzinfo=tzutc()), end_at=datetime.datetime(2026, 6, 3, 7, 11, 39, 535019, tzinfo=tzutc()), state=<SandboxState.RUNNING: 'running'>, cpu_count=2, memory_mb=2000, envd_version='0.2.0', _envd_access_token=None, allow_internet_access=None, network=None, lifecycle=None, volume_mounts=[]) http://127.0.0.1:12580/sandboxes/router/5494701e3b48486a97ff9579570f89f3/49983
<html>
<head><title>502 Bad Gateway</title></head>
<body>
<center><h1>502 Bad Gateway</h1></center>
<hr><center>openresty</center>
</body>
</html>
: This error is likely due to sandbox timeout. You can modify the sandbox timeout by passing 'timeout' when starting the sandbox or calling '.set_timeout' on the sandbox with the desired timeout.
`
(.venv) mercury@mercury-X99:/mnt/nvme2/develope/develope/code/CubeSandbox/examples/e2b-dev-sidecar$ source .env 

### 命令行创建的模板【可用】
```shell
(.venv) mercury@mercury-X99:/mnt/nvme2/develope/develope/code/CubeSandbox/examples/e2b-dev-sidecar$ python demo.py
SandboxInfo(sandbox_id='c7eccbb18af445d2af2d2ff2a19c4711', sandbox_domain='cube.app', template_id='tpl-3a05aafec23c4d928cfa1850', name=None, metadata={'X-Caller': 'X-Caller', 'cube.master.runtime.restore.snapshot.id': 'tpl-3a05aafec23c4d928cfa1850', 'cube.master.runtime.restore.snapshot.attached_at': '2026-06-03T07:12:31.68806739Z', 'cube.master.runtime.snapshot.attached_at': '2026-06-03T07:12:31.68806739Z', 'cube.product': 'cubebox', 'cube.master.instance.type': 'cubebox', 'cube.master.runtime.snapshot.id': 'tpl-3a05aafec23c4d928cfa1850', 'cube.numa_node': '0', 'cube.master.appsnapshot.template.id': 'tpl-3a05aafec23c4d928cfa1850'}, started_at=datetime.datetime(2026, 6, 3, 7, 12, 31, 688109, tzinfo=tzutc()), end_at=datetime.datetime(2026, 6, 3, 7, 12, 31, 688109, tzinfo=tzutc()), state=<SandboxState.RUNNING: 'running'>, cpu_count=2, memory_mb=2000, envd_version='0.2.0', _envd_access_token=None, allow_internet_access=None, network=None, lifecycle=None, volume_mounts=[]) http://127.0.0.1:12580/sandboxes/router/c7eccbb18af445d2af2d2ff2a19c4711/49983
Hello world Cube！
```

## cli【虚拟环境内部】创建与api【宿主机】创建的模板对比
```shell
# 通过API【宿主机】创建的模板
cubemastercli tpl info \
  --template-id tpl-698a1ab5376f47e883639fa1 \
  --include-request \
  --json
{
    "ret": {
        "ret_code": 200,
        "ret_msg": "success"
    },
    "template_id": "tpl-698a1ab5376f47e883639fa1",
    "instance_type": "cubebox",
    "version": "v2",
    "status": "READY",
    "replicas": [
        {
            "node_id": "10.0.2.15",
            "node_ip": "10.0.2.15",
            "instance_type": "cubebox",
            "spec": "cpu=2000m,mem=2000Mi",
            "status": "READY",
            "phase": "READY"
        }
    ],
    "create_request": {
        "requestID": "c61a8e98-3aa3-4a66-8312-007998fac077",
        "volumes": [
            {
                "name": "cube_rootfs_rw",
                "volume_source": {
                    "empty_dir": {
                        "size_limit": "1G"
                    }
                }
            }
        ],
        "containers": [
            {
                "name": "cubebox-name-0",
                "image": {
                    "image": "rfs-94b0e4b84a45b0f3b630eedf",
                    "annotations": {
                        "cube.master.rootfs.artifact.id": "rfs-94b0e4b84a45b0f3b630eedf",
                        "cube.master.rootfs.artifact.sha256": "5e5c39c2f56f0e2f5599e2e8f59e0c63d00eb54bad1bf273481ef0debc082cf5",
                        "cube.master.rootfs.artifact.size_bytes": "1073741824",
                        "cube.master.rootfs.artifact.token": "49cddce5-e68c-48f6-8e70-e6254bfbcfce",
                        "cube.master.rootfs.artifact.url": "http://127.0.0.1:8089/cube/template/artifact/download?artifact_id=rfs-94b0e4b84a45b0f3b630eedf\u0026token=49cddce5-e68c-48f6-8e70-e6254bfbcfce",
                        "cube.master.rootfs.writable_layer_size": "1G",
                        "cube.master.template.spec_fingerprint": "94b0e4b84a45b0f3b630eedf5fd056ac2815672515a25b79d61326f27681cbc4"
                    },
                    "storage_media": "ext4",
                    "writable_layer_size": "1G"
                },
                "args": [
                    "/usr/local/bin/start-lightweight-code-interpreter.sh"
                ],
                "working_dir": "/workspace",
                "envs": [
                    {
                        "key": "CODE_INTERPRETER_HOST",
                        "value": "0.0.0.0"
                    },
                    {
                        "key": "CODE_INTERPRETER_PORT",
                        "value": "49999"
                    },
                    {
                        "key": "CODE_INTERPRETER_WORKDIR",
                        "value": "/workspace"
                    },
                    {
                        "key": "DEBIAN_FRONTEND",
                        "value": "noninteractive"
                    },
                    {
                        "key": "GPG_KEY",
                        "value": "7169605F62C751356D054A26A821E680E5FA6305"
                    },
                    {
                        "key": "LANG",
                        "value": "C.UTF-8"
                    },
                    {
                        "key": "PATH",
                        "value": "/usr/local/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
                    },
                    {
                        "key": "PYTHONDONTWRITEBYTECODE",
                        "value": "1"
                    },
                    {
                        "key": "PYTHONUNBUFFERED",
                        "value": "1"
                    },
                    {
                        "key": "PYTHON_SHA256",
                        "value": "fb85a13414b028c49ba18bbd523c2d055a30b56b18b92ce454ea2c51edc656c4"
                    },
                    {
                        "key": "PYTHON_VERSION",
                        "value": "3.12.12"
                    },
                    {
                        "key": "UV_TOOL_BIN_DIR",
                        "value": "/usr/local/bin"
                    }
                ],
                "volume_mounts": [
                    {
                        "name": "cube_rootfs_rw",
                        "container_path": "/"
                    }
                ],
                "r_limit": {
                    "no_file": 1000000
                },
                "resources": {
                    "cpu": "2000m",
                    "mem": "2000Mi"
                },
                "security_context": {
                    "privileged": true
                }
            }
        ],
        "annotations": {
            "cube.master.appsnapshot.template.id": "tpl-698a1ab5376f47e883639fa1",
            "cube.master.appsnapshot.template.version": "v2",
            "cube.master.appsnapshot.version": "v2",
            "cube.master.rootfs.artifact.id": "rfs-94b0e4b84a45b0f3b630eedf",
            "cube.master.rootfs.writable_layer_size": "1G",
            "cube.master.system_disk_size": "1",
            "cube.master.template.spec_fingerprint": "94b0e4b84a45b0f3b630eedf5fd056ac2815672515a25b79d61326f27681cbc4"
        },
        "instance_type": "cubebox",
        "network_type": "tap"
    }
}

# 通过cli【虚拟主机内部】创建的模板
[root@localhost ~]# cubemastercli tpl info \
  --template-id tpl-3a05aafec23c4d928cfa1850 \
  --include-request \
  --json
{
    "ret": {
        "ret_code": 200,
        "ret_msg": "success"
    },
    "template_id": "tpl-3a05aafec23c4d928cfa1850",
    "instance_type": "cubebox",
    "version": "v2",
    "status": "READY",
    "replicas": [
        {
            "node_id": "10.0.2.15",
            "node_ip": "10.0.2.15",
            "instance_type": "cubebox",
            "spec": "cpu=2000m,mem=2000Mi",
            "status": "READY",
            "phase": "READY"
        }
    ],
    "create_request": {
        "requestID": "1832cda2-e04d-4937-aeb6-305bcbb592b6",
        "volumes": [
            {
                "name": "cube_rootfs_rw",
                "volume_source": {
                    "empty_dir": {
                        "size_limit": "1G"
                    }
                }
            }
        ],
        "containers": [
            {
                "name": "cubebox-name-0",
                "image": {
                    "image": "rfs-3b28d541be85881ec725d009",
                    "annotations": {
                        "cube.master.rootfs.artifact.id": "rfs-3b28d541be85881ec725d009",
                        "cube.master.rootfs.artifact.sha256": "07bc00907b29e754c9ec348e3cedb44e06b049c9c575abcf1c22053b361c847e",
                        "cube.master.rootfs.artifact.size_bytes": "1073741824",
                        "cube.master.rootfs.artifact.token": "4e0851a1-9557-455e-85a3-f575d293b0b4",
                        "cube.master.rootfs.artifact.url": "http://0.0.0.0:8089/cube/template/artifact/download?artifact_id=rfs-3b28d541be85881ec725d009\u0026token=4e0851a1-9557-455e-85a3-f575d293b0b4",
                        "cube.master.rootfs.writable_layer_size": "1G",
                        "cube.master.template.spec_fingerprint": "3b28d541be85881ec725d00925793852a41f6ec3325b64c4cae5438226da80bf"
                    },
                    "storage_media": "ext4",
                    "writable_layer_size": "1G"
                },
                "args": [
                    "/usr/local/bin/start-lightweight-code-interpreter.sh"
                ],
                "working_dir": "/workspace",
                "envs": [
                    {
                        "key": "CODE_INTERPRETER_HOST",
                        "value": "0.0.0.0"
                    },
                    {
                        "key": "CODE_INTERPRETER_PORT",
                        "value": "49999"
                    },
                    {
                        "key": "CODE_INTERPRETER_WORKDIR",
                        "value": "/workspace"
                    },
                    {
                        "key": "DEBIAN_FRONTEND",
                        "value": "noninteractive"
                    },
                    {
                        "key": "GPG_KEY",
                        "value": "7169605F62C751356D054A26A821E680E5FA6305"
                    },
                    {
                        "key": "LANG",
                        "value": "C.UTF-8"
                    },
                    {
                        "key": "PATH",
                        "value": "/usr/local/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
                    },
                    {
                        "key": "PYTHONDONTWRITEBYTECODE",
                        "value": "1"
                    },
                    {
                        "key": "PYTHONUNBUFFERED",
                        "value": "1"
                    },
                    {
                        "key": "PYTHON_SHA256",
                        "value": "fb85a13414b028c49ba18bbd523c2d055a30b56b18b92ce454ea2c51edc656c4"
                    },
                    {
                        "key": "PYTHON_VERSION",
                        "value": "3.12.12"
                    },
                    {
                        "key": "UV_TOOL_BIN_DIR",
                        "value": "/usr/local/bin"
                    }
                ],
                "volume_mounts": [
                    {
                        "name": "cube_rootfs_rw",
                        "container_path": "/"
                    }
                ],
                "r_limit": {
                    "no_file": 1000000
                },
                "resources": {
                    "cpu": "2000m",
                    "mem": "2000Mi"
                },
                "security_context": {
                    "privileged": true
                },
                "probe": {
                    "probe_handler": {
                        "http_get": {
                            "path": "/health",
                            "port": 49999,
                            "host": ""
                        }
                    },
                    "timeout_ms": 30000,
                    "period_ms": 500,
                    "success_threshold": 1,
                    "failure_threshold": 60
                }
            }
        ],
        "annotations": {
            "com.exposed_ports": "49983:49999",
            "cube.master.appsnapshot.template.id": "tpl-3a05aafec23c4d928cfa1850",
            "cube.master.appsnapshot.template.version": "v2",
            "cube.master.appsnapshot.version": "v2",
            "cube.master.rootfs.artifact.id": "rfs-3b28d541be85881ec725d009",
            "cube.master.rootfs.writable_layer_size": "1G",
            "cube.master.system_disk_size": "1",
            "cube.master.template.spec_fingerprint": "3b28d541be85881ec725d00925793852a41f6ec3325b64c4cae5438226da80bf"
        },
        "instance_type": "cubebox",
        "network_type": "tap"
    }
}
```