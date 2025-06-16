package chain

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"slices"

	"github.com/google/uuid"
	"github.com/jlrosende/go-agents/agents"
	"github.com/jlrosende/go-agents/agents/workflows/base"
	pb "github.com/jlrosende/go-agents/proto/a2a/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
)

type ChainAgent struct {
	base.BaseAgent

	AgentsChain []string
	Cumulative  bool

	agents map[string]agents.Agent
}

var _ agents.Agent = (*ChainAgent)(nil)

func NewChainAgent(name string, agentNames []string, cumulative bool) *ChainAgent {
	return &ChainAgent{
		BaseAgent: base.BaseAgent{
			Name: name,
		},
		AgentsChain: agentNames,
		Cumulative:  cumulative,
		agents:      make(map[string]agents.Agent),
	}
}

func (a *ChainAgent) Initialize() error {
	a.Logger = slog.Default().With(
		slog.String("agent", a.Name),
		slog.String("type", "CHAIN_AGENT"),
	)

	if a.Url == "" {
		a.Url = fmt.Sprintf("unix:///tmp/go-agent-%s.sock", a.Name)
	}

	if a.Protocol == "" {
		a.Protocol = agents.PROTOCOL_UNIX
	}

	grpcPanicRecoveryHandler := func(p any) (err error) {
		a.Logger.Error("recovered from panic")
		return status.Errorf(codes.Internal, "%s", p)
	}

	interceptorLogger := func(l *slog.Logger) logging.Logger {
		return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
			l.Log(ctx, slog.Level(lvl), msg, fields...)
		})
	}

	a.Server = grpc.NewServer(
		grpc.ChainStreamInterceptor(
			logging.StreamServerInterceptor(
				interceptorLogger(a.Logger),
				// , logging.WithFieldsFromContext(logTraceID)
			),
			recovery.StreamServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
		),

		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(
				interceptorLogger(a.Logger),
				// logging.WithFieldsFromContext(logTraceID),
			),
			recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
		),
	)

	return nil
}

func (a *ChainAgent) AttachAgents(agentMap map[string]agents.Agent) {

	if a.agents == nil {
		a.agents = make(map[string]agents.Agent)
	}

	for name, agent := range agentMap {
		if slices.Contains(a.AgentsChain, name) {
			a.agents[agent.GetName()] = agent
		}
	}
}

func (a *ChainAgent) Start() error {

	a.Logger.Debug("start CLIENT", "url", a.Url)

	if err := a.StartClient(); err != nil {
		return fmt.Errorf("error start agent %s client, %w", a.GetName(), err)
	}

	a.Logger.Debug("start SERVER")

	err := a.StartServer(func(server *grpc.Server) {
		pb.RegisterA2AServiceServer(a.Server, a)
	})

	if err != nil {
		return fmt.Errorf("error start agent %s server, %w", a.GetName(), err)
	}

	return nil
}

// func (a *ChainAgent) Generate(message string) ([]mcp_tool.Content, error) {

// 	// response := []mcp_tool.Content{}

// 	msg := message

// 	a.Logger.Debug("AAAAAAAAAAAAAAAAAAAAAAAAAAA", "msg", msg)

// 	// Iterate over agents
// 	for _, step := range a.AgentsChain {
// 		if agent, ok := a.agents[step]; ok {

// 			// stepCli := agent.GetClient()

// 			a.Logger.Debug(agent.GetName())

// 			// r, err := stepCli.SendMessage(context.Background(), &pb.SendMessageRequest{
// 			// 	Request: &pb.Message{
// 			// 		MessageId: uuid.NewString(),
// 			// 		Role:      pb.Role_ROLE_USER,
// 			// 		Content: []*pb.Part{
// 			// 			{
// 			// 				Part: &pb.Part_Text{
// 			// 					Text: msg,
// 			// 				},
// 			// 			},
// 			// 		},
// 			// 	},
// 			// })

// 			// if err != nil {
// 			// 	return nil, fmt.Errorf("could send message: %w", err)
// 			// }

// 			// log.Printf("MSG Response: %s", r.GetMsg())
// 		}
// 	}

// 	// result := mcp.Result(response)

// 	// response, err := a.llm.Generate(message)

// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// 	// return response, nil

// 	return nil, nil
// }

func (a *ChainAgent) SendMessage(ctx context.Context, in *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {

	a.Logger.Debug(fmt.Sprintf("Received Chain: %v", in.GetRequest()))

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

	// msg, err := a.Generate(buffer.String())

	// Iterate over agents
	for _, step := range a.AgentsChain {
		if agent, ok := a.agents[step]; ok {

			// stepCli := agent.GetClient()

			a.Logger.Debug(agent.GetName())

			// r, err := stepCli.SendMessage(context.Background(), &pb.SendMessageRequest{
			// 	Request: &pb.Message{
			// 		MessageId: uuid.NewString(),
			// 		Role:      pb.Role_ROLE_USER,
			// 		Content: []*pb.Part{
			// 			{
			// 				Part: &pb.Part_Text{
			// 					Text: msg,
			// 				},
			// 			},
			// 		},
			// 	},
			// })

			// if err != nil {
			// 	return nil, fmt.Errorf("could send message: %w", err)
			// }

			// log.Printf("MSG Response: %s", r.GetMsg())
		}
	}

	// result := mcp.Result(response)

	// response, err := a.llm.Generate(message)

	// if err != nil {
	// 	return nil, err
	// }

	// return response, nil

	// if err != nil {
	// 	return nil, fmt.Errorf("error sending message to agent %w", err)
	// }

	// TODO Need more logic to add more interactions
	// - Tasks
	//   - Support for Artifacts
	// - Multi-Turn Intecraction
	// - Push notifications
	// - File exchange support
	// - Structured Responses

	return &pb.SendMessageResponse{
		Payload: &pb.SendMessageResponse_Msg{
			Msg: &pb.Message{
				MessageId: uuid.NewString(),
				ContextId: uuid.NewString(),
				Role:      pb.Role_ROLE_AGENT,
				Content: []*pb.Part{
					{
						Part: &pb.Part_Text{
							Text: "a",
							// Text: mcp.Result("a").AllText(),
						},
					},
				},
			},
		},
	}, nil
}
