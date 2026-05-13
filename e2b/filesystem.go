package e2b

import (
	"context"
	"fmt"
	"net/http"
)

type DownloadFileRequest struct {
	Path string `json:"path,omitempty"`
}
type DownloadFileResponse struct {
	Path    string `json:"path,omitempty"`
	Content string `json:"content,omitempty"`
}

type UploadFileRequest struct {
	Path    string `json:"path,omitempty"`
	Content string `json:"content,omitempty"`
}
type UploadFileResponse struct{}

type ComposeFilesRequest struct {
	Sources []string `json:"sources,omitempty"`
	Target  string   `json:"target,omitempty"`
}
type ComposeFilesResponse struct{}

type CreateWatcherRequest struct {
	Path      string `json:"path,omitempty"`
	Recursive bool   `json:"recursive,omitempty"`
}
type CreateWatcherResponse struct {
	WatcherID string `json:"watcherID,omitempty"`
}

type GetWatcherEventsRequest struct {
	WatcherID string `json:"watcherID,omitempty"`
}
type WatcherEvent struct {
	Type string `json:"type,omitempty"`
	Path string `json:"path,omitempty"`
}
type GetWatcherEventsResponse struct {
	Events []WatcherEvent `json:"events,omitempty"`
}

type ListDirRequest struct {
	Path string `json:"path,omitempty"`
}
type DirEntry struct {
	Name  string `json:"name,omitempty"`
	Path  string `json:"path,omitempty"`
	IsDir bool   `json:"isDir,omitempty"`
}
type ListDirResponse struct {
	Entries []DirEntry `json:"entries,omitempty"`
}

type MakeDirRequest struct {
	Path      string `json:"path,omitempty"`
	Recursive bool   `json:"recursive,omitempty"`
}
type MakeDirResponse struct{}

type MoveRequest struct {
	Source string `json:"source,omitempty"`
	Target string `json:"target,omitempty"`
}
type MoveResponse struct{}

type RemoveRequest struct {
	Path      string `json:"path,omitempty"`
	Recursive bool   `json:"recursive,omitempty"`
}
type RemoveResponse struct{}

type RemoveWatcherRequest struct {
	WatcherID string `json:"watcherID,omitempty"`
}
type RemoveWatcherResponse struct{}

type StatRequest struct {
	Path string `json:"path,omitempty"`
}
type StatResponse struct {
	Name    string `json:"name,omitempty"`
	Path    string `json:"path,omitempty"`
	Size    int64  `json:"size,omitempty"`
	IsDir   bool   `json:"isDir,omitempty"`
	Mode    string `json:"mode,omitempty"`
	ModTime string `json:"modTime,omitempty"`
}

type WatchDirRequest struct {
	Path      string `json:"path,omitempty"`
	Recursive bool   `json:"recursive,omitempty"`
}
type WatchDirResponse struct {
	WatcherID string `json:"watcherID,omitempty"`
}

func (c *Client) DownloadFile(ctx context.Context, req DownloadFileRequest) (*DownloadFileResponse, error) {
	var out DownloadFileResponse
	err := c.doJSON(ctx, http.MethodGet, "/filesystem/download", nil, req, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) UploadFile(ctx context.Context, req UploadFileRequest) (*UploadFileResponse, error) {
	var out UploadFileResponse
	err := c.doJSON(ctx, http.MethodPost, "/filesystem/upload", nil, req, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ComposeFiles(ctx context.Context, req ComposeFilesRequest) (*ComposeFilesResponse, error) {
	var out ComposeFilesResponse
	err := c.doJSON(ctx, http.MethodPost, "/filesystem/compose", nil, req, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) CreateWatcher(ctx context.Context, req CreateWatcherRequest) (*CreateWatcherResponse, error) {
	var out CreateWatcherResponse
	err := c.doJSON(ctx, http.MethodPost, "/filesystem/createwatcher", nil, req, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetWatcherEvents(ctx context.Context, req GetWatcherEventsRequest) (*GetWatcherEventsResponse, error) {
	var out GetWatcherEventsResponse
	err := c.doJSON(ctx, http.MethodPost, "/filesystem/getwatcherevents", nil, req, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ListDir(ctx context.Context, req ListDirRequest) (*ListDirResponse, error) {
	var out ListDirResponse
	err := c.doJSON(ctx, http.MethodPost, "/filesystem/listdir", nil, req, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) MakeDir(ctx context.Context, req MakeDirRequest) (*MakeDirResponse, error) {
	var out MakeDirResponse
	err := c.doJSON(ctx, http.MethodPost, "/filesystem/makedir", nil, req, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) Move(ctx context.Context, req MoveRequest) (*MoveResponse, error) {
	var out MoveResponse
	err := c.doJSON(ctx, http.MethodPost, "/filesystem/move", nil, req, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) Remove(ctx context.Context, req RemoveRequest) (*RemoveResponse, error) {
	var out RemoveResponse
	err := c.doJSON(ctx, http.MethodPost, "/filesystem/remove", nil, req, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) RemoveWatcher(ctx context.Context, req RemoveWatcherRequest) (*RemoveWatcherResponse, error) {
	var out RemoveWatcherResponse
	err := c.doJSON(ctx, http.MethodPost, "/filesystem/removewatcher", nil, req, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) Stat(ctx context.Context, req StatRequest) (*StatResponse, error) {
	var out StatResponse
	err := c.doJSON(ctx, http.MethodPost, "/filesystem/stat", nil, req, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) WatchDir(ctx context.Context, req WatchDirRequest) (*WatchDirResponse, error) {
	var out WatchDirResponse
	err := c.doJSON(ctx, http.MethodPost, "/filesystem/watchdir", nil, req, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) FilesystemRPC(ctx context.Context, method string, body interface{}) (map[string]interface{}, error) {
	out := map[string]interface{}{}
	err := c.doJSON(ctx, http.MethodPost, fmt.Sprintf("/filesystem/%s", method), nil, body, &out)
	return out, err
}
