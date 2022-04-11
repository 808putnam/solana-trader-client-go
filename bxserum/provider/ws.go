package provider

import (
	"context"
	"fmt"
	"github.com/bloXroute-Labs/serum-api/bxserum/helpers"
	pb "github.com/bloXroute-Labs/serum-api/proto"
	"github.com/gorilla/websocket"
	"github.com/sourcegraph/jsonrpc2"
	"sync"
)

type WSClient struct {
	pb.UnsafeApiServer

	addr string
	conn *websocket.Conn
}

var requestID uint64 = 1
var requestIDLock sync.Mutex

func getRequestID() uint64 {
	var val uint64

	requestIDLock.Lock()
	val = requestID
	requestID++
	requestIDLock.Unlock()

	return val
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

	return &WSClient{addr: addr, conn: conn}, nil
}

func (w *WSClient) GetOrderbook(market string) (*pb.GetOrderbookResponse, error) {
	request, err := jsonRPCRequest("GetOrderbook", map[string]string{"market": market})
	if err != nil {
		return nil, err
	}
	return helpers.UnaryWSRequest[pb.GetOrderbookResponse](w.conn, request)
}

func (w *WSClient) GetOrderbookStream(ctx context.Context, market string, orderbookChan chan *pb.GetOrderbookStreamResponse) error {
	request, err := jsonRPCRequest("GetOrderbookStream", map[string]string{"market": market})
	if err != nil {
		return err
	}
	return helpers.UnaryWSStream[pb.GetOrderbookStreamResponse](ctx, w.conn, request, orderbookChan)
}

func (w *WSClient) Close() error {
	err := w.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		// TODO close conn harshly?
		return fmt.Errorf("error writing close msg -  %v", err)
	}
	return nil
}

func jsonRPCRequest(method string, params map[string]string) ([]byte, error) {
	id := getRequestID()
	req := jsonrpc2.Request{
		Method: method,
		ID:     jsonrpc2.ID{Num: id},
	}
	if err := req.SetParams(params); err != nil {
		return nil, err
	}

	return req.MarshalJSON()
}
