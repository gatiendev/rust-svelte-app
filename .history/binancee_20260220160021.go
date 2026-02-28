package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

type Broadcaster struct {
	clients   map[*websocket.Conn]bool
	broadcast chan []byte
	mu        sync.Mutex
}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan []byte),
	}
}

func (b *Broadcaster) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true }, // allow all origins
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	b.mu.Lock()
	b.clients[conn] = true
	b.mu.Unlock()

	// Wait for client disconnect (or read loop if you want to handle incoming messages)
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}

	b.mu.Lock()
	delete(b.clients, conn)
	b.mu.Unlock()
}

func (b *Broadcaster) BroadcastLoop() {
	for msg := range b.broadcast {
		b.mu.Lock()
		for client := range b.clients {
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				client.Close()
				delete(b.clients, client)
			}
		}
		b.mu.Unlock()
	}
}

// Binance WebSocket endpoint for combined streams
const binanceWS = "wss://stream.binance.com:9443/ws"

// Ticker message structure from Binance
type TickerData struct {
	EventType          string `json:"e"` // Event type
	EventTime          int64  `json:"E"` // Event time
	Symbol             string `json:"s"` // Symbol
	PriceChange        string `json:"p"` // Price change
	PriceChangePercent string `json:"P"` // Price change percent
	WeightedAvgPrice   string `json:"w"` // Weighted average price
	LastPrice          string `json:"c"` // Last price
	LastQuantity       string `json:"Q"` // Last quantity
	OpenPrice          string `json:"o"` // Open price
	HighPrice          string `json:"h"` // High price
	LowPrice           string `json:"l"` // Low price
	Volume             string `json:"v"` // Total traded base asset volume
	QuoteVolume        string `json:"q"` // Total traded quote asset volume
	OpenTime           int64  `json:"O"` // Statistics open time
	CloseTime          int64  `json:"C"` // Statistics close time
	FirstTradeID       int64  `json:"F"` // First trade ID
	LastTradeID        int64  `json:"L"` // Last trade ID
	TotalTrades        int    `json:"n"` // Total number of trades
}

// Subscription message to send to Binance
type SubscriptionMessage struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
	ID     int      `json:"id"`
}

func main() {
	broadcaster := NewBroadcaster()
	go broadcaster.BroadcastLoop()

	// Serve static files (your Svelte build) and WebSocket endpoint
	http.Handle("/", http.FileServer(http.Dir("./static"))) // Svelte build output
	http.HandleFunc("/ws", broadcaster.HandleWebSocket)

	// Start HTTP server in a goroutine
	go func() {
		log.Println("HTTP server listening on :8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal("HTTP server error:", err)
		}
	}()

	// Your existing Binance streaming code...
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		cancel()
	}()

	if err := streamPrices(ctx, broadcaster); err != nil { // pass broadcaster
		log.Fatal("Error streaming prices:", err)
	}
}

func streamPrices(ctx context.Context) error {
	// Connect to Binance WebSocket
	conn, _, err := websocket.DefaultDialer.Dial(binanceWS, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to Binance: %w", err)
	}
	defer conn.Close()

	log.Println("Connected to Binance WebSocket")

	// Subscribe to BTCUSDT ticker stream
	// Stream name format: <symbol>@ticker (lowercase)
	// According to Binance docs: single connection can handle up to 1024 streams [citation:1]
	subscribeMsg := SubscriptionMessage{
		Method: "SUBSCRIBE",
		Params: []string{"btcusdt@ticker"},
		ID:     1,
	}

	if err := conn.WriteJSON(subscribeMsg); err != nil {
		return fmt.Errorf("failed to subscribe: %w", err)
	}

	log.Println("Subscribed to BTCUSDT ticker stream")

	// Set up ping/pong handling
	// Binance sends ping frames every 3 minutes; connection closes if no pong within 10 minutes [citation:1]
	conn.SetPingHandler(func(appData string) error {
		log.Println("Received ping, sending pong")
		return conn.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(time.Second))
	})

	// Message processing loop
	for {
		select {
		case <-ctx.Done():
			log.Println("Context cancelled, closing connection")
			// Send unsubscribe message before closing
			unsubscribeMsg := SubscriptionMessage{
				Method: "UNSUBSCRIBE",
				Params: []string{"btcusdt@ticker"},
				ID:     2,
			}
			conn.WriteJSON(unsubscribeMsg)
			return nil

		default:
			// Read message from WebSocket
			_, message, err := conn.ReadMessage()
			if err != nil {
				// Check if this is a normal closure
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					log.Println("Connection closed normally")
					return nil
				}
				return fmt.Errorf("error reading message: %w", err)
			}

			// Parse and process the message
			if err := processMessage(message); err != nil {
				log.Printf("Error processing message: %v", err)
				// Continue running despite individual message errors
			}
		}
	}
}

func processMessage(message []byte) error {
	// First, check if this is a subscription confirmation
	var raw map[string]interface{}
	if err := json.Unmarshal(message, &raw); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Check for subscription response
	if result, ok := raw["result"]; ok {
		log.Printf("Subscription response: %v", result)
		return nil
	}

	// Parse as ticker data
	var ticker TickerData
	if err := json.Unmarshal(message, &ticker); err != nil {
		return fmt.Errorf("failed to parse ticker: %w", err)
	}

	// Check if this is a ticker message
	if ticker.EventType == "24hrTicker" {
		// Format and display the price data
		eventTime := time.UnixMilli(ticker.EventTime)
		log.Printf("[%s] %s - Last Price: %s | 24h Change: %s%% | Volume: %s",
			eventTime.Format("15:04:05"),
			ticker.Symbol,
			ticker.LastPrice,
			ticker.PriceChangePercent,
			ticker.Volume,
		)

		// Here you can add your custom processing:
		// - Store in database
		// - Calculate indicators (pivot points, RSI, etc.)
		// - Trigger trading signals
		// - Forward to WebSocket clients
		// - etc.
	}

	return nil
}
