package test

import (
	"context"
	"fmt"
	"os"
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
		//if s.State == "running" {
		//	err = client.DeleteSandbox(context.Background(), s.SandboxID)
		//	if err != nil {
		//		t.Fatal(err)
		//	}
		//}
	}
}

func TestCreateSandbox(t *testing.T) {
	client := e2b.NewClient()
	sandbox, err := client.CreateSandbox(context.Background(), e2b.CreateSandboxRequest{
		TemplateID: "tpl-3a864cb982224e97ac2168b5",
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
	templateID := fmt.Sprintf("sdk-probe-%d", time.Now().UnixNano())

	result, err := client.CreateTemplate(ctx, e2b.JSONMap{
		"image":             "cube-sandbox-cn.tencentcloudcr.com/cube-sandbox/sandbox-code:latest",
		"templateID":        templateID,
		"exposedPorts":      []int{49999, 49983},
		"probePort":         49999,
		"writableLayerSize": "1G",
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

func TestCreateSandboxWithMountedExec(t *testing.T) {
	client := e2b.NewClient()

	sandbox, err := client.CreateSandbox(context.Background(), e2b.CreateSandboxRequest{
		TemplateID: "tpl-3a864cb982224e97ac2168b5",
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

	unique := fmt.Sprintf("e2b-sdk-%d", time.Now().UnixNano())
	writeReq := e2b.StartProcessRequest{
		SandboxID: sandbox.SandboxID,
		Args: []string{
			"sh",
			"-lc",
			fmt.Sprintf("mkdir -p /workspace && echo %s >/workspace/flag && test -f /workspace/flag", unique),
		},
	}
	start, err := client.StartProcess(context.Background(), writeReq)
	if err != nil {
		t.Fatal(err)
	}
	if start == nil || start.PID == "" {
		t.Fatalf("expected non-empty process PID, got %+v", start)
	}

	checkReq := e2b.StartProcessRequest{
		SandboxID: sandbox.SandboxID,
		Args: []string{
			"sh",
			"-lc",
			"test -f /workspace/flag",
		},
	}
	check, err := client.StartProcess(context.Background(), checkReq)
	if err != nil {
		t.Fatal(err)
	}
	if check == nil || check.PID == "" {
		t.Fatalf("expected non-empty check PID, got %+v", check)
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
// the data plane (CubeProxy → envd on port 49983).
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

	runner := client.Commands(sandboxID)
	result, err := runner.RunSimple(ctx, "echo hello cube")
	if err != nil {
		t.Fatalf("RunSimple failed: %v", err)
	}

	t.Logf("Command stdout: %q", result.Stdout)
	t.Logf("Command stderr: %q", result.Stderr)
	t.Logf("Command exit code: %d", result.ExitCode)

	if result.StdoutText() != "hello cube" {
		t.Errorf("expected 'hello cube', got %q", result.StdoutText())
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
