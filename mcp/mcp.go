package mcp

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/client"
	mcp_transport "github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
)

type Transport string

const (
	TRANSPORT_HTTP  Transport = "http"
	TRANSPORT_STDIO Transport = "stdio"
	TRANSPORT_SSE   Transport = "sse"
)

type MCPServer struct {
	ctx    context.Context
	Name   string
	client *client.Client
}

func NewMCPServer(ctx context.Context, name string, transport Transport, url, command string, environments map[string]string, args ...string) (*MCPServer, error) {

	var t mcp_transport.Interface
	var err error

	switch transport {

	case TRANSPORT_HTTP:
		t, err = mcp_transport.NewStreamableHTTP(url)

		if err != nil {
			return nil, fmt.Errorf("error create client http %s, %w", name, err)
		}
	case TRANSPORT_SSE:
		t, err = mcp_transport.NewSSE(url)

		if err != nil {
			return nil, fmt.Errorf("error create client sse %s, %w", name, err)
		}
	case TRANSPORT_STDIO:
		fallthrough
	default:
		envs := []string{}

		for key, value := range environments {
			envs = append(envs, fmt.Sprintf("%s=%s", key, value))
		}

		t = mcp_transport.NewStdio(command, envs, args...)
	}

	return &MCPServer{
		ctx:    ctx,
		Name:   name,
		client: client.NewClient(t),
	}, nil
}

func (server *MCPServer) Start() error {
	err := server.client.Start(server.ctx)
	if err != nil {
		return fmt.Errorf("error start mcp server %s, %w", server.Name, err)
	}
	_, err = server.client.Initialize(server.ctx, mcp.InitializeRequest{})

	if err != nil {
		return fmt.Errorf("error initialize mcp server %s, %w", server.Name, err)
	}
	return nil
}

func (server *MCPServer) ListTools() (*mcp.ListToolsResult, error) {
	tools, err := server.client.ListTools(server.ctx, mcp.ListToolsRequest{})

	if err != nil {
		return nil, fmt.Errorf("error list tools mcp server %s, %w", server.Name, err)
	}

	return tools, nil
}

func (server *MCPServer) CallTool(name string, args any) (*mcp.CallToolResult, error) {
	result, err := server.client.CallTool(server.ctx, mcp.CallToolRequest{
		Params: struct {
			Name      string    "json:\"name\""
			Arguments any       "json:\"arguments,omitempty\""
			Meta      *mcp.Meta "json:\"_meta,omitempty\""
		}{
			Name:      name,
			Arguments: args,
		},
	})

	if err != nil {
		return nil, fmt.Errorf("error list tools mcp server %s, %w", server.Name, err)
	}

	return result, nil
}
