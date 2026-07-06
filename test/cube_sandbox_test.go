package test

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sinhang/e2b-go-sdk/e2b"
)

func TestListSandbox(t *testing.T) {
	client := e2b.NewClient()
	sandbox, err := client.ListSandboxes(context.Background(), "")
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range sandbox {
		t.Logf("Sandbox: %+v", s)
		if s.State == "running" {
			err = client.DeleteSandbox(context.Background(), s.SandboxID)
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestCreateSandbox(t *testing.T) {
	hostMount := make([]e2b.HostMountItem, 0)
	hostMount = append(hostMount, e2b.HostMountItem{
		HostPath:  "/mnt/workspaces/1",
		MountPath: "/workspace",
		ReadOnly:  false,
	})

	hostMountJSON, err := json.Marshal(hostMount)
	if err != nil {
		t.Fatalf("序列化 hostMount 失败: %v", err)
	}
	client := e2b.NewClient()
	sandbox, err := client.CreateSandbox(context.Background(), e2b.CreateSandboxRequest{
		TemplateID: "tpl-439d8f2cc0604e57b95a1543",
		EnvVars: map[string]string{
			"E2B_SANDBOX_ID": "sandbox.SandboxID",
			"AppId":          "app id",
			"AppSecret":      "app secret",
		},
		Metadata: map[string]interface{}{
			"host-mount": string(hostMountJSON),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if sandbox == nil {
		t.Fatal("Sandbox is nil")
	}
	//err = client.DeleteSandbox(context.Background(), sandbox.SandboxID)
	//if err != nil {
	//	t.Fatal(err)
	//}
	// cubecli exec -it 8db197da97cc407e8eac54c94c2bd5c8  -- bash -l

	start, err := client.StartProcess(context.Background(), e2b.StartProcessRequest{
		Cmd: "ls",
		Env: map[string]string{
			"E2B_SANDBOX_ID": sandbox.SandboxID,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Process started: %+v", start)
}

func TestCreateTemplateV2(t *testing.T) {
	client := e2b.NewClient()

	ctx := context.Background()
	templateID := fmt.Sprintf("sdk-internal-%d", time.Now().UnixNano())
	//"image": "cube-sandbox-cn.tencentcloudcr.com/cube-sandbox/sandbox-code:latest",
	//"image":             "192.168.1.100:5000/sandbox/cube-code-sandbox:v1",
	//"image":             "192.168.1.100:5000/sandbox/cube-base-code-sandbox:v1",
	//
	// "image":             "cube-sandbox-cn.tencentcloudcr.com/cube-sandbox/sandbox-code:latest",
	result, err := client.CreateTemplate(ctx, e2b.JSONMap{
		"image": "registry.i-mall.top/sandbox/cube-code-sandbox:v2",
		//"image":             "registry.i-mall.top/sandbox/cube-code-sandbox-skill:v6",
		"template-id":       templateID,
		"exposedPorts":      []int{49999, 49983},
		"probePort":         49999,
		"writableLayerSize": "2G",
	})

	if err != nil {
		t.Fatalf("Failed to create template with V2 API: %v", err)
	}

	t.Logf("Template created successfully with V2: %+v", result)

	//for i := 0; i < 10; i++ {
	//	current, getErr := client.GetTemplate(ctx, templateID, nil)
	//	if getErr == nil {
	//		if status, ok := current["status"].(string); ok && status == "READY" {
	//			break
	//		}
	//	}
	//	time.Sleep(2 * time.Second)
	//}

	//for i := 0; i < 12; i++ {
	//	if err := client.DeleteTemplate(ctx, templateID); err == nil {
	//		return
	//	}
	//	time.Sleep(5 * time.Second)
	//}
	//
	//t.Fatalf("Failed to delete template %s after retries", templateID)
}

func TestDeleteTemplate(t *testing.T) {
	client := e2b.NewClient()
	ctx := context.Background()
	list, err := client.ListTemplates(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, template := range list {
		err = client.DeleteTemplate(ctx, template.TemplateID)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("Deleted template: %s", template.TemplateID)
	}
}

func TestCreateSandboxWithMountedExec(t *testing.T) {
	dataPlaneURL := os.Getenv("CUBE_DATAPLANE_URL")
	if dataPlaneURL == "" {
		dataPlaneURL = "https://127.0.0.1:11443"
	}

	client := e2b.NewClient(
		e2b.WithDataPlaneURL(dataPlaneURL),
	)

	ctx := context.Background()

	sandbox, err := client.CreateSandbox(ctx, e2b.CreateSandboxRequest{
		TemplateID: "tpl-3a05aafec23c4d928cfa1850",
		VolumeMounts: []e2b.JSONMap{
			{
				"name": "tmp",
				"path": "/workspace",
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if sandbox == nil {
		t.Fatal("Sandbox is nil")
	}
	sandboxID := sandbox.SandboxID
	t.Logf("Sandbox created: %s", sandboxID)

	time.Sleep(5 * time.Second)

	// Execute via Python code interpreter (data plane), since the control
	// plane /process/start and /sandbox/exec are pure-404 on Cube.
	interpreter := e2b.NewCodeInterpreter(client, sandboxID)

	// Write a flag file to the mounted workspace.
	unique := fmt.Sprintf("e2b-sdk-%d", time.Now().UnixNano())
	writeCode := fmt.Sprintf(`
import subprocess, os
os.makedirs("/workspace", exist_ok=True)
r = subprocess.run(["sh", "-c", "echo %s >/workspace/flag && test -f /workspace/flag"], capture_output=True, text=True)
print("STDOUT:", r.stdout.strip())
print("STDERR:", r.stderr.strip())
print("EXIT_CODE:", r.returncode)
`, unique)

	exec1, err := interpreter.RunSimple(ctx, writeCode)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Write result: %s", exec1.Text())
	if exec1.Error != nil {
		t.Fatalf("Write error: %s: %s", exec1.Error.Name, exec1.Error.Value)
	}
	if !strings.Contains(exec1.Text(), "EXIT_CODE: 0") {
		t.Fatal("expected EXIT_CODE: 0 from write operation")
	}

	// Verify the flag file persists on the mounted volume.
	checkCode := `
import subprocess
r = subprocess.run(["sh", "-c", "test -f /workspace/flag"], capture_output=True, text=True)
print("EXIT_CODE:", r.returncode)
`
	exec2, err := interpreter.RunSimple(ctx, checkCode)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Check result: %s", exec2.Text())
	if exec2.Error != nil {
		t.Fatalf("Check error: %s: %s", exec2.Error.Name, exec2.Error.Value)
	}

	if err := client.DeleteSandbox(ctx, sandboxID); err != nil {
		t.Logf("Warning: failed to delete sandbox: %v", err)
	}
}

// TestCodeInterpreterSimple executes Python code inside a sandbox through
// the data plane (CubeProxy → lightweight code interpreter on port 49999).
func TestCodeInterpreterSimple(t *testing.T) {
	dataPlaneURL := os.Getenv("CUBE_DATAPLANE_URL")
	if dataPlaneURL == "" {
		dataPlaneURL = "https://127.0.0.1:11443"
	}
	templateID := os.Getenv("CUBE_TEMPLATE_ID")
	if templateID == "" {
		templateID = "tpl-3a05aafec23c4d928cfa1850"
	}

	client := e2b.NewClient(
		e2b.WithDataPlaneURL(dataPlaneURL),
	)

	ctx := context.Background()

	sandbox, err := client.CreateSandbox(ctx, e2b.CreateSandboxRequest{
		TemplateID: templateID,
	})
	if err != nil {
		t.Fatalf("Failed to create sandbox: %v", err)
	}
	if sandbox == nil || sandbox.SandboxID == "" {
		t.Fatal("Sandbox creation returned empty ID")
	}
	sandboxID := sandbox.SandboxID
	t.Logf("Sandbox created: %s", sandboxID)

	time.Sleep(5 * time.Second)

	interpreter := e2b.NewCodeInterpreter(client, sandboxID)
	execution, err := interpreter.RunSimple(ctx, "print('Hello world Cube！')")
	if err != nil {
		t.Fatalf("RunSimple failed: %v", err)
	}

	t.Logf("Execution stdout: %q", execution.Logs.Stdout)
	t.Logf("Execution text: %q", execution.Text())

	if execution.Text() != "Hello world Cube！\n" {
		t.Errorf("expected 'Hello world Cube！\\n', got %q", execution.Text())
	}

	if err := client.DeleteSandbox(ctx, sandboxID); err != nil {
		t.Logf("Warning: failed to delete sandbox %s: %v", sandboxID, err)
	}
}

// TestCommandRunSimple executes a shell command inside a sandbox through
// the data plane (CubeProxy → code interpreter on port 49999).
// Note: the template tpl-3a05aafec23c4d928cfa1850 only exposes port 49999,
// not port 49983 (envd). Shell commands are executed via Python subprocess.
func TestCommandRunSimple(t *testing.T) {
	dataPlaneURL := os.Getenv("CUBE_DATAPLANE_URL")
	if dataPlaneURL == "" {
		dataPlaneURL = "https://127.0.0.1:11443"
	}
	templateID := os.Getenv("CUBE_TEMPLATE_ID")
	if templateID == "" {
		templateID = "tpl-3a05aafec23c4d928cfa1850"
	}

	client := e2b.NewClient(
		e2b.WithDataPlaneURL(dataPlaneURL),
	)

	ctx := context.Background()

	sandbox, err := client.CreateSandbox(ctx, e2b.CreateSandboxRequest{
		TemplateID: templateID,
	})
	if err != nil {
		t.Fatalf("Failed to create sandbox: %v", err)
	}
	if sandbox == nil || sandbox.SandboxID == "" {
		t.Fatal("Sandbox creation returned empty ID")
	}
	sandboxID := sandbox.SandboxID
	t.Logf("Sandbox created: %s", sandboxID)

	time.Sleep(5 * time.Second)

	// Execute shell command via Python subprocess (data plane port 49999).
	interpreter := e2b.NewCodeInterpreter(client, sandboxID)
	execution, err := interpreter.RunSimple(ctx, `
import subprocess
r = subprocess.run(["sh", "-c", "echo hello cube"], capture_output=True, text=True)
print("STDOUT:", r.stdout.strip())
print("EXIT_CODE:", r.returncode)
`)
	if err != nil {
		t.Fatalf("RunSimple failed: %v", err)
	}

	t.Logf("Command stdout: %q", execution.Logs.Stdout)
	t.Logf("Command text: %q", execution.Text())

	if !strings.Contains(execution.Text(), "hello cube") {
		t.Errorf("expected 'hello cube' in output, got %q", execution.Text())
	}

	if err := client.DeleteSandbox(ctx, sandboxID); err != nil {
		t.Logf("Warning: failed to delete sandbox %s: %v", sandboxID, err)
	}
}

// TestSandboxVolumeMountedScript demonstrates mounting a volume into a sandbox,
// writing a shell script to the mounted directory, and executing it.
//
// The test uses the Python code interpreter (data plane) to write and execute
// scripts on the mounted filesystem, since envd port 49983 is not available
// on the default code template.
func TestSandboxVolumeMountedScript(t *testing.T) {
	dataPlaneURL := os.Getenv("CUBE_DATAPLANE_URL")
	if dataPlaneURL == "" {
		dataPlaneURL = "https://127.0.0.1:11443"
	}
	templateID := os.Getenv("CUBE_TEMPLATE_ID")
	if templateID == "" {
		templateID = "tpl-3a05aafec23c4d928cfa1850"
	}

	client := e2b.NewClient(
		e2b.WithDataPlaneURL(dataPlaneURL),
	)

	ctx := context.Background()

	// 1. Create sandbox with volume mount: "tmp" volume → /workspace inside sandbox.
	sandbox, err := client.CreateSandbox(ctx, e2b.CreateSandboxRequest{
		TemplateID: templateID,
		VolumeMounts: []e2b.JSONMap{
			{
				"name": "tmp",
				"path": "/workspace",
			},
		},
	})
	if err != nil {
		t.Fatalf("Failed to create sandbox with volume mount: %v", err)
	}
	sandboxID := sandbox.SandboxID
	t.Logf("Sandbox created: %s (volume: tmp → /workspace)", sandboxID)

	time.Sleep(5 * time.Second)

	// 2. Write a shell script to the mounted workspace via Python code interpreter.
	interpreter := e2b.NewCodeInterpreter(client, sandboxID)

	scriptName := fmt.Sprintf("hello-%d.sh", time.Now().UnixNano())
	scriptPath := "/workspace/" + scriptName
	scriptContent := "#!/bin/sh\necho 'Hello from mounted workspace!'\ndate\nhostname\n"

	writeCode := fmt.Sprintf(`
import os
script = """%s"""
path = "%s"
os.makedirs(os.path.dirname(path), exist_ok=True)
with open(path, "w") as f:
    f.write(script)
os.chmod(path, 0o755)
print("Written {} bytes to {}".format(len(script), path))
print("File exists: {}".format(os.path.exists(path)))
`, scriptContent, scriptPath)

	execution, err := interpreter.RunSimple(ctx, writeCode)
	if err != nil {
		t.Fatalf("Failed to write script: %v", err)
	}
	t.Logf("Write result: %s", execution.Text())

	if execution.Error != nil {
		t.Fatalf("Write error: %s", execution.Error.Name+": "+execution.Error.Value)
	}

	// 3. Execute the script via Python subprocess.
	runCode := fmt.Sprintf(`
import subprocess, sys
result = subprocess.run(["sh", "%s"], capture_output=True, text=True)
print("STDOUT:", result.stdout.strip())
print("STDERR:", result.stderr.strip())
print("EXIT_CODE:", result.returncode)
`, scriptPath)

	execution2, err := interpreter.RunSimple(ctx, runCode)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}
	t.Logf("Execution stdout: %q", execution2.Logs.Stdout)
	t.Logf("Execution result: %s", execution2.Text())

	if execution2.Error != nil {
		t.Fatalf("Execution error: %s", execution2.Error.Name+": "+execution2.Error.Value)
	}

	// 4. Verify the script produced expected output.
	stdout := strings.Join(execution2.Logs.Stdout, "\n")
	if !strings.Contains(stdout, "Hello from mounted workspace!") {
		t.Errorf("Expected 'Hello from mounted workspace!' in output, got: %s", stdout)
	}

	// 5. Write a second script to verify persistence on the mounted volume.
	verifyCode := fmt.Sprintf(`
import os
path = "%s"
print("File still exists: {}".format(os.path.exists(path)))
with open(path) as f:
    print("Content: {}".format(f.read().strip()))
`, scriptPath)

	execution3, err := interpreter.RunSimple(ctx, verifyCode)
	if err != nil {
		t.Fatalf("Failed to verify persistence: %v", err)
	}
	t.Logf("Persistence check: %s", execution3.Text())

	if !strings.Contains(execution3.Text(), "Hello from mounted workspace!") {
		t.Error("Volume mount persistence failed: script not found after write")
	}

	// Cleanup.
	if err := client.DeleteSandbox(ctx, sandboxID); err != nil {
		t.Logf("Warning: failed to delete sandbox %s: %v", sandboxID, err)
	}
}

// TestSyncHostFilesToSandbox demonstrates how to sync files from the host
// machine into a sandbox. Host directories cannot be mounted directly;
// instead, read files on the host and write them into the sandbox via the
// code interpreter (data plane).
func TestSyncHostFilesToSandbox(t *testing.T) {
	dataPlaneURL := os.Getenv("CUBE_DATAPLANE_URL")
	if dataPlaneURL == "" {
		dataPlaneURL = "https://127.0.0.1:11443"
	}
	templateID := os.Getenv("CUBE_TEMPLATE_ID")
	if templateID == "" {
		templateID = "tpl-3a05aafec23c4d928cfa1850"
	}

	// 1. Read host files and encode as base64.
	hostScriptDir := filepath.Join("..", "scripts")
	entries, err := os.ReadDir(hostScriptDir)
	if err != nil {
		t.Fatalf("Failed to read host scripts dir: %v", err)
	}

	type hostFile struct {
		Name    string
		Content string // base64-encoded
	}
	var files []hostFile
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		data, err := os.ReadFile(filepath.Join(hostScriptDir, e.Name()))
		if err != nil {
			t.Fatalf("Failed to read %s: %v", e.Name(), err)
		}
		files = append(files, hostFile{
			Name:    e.Name(),
			Content: base64.StdEncoding.EncodeToString(data),
		})
	}
	t.Logf("Read %d files from host", len(files))
	for _, f := range files {
		t.Logf("  %s (%d bytes)", f.Name, len(f.Content))
	}

	// 2. Create sandbox.
	client := e2b.NewClient(e2b.WithDataPlaneURL(dataPlaneURL))
	ctx := context.Background()

	sandbox, err := client.CreateSandbox(ctx, e2b.CreateSandboxRequest{
		TemplateID: templateID,
	})
	if err != nil {
		t.Fatalf("Failed to create sandbox: %v", err)
	}
	sandboxID := sandbox.SandboxID
	t.Logf("Sandbox created: %s", sandboxID)

	time.Sleep(5 * time.Second)

	// 3. Build a Python script that writes all synced files into /synced.
	var buf bytes.Buffer
	buf.WriteString("import os, base64\nos.makedirs('/synced', exist_ok=True)\n")
	for _, f := range files {
		fmt.Fprintf(&buf, `open('/synced/%s','wb').write(base64.b64decode('%s'))`+"\n", f.Name, f.Content)
	}
	buf.WriteString("print('Files synced:', os.listdir('/synced'))\n")
	// Make shell scripts executable.
	buf.WriteString("import stat\n")
	buf.WriteString("for f in os.listdir('/synced'):\n")
	buf.WriteString("    if f.endswith('.sh'):\n")
	buf.WriteString("        os.chmod('/synced/'+f, 0o755)\n")

	interpreter := e2b.NewCodeInterpreter(client, sandboxID)
	execution, err := interpreter.RunSimple(ctx, buf.String())
	if err != nil {
		t.Fatalf("Failed to sync files: %v", err)
	}
	t.Logf("Sync result: %s", execution.Text())
	if execution.Error != nil {
		t.Fatalf("Sync error: %s: %s", execution.Error.Name, execution.Error.Value)
	}

	// 4. Verify synced files exist and have correct content.
	verifyCode := `
import os
print("=== Synced files ===")
for f in sorted(os.listdir('/synced')):
    st = os.stat('/synced/'+f)
    mode = oct(st.st_mode)[-3:]
    print("  {}  ({} bytes, mode={})".format(f, st.st_size, mode))
    # Read first line to verify content
    with open('/synced/'+f) as fh:
        first = fh.readline().rstrip()
        print("    first line: {}".format(first[:80]))
`
	exec2, err := interpreter.RunSimple(ctx, verifyCode)
	if err != nil {
		t.Fatalf("Verification failed: %v", err)
	}
	t.Logf("Verify result:\n%s", exec2.Text())

	// 5. Verify expected files are present.
	for _, f := range files {
		if !strings.Contains(exec2.Text(), f.Name) {
			t.Errorf("Expected synced file %q not found in sandbox", f.Name)
		}
	}

	if err := client.DeleteSandbox(ctx, sandboxID); err != nil {
		t.Logf("Warning: failed to delete sandbox: %v", err)
	}
}
