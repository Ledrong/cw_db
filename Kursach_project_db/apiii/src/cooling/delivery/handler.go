package delivery

import (
	"Kursach_project/apiii/models"
	coolingUsecase "Kursach_project/apiii/src/cooling/usecase"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type Delivery struct {
	CoolingUC coolingUsecase.UseCaseI
}

func (delivery *Delivery) CreateCooling(c echo.Context) error {
	var cooling models.Cooling
	e := c.Bind(&cooling)
	if e != nil {
		//c.Logger().Error(err) //а мне все это надо ли надо ли 50 на 50
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	e = delivery.CoolingUC.CreateCooling(&cooling)
	if e != nil {
		if e == models.ErrConflict {
			//c.Logger().Error(err)
			return c.JSON(http.StatusConflict, cooling)
		} else {
			//c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
		}
	}

	return c.JSON(http.StatusCreated, cooling)
}

func (delivery *Delivery) UpdateCooling(c echo.Context) error {
	var cooling models.Cooling
	e := c.Bind(&cooling)
	if e != nil {
		fmt.Println(e.Error())
		//c.Logger().Error(err) //а мне все это надо ли надо ли 50 на 50
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	id, e := strconv.ParseInt(c.Param("coolingid"), 10, 64)

	if e != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	e = delivery.CoolingUC.UpdateCooling(&cooling, id)
	if e != nil {
		if e == models.ErrNotFound { //возвращает ли err not found возвращает
			//c.Logger().Error(e)
			return echo.NewHTTPError(http.StatusNotFound, models.ErrNotFound.Error())
		} else {
			//c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
		}
	}

	return c.JSON(http.StatusOK, cooling)

}

// coolingid с маленькой? да
func (delivery *Delivery) SelectCooling(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("coolingid"), 10, 64)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	cooling, e := delivery.CoolingUC.SelectCooling(id)
	if e != nil {
		if e == models.ErrNotFound { //возвращает ли err not found возвращает
			//c.Logger().Error(e)
			return echo.NewHTTPError(http.StatusNotFound, models.ErrNotFound.Error())
		} else {
			//c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
		}
	}

	return c.JSON(http.StatusOK, cooling)
}

func (delivery *Delivery) ShowFullCooling(c echo.Context) error {
	cooling, e := delivery.CoolingUC.ShowFullCooling()
	if e != nil {
		//c.Logger().Error(e)
		return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
	}

	return c.JSON(http.StatusOK, cooling)
}

func (delivery *Delivery) ShowPartCooling(c echo.Context) error {
	cooling, e := delivery.CoolingUC.ShowPartCooling()
	if e != nil {
		//c.Logger().Error(e)
		return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
	}

	return c.JSON(http.StatusOK, cooling)

} //вот тут что делать с двумя переменными и что returнить
func (delivery *Delivery) DeleteCooling(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("coolingid"), 10, 64)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}
	e := delivery.CoolingUC.DeleteCoolingById(id)
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

func NewDelivery(e *echo.Echo, coolingUC coolingUsecase.UseCaseI) {
	handler := &Delivery{
		CoolingUC: coolingUC,
	}
	//вот тут что
	e.POST("/cooling/create", handler.CreateCooling)
	e.PUT("/cooling/:coolingid/update", handler.UpdateCooling)
	e.GET("/cooling/:coolingid/select", handler.SelectCooling)
	e.GET("/cooling/showf", handler.ShowFullCooling)
	e.GET("/cooling/showp", handler.ShowPartCooling)
	e.DELETE("/cooling/:coolingid/delete", handler.DeleteCooling)
	//e.POST("/cooling/create", middl.AuthMiddleware(handler.CreateCooling, "admin"))
	//e.PUT("/cooling/:coolingid/update", middl.AuthMiddleware(handler.UpdateCooling, "admin"))
	//e.GET("/cooling/:coolingid/select", middl.AuthMiddleware(handler.SelectCooling, "admin", "user"))
	//e.GET("/cooling/showf", middl.AuthMiddleware(handler.ShowFullCooling, "admin", "user"))
	//e.GET("/cooling/showp", handler.ShowPartCooling)
	//e.DELETE("/cooling/:coolingid/delete", middl.AuthMiddleware(handler.DeleteCooling, "admin"))
}

//func GuestNewDelivery(e *echo.Echo, coolingUC coolingUsecase.UseCaseI) {
//	handler := &Delivery{
//		CoolingUC: coolingUC,
//	}
//	guest := e.Group("/guest")
//	//вот тут что
//	guest.GET("/cooling/:coolingid/select", handler.SelectCooling)
//	guest.GET("/cooling/showp", handler.ShowPartCooling)
//}
//
//func UserNewDelivery(e *echo.Echo, coolingUC coolingUsecase.UseCaseI) {
//	handler := &Delivery{
//		CoolingUC: coolingUC,
//	}
//	//вот тут что
//	user := e.Group("/user")
//	user.GET("/cooling/:coolingid/select", handler.SelectCooling)
//	user.GET("/cooling/showf", handler.ShowFullCooling)
//	user.GET("/cooling/showp", handler.ShowPartCooling)
//}
//
//func AdminNewDelivery(e *echo.Echo, coolingUC coolingUsecase.UseCaseI) {
//	handler := &Delivery{
//		CoolingUC: coolingUC,
//	}
//	//admin := e.Group("/admin")
//	//вот тут что
//	e.POST("/cooling/create", handler.CreateCooling)
//	e.PUT("/cooling/:coolingid/update", handler.UpdateCooling)
//	e.GET("/cooling/:coolingid/select", handler.SelectCooling)
//	e.GET("/cooling/showf", handler.ShowFullCooling)
//	e.GET("/cooling/showp", handler.ShowPartCooling)
//	e.DELETE("/cooling/:coolingid/delete", handler.DeleteCooling)
//}
