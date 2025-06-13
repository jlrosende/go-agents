package base

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"slices"
	"sync"
	"syscall"

	"github.com/invopop/jsonschema"
	"github.com/jlrosende/go-agents/agents"
	"github.com/jlrosende/go-agents/llm/providers"
	"github.com/jlrosende/go-agents/mcp"
	"google.golang.org/grpc"

	mcp_tool "github.com/mark3labs/mcp-go/mcp"

	pb "github.com/jlrosende/go-agents/proto/a2a/v1"
)

type BaseAgent struct {
	ctx context.Context

	Name        string
	Description string
	// MCP
	Servers      []string
	IncludeTools []string
	ExcludeTools []string
	mcpServers   map[string]*mcp.MCPServer

	Logger *slog.Logger

	// LLM
	Model        string
	Instructions string
	llm          providers.LLM

	RequestParams *providers.RequestParams

	// GRCP Server
	pb.UnimplementedA2AServiceServer
}

var _ agents.Agent = (*BaseAgent)(nil)

func NewBaseAgent(ctx context.Context, name, description, model, instructions string, servers, includeTools, excludeTools []string, reqParams *providers.RequestParams) agents.Agent {

	// Init LLM factory with model and tools
	return &BaseAgent{
		ctx: ctx,

		Name:        name,
		Description: description,

		Servers:      servers,
		IncludeTools: includeTools,
		ExcludeTools: excludeTools,

		Model:        model,
		Instructions: instructions,
		mcpServers:   map[string]*mcp.MCPServer{},

		RequestParams: reqParams,
	}
}

func (a *BaseAgent) Start() error {

	lis, err := net.Listen("unix", fmt.Sprintf("/tmp/%s.sock", a.Name))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	server := grpc.NewServer()

	pb.RegisterA2AServiceServer(server, a)

	a.Logger.Info(fmt.Sprintf("agent %s listening at %v", a.Name, lis.Addr()))

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		s := <-sigCh
		a.Logger.Info(fmt.Sprintf("got signal %v, attempting graceful shutdown", s))

		server.GracefulStop()
		// grpc.Stop() // leads to error while receiving stream response: rpc error: code = Unavailable desc = transport is closing
		wg.Done()
	}()

	if err := server.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %@", err)
	}

	wg.Wait()

	a.Logger.Info(fmt.Sprintf("clean agent %s shutdown", a.Name))

	return nil
}

func (a BaseAgent) GetName() string {
	return a.Name
}

func (a BaseAgent) GetModel() string {
	return a.Model
}

func (a BaseAgent) GetInstructions() string {
	return a.Instructions
}

func (a BaseAgent) GetRequestParams() *providers.RequestParams {
	return a.RequestParams
}

func (a *BaseAgent) AttachLLM(llm providers.LLM) {
	a.llm = llm
}

func (a *BaseAgent) Initialize() error {

	a.Logger = slog.Default().With(
		slog.String("agent", a.Name),
		slog.String("type", "BASE_AGENT"),
		slog.String("model", a.Model),
	)

	err := a.llm.Initialize()

	if err != nil {
		return fmt.Errorf("error initialize llm %s in agent %s, %w", a.Model, a.Name, err)
	}

	// Init clients and create missing configurations
	a.llm.AttachTools(a.mcpServers, a.IncludeTools, a.ExcludeTools)

	return nil
}

func (a *BaseAgent) AttachMCPServers(servers map[string]*mcp.MCPServer) {

	if a.mcpServers == nil {
		a.mcpServers = map[string]*mcp.MCPServer{}
	}

	for name, server := range servers {
		if slices.Contains(a.Servers, name) {
			a.mcpServers[name] = server
		}
	}
}

func (a *BaseAgent) Send(message string) (string, error) {
	response, err := a.Generate(message)
	if err != nil {
		return "", err
	}

	// Join response text

	result := mcp.Result(response)

	return result.AllText(), nil
}

func (a BaseAgent) Generate(message string) ([]mcp_tool.Content, error) {
	response, err := a.llm.Generate(message)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (a BaseAgent) Structured(message string, responseStruct any) ([]mcp_tool.Content, error) {

	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}

	schema := reflector.Reflect(responseStruct)
	// return schema

	response, err := a.llm.Structured(message, schema)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (a *BaseAgent) GetAgentCard(ctx context.Context, in *pb.GetAgentCardRequest) (*pb.AgentCard, error) {
	return &pb.AgentCard{
		Name:        a.Name,
		Description: a.Description,
		Url:         "http:://<url>:<port>",
		Version:     "0.0.2",
		Capabilities: &pb.AgentCapabilities{
			Streaming:         true,
			PushNotifications: true,
		},
		DefaultInputModes:  []string{"text", "audio"},
		DefaultOutputModes: []string{"text", "audio"},
		Skills: []*pb.AgentSkill{
			{
				Id:          "filesystem",
				Name:        "FileSystem",
				Description: "hj0",
				Tags:        []string{"jhgbvkjhbvkjhv"},
			},
		},
	}, nil
}

func (a *BaseAgent) SendMessage(ctx context.Context, in *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {

	a.Logger.Debug(fmt.Sprintf("Received: %v", in.GetRequest()))

	var buffer bytes.Buffer
	for _, part := range in.GetRequest().GetContent() {

		switch p := part.GetPart().(type) {
		case *pb.Part_Text:
			buffer.WriteString(p.Text + "\n")
		case *pb.Part_Data:
			buffer.WriteString(p.Data.GetData().String() + "\n")
		case *pb.Part_File:
			buffer.WriteString(p.File.GetFileWithUri() + "\n")
		}
	}

	msg, err := a.Send(buffer.String())

	if err != nil {
		return nil, fmt.Errorf("error sending message to agent %w", err)
	}

	return &pb.SendMessageResponse{
		Payload: &pb.SendMessageResponse_Msg{
			Msg: &pb.Message{
				Role: pb.Role_ROLE_AGENT,
				Content: []*pb.Part{
					&pb.Part{
						Part: &pb.Part_Text{
							Text: msg,
						},
					},
				},
			},
		},
	}, nil
}
