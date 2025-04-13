let ws;
let lastPair = localStorage.getItem("lastPair") || "btcusdt"; // Default to BTC/USDT

// Connect to WebSocket and send symbol
function connectToWebSocket() {
    ws = new WebSocket("ws://3.123.41.68/8080/ws");

    ws.onopen = () => {
        console.log("WebSocket connected");
        ws.send(JSON.stringify(lastPair));
    };

    ws.onmessage = (event) => {
        const data = JSON.parse(event.data);
        document.getElementById("price").innerText = `${data.symbol.toUpperCase()}: $${data.price}`;
    };

    ws.onerror = (err) => {
        console.log("WebSocket error", err);
    };

    ws.onclose = () => {
        console.log("WebSocket disconnected");
        setTimeout(connectToWebSocket, 1000); // Reconnect
    };
}

// Change trading pair
function changePair() {
    const symbol = document.getElementById("symbol").value.toLowerCase();
    localStorage.setItem("lastPair", symbol);
    lastPair = symbol;
    if (ws) {
        ws.close(); // Close existing WebSocket connection
    }
    connectToWebSocket();
}

window.onload = () => {
    document.getElementById("symbol").value = lastPair;
    document.getElementById("priceLabel").value = lastPair;
    connectToWebSocket();
};