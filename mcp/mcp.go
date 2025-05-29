package mcp

import (
	"fmt"

	"github.com/mark3labs/mcp-go/client"
	mcp_transport "github.com/mark3labs/mcp-go/client/transport"
)

type Transport string

const (
	TRANSPORT_HTTP  Transport = "http"
	TRANSPORT_STDIO Transport = "stdio"
	TRANSPORT_SSE   Transport = "sse"
)

type MCPServer struct {
	name   string
	client *client.Client
}

func NewMCPServer(name string, transport Transport, url, command string, env []string, args ...string) (*client.Client, error) {
	switch transport {

	case TRANSPORT_HTTP:
		t, err := mcp_transport.NewStreamableHTTP(url)

		if err != nil {
			return nil, fmt.Errorf("error create client sse %s, %w", name, err)
		}

		return client.NewClient(t), nil
	case TRANSPORT_SSE:
		t, err := mcp_transport.NewSSE(url)

		if err != nil {
			return nil, fmt.Errorf("error create client sse %s, %w", name, err)
		}

		return client.NewClient(t), nil
	default:

		return client.NewClient(
			mcp_transport.NewStdio(command, env, args...),
		), nil
	}
}
