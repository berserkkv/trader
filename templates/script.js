let ws = null;
let lastPair = localStorage.getItem("lastPair") || "BTCUSDT"; // Default to BTC/USDT
let lastPrice = null;

// Connect to WebSocket and send symbol
function connectToWebSocket() {
    //ws = new WebSocket("ws://3.123.41.68:8080/ws");
    ws = new WebSocket("ws://localhost:8080/ws");
    ws.onopen = () => {
        console.log("WebSocket connected");
        ws.send(JSON.stringify(lastPair.toLowerCase()));
    };

    ws.onmessage = (event) => {
        const data = JSON.parse(event.data);
        curPrice = parseFloat(data.price).toFixed(2);

        document.getElementById("price").innerText = curPrice
        if(lastPrice !== null && lastPrice > curPrice){
            document.getElementById("price").className = "text-red-500";
        }else{
            document.getElementById("price").className = "text-green-500";
        }
        lastPrice = curPrice;

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
    document.getElementById("price").innerText = "Loading...";
    const symbol = document.getElementById("symbol").value.toLowerCase();
    localStorage.setItem("lastPair", symbol.toUpperCase());
    lastPair = symbol;
    document.getElementById("priceLabel").innerText = lastPair.toUpperCase();
    if (ws) {
        ws.close(); // triggers `onclose`, where new connection is created
        ws.onclose = () => connectToWebSocket(); // wait for clean close
    } else {
        connectToWebSocket();
    }
}

const buyOrderType = document.getElementById("buyOrderType");
const buyLimitWrapper = document.getElementById("buyLimitWrapper");

const sellOrderType = document.getElementById("sellOrderType");
const sellLimitWrapper = document.getElementById("sellLimitWrapper");

const sellSlider = document.getElementById("sellSlider");
const sellPercent = document.getElementById("sellPercent");

// Toggle Buy Limit input
buyOrderType.addEventListener("change", () => {
    if (buyOrderType.value === "limit") {
        buyLimitWrapper.classList.remove("h-0", "opacity-0");
        buyLimitWrapper.classList.add("h-20", "opacity-100");
    } else {
        buyLimitWrapper.classList.remove("h-20", "opacity-100");
        buyLimitWrapper.classList.add("h-0", "opacity-0");
    }
});

// Toggle Sell Limit input
sellOrderType.addEventListener("change", () => {
    if (sellOrderType.value === "limit") {
        sellLimitWrapper.classList.remove("h-0", "opacity-0");
        sellLimitWrapper.classList.add("h-20", "opacity-100");
    } else {
        sellLimitWrapper.classList.remove("h-20", "opacity-100");
        sellLimitWrapper.classList.add("h-0", "opacity-0");
    }
});

// Update Sell slider %
sellSlider.addEventListener("input", () => {
    sellPercent.innerText = `${sellSlider.value}%`;
});

window.onload = () => {
    document.getElementById("symbol").value = lastPair;
    document.getElementById("priceLabel").innerText = lastPair.toUpperCase();
    connectToWebSocket();
};