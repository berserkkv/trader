package controllerImpl

import (
	"github.com/berserkkv/trader/httpctx"
	"github.com/berserkkv/trader/model"
	service "github.com/berserkkv/trader/service/interface"
	"log/slog"
	"net/http"
	"strconv"
)

type BotControllerImpl struct {
	s service.BotService
}

func NewBotController(svc service.BotService) *BotControllerImpl {
	return &BotControllerImpl{s: svc}
}

func (ctrl *BotControllerImpl) GetAllBots(ctx httpctx.Context) {
	fields := map[string]interface{}{}

	if val := ctx.Query("isNotActive"); val != "" {
		if b, err := strconv.ParseBool(val); err == nil {
			fields["is_not_active"] = b
		}
	}
	if v := ctx.Query("strategy"); v != "" {
		fields["strategy_name"] = v
	}
	if v := ctx.Query("timeframe"); v != "" {
		fields["timeframe"] = v
	}
	if v := ctx.Query("symbol"); v != "" {
		fields["symbol"] = v
	}
	if val := ctx.Query("inPos"); val != "" {
		if b, err := strconv.ParseBool(val); err == nil {
			fields["in_pos"] = b
		}
	}

	bots := ctrl.s.GetAll(fields)
	ctx.JSON(http.StatusOK, bots)
}

func (ctrl *BotControllerImpl) GetBotByID(ctx httpctx.Context) {
	id, ok := extractID(ctx)
	if !ok {
		return
	}
	b := ctrl.s.GetById(id)
	ctx.JSON(http.StatusOK, b)
}

func (ctrl *BotControllerImpl) CreateBot(ctx httpctx.Context) {
	var req model.CreateBotRequest

	if err := ctx.BindJSON(&req); err != nil {
		slog.Error("Invalid input for creating bot request", "error", err)
		ctx.JSON(http.StatusBadRequest, map[string]any{"error": "Invalid input", "details": err.Error()})
		return
	}

	bot, err := ctrl.s.CreateAndAddToBotFather(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]any{"error": "Could not create bot", "details": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, bot)
}

func (ctrl *BotControllerImpl) DeleteBot(ctx httpctx.Context) {
	id, ok := extractID(ctx)
	if !ok {
		return
	}

	err := ctrl.s.Delete(id)
	if err != nil {
		slog.Error("Error deleting bot", "id", id, "error", err)
	}
	ctx.JSON(http.StatusOK, map[string]any{"message": "Deleted bot", "id": id})
}

func (ctrl *BotControllerImpl) StartBot(ctx httpctx.Context) {
	id, ok := extractID(ctx)
	if !ok {
		return
	}

	err := ctrl.s.Start(id)
	if err != nil {
		slog.Error("Error starting bot", "id", id, "error", err)
	}
	ctx.JSON(http.StatusOK, map[string]any{"message": "Started bot", "id": id})
}

func (ctrl *BotControllerImpl) StopBot(ctx httpctx.Context) {
	id, ok := extractID(ctx)
	if !ok {
		return
	}

	err := ctrl.s.Stop(id)
	if err != nil {
		slog.Error("Error stopping bot", "id", id, "error", err)
	}
	ctx.JSON(http.StatusOK, map[string]any{"message": "Stopped bot", "id": id})
}

func (ctrl *BotControllerImpl) ClosePosition(ctx httpctx.Context) {
	id, ok := extractID(ctx)
	if !ok {
		return
	}

	err := ctrl.s.ClosePosition(id)
	if err != nil {
		ctx.JSON(http.StatusOK, map[string]any{"message": "Error closing position", "id": id})
		return
	}
	ctx.JSON(http.StatusOK, map[string]any{"message": "Closed position", "id": id})
}

func extractID(ctx httpctx.Context) (int64, bool) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		slog.Error("Invalid bot ID", "id", idStr, "error", err)
		ctx.JSON(http.StatusBadRequest, map[string]any{"error": "Invalid bot ID"})
		return 0, false
	}
	return id, true
}
