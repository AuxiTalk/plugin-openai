package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/auxitalk/plugin-openai/internal/openai"
)

type Runtime struct {
	rpc    *RPC
	logs   io.Writer
	client *openai.Client
	cfg    openai.Config
}

func NewRuntime(input io.Reader, output io.Writer, logs io.Writer, cfg openai.Config) *Runtime {
	r := &Runtime{
		rpc:    NewRPC(input, output),
		logs:   logs,
		client: openai.NewClient(cfg),
		cfg:    cfg,
	}
	r.registerHandlers()
	return r
}

func (r *Runtime) Listen() error {
	fmt.Fprintf(r.logs, "[plugin-openai] ready model=%s base_url=%s\n", r.cfg.Model, r.cfg.BaseURL)
	return r.rpc.Listen()
}

func (r *Runtime) registerHandlers() {
	r.rpc.Handle("plugin.handshake", r.handshake)
	r.rpc.Handle("plugin.start", r.start)
	r.rpc.Handle("plugin.stop", r.stop)
	r.rpc.Handle("plugin.health", r.health)
	r.rpc.Handle("capability.call", r.capabilityCall)
}

func (r *Runtime) handshake(_ json.RawMessage) (any, error) {
	return map[string]any{
		"pluginId":        "openai",
		"protocolVersion": "0.1",
		"capabilities":     []string{"ai.complete"},
	}, nil
}

func (r *Runtime) start(_ json.RawMessage) (any, error) {
	fmt.Fprintln(r.logs, "[plugin-openai] started")
	return map[string]any{"started": true}, nil
}

func (r *Runtime) stop(_ json.RawMessage) (any, error) {
	fmt.Fprintln(r.logs, "[plugin-openai] stopped")
	return map[string]any{"stopped": true}, nil
}

func (r *Runtime) health(_ json.RawMessage) (any, error) {
	return map[string]any{"ok": true, "pluginId": "openai", "model": r.cfg.Model}, nil
}

func (r *Runtime) capabilityCall(params json.RawMessage) (any, error) {
	var req struct {
		Name  string          `json:"name"`
		Input json.RawMessage `json:"input"`
	}
	if err := json.Unmarshal(params, &req); err != nil {
		return nil, err
	}
	if req.Name != "ai.complete" {
		return nil, fmt.Errorf("capability not found: %s", req.Name)
	}

	var input openai.CompleteInput
	if err := json.Unmarshal(req.Input, &input); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.Timeout)
	defer cancel()

	started := time.Now()
	result, err := r.client.Complete(ctx, input)
	if err != nil {
		return nil, err
	}
	fmt.Fprintf(r.logs, "[plugin-openai] ai.complete completed duration=%s model=%s\n", time.Since(started), result.Model)
	return result, nil
}
