package delivery

import (
	"Kursach_project/apiii/models"
	powerboxUsecase "Kursach_project/apiii/src/powerbox/usecase"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type Delivery struct {
	PowerboxUC powerboxUsecase.UseCaseI
}

func (delivery *Delivery) CreatePowerbox(c echo.Context) error {
	var powerbox models.Powerbox
	e := c.Bind(&powerbox)
	if e != nil {
		//c.Logger().Error(err) //а мне все это надо ли надо ли 50 на 50
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	e = delivery.PowerboxUC.CreatePowerbox(&powerbox)
	if e != nil {
		if e == models.ErrConflict {
			//c.Logger().Error(err)
			return c.JSON(http.StatusConflict, powerbox)
		} else {
			//c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
		}
	}

	return c.JSON(http.StatusCreated, powerbox)
}

func (delivery *Delivery) UpdatePowerbox(c echo.Context) error {
	var powerbox models.Powerbox
	e := c.Bind(&powerbox)
	if e != nil {
		fmt.Println(e.Error())
		//c.Logger().Error(err) //а мне все это надо ли надо ли 50 на 50
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	id, e := strconv.ParseInt(c.Param("powerboxid"), 10, 64)

	if e != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	e = delivery.PowerboxUC.UpdatePowerbox(&powerbox, id)
	if e != nil {
		if e == models.ErrNotFound { //возвращает ли err not found возвращает
			//c.Logger().Error(e)
			return echo.NewHTTPError(http.StatusNotFound, models.ErrNotFound.Error())
		} else {
			//c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
		}
	}

	return c.JSON(http.StatusOK, powerbox)

}

// powerboxid с маленькой? да
func (delivery *Delivery) SelectPowerbox(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("powerboxid"), 10, 64)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}
	powerbox, e := delivery.PowerboxUC.SelectPowerbox(id)
	if e != nil {
		if e == models.ErrNotFound { //возвращает ли err not found возвращает
			//c.Logger().Error(e)
			return echo.NewHTTPError(http.StatusNotFound, models.ErrNotFound.Error())
		} else {
			//c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
		}
	}

	return c.JSON(http.StatusOK, powerbox)
}

func (delivery *Delivery) ShowFullPowerbox(c echo.Context) error {
	powerbox, e := delivery.PowerboxUC.ShowFullPowerbox()
	if e != nil {
		//c.Logger().Error(e)
		return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
	}

	return c.JSON(http.StatusOK, powerbox)
}

func (delivery *Delivery) ShowPartPowerbox(c echo.Context) error {
	powerbox, e := delivery.PowerboxUC.ShowPartPowerbox()
	if e != nil {
		//c.Logger().Error(e)
		return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
	}

	return c.JSON(http.StatusOK, powerbox)

} //вот тут что делать с двумя переменными и что returнить
func (delivery *Delivery) DeletePowerbox(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("powerboxid"), 10, 64)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}
	e := delivery.PowerboxUC.DeletePowerboxById(id)
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

func NewDelivery(e *echo.Echo, powerboxUC powerboxUsecase.UseCaseI) {
	handler := &Delivery{
		PowerboxUC: powerboxUC,
	}
	//вот тут что
	e.POST("/powerbox/create", handler.CreatePowerbox)
	e.PUT("/powerbox/:powerboxid/update", handler.UpdatePowerbox)
	e.GET("/powerbox/:powerboxid/select", handler.SelectPowerbox)
	e.GET("/powerbox/showf", handler.ShowFullPowerbox)
	e.GET("/powerbox/showp", handler.ShowPartPowerbox)
	e.DELETE("/powerbox/:powerboxid/delete", handler.DeletePowerbox)
	//e.POST("/powerbox/create", middl.AuthMiddleware(handler.CreatePowerbox, "admin"))
	//e.PUT("/powerbox/:powerboxid/update", middl.AuthMiddleware(handler.UpdatePowerbox, "admin"))
	//e.GET("/powerbox/:powerboxid/select", middl.AuthMiddleware(handler.SelectPowerbox, "admin", "user"))
	//e.GET("/powerbox/showf", middl.AuthMiddleware(handler.ShowFullPowerbox, "admin", "user"))
	//e.GET("/powerbox/showp", handler.ShowPartPowerbox)
	//e.DELETE("/powerbox/:powerboxid/delete", middl.AuthMiddleware(handler.DeletePowerbox, "admin"))
}
