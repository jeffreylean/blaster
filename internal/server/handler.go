package server

import (
	"net/http"

	"github.com/jeffreylean/blaster/internal/blast"
	"github.com/labstack/echo"
)

func HealthCheck(c echo.Context) error {
	return c.JSON(200, "working fine....")
}

func RunLoadScheduler(c echo.Context) error {
	var i struct {
		TargetURL string `json:"targetUrl" validate:"required"`
		Workers   uint32 `json:"workers" validate:"required"`
		Rampup    uint32 `json:"rampup" validate:"required"`
		Payload   string `json:"payload"`
	}

	if err := c.Bind(&i); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	if err := c.Validate(&i); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	go blast.ScheduleBlast(i.TargetURL, []byte(i.Payload), i.Workers, i.Rampup)

	return c.JSON(200, nil)
}
