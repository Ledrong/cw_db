package delivery

import (
	"Kursach_project/apiii/models"
	motherboardUsecase "Kursach_project/apiii/src/motherboard/usecase"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type Delivery struct {
	MotherboardUC motherboardUsecase.UseCaseI
}

func (delivery *Delivery) CreateMotherboard(c echo.Context) error {
	var motherboard models.Motherboard
	e := c.Bind(&motherboard)
	if e != nil {
		//c.Logger().Error(err) //а мне все это надо ли надо ли 50 на 50
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	e = delivery.MotherboardUC.CreateMotherboard(&motherboard)
	if e != nil {
		if e == models.ErrConflict {
			//c.Logger().Error(err)
			return c.JSON(http.StatusConflict, motherboard)
		} else {
			//c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
		}
	}

	return c.JSON(http.StatusCreated, motherboard)
}

func (delivery *Delivery) UpdateMotherboard(c echo.Context) error {
	var motherboard models.Motherboard
	e := c.Bind(&motherboard)
	if e != nil {
		fmt.Println(e.Error())
		//c.Logger().Error(err) //а мне все это надо ли надо ли 50 на 50
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	id, e := strconv.ParseInt(c.Param("motherboardid"), 10, 64)

	if e != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	e = delivery.MotherboardUC.UpdateMotherboard(&motherboard, id)
	if e != nil {
		if e == models.ErrNotFound { //возвращает ли err not found возвращает
			//c.Logger().Error(e)
			return echo.NewHTTPError(http.StatusNotFound, models.ErrNotFound.Error())
		} else {
			//c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
		}
	}

	return c.JSON(http.StatusOK, motherboard)

}

// motherboardid с маленькой? да
func (delivery *Delivery) SelectMotherboard(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("motherboardid"), 10, 64)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}
	motherboard, e := delivery.MotherboardUC.SelectMotherboard(id)
	if e != nil {
		if e == models.ErrNotFound { //возвращает ли err not found возвращает
			//c.Logger().Error(e)
			return echo.NewHTTPError(http.StatusNotFound, models.ErrNotFound.Error())
		} else {
			//c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
		}
	}

	return c.JSON(http.StatusOK, motherboard)
}

func (delivery *Delivery) ShowFullMotherboard(c echo.Context) error {
	motherboard, e := delivery.MotherboardUC.ShowFullMotherboard()
	if e != nil {
		//c.Logger().Error(e)
		return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
	}

	return c.JSON(http.StatusOK, motherboard)
}

func (delivery *Delivery) ShowCompatibilityRam(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("motherboardid"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	ram, e := delivery.MotherboardUC.ShowCompatibilityRam(id)
	if e != nil {
		//c.Logger().Error(e)
		return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
	}

	return c.JSON(http.StatusOK, ram)
}

func (delivery *Delivery) ShowCompatibilityCpu(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("motherboardid"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	cpu, e := delivery.MotherboardUC.ShowCompatibilityCpu(id)
	if e != nil {
		//c.Logger().Error(e)
		return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
	}

	return c.JSON(http.StatusOK, cpu)
}

func (delivery *Delivery) ShowPartMotherboard(c echo.Context) error {
	motherboard, e := delivery.MotherboardUC.ShowPartMotherboard()
	if e != nil {
		//c.Logger().Error(e)
		return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
	}

	return c.JSON(http.StatusOK, motherboard)

} //вот тут что делать с двумя переменными и что returнить
func (delivery *Delivery) DeleteMotherboard(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("motherboardid"), 10, 64)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	e := delivery.MotherboardUC.DeleteMotherboardById(id)
	if e != nil {
		if e == models.ErrNotFound { //проверить
			//c.Logger().Error(e)
			return echo.NewHTTPError(http.StatusNotFound, models.ErrNotFound.Error())
		} else {
			//`c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
		}
	}

	return c.NoContent(http.StatusOK) //мб заменить на http.StatusNoContent
}

func NewDelivery(e *echo.Echo, motherboardUC motherboardUsecase.UseCaseI) {
	handler := &Delivery{
		MotherboardUC: motherboardUC,
	}
	//вот тут что
	e.POST("/motherboard/create", handler.CreateMotherboard)
	e.PUT("/motherboard/:motherboardid/update", handler.UpdateMotherboard)
	e.GET("/motherboard/:motherboardid/select", handler.SelectMotherboard)
	e.GET("/motherboard/showf", handler.ShowFullMotherboard)
	e.GET("/motherboard/showp", handler.ShowPartMotherboard)
	e.GET("/motherboard/:motherboardid/showram", handler.ShowCompatibilityRam)
	e.GET("/motherboard/:motherboardid/showcpu", handler.ShowCompatibilityCpu)
	e.DELETE("/motherboard/:motherboardid/delete", handler.DeleteMotherboard)
	//e.POST("/motherboard/create", middl.AuthMiddleware(handler.CreateMotherboard, "admin"))
	//e.PUT("/motherboard/:motherboardid/update", middl.AuthMiddleware(handler.UpdateMotherboard, "admin"))
	//e.GET("/motherboard/:motherboardid/select", middl.AuthMiddleware(handler.SelectMotherboard, "admin", "user"))
	//e.GET("/motherboard/showf", middl.AuthMiddleware(handler.ShowFullMotherboard, "admin", "user"))
	//e.GET("/motherboard/:motherboardid/showram", middl.AuthMiddleware(handler.ShowCompatibilityRam, "admin", "user"))
	//e.GET("/motherboard/:motherboardid/showcpu", middl.AuthMiddleware(handler.ShowCompatibilityCpu, "admin", "user"))
	//e.GET("/motherboard/showp", handler.ShowPartMotherboard)
	//e.DELETE("/motherboard/:motherboardid/delete", middl.AuthMiddleware(handler.DeleteMotherboard, "admin"))
}
