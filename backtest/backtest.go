package backtest

import (
	"github.com/berserkkv/trader/bot"
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/model/enum/order"
	"github.com/berserkkv/trader/service/tools"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"image/color"
	"log"
	"log/slog"
	"math/rand"
	"time"
)

type BacktestBotFather struct {
	candles []model.Candle
	length  int
	Index   int
	Capital []float64
}

func NewBacktestBotFather(length int) *BacktestBotFather {
	return &BacktestBotFather{
		length: length,
	}
}

func (b *BacktestBotFather) Run(bot *bot.Bot) {
	for b.next() {
		c := b.getCandles()
		if len(c) == 0 {
			break
		}
		//if bot.CurrentCapital > bot.InitialCapital*10 {
		//	slog.Info("Excellent")
		//	return
		//}

		if bot.InPos {

			if b.stopLossInRange(bot.OrderStopLoss, c[len(c)-1], bot.OrderType) {
				_, err := bot.ClosePosition(bot.OrderStopLoss)
				if err != nil {
					slog.Error("Error closing position", "err", err)
				}
			}

			if b.takeProfitInRange(bot.OrderTakeProfit, c[len(c)-1], bot.OrderType) {
				_, err := bot.ClosePosition(bot.OrderTakeProfit)
				if err != nil {
					slog.Error("Error closing position", "err", err)
				}
			}

		} else {
			command, _ := bot.Strategy.Run(c)

			if command == order.LONG || command == order.SHORT {

				if bot.CurrentCapital != 0.0 {
					b.Capital = append(b.Capital, bot.CurrentCapital)
				} else {
					b.Capital = append(b.Capital, bot.OrderCapital)
				}
				err := bot.BacktestOpenPosition(command, c[len(c)-1].Close)
				if err != nil {
					slog.Error("Error opening position", "error", err)
					return
				}
			}
		}

	}

}

func (b *BacktestBotFather) LoadData(name string) error {
	candles, err := tools.LoadCandlesFromCSV(name)
	if err != nil {
		slog.Error("Error loading candles from CSV", "error", err)
		return err
	}
	b.candles = candles
	return nil
}

func (b *BacktestBotFather) getCandles() []model.Candle {
	end := b.Index + b.length
	if end > len(b.candles) {
		slog.Info("End of candles")
		return []model.Candle{}
	}
	return b.candles[b.Index:end]
}

func (b *BacktestBotFather) next() bool {
	if b.Index+1 >= len(b.candles) {
		return false
	}
	b.Index++
	return true
}

func (b *BacktestBotFather) stopLossInRange(stopLoss float64, c model.Candle, command order.Command) bool {
	if command == order.LONG {
		return c.Close <= stopLoss || c.Open <= stopLoss || c.High <= stopLoss || c.Low <= stopLoss
	} else if command == order.SHORT {
		return c.Close >= stopLoss || c.Open >= stopLoss || c.High >= stopLoss || c.Low >= stopLoss
	}
	return false
}

func (b *BacktestBotFather) takeProfitInRange(p float64, candle model.Candle, command order.Command) bool {
	if command == order.LONG {
		return candle.High >= p || candle.Open >= p || candle.Close >= p || candle.Low >= p
	} else if command == order.SHORT {
		return candle.Low <= p || candle.Open <= p || candle.Close <= p || candle.High <= p
	}
	return false
}

func (b *BacktestBotFather) Randomize() {
	rand.Seed(time.Now().UnixNano())

	b.Index = rand.Intn(len(b.candles))
	slog.Info("Randomized", "Index", b.Index)
}

func (b *BacktestBotFather) PrintChart(values []float64, info string) {
	if len(values) == 0 {
		slog.Info("Empty values, nothing to print")
		return
	}
	pts := make(plotter.XYs, len(values))
	for i, v := range values {
		pts[i].X = float64(i)
		pts[i].Y = v
	}

	// Create new plot
	p := plot.New()
	p.Title.Text = info
	p.X.Label.Text = "Index"
	p.Y.Label.Text = "Value"
	if values[len(values)-1] < 100 {
		p.BackgroundColor = color.RGBA{R: 255, G: 10, B: 10, A: 255}
	}

	line, err := plotter.NewLine(pts)
	if err != nil {
		log.Fatal(err)
	}
	line.Color = color.RGBA{R: 10, G: 10, B: 10, A: 255}

	//line.FillColor = color.White
	p.Add(line)

	// Save to PNG
	if err := p.Save(8*vg.Inch, 4*vg.Inch, "chart.png"); err != nil {
		log.Fatal(err)
	}
}
