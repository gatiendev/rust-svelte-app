<script>
  import { onMount } from 'svelte';

  let priceData = { symbol: 'BTCUSDT', price: '0', change: '0', volume: '0' };
  let connectionStatus = 'Connecting...';

  onMount(() => {
    const ws = new WebSocket('ws://localhost:8080/ws');

    ws.onopen = () => {
      connectionStatus = 'Connected';
    };

    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      priceData = data;
    };

    ws.onclose = () => {
      connectionStatus = 'Disconnected. Reconnecting...';
      // Attempt reconnect after 2 seconds
      setTimeout(() => {
        window.location.reload(); // simple reload for demo
      }, 2000);
    };

    return () => {
      ws.close();
    };
  });
</script>

<main>
  <h1>Crypto Price Feed</h1>
  <p>Status: {connectionStatus}</p>

  <div class="ticker">
    <h2>{priceData.symbol}</h2>
    <div class="price">${Number(priceData.price).toLocaleString()}</div>
    <div class="change" class:positive={priceData.change >= 0} class:negative={priceData.change < 0}>
      24h Change: {Number(priceData.change).toFixed(2)}%
    </div>
    <div class="volume">Volume: {Number(priceData.volume).toLocaleString()}</div>
  </div>
</main>

<style>
  main {
    font-family: sans-serif;
    max-width: 600px;
    margin: 0 auto;
    padding: 2rem;
  }
  .ticker {
    background: #f5f5f5;
    border-radius: 8px;
    padding: 2rem;
    text-align: center;
  }
  .price {
    font-size: 3rem;
    font-weight: bold;
    margin: 1rem 0;
  }
  .change {
    font-size: 1.2rem;
    margin: 0.5rem 0;
  }
  .positive { color: green; }
  .negative { color: red; }
  .volume {
    color: #666;
  }
</style>