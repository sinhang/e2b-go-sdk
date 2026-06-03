package test

import (
	"context"
	"fmt"
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
		"image": "cube-sandbox-cn.tencentcloudcr.com/cube-sandbox/sandbox-code:latest",
		//"image":      "cube-sandbox-cn.tencentcloudcr.com/cube-sandbox/sandbox-browser:latest",
		"templateID": templateID,
		//"sourceImageRef": "cube-sandbox-cn.tencentcloudcr.com/cube-sandbox/sandbox-code:latest",
		//"exposedPorts":      []int{49999, 49983},
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
