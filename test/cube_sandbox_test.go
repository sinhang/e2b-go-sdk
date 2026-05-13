package test

import (
	"context"
	"testing"

	"github.com/sinhang/e2b-go-sdk/e2b"
)

func TestListSandbox(t *testing.T) {
	client := e2b.NewClient("API_KEY")
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
	client := e2b.NewClient("API_KEY")
	sandbox, err := client.CreateSandbox(context.Background(), e2b.CreateSandboxRequest{
		TemplateID: "tpl-3a864cb982224e97ac2168b5",
	})
	if err != nil {
		t.Fatal(err)
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
	client := e2b.NewClient("API_KEY")

	ctx := context.Background()

	result, err := client.CreateTemplate(ctx, e2b.JSONMap{
		"image":               "cube-sandbox-cn.tencentcloudcr.com/cube-sandbox/sandbox-browser:latest",
		"templateID":          "sandbox-browser",
		"source_image_ref":    "cube-sandbox-cn.tencentcloudcr.com/cube-sandbox/sandbox-browser:latest",
		"exposePorts":         []int{49998, 49982},
		"writable_layer_size": "1G",
	})

	if err != nil {
		t.Fatalf("Failed to create template with V2 API: %v", err)
	}

	t.Logf("Template created successfully with V2: %+v", result)
}
