package symbol

import "fmt"

type Symbol string

const (
	// USDT-quoted pairs
	BTCUSDT   Symbol = "BTCUSDT"
	ETHUSDT   Symbol = "ETHUSDT"
	SOLUSDT   Symbol = "SOLUSDT"
	XRPUSDT   Symbol = "XRPUSDT"
	DOGEUSDT  Symbol = "DOGEUSDT"
	BNBUSDT   Symbol = "BNBUSDT"
	ADAUSDT   Symbol = "ADAUSDT"
	LTCUSDT   Symbol = "LTCUSDT"
	AVAXUSDT  Symbol = "AVAXUSDT"
	MATICUSDT Symbol = "MATICUSDT"
	DOTUSDT   Symbol = "DOTUSDT"
	TRXUSDT   Symbol = "TRXUSDT"
	SHIBUSDT  Symbol = "SHIBUSDT"
	UNIUSDT   Symbol = "UNIUSDT"
	ETCUSDT   Symbol = "ETCUSDT"
	ATOMUSDT  Symbol = "ATOMUSDT"
	NEARUSDT  Symbol = "NEARUSDT"
	APTUSDT   Symbol = "APTUSDT"
	OPUSDT    Symbol = "OPUSDT"
	FTMUSDT   Symbol = "FTMUSDT"
	CAKEUSDT  Symbol = "CAKEUSDT"
	ALGOUSDT  Symbol = "ALGOUSDT"
	PEPEUSDT  Symbol = "PEPEUSDT"
	SUIUSDT   Symbol = "SUIUSDT"

	// BTC-quoted pairs
	ETHBTC   Symbol = "ETHBTC"
	BNBBTC   Symbol = "BNBBTC"
	SOLBTC   Symbol = "SOLBTC"
	XRPBTC   Symbol = "XRPBTC"
	DOGEBTC  Symbol = "DOGEBTC"
	ADABTC   Symbol = "ADABTC"
	LTCBTC   Symbol = "LTCBTC"
	AVAXBTC  Symbol = "AVAXBTC"
	MATICBTC Symbol = "MATICBTC"
	DOTBTC   Symbol = "DOTBTC"
	TRXBTC   Symbol = "TRXBTC"
	SHIBBTC  Symbol = "SHIBBTC"
	UNIBTC   Symbol = "UNIBTC"
	ETCBTC   Symbol = "ETCBTC"
	ATOMBTC  Symbol = "ATOMBTC"
	NEARBTC  Symbol = "NEARBTC"
	APTBTC   Symbol = "APTBTC"
	OPBTC    Symbol = "OPBTC"
	FTMBTC   Symbol = "FTMBTC"
	CAKEBTC  Symbol = "CAKEBTC"
	ALGOBTC  Symbol = "ALGOBTC"
	PEPEBTC  Symbol = "PEPEBTC"
	SUIBTC   Symbol = "SUIBTC"
)

var validSymbols = map[string]Symbol{
	// USDT pairs
	string(BTCUSDT):   BTCUSDT,
	string(ETHUSDT):   ETHUSDT,
	string(SOLUSDT):   SOLUSDT,
	string(XRPUSDT):   XRPUSDT,
	string(DOGEUSDT):  DOGEUSDT,
	string(BNBUSDT):   BNBUSDT,
	string(ADAUSDT):   ADAUSDT,
	string(LTCUSDT):   LTCUSDT,
	string(AVAXUSDT):  AVAXUSDT,
	string(MATICUSDT): MATICUSDT,
	string(DOTUSDT):   DOTUSDT,
	string(TRXUSDT):   TRXUSDT,
	string(SHIBUSDT):  SHIBUSDT,
	string(UNIUSDT):   UNIUSDT,
	string(ETCUSDT):   ETCUSDT,
	string(ATOMUSDT):  ATOMUSDT,
	string(NEARUSDT):  NEARUSDT,
	string(APTUSDT):   APTUSDT,
	string(OPUSDT):    OPUSDT,
	string(FTMUSDT):   FTMUSDT,
	string(CAKEUSDT):  CAKEUSDT,
	string(ALGOUSDT):  ALGOUSDT,
	string(PEPEUSDT):  PEPEUSDT,
	string(SUIUSDT):   SUIUSDT,

	// BTC pairs
	string(ETHBTC):   ETHBTC,
	string(BNBBTC):   BNBBTC,
	string(SOLBTC):   SOLBTC,
	string(XRPBTC):   XRPBTC,
	string(DOGEBTC):  DOGEBTC,
	string(ADABTC):   ADABTC,
	string(LTCBTC):   LTCBTC,
	string(AVAXBTC):  AVAXBTC,
	string(MATICBTC): MATICBTC,
	string(DOTBTC):   DOTBTC,
	string(TRXBTC):   TRXBTC,
	string(SHIBBTC):  SHIBBTC,
	string(UNIBTC):   UNIBTC,
	string(ETCBTC):   ETCBTC,
	string(ATOMBTC):  ATOMBTC,
	string(NEARBTC):  NEARBTC,
	string(APTBTC):   APTBTC,
	string(OPBTC):    OPBTC,
	string(FTMBTC):   FTMBTC,
	string(CAKEBTC):  CAKEBTC,
	string(ALGOBTC):  ALGOBTC,
	string(PEPEBTC):  PEPEBTC,
	string(SUIBTC):   SUIBTC,
}

func Parse(s string) (Symbol, error) {
	if sym, ok := validSymbols[s]; ok {
		return sym, nil
	}
	return "", fmt.Errorf("invalid symbol: %s", s)
}
