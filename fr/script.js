document.getElementById("trade-form").addEventListener("submit", function (e) {
    e.preventDefault();

    const symbol = document.getElementById("symbol").value;
    const action = document.getElementById("action").value;
    const price = document.getElementById("price").value;
    const amount = document.getElementById("amount").value;

    const trade = { symbol, action, price, amount };

    fetch("/api/add_trade", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(trade)
    })
        .then(response => response.text())
        .then(() => {
            loadTrades();
            document.getElementById("trade-form").reset();
        });
});

function loadTrades() {
    fetch("/api/trades")
        .then(response => response.json())
        .then(trades => {
            const tradeList = document.getElementById("trade-list");
            tradeList.innerHTML = "";
            trades.forEach(trade => {
                const li = document.createElement("li");
                li.textContent = `${trade.symbol} - ${trade.action} - $${trade.price} x ${trade.amount}`;
                tradeList.appendChild(li);
            });
        });
}

// Load trades on page load
window.onload = loadTrades;
