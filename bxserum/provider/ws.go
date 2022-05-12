package provider

import (
	"context"
	"fmt"
	"github.com/bloXroute-Labs/serum-api/bxserum/connections"
	pb "github.com/bloXroute-Labs/serum-api/proto"
	"github.com/bloXroute-Labs/serum-api/utils"
	"github.com/gorilla/websocket"
	"github.com/sourcegraph/jsonrpc2"
)

type WSClient struct {
	pb.UnimplementedApiServer

	addr      string
	conn      *websocket.Conn
	requestID utils.RequestID
}

// Connects to Mainnet Serum API
func NewWSClient() (*WSClient, error) {
	return NewWSClientWithEndpoint("ws://174.129.154.164:1810/ws")
}

// Connects to Testnet Serum API
func NewWSClientTestnet() (*WSClient, error) {
	panic("implement me")
}

// Connects to custom Serum API
func NewWSClientWithEndpoint(addr string) (*WSClient, error) {
	conn, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		return nil, err
	}

	return &WSClient{
		addr:      addr,
		conn:      conn,
		requestID: utils.NewRequestID(),
	}, nil
}

// Set limit to 0 to get all bids/asks
func (w *WSClient) GetOrderbook(market string, limit uint32) (*pb.GetOrderbookResponse, error) {
	request, err := w.jsonRPCRequest("GetOrderbook", map[string]interface{}{"market": market, "limit": limit})
	if err != nil {
		return nil, err
	}
	return connections.WSRequest[pb.GetOrderbookResponse](w.conn, request)
}

func (w *WSClient) GetOrderbookStream(ctx context.Context, market string, limit uint32, orderbookChan chan *pb.GetOrderbookStreamResponse) error {
	request, err := w.jsonRPCRequest("GetOrderbookStream", map[string]interface{}{"market": market, "limit": limit})
	if err != nil {
		return err
	}
	return connections.WSStream[pb.GetOrderbookStreamResponse](ctx, w.conn, request, orderbookChan)
}

// Set limit to 0 to get all trades
func (w *WSClient) GetTrades(market string, limit uint32) (*pb.GetTradesResponse, error) {
	request, err := w.jsonRPCRequest("GetTrades", map[string]interface{}{"market": market, "limit": limit})
	if err != nil {
		return nil, err
	}
	return connections.WSRequest[pb.GetTradesResponse](w.conn, request)
}

func (w *WSClient) GetTradesStream(ctx context.Context, market string, limit uint32, tradesChan chan *pb.GetTradesStreamResponse) error {
	request, err := w.jsonRPCRequest("GetTradeStream", map[string]interface{}{"market": market, "limit": limit})
	if err != nil {
		return err
	}
	return connections.WSStream[pb.GetTradesStreamResponse](ctx, w.conn, request, tradesChan)
}

// GetOrders returns all opened orders by owner address and market
func (w *WSClient) GetOrders(market string, owner string) (*pb.GetOrdersResponse, error) {
	request, err := w.jsonRPCRequest("GetOrders", map[string]interface{}{"market": market, "address": owner})
	if err != nil {
		return nil, err
	}
	return connections.WSRequest[pb.GetOrdersResponse](w.conn, request)
}

// Set market to empty string to get all tickers
func (w *WSClient) GetTickers(market string) (*pb.GetTickersResponse, error) {
	request, err := w.jsonRPCRequest("GetTickers", map[string]interface{}{"market": market})
	if err != nil {
		return nil, err
	}
	return connections.WSRequest[pb.GetTickersResponse](w.conn, request)
}

func (w *WSClient) GetMarkets() (*pb.GetMarketsResponse, error) {
	request, err := w.jsonRPCRequest("GetMarkets", nil)
	if err != nil {
		return nil, err
	}
	return connections.WSRequest[pb.GetMarketsResponse](w.conn, request)
}

func (w *WSClient) Close() error {
	err := w.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		return fmt.Errorf("error writing close msg -  %v", err)
	}
	return nil
}

func (w *WSClient) jsonRPCRequest(method string, params map[string]interface{}) ([]byte, error) {
	id := w.requestID.Next()
	req := jsonrpc2.Request{
		Method: method,
		ID:     jsonrpc2.ID{Num: id},
	}
	if err := req.SetParams(params); err != nil {
		return nil, err
	}

	return req.MarshalJSON()
}
