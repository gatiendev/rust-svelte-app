package main

// Import statements group all dependencies. Standard library packages come first,
// then third-party packages (separated by a blank line for readability).
import (
	"context"       // Provides cancellation and deadlines across API boundaries
	"encoding/json" // For parsing and generating JSON data
	"fmt"           // Formatted I/O (used for error wrapping)
	"log"           // Simple logging to console
	"net/http"      // HTTP client and server implementations
	"os"            // Operating system functions (signals, exit)
	"os/signal"     // For handling OS signals (graceful shutdown)
	"sync"          // Provides mutexes for safe concurrent access
	"syscall"       // Low-level system calls (for signal constants)
	"time"          // Time handling (timestamps, durations)

	"github.com/gorilla/websocket" // Fast, well-tested WebSocket library
	"github.com/rs/cors"           // CORS middleware for HTTP handlers
)

// Broadcaster manages WebSocket clients and broadcasts messages to them.
// It contains:
// - clients: a set of active WebSocket connections (map with bool as placeholder)
// - broadcast: a channel that accepts messages to be sent to all clients
// - mu: a mutex to protect concurrent access to the clients map
type Broadcaster struct {
	clients   map[*websocket.Conn]bool
	broadcast chan []byte
	mu        sync.Mutex
}

// NewBroadcaster is a constructor that creates and initializes a Broadcaster.
// It returns a pointer so that the same instance can be shared across goroutines.
// The clients map is initialized (empty), and the broadcast channel is created.
// Note: The channel is unbuffered, meaning sends will block until a receiver is ready.
func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan []byte),
	}
}

// HandleWebSocket upgrades an HTTP connection to WebSocket and registers the client.
// This method is called when a client connects to the /ws endpoint.
// Parameters:
//
//	w: http.ResponseWriter to write the HTTP response
//	r: *http.Request containing the client's request
func (b *Broadcaster) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrader defines parameters for upgrading from HTTP to WebSocket.
	// CheckOrigin is set to allow all origins (for development; restrict in production).
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	// defer ensures the connection is closed when this function returns,
	// even if an error occurs later.
	defer conn.Close()

	// Lock the mutex before modifying the clients map to prevent data races.
	b.mu.Lock()
	b.clients[conn] = true
	b.mu.Unlock()

	// This loop blocks until the client disconnects or an error occurs.
	// ReadMessage returns the message type, message data, and error.
	// We ignore the actual messages because we only broadcast from server to client,
	// but we need to read to detect when the client closes the connection.
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}

	// When the loop exits, remove the client from the map.
	b.mu.Lock()
	delete(b.clients, conn)
	b.mu.Unlock()
}

// BroadcastLoop runs as a separate goroutine and listens for messages on the broadcast channel.
// For each message received, it sends the message to all currently connected clients.
// If writing to a client fails, the client is assumed dead and is removed.
func (b *Broadcaster) BroadcastLoop() {
	// Ranging over a channel continues until the channel is closed.
	// Here we never close broadcast, so this loop runs forever.
	for msg := range b.broadcast {
		// Lock the mutex while iterating over clients to prevent concurrent modifications.
		b.mu.Lock()
		for client := range b.clients {
			// WriteMessage sends a WebSocket message with the given message type and data.
			// websocket.TextMessage indicates UTF-8 encoded text data.
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				// If write fails, close the connection and remove it from the map.
				client.Close()
				delete(b.clients, client)
			}
		}
		b.mu.Unlock()
	}
}

// binanceWS is the WebSocket endpoint for Binance combined streams.
// Constants are declared at package level and are immutable.
const binanceWS = "wss://stream.binance.com:9443/ws"

// TickerData matches the JSON structure of Binance's 24hr ticker stream.
// Each field has a JSON struct tag that tells encoding/json how to map
// JSON object keys to struct fields.
type TickerData struct {
	EventType          string `json:"e"` // e.g., "24hrTicker"
	EventTime          int64  `json:"E"` // Event timestamp in milliseconds
	Symbol             string `json:"s"` // Trading pair symbol (e.g., "BTCUSDT")
	PriceChange        string `json:"p"` // Absolute price change
	PriceChangePercent string `json:"P"` // Price change percent
	WeightedAvgPrice   string `json:"w"` // Weighted average price
	LastPrice          string `json:"c"` // Last trade price
	LastQuantity       string `json:"Q"` // Last trade quantity
	OpenPrice          string `json:"o"` // Open price in the period
	HighPrice          string `json:"h"` // Highest price in the period
	LowPrice           string `json:"l"` // Lowest price in the period
	Volume             string `json:"v"` // Total traded base asset volume
	QuoteVolume        string `json:"q"` // Total traded quote asset volume
	OpenTime           int64  `json:"O"` // Period open time
	CloseTime          int64  `json:"C"` // Period close time
	FirstTradeID       int64  `json:"F"` // First trade ID in the period
	LastTradeID        int64  `json:"L"` // Last trade ID in the period
	TotalTrades        int    `json:"n"` // Number of trades in the period
}

// SubscriptionMessage is used to subscribe/unsubscribe to Binance streams.
// It matches the JSON format required by Binance WebSocket API.
type SubscriptionMessage struct {
	Method string   `json:"method"` // "SUBSCRIBE" or "UNSUBSCRIBE"
	Params []string `json:"params"` // List of stream names (e.g., "btcusdt@ticker")
	ID     int      `json:"id"`     // Request ID (used to match responses)
}

func main() {
	// Create a broadcaster instance to manage WebSocket clients.
	broadcaster := NewBroadcaster()
	// Start the broadcast loop in its own goroutine so it runs concurrently.
	go broadcaster.BroadcastLoop()

	// Set up HTTP routes.
	// FileServer serves static files from the "./static" directory.
	http.Handle("/", http.FileServer(http.Dir("./static")))
	// HandleWebSocket will be called for requests to "/ws".
	http.HandleFunc("/ws", broadcaster.HandleWebSocket)

	// Wrap the default HTTP handler with CORS middleware to allow cross-origin requests.
	handler := cors.Default().Handler(http.DefaultServeMux)

	// Start the HTTP server in a goroutine so it doesn't block the main thread.
	// ListenAndServe binds to port 8000 on all network interfaces.
	go func() {
		log.Println("HTTP server listening on :8000")
		if err := http.ListenAndServe(":8000", handler); err != nil {
			// Fatal will log the error and exit the program.
			log.Fatal("HTTP server error:", err)
		}
	}()

	// --- Setup for Binance streaming ---

	// Create a context that can be cancelled, used to signal shutdown to goroutines.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure cancel is called when main exits.

	// Create a channel to receive OS signals (interrupt, terminate).
	sigs := make(chan os.Signal, 1)
	// Notify the channel when SIGINT (Ctrl+C) or SIGTERM (kill) is received.
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	// Start a goroutine that waits for a signal and then cancels the context.
	go func() {
		<-sigs
		log.Println("Received shutdown signal, cleaning up...")
		cancel()
	}()

	// Start streaming prices from Binance. This function blocks until the context is cancelled
	// or an unrecoverable error occurs.
	if err := streamPrices(ctx, broadcaster); err != nil {
		log.Fatal("Error streaming prices:", err)
	}
}

// streamPrices connects to Binance WebSocket, subscribes to ticker data,
// and forwards messages to the broadcaster for distribution to clients.
// It uses the provided context for cancellation and cleanup.
func streamPrices(ctx context.Context, broadcaster *Broadcaster) error {
	// Dial establishes a WebSocket connection to Binance.
	conn, _, err := websocket.DefaultDialer.Dial(binanceWS, nil)
	if err != nil {
		// fmt.Errorf with %w wraps the error, preserving the original error type.
		return fmt.Errorf("failed to connect to Binance: %w", err)
	}
	// Ensure the connection is closed when this function exits.
	defer conn.Close()

	log.Println("Connected to Binance WebSocket")

	// Create and send subscription message.
	subscribeMsg := SubscriptionMessage{
		Method: "SUBSCRIBE",
		Params: []string{"btcusdt@ticker"},
		ID:     1,
	}
	if err := conn.WriteJSON(subscribeMsg); err != nil {
		return fmt.Errorf("failed to subscribe: %w", err)
	}
	log.Println("Subscribed to BTCUSDT ticker stream")

	// Set a handler for ping frames. Binance sends pings periodically; we must respond
	// with a pong to keep the connection alive.
	conn.SetPingHandler(func(appData string) error {
		log.Println("Received ping, sending pong")
		// WriteControl sends a control message (pong) with a deadline.
		return conn.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(time.Second))
	})

	// Main message processing loop.
	for {
		select {
		case <-ctx.Done():
			// Context was cancelled (e.g., due to SIGINT). Clean up and exit.
			log.Println("Context cancelled, closing connection")
			// Send unsubscribe message before closing.
			unsubscribeMsg := SubscriptionMessage{
				Method: "UNSUBSCRIBE",
				Params: []string{"btcusdt@ticker"},
				ID:     2,
			}
			conn.WriteJSON(unsubscribeMsg) // Ignore error – we're closing anyway.
			return nil

		default:
			// Read a message from Binance.
			_, message, err := conn.ReadMessage()
			if err != nil {
				// Check if the error is a normal closure (e.g., Binance closed the connection).
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					log.Println("Connection closed normally")
					return nil
				}
				return fmt.Errorf("error reading message: %w", err)
			}

			// Parse the JSON message into our TickerData struct.
			var ticker TickerData
			if err := json.Unmarshal(message, &ticker); err != nil {
				log.Printf("Failed to parse ticker: %v", err)
				continue // Skip this message and continue reading.
			}

			// Check if this is a ticker event (Binance sends other types too).
			if ticker.EventType == "24hrTicker" {
				// Log the price update to the console.
				eventTime := time.UnixMilli(ticker.EventTime)
				log.Printf("[%s] %s - Last Price: %s | 24h Change: %s%% | Volume: %s",
					eventTime.Format("15:04:05"),
					ticker.Symbol,
					ticker.LastPrice,
					ticker.PriceChangePercent,
					ticker.Volume,
				)

				// Prepare a simplified JSON object to send to frontend clients.
				simplified := map[string]interface{}{
					"symbol":    ticker.Symbol,
					"price":     ticker.LastPrice,
					"change":    ticker.PriceChangePercent,
					"volume":    ticker.Volume,
					"timestamp": ticker.EventTime,
				}
				// Marshal the map to JSON. Error is ignored because this map always marshals successfully.
				broadcastMsg, _ := json.Marshal(simplified)

				// Send the message to the broadcaster's channel. This will be picked up by
				// BroadcastLoop and sent to all connected WebSocket clients.
				broadcaster.broadcast <- broadcastMsg
			}
		}
	}
}

// processMessageForBroadcast is a helper that extracts ticker data and prepares it for broadcasting.
// It returns the JSON bytes and a boolean indicating whether broadcasting should occur.
// (Currently not used – streamPrices handles broadcasting directly.)
func processMessageForBroadcast(message []byte) ([]byte, bool) {
	var ticker TickerData
	if err := json.Unmarshal(message, &ticker); err != nil {
		return nil, false
	}
	if ticker.EventType == "24hrTicker" {
		simplified := map[string]interface{}{
			"symbol":    ticker.Symbol,
			"price":     ticker.LastPrice,
			"change":    ticker.PriceChangePercent,
			"volume":    ticker.Volume,
			"timestamp": ticker.EventTime,
		}
		data, _ := json.Marshal(simplified)
		return data, true
	}
	return nil, false
}

// processMessage is another helper that logs ticker data.
// (Currently not used – streamPrices logs directly.)
func processMessage(message []byte) error {
	// Unmarshal into a generic map to check for subscription confirmation.
	var raw map[string]interface{}
	if err := json.Unmarshal(message, &raw); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Check for a "result" field – Binance sends this for subscription acknowledgements.
	if result, ok := raw["result"]; ok {
		log.Printf("Subscription response: %v", result)
		return nil
	}

	// Otherwise, parse as ticker data.
	var ticker TickerData
	if err := json.Unmarshal(message, &ticker); err != nil {
		return fmt.Errorf("failed to parse ticker: %w", err)
	}

	if ticker.EventType == "24hrTicker" {
		eventTime := time.UnixMilli(ticker.EventTime)
		log.Printf("[%s] %s - Last Price: %s | 24h Change: %s%% | Volume: %s",
			eventTime.Format("15:04:05"),
			ticker.Symbol,
			ticker.LastPrice,
			ticker.PriceChangePercent,
			ticker.Volume,
		)
	}
	return nil
}
