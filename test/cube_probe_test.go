package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

const probeBaseURL = "http://127.0.0.1:13000"
const probeAPIKey = "API_KEY"

type probeResult struct {
	Method string
	API    string
	Cube   string
	Code   int
}

// rawProbeResult includes both status code and whether the response body indicates
// the route handler processed the request (even on 404).
type rawProbeResult struct {
	Code        int
	BodyHasJSON bool // true if response body is valid JSON (handler processed request)
}

// rawProbe makes a direct HTTP request and returns detailed result.
func rawProbe(method, path string, body any) rawProbeResult {
	client := &http.Client{Timeout: 10 * time.Second}

	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return rawProbeResult{Code: 0}
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, probeBaseURL+path, reqBody)
	if err != nil {
		return rawProbeResult{Code: 0}
	}
	req.Header.Set("X-API-Key", probeAPIKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.Do(req)
	if err != nil {
		return rawProbeResult{Code: 0}
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	bodyHasJSON := false
	if len(respBytes) > 0 {
		var js interface{}
		bodyHasJSON = json.Unmarshal(respBytes, &js) == nil
	}
	return rawProbeResult{Code: resp.StatusCode, BodyHasJSON: bodyHasJSON}
}

// rawProbeCode is a convenience wrapper that returns just the status code.
func rawProbeCode(method, path string, body any) int {
	return rawProbe(method, path, body).Code
}

func makeResult(method, api string, code int, _ bool, compatCodes ...int) probeResult {
	// Check all codes: primary + compat fallback codes.
	allCodes := append([]int{code}, compatCodes...)
	effectiveCode := code
	cube := "否"

	// Route exists if status code is not 404 and not 0 (connection error).
	for _, c := range allCodes {
		if c != 404 && c != 0 {
			cube = "是"
			effectiveCode = c
			break
		}
	}

	// If primary failed and compat succeeded, use compat code for display.
	if cube == "是" && (code == 404 || code == 0) {
		for _, c := range compatCodes {
			if c != 404 && c != 0 {
				effectiveCode = c
				break
			}
		}
	}
	return probeResult{Method: method, API: api, Cube: cube, Code: effectiveCode}
}

func probeDirect(method, api, path string, body any) probeResult {
	r := rawProbe(method, path, body)
	return makeResult(method, api, r.Code, r.BodyHasJSON)
}

// probeWithCompat probes an endpoint with compat fallback routes.
func probeWithCompat(method, api, primaryPath string, primaryBody any, fallbackPaths ...string) probeResult {
	r := rawProbe(method, primaryPath, primaryBody)
	if r.Code == 404 {
		compatCodes := make([]int, 0, len(fallbackPaths))
		for _, fbPath := range fallbackPaths {
			fbResult := rawProbe(method, fbPath, primaryBody)
			compatCodes = append(compatCodes, fbResult.Code)
		}
		return makeResult(method, api, r.Code, r.BodyHasJSON, compatCodes...)
	}
	return makeResult(method, api, r.Code, r.BodyHasJSON)
}

func TestProbeAllAPIs(t *testing.T) {
	// -----------------------------------------------------------------------
	// Phase 1: Create resources needed for probing
	// -----------------------------------------------------------------------

	// Find working template for sandbox creation.
	templates := listTemplatesRaw(t)
	if len(templates) == 0 {
		t.Fatal("No templates available")
	}
	sandboxTemplateID := templates[0]
	t.Logf("Template for sandbox: %s", sandboxTemplateID)

	// Use an existing template for template-specific probes.
	existingTemplateID := templates[0]
	if len(templates) > 1 {
		existingTemplateID = templates[1] // use a different one to avoid conflicts
	}
	t.Logf("Template for template probes: %s", existingTemplateID)

	// Create a fresh sandbox for probing.
	sandboxID := createSandboxRaw(t, sandboxTemplateID)
	t.Logf("Probe sandbox: %s", sandboxID)

	// Create a fresh template for probing template endpoints (POST, GET, etc.)
	templateID := fmt.Sprintf("probe-%d", time.Now().UnixNano())
	createTemplateRaw(t, templateID, false)
	t.Logf("Probe template: %s", templateID)

	// Create a disposable template specifically for DELETE test.
	deleteTemplateID := fmt.Sprintf("probe-del-%d", time.Now().UnixNano())
	createTemplateRaw(t, deleteTemplateID, true)

	dummyID := "probe-dummy-0000000000000000"
	dummyBuildID := "probe-build-dummy"

	results := []probeResult{}
	add := func(r probeResult) { results = append(results, r) }

	// ---------- Sandboxes ----------
	t.Log("Probing sandbox endpoints...")

	add(probeDirect("GET", "/sandboxes", "/sandboxes", nil))
	add(probeDirect("POST", "/sandboxes", "/sandboxes", map[string]any{"templateID": sandboxTemplateID}))
	add(probeDirect("GET", "/v2/sandboxes", "/v2/sandboxes", nil))
	add(probeDirect("GET", "/sandboxes/metrics", "/sandboxes/metrics", nil))
	add(probeDirect("GET", "/sandboxes/{sandboxID}/logs", "/sandboxes/"+sandboxID+"/logs", nil))
	add(probeDirect("GET", "/v2/sandboxes/{sandboxID}/logs", "/v2/sandboxes/"+sandboxID+"/logs", nil))
	add(probeDirect("GET", "/sandboxes/{sandboxID}", "/sandboxes/"+sandboxID, nil))
	add(probeDirect("DELETE", "/sandboxes/{sandboxID}", "/sandboxes/"+dummyID, nil))
	add(probeDirect("GET", "/sandboxes/{sandboxID}/metrics", "/sandboxes/"+sandboxID+"/metrics", nil))
	add(probeDirect("POST", "/sandboxes/{sandboxID}/snapshots", "/sandboxes/"+sandboxID+"/snapshots", map[string]any{"name": "probe-snap"}))
	add(probeDirect("POST", "/sandboxes/{sandboxID}/pause", "/sandboxes/"+sandboxID+"/pause", nil))
	add(probeDirect("POST", "/sandboxes/{sandboxID}/resume", "/sandboxes/"+sandboxID+"/resume", nil))
	add(probeDirect("POST", "/sandboxes/{sandboxID}/connect", "/sandboxes/"+sandboxID+"/connect", map[string]any{"timeout": 300}))
	add(probeDirect("POST", "/sandboxes/{sandboxID}/timeout", "/sandboxes/"+sandboxID+"/timeout", map[string]any{"timeout": 300}))
	add(probeDirect("PUT", "/sandboxes/{sandboxID}/network", "/sandboxes/"+sandboxID+"/network", map[string]any{}))
	add(probeDirect("POST", "/sandboxes/{sandboxID}/refresh", "/sandboxes/"+sandboxID+"/refresh", map[string]any{}))
	add(probeDirect("GET", "/snapshots", "/snapshots", nil))

	// ---------- Templates ----------
	t.Log("Probing template endpoints...")

	// /v3/templates: direct 404 → compat fallback to /templates
	v3Body := map[string]any{
		"image":             "cube-sandbox-cn.tencentcloudcr.com/cube-sandbox/sandbox-code:latest",
		"templateID":        fmt.Sprintf("probe-v3-%d", time.Now().UnixNano()),
		"writableLayerSize": "1G",
	}
	add(probeWithCompat("POST", "/v3/templates", "/v3/templates", v3Body, "/templates"))

	// /v2/templates: direct 404 → compat fallback to /templates
	v2Body := map[string]any{
		"image":             "cube-sandbox-cn.tencentcloudcr.com/cube-sandbox/sandbox-code:latest",
		"templateID":        fmt.Sprintf("probe-v2-%d", time.Now().UnixNano()),
		"writableLayerSize": "1G",
	}
	add(probeWithCompat("POST", "/v2/templates", "/v2/templates", v2Body, "/templates"))

	add(probeDirect("GET", "/templates/upload-link", "/templates/upload-link", nil))
	add(probeDirect("GET", "/templates", "/templates", nil))

	tplBody := map[string]any{
		"image":             "cube-sandbox-cn.tencentcloudcr.com/cube-sandbox/sandbox-code:latest",
		"templateID":        fmt.Sprintf("probe-legacy-%d", time.Now().UnixNano()),
		"writableLayerSize": "1G",
	}
	add(probeDirect("POST", "/templates", "/templates", tplBody))

	// Use the EXISTING template ID for GET (newly created template may not be immediately queryable).
	add(probeDirect("GET", "/templates/{templateID}", "/templates/"+existingTemplateID, nil))
	add(probeDirect("POST", "/templates/{templateID}", "/templates/"+templateID, map[string]any{}))
	// DELETE: use the disposable template we created specifically for this.
	add(probeDirect("DELETE", "/templates/{templateID}", "/templates/"+deleteTemplateID, nil))
	add(probeDirect("PATCH", "/templates/{templateID}", "/templates/"+existingTemplateID, map[string]any{}))
	add(probeDirect("POST", "/templates/{templateID}/build", "/templates/"+existingTemplateID+"/build", map[string]any{}))
	add(probeDirect("POST", "/v2/templates/{templateID}/build", "/v2/templates/"+existingTemplateID+"/build", map[string]any{}))
	add(probeDirect("PATCH", "/v2/templates/{templateID}", "/v2/templates/"+existingTemplateID, map[string]any{}))
	add(probeDirect("GET", "/templates/{templateID}/builds/{buildID}", "/templates/"+existingTemplateID+"/builds/"+dummyBuildID, nil))
	add(probeDirect("GET", "/templates/{templateID}/builds/{buildID}/logs", "/templates/"+existingTemplateID+"/builds/"+dummyBuildID+"/logs", nil))
	add(probeDirect("GET", "/templates/aliases/{alias}", "/templates/aliases/"+dummyID, nil))

	// ---------- Tags ----------
	t.Log("Probing tag endpoints...")
	add(probeDirect("POST", "/templates/{templateID}/tags", "/templates/"+existingTemplateID+"/tags", map[string]any{"tag": "probe"}))
	add(probeDirect("DELETE", "/templates/{templateID}/tags", "/templates/"+existingTemplateID+"/tags", map[string]any{}))
	add(probeDirect("GET", "/templates/{templateID}/tags", "/templates/"+existingTemplateID+"/tags", nil))

	// ---------- Volumes ----------
	t.Log("Probing volume endpoints...")
	add(probeDirect("GET", "/volumes", "/volumes", nil))
	add(probeDirect("POST", "/volumes", "/volumes", map[string]any{}))
	add(probeDirect("GET", "/volumes/{volumeID}", "/volumes/"+dummyID, nil))
	add(probeDirect("DELETE", "/volumes/{volumeID}", "/volumes/"+dummyID, nil))

	// ---------- Envd ----------
	t.Log("Probing envd endpoints...")
	add(probeDirect("GET", "/envd/health", "/envd/health", nil))
	add(probeDirect("GET", "/envd/stats", "/envd/stats", nil))
	add(probeDirect("GET", "/envd/envs", "/envd/envs", nil))

	// ---------- Filesystem ----------
	t.Log("Probing filesystem endpoints...")
	add(probeDirect("GET", "/filesystem/download", "/filesystem/download?path=/tmp", nil))
	add(probeDirect("POST", "/filesystem/upload", "/filesystem/upload", map[string]any{"path": "/tmp/probe", "content": "test"}))
	add(probeDirect("POST", "/filesystem/compose", "/filesystem/compose", map[string]any{"sources": []string{"/tmp/a"}, "target": "/tmp/b"}))
	add(probeDirect("POST", "/filesystem/createwatcher", "/filesystem/createwatcher", map[string]any{"path": "/tmp"}))
	add(probeDirect("POST", "/filesystem/getwatcherevents", "/filesystem/getwatcherevents", map[string]any{"watcherID": "probe"}))
	add(probeDirect("POST", "/filesystem/listdir", "/filesystem/listdir", map[string]any{"path": "/"}))
	add(probeDirect("POST", "/filesystem/makedir", "/filesystem/makedir", map[string]any{"path": "/tmp/probe-dir"}))
	add(probeDirect("POST", "/filesystem/move", "/filesystem/move", map[string]any{"source": "/tmp/a", "target": "/tmp/b"}))
	add(probeDirect("POST", "/filesystem/remove", "/filesystem/remove", map[string]any{"path": "/tmp/probe"}))
	add(probeDirect("POST", "/filesystem/removewatcher", "/filesystem/removewatcher", map[string]any{"watcherID": "probe"}))
	add(probeDirect("POST", "/filesystem/stat", "/filesystem/stat", map[string]any{"path": "/tmp"}))
	add(probeDirect("POST", "/filesystem/watchdir", "/filesystem/watchdir", map[string]any{"path": "/tmp"}))

	// ---------- Process ----------
	t.Log("Probing process endpoints...")
	add(probeDirect("POST", "/process/closestdin", "/process/closestdin", map[string]any{"pid": "probe"}))
	add(probeDirect("POST", "/process/connect", "/process/connect", map[string]any{"pid": "probe"}))
	add(probeDirect("POST", "/process/list", "/process/list", map[string]any{}))
	add(probeDirect("POST", "/process/sendinput", "/process/sendinput", map[string]any{"pid": "probe", "input": "test"}))
	add(probeDirect("POST", "/process/sendsignal", "/process/sendsignal", map[string]any{"pid": "probe", "signal": "SIGTERM"}))

	// /process/start: direct 404, compat /sandbox/exec also 404
	startBody := map[string]any{"cmd": "echo hello", "sandboxID": sandboxID}
	add(probeWithCompat("POST", "/process/start", "/process/start", startBody, "/sandbox/exec"))

	add(probeDirect("POST", "/process/streaminput", "/process/streaminput", map[string]any{"pid": "probe", "input": "test"}))
	add(probeDirect("POST", "/process/update", "/process/update", map[string]any{"pid": "probe"}))

	// ---------- Teams ----------
	t.Log("Probing team endpoints...")
	add(probeDirect("GET", "/teams", "/teams", nil))
	add(probeDirect("GET", "/teams/metrics", "/teams/metrics", nil))
	add(probeDirect("GET", "/teams/metrics/max", "/teams/metrics/max", nil))

	// -----------------------------------------------------------------------
	// Phase 3: Cleanup
	// -----------------------------------------------------------------------
	t.Log("Cleaning up...")
	rawProbeCode("DELETE", "/sandboxes/"+sandboxID, nil)
	// deleteTemplateID was probed during DELETE test, templateID may need cleanup
	rawProbeCode("DELETE", "/templates/"+templateID, nil)

	// -----------------------------------------------------------------------
	// Phase 4: Output results
	// -----------------------------------------------------------------------
	t.Log("")
	t.Log("========== COMPATIBILITY MATRIX RESULTS ==========")
	t.Log("")

	oldResults := getOldResults()
	changes := 0

	t.Log("| method | api | e2b | cube | code | change |")
	t.Log("|---|---|---|---|---|---|")
	for _, r := range results {
		change := ""
		lookupKey := r.Method + " " + r.API
		if old, ok := oldResults[lookupKey]; ok {
			if old != r.Cube {
				change = fmt.Sprintf("⚠ %s→%s", old, r.Cube)
				changes++
			}
		}
		t.Logf("| %s | %s | 是 | %s | %d | %s |",
			r.Method, r.API, r.Cube, r.Code, change)
	}

	yesCount := 0
	for _, r := range results {
		if r.Cube == "是" {
			yesCount++
		}
	}
	t.Logf("Total: %d endpoints | cube=是: %d | cube=否: %d | changes: %d",
		len(results), yesCount, len(results)-yesCount, changes)

	// Print simple markdown table for README
	t.Log("")
	t.Log("========== README TABLE (copy-paste) ==========")
	t.Log("")
	t.Log("| method | api | e2b | cube |")
	t.Log("|---|---|---|---|")
	for _, r := range results {
		t.Logf("| %s | %s | 是 | %s |", r.Method, r.API, r.Cube)
	}

	// Write details to file
	f, _ := os.Create("/tmp/cube_probe_results.md")
	fmt.Fprintf(f, "## Cube Compatibility Matrix (%s)\n\n", time.Now().Format("2006-01-02"))
	fmt.Fprintf(f, "Probe target: `http://127.0.0.1:13000`  \n")
	fmt.Fprintf(f, "Rule: direct or SDK compat-fallback available => `cube=是`; otherwise => `cube=否`\n\n")
	fmt.Fprintf(f, "| method | api | e2b | cube |\n")
	fmt.Fprintf(f, "|---|---|---|---|\n")
	for _, r := range results {
		fmt.Fprintf(f, "| %s | %s | 是 | %s |\n", r.Method, r.API, r.Cube)
	}
	f.Close()
}

// --- Raw HTTP helpers ---

func listTemplatesRaw(t *testing.T) []string {
	t.Helper()
	client := &http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest("GET", probeBaseURL+"/templates", nil)
	req.Header.Set("X-API-Key", probeAPIKey)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to list templates: %v", err)
	}
	defer resp.Body.Close()
	var templates []struct {
		TemplateID string `json:"templateID"`
		Status     string `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&templates); err != nil {
		t.Fatalf("Failed to decode templates: %v", err)
	}
	ids := make([]string, 0, len(templates))
	for _, tpl := range templates {
		if tpl.Status == "READY" {
			ids = append(ids, tpl.TemplateID)
		}
	}
	return ids
}

func createSandboxRaw(t *testing.T, templateID string) string {
	t.Helper()
	body, _ := json.Marshal(map[string]any{"templateID": templateID})
	client := &http.Client{Timeout: 15 * time.Second}
	req, _ := http.NewRequest("POST", probeBaseURL+"/sandboxes", bytes.NewReader(body))
	req.Header.Set("X-API-Key", probeAPIKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to create sandbox: %v", err)
	}
	defer resp.Body.Close()
	var result struct {
		SandboxID string `json:"sandboxID"`
	}
	respBytes, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(respBytes, &result); err != nil {
		t.Fatalf("Failed to decode: %s", string(respBytes))
	}
	if result.SandboxID == "" {
		t.Fatalf("Empty sandbox ID: %s", string(respBytes))
	}
	return result.SandboxID
}

func createTemplateRaw(t *testing.T, templateID string, failOk bool) {
	t.Helper()
	body, _ := json.Marshal(map[string]any{
		"image":             "cube-sandbox-cn.tencentcloudcr.com/cube-sandbox/sandbox-code:latest",
		"templateID":        templateID,
		"writableLayerSize": "1G",
	})
	client := &http.Client{Timeout: 15 * time.Second}
	req, _ := http.NewRequest("POST", probeBaseURL+"/templates", bytes.NewReader(body))
	req.Header.Set("X-API-Key", probeAPIKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		if failOk {
			t.Logf("Warning: create template %s failed: %v", templateID, err)
		} else {
			t.Fatalf("Failed to create template %s: %v", templateID, err)
		}
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		msg := fmt.Sprintf("Warning: template %s creation returned %d", templateID, resp.StatusCode)
		if failOk {
			t.Log(msg)
		} else {
			t.Error(msg)
		}
	}
}

// getOldResults returns the existing compatibility matrix for change detection.
func getOldResults() map[string]string {
	return map[string]string{
		"GET /sandboxes":                                    "是",
		"POST /sandboxes":                                   "是",
		"GET /v2/sandboxes":                                 "是",
		"GET /sandboxes/metrics":                            "否",
		"GET /sandboxes/{sandboxID}/logs":                   "是",
		"GET /v2/sandboxes/{sandboxID}/logs":                "是",
		"GET /sandboxes/{sandboxID}":                        "是",
		"DELETE /sandboxes/{sandboxID}":                     "是",
		"GET /sandboxes/{sandboxID}/metrics":                "否",
		"POST /sandboxes/{sandboxID}/pause":                 "是",
		"POST /sandboxes/{sandboxID}/resume":                "是",
		"POST /sandboxes/{sandboxID}/connect":               "是",
		"POST /sandboxes/{sandboxID}/timeout":               "否",
		"PUT /sandboxes/{sandboxID}/network":                "否",
		"POST /sandboxes/{sandboxID}/refresh":               "否",
		"POST /sandboxes/{sandboxID}/snapshots":             "是",
		"GET /snapshots":                                    "否",
		"POST /v3/templates":                                "是",
		"POST /v2/templates":                                "是",
		"GET /templates/upload-link":                        "否",
		"GET /templates":                                    "是",
		"POST /templates":                                   "是",
		"GET /templates/{templateID}":                       "是",
		"POST /templates/{templateID}":                      "是",
		"DELETE /templates/{templateID}":                    "是",
		"PATCH /templates/{templateID}":                     "否",
		"POST /templates/{templateID}/build":                "否",
		"POST /v2/templates/{templateID}/build":             "否",
		"PATCH /v2/templates/{templateID}":                  "否",
		"GET /templates/{templateID}/builds/{buildID}":      "否",
		"GET /templates/{templateID}/builds/{buildID}/logs": "是",
		"GET /templates/aliases/{alias}":                    "否",
		"POST /templates/{templateID}/tags":                 "否",
		"DELETE /templates/{templateID}/tags":               "否",
		"GET /templates/{templateID}/tags":                  "否",
		"GET /volumes":                                      "否",
		"POST /volumes":                                     "否",
		"GET /volumes/{volumeID}":                           "否",
		"DELETE /volumes/{volumeID}":                        "否",
		"GET /envd/health":                                  "否",
		"GET /envd/stats":                                   "否",
		"GET /envd/envs":                                    "否",
		"GET /filesystem/download":                          "否",
		"POST /filesystem/upload":                           "否",
		"POST /filesystem/compose":                          "否",
		"POST /filesystem/createwatcher":                    "否",
		"POST /filesystem/getwatcherevents":                 "否",
		"POST /filesystem/listdir":                          "否",
		"POST /filesystem/makedir":                          "否",
		"POST /filesystem/move":                             "否",
		"POST /filesystem/remove":                           "否",
		"POST /filesystem/removewatcher":                    "否",
		"POST /filesystem/stat":                             "否",
		"POST /filesystem/watchdir":                         "否",
		"POST /process/closestdin":                          "否",
		"POST /process/connect":                             "否",
		"POST /process/list":                                "否",
		"POST /process/sendinput":                           "否",
		"POST /process/sendsignal":                          "否",
		"POST /process/start":                               "是",
		"POST /process/streaminput":                         "否",
		"POST /process/update":                              "否",
		"GET /teams":                                        "否",
		"GET /teams/metrics":                                "否",
		"GET /teams/metrics/max":                            "否",
	}
}
