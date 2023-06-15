package delivery

import (
	"Kursach_project/apiii/models"
	ramUsecase "Kursach_project/apiii/src/ram/usecase"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type Delivery struct {
	RamUC ramUsecase.UseCaseI
}

func (delivery *Delivery) CreateRam(c echo.Context) error {
	var ram models.Ram
	e := c.Bind(&ram)
	if e != nil {
		//c.Logger().Error(err) //а мне все это надо ли надо ли 50 на 50
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	e = delivery.RamUC.CreateRam(&ram)
	if e != nil {
		if e == models.ErrConflict {
			//c.Logger().Error(err)
			return c.JSON(http.StatusConflict, ram)
		} else {
			//c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
		}
	}

	return c.JSON(http.StatusCreated, ram)
}

func (delivery *Delivery) UpdateRam(c echo.Context) error {
	var ram models.Ram
	e := c.Bind(&ram)
	if e != nil {
		fmt.Println(e.Error())
		//c.Logger().Error(err) //а мне все это надо ли надо ли 50 на 50
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	id, e := strconv.ParseInt(c.Param("ramid"), 10, 64)

	if e != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	e = delivery.RamUC.UpdateRam(&ram, id)
	if e != nil {
		if e == models.ErrNotFound { //возвращает ли err not found возвращает
			//c.Logger().Error(e)
			return echo.NewHTTPError(http.StatusNotFound, models.ErrNotFound.Error())
		} else {
			//c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
		}
	}

	return c.JSON(http.StatusOK, ram)

}

// ramid с маленькой? да
func (delivery *Delivery) SelectRam(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("ramid"), 10, 64)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}
	ram, e := delivery.RamUC.SelectRam(id)
	if e != nil {
		if e == models.ErrNotFound { //возвращает ли err not found возвращает
			//c.Logger().Error(e)
			return echo.NewHTTPError(http.StatusNotFound, models.ErrNotFound.Error())
		} else {
			//c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
		}
	}

	return c.JSON(http.StatusOK, ram)
}

func (delivery *Delivery) ShowFullRam(c echo.Context) error {
	ram, e := delivery.RamUC.ShowFullRam()
	if e != nil {
		//c.Logger().Error(e)
		return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
	}

	return c.JSON(http.StatusOK, ram)
}

func (delivery *Delivery) ShowCompatibilityCpu(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("ramid"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	cpu, e := delivery.RamUC.ShowCompatibilityCpu(id)
	if e != nil {
		//c.Logger().Error(e)
		return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
	}

	return c.JSON(http.StatusOK, cpu)
}

func (delivery *Delivery) ShowCompatibilityMotherboard(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("ramid"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	motherboard, e := delivery.RamUC.ShowCompatibilityMotherboard(id)
	if e != nil {
		//c.Logger().Error(e)
		return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
	}

	return c.JSON(http.StatusOK, motherboard)
}

func (delivery *Delivery) ShowPartRam(c echo.Context) error {
	ram, e := delivery.RamUC.ShowPartRam()
	if e != nil {
		//c.Logger().Error(e)
		return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
	}

	return c.JSON(http.StatusOK, ram)

} //вот тут что делать с двумя переменными и что returнить
func (delivery *Delivery) DeleteRam(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("ramid"), 10, 64)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}
	e := delivery.RamUC.DeleteRamById(id)
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

func NewDelivery(e *echo.Echo, ramUC ramUsecase.UseCaseI) {
	handler := &Delivery{
		RamUC: ramUC,
	}
	//вот тут что
	e.POST("ram/create", handler.CreateRam)
	e.PUT("/ram/:ramid/update", handler.UpdateRam)
	e.GET("ram/:ramid/select", handler.SelectRam)
	e.GET("ram/showf", handler.ShowFullRam)
	e.GET("ram/showp", handler.ShowPartRam)
	e.GET("/ram/:ramid/showcpu", handler.ShowCompatibilityCpu)
	e.GET("/ram/:ramid/showmotherboard", handler.ShowCompatibilityMotherboard)
	e.DELETE("ram/:ramid/delete", handler.DeleteRam)
	//e.POST("/ram/create", middl.AuthMiddleware(handler.CreateRam, "admin"))
	//e.PUT("/ram/:ramid/update", middl.AuthMiddleware(handler.UpdateRam, "admin"))
	//e.GET("/ram/:ramid/select", middl.AuthMiddleware(handler.SelectRam, "admin", "user"))
	//e.GET("/ram/showf", middl.AuthMiddleware(handler.ShowFullRam, "admin", "user"))
	//e.GET("/ram/:ramid/showcpu", middl.AuthMiddleware(handler.ShowCompatibilityCpu, "admin", "user"))
	//e.GET("/ram/:ramid/showmotherboard", middl.AuthMiddleware(handler.ShowCompatibilityMotherboard, "admin", "user"))
	//e.GET("/ram/showp", handler.ShowPartRam)
	//e.DELETE("/ram/:ramid/delete", middl.AuthMiddleware(handler.DeleteRam, "admin"))
}
