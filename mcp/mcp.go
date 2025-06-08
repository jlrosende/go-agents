package mcp

import (
	"bytes"
	"context"
	"fmt"
	"slices"
	"strings"

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
	tools  map[string]struct{}
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
			envs = append(envs, fmt.Sprintf("%s=%s", strings.ToUpper(key), value))
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

func (server *MCPServer) ListTools() ([]mcp.Tool, error) {
	tools, err := server.client.ListTools(server.ctx, mcp.ListToolsRequest{})

	if err != nil {
		return nil, fmt.Errorf("error list tools mcp server %s, %w", server.Name, err)
	}

	return tools.Tools, nil
}

func (server *MCPServer) CallTool(name string, args any) (*mcp.CallToolResult, error) {
	result, err := server.client.CallTool(server.ctx, mcp.CallToolRequest{
		Params: struct {
			Name      string    `json:"name"`
			Arguments any       `json:"arguments,omitempty"`
			Meta      *mcp.Meta `json:"_meta,omitempty"`
		}{
			Name:      name,
			Arguments: args,
		},
	})

	if err != nil {
		return nil, fmt.Errorf("error call tool %s on mcp server %s, %w", name, server.Name, err)
	}

	return result, nil
}

type Result []mcp.Content

func (r Result) AllText() string {
	var buffer bytes.Buffer
	for _, content := range r {
		if ok, text := GetText(content); ok {
			buffer.WriteString(text + "\n")
		}
	}
	return buffer.String()
}

func (r Result) FirstText() string {
	for _, content := range r {
		if ok, text := GetText(content); ok {
			return text
		}
	}
	return ""
}

func (r Result) LastText() string {
	for _, content := range slices.Backward(r) {
		if ok, text := GetText(content); ok {
			return text
		}
	}
	return ""
}

func GetText(content mcp.Content) (bool, string) {
	switch c := content.(type) {
	case mcp.TextContent:
		return true, c.Text
	case mcp.EmbeddedResource:
		switch embed := c.Resource.(type) {
		case mcp.TextResourceContents:
			return true, embed.Text
		}
	}
	return false, ""
}

func GetUri(content mcp.Content) (bool, string) {
	switch c := content.(type) {

	case mcp.EmbeddedResource:
		switch embed := c.Resource.(type) {
		case mcp.TextResourceContents:
			return true, embed.URI
		case mcp.BlobResourceContents:
			return true, embed.URI
		}
	}
	return false, ""
}

func GetImageData(content mcp.Content) (bool, string) {
	switch c := content.(type) {
	case mcp.ImageContent:
		return true, c.Data
	case mcp.EmbeddedResource:
		switch embed := c.Resource.(type) {
		case mcp.BlobResourceContents:
			return true, embed.Blob
		}
	}
	return false, ""
}
