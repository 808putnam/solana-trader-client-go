package provider

import (
	"context"
	"github.com/bloXroute-Labs/serum-api/bxserum/connections"
	pb "github.com/bloXroute-Labs/serum-api/proto"
	"github.com/bloXroute-Labs/serum-api/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	pb.UnimplementedApiServer

	apiClient pb.ApiClient
	requestID utils.RequestID
}

// Connects to Mainnet Serum API
func NewGRPCClient() (*GRPCClient, error) {
	return NewGRPCClientWithEndpoint("174.129.154.164:1811")
}

// Connects to Testnet Serum API
func NewGRPCTestnet() (*GRPCClient, error) {
	panic("implement me")
}

// Connects to custom Serum API
func NewGRPCClientWithEndpoint(endpoint string) (*GRPCClient, error) {
	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &GRPCClient{
		apiClient: pb.NewApiClient(conn),
		requestID: utils.NewRequestID(),
	}, nil
}

// Set limit to 0 to get all bids/asks
func (g *GRPCClient) GetOrderbook(ctx context.Context, market string, limit uint32) (*pb.GetOrderbookResponse, error) {
	return g.apiClient.GetOrderbook(ctx, &pb.GetOrderBookRequest{Market: market, Limit: limit})
}

func (g *GRPCClient) GetOrderbookStream(ctx context.Context, market string, limit uint32, outputChan chan *pb.GetOrderbookStreamResponse) error {
	stream, err := g.apiClient.GetOrderbookStream(ctx, &pb.GetOrderBookRequest{Market: market, Limit: limit})
	if err != nil {
		return err
	}

	return connections.GRPCStream[pb.GetOrderbookStreamResponse](stream, market, outputChan)
}

func (g *GRPCClient) GetMarkets(ctx context.Context) (*pb.GetMarketsResponse, error) {
	return g.apiClient.GetMarkets(ctx, &pb.GetMarketsRequest{})
}
