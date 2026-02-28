<script>
  import { onMount } from 'svelte';

  let priceData = { symbol: 'BTCUSDT', price: '0', change: '0', volume: '0' };
  let connectionStatus = 'Connecting...';
  let prevPrice = 0;

  onMount(() => {
    const wsUrl = import.meta.env.VITE_WS_URL || 'ws://localhost:8000/ws';
    const ws = new WebSocket(wsUrl);

    ws.onopen = () => {
      connectionStatus = 'Connected';
    };

    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      prevPrice = parseFloat(priceData.price) || 0;
      priceData = data;
    };

    ws.onclose = () => {
      connectionStatus = 'Disconnected';
      // Attempt reconnect after 2 seconds
      setTimeout(() => {
        window.location.reload();
      }, 2000);
    };

    return () => {
      ws.close();
    };
  });
</script>

<main class="min-h-screen bg-gradient-to-br from-gray-900 to-gray-800 flex items-center justify-center p-4">
  <div class="w-full max-w-md">
    <!-- Status Badge -->
    <div class="mb-4 flex justify-end">
      <div class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium
        {connectionStatus === 'Connected' ? 'bg-green-100 text-green-800' : 
          connectionStatus === 'Connecting...' ? 'bg-yellow-100 text-yellow-800' : 
          'bg-red-100 text-red-800'}">
        <span class="w-2 h-2 rounded-full mr-2
          {connectionStatus === 'Connected' ? 'bg-green-500' : 
            connectionStatus === 'Connecting...' ? 'bg-yellow-500' : 
            'bg-red-500'}">
        </span>
        {connectionStatus}
      </div>
    </div>

    <!-- Main Card -->
    <div class="bg-white/10 backdrop-blur-lg rounded-2xl shadow-2xl border border-white/20 p-8 text-white">
      <!-- Symbol Header -->
      <div class="flex items-center justify-between mb-6">
        <h2 class="text-3xl font-bold tracking-tight">{priceData.symbol}</h2>
        <div class="px-3 py-1 bg-white/20 rounded-full text-sm font-medium">
          LIVE
        </div>
      </div>

      <!-- Price Display -->
      <div class="mb-6">
        <div class="text-sm text-white/60 mb-1">Current Price</div>
        <div class="text-5xl font-bold tabular-nums tracking-tight">
          ${Number(priceData.price).toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
        </div>
      </div>

      <!-- 24h Change & Volume -->
      <div class="grid grid-cols-2 gap-4">
        <div class="bg-white/5 rounded-xl p-4">
          <div class="text-sm text-white/60 mb-1">24h Change</div>
          <div class="text-2xl font-semibold flex items-center
            {parseFloat(priceData.change) >= 0 ? 'text-green-400' : 'text-red-400'}">
            {parseFloat(priceData.change) >= 0 ? '▲' : '▼'}
            {Math.abs(parseFloat(priceData.change)).toFixed(2)}%
          </div>
        </div>
        <div class="bg-white/5 rounded-xl p-4">
          <div class="text-sm text-white/60 mb-1">Volume</div>
          <div class="text-2xl font-semibold text-white">
            {Number(priceData.volume).toLocaleString(undefined, { maximumFractionDigits: 0 })}
          </div>
        </div>
      </div>

      <!-- Animated price change indicator (optional) -->
      {#if prevPrice !== parseFloat(priceData.price) && prevPrice !== 0}
        <div class="mt-4 text-center text-sm text-white/40 animate-pulse">
          {parseFloat(priceData.price) > prevPrice ? '▲ Up' : '▼ Down'} from previous
        </div>
      {/if}
    </div>

    <!-- Footer Note -->
    <p class="mt-4 text-center text-sm text-white/40">
      Real-time data via Binance WebSocket
    </p>
  </div>
</main>

<style>
  /* Optional: Add a subtle pulsing effect for the status dot */
  @keyframes pulse {
    0%, 100% { opacity: 1; }
    50% { opacity: 0.5; }
  }
  .animate-pulse {
    animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
  }
</style>