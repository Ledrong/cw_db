package delivery

import (
	"Kursach_project/apiii/middl"
	"Kursach_project/apiii/models"
	videocardUsecase "Kursach_project/apiii/src/videocard/usecase"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type Delivery struct {
	VideocardUC videocardUsecase.UseCaseI
}

func (delivery *Delivery) CreateVideocard(c echo.Context) error {
	var videocard models.Videocard
	e := c.Bind(&videocard)
	if e != nil {
		//c.Logger().Error(err) //а мне все это надо ли надо ли 50 на 50
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	e = delivery.VideocardUC.CreateVideocard(&videocard)
	if e != nil {
		if e == models.ErrConflict {
			//c.Logger().Error(err)
			return c.JSON(http.StatusConflict, videocard)
		} else {
			//c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
		}
	}

	return c.JSON(http.StatusCreated, videocard)
}

func (delivery *Delivery) UpdateVideocard(c echo.Context) error {
	var videocard models.Videocard
	e := c.Bind(&videocard)
	if e != nil {
		fmt.Println(e.Error())
		//c.Logger().Error(err) //а мне все это надо ли надо ли 50 на 50
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	id, e := strconv.ParseInt(c.Param("videocardid"), 10, 64)

	if e != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	e = delivery.VideocardUC.UpdateVideocard(&videocard, id)
	if e != nil {
		if e == models.ErrNotFound { //возвращает ли err not found возвращает
			//c.Logger().Error(e)
			return echo.NewHTTPError(http.StatusNotFound, models.ErrNotFound.Error())
		} else {
			//c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
		}
	}

	return c.JSON(http.StatusOK, videocard)

}

// videocardid с маленькой? да
func (delivery *Delivery) SelectVideocard(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("videocardid"), 10, 64)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}
	videocard, e := delivery.VideocardUC.SelectVideocard(id)
	if e != nil {
		if e == models.ErrNotFound { //возвращает ли err not found возвращает
			//c.Logger().Error(e)
			return echo.NewHTTPError(http.StatusNotFound, models.ErrNotFound.Error())
		} else {
			//c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
		}
	}

	return c.JSON(http.StatusOK, videocard)
}

func (delivery *Delivery) ShowFullVideocard(c echo.Context) error {
	videocard, e := delivery.VideocardUC.ShowFullVideocard()
	if e != nil {
		//c.Logger().Error(e)
		return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
	}

	return c.JSON(http.StatusOK, videocard)
}

func (delivery *Delivery) ShowPartVideocard(c echo.Context) error {
	videocard, e := delivery.VideocardUC.ShowPartVideocard()
	if e != nil {
		//c.Logger().Error(e)
		return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
	}

	return c.JSON(http.StatusOK, videocard)

} //вот тут что делать с двумя переменными и что returнить
func (delivery *Delivery) DeleteVideocard(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("videocardid"), 10, 64)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}
	e := delivery.VideocardUC.DeleteVideocardById(id)
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

func NewDelivery(e *echo.Echo, videocardUC videocardUsecase.UseCaseI) {
	handler := &Delivery{
		VideocardUC: videocardUC,
	}
	//вот тут что
	//e.POST("/videocard/create", handler.CreateVideocard)
	//e.PUT("/videocard/:videocardid/update", handler.UpdateVideocard)
	//e.GET("/videocard/:videocardid/select", handler.SelectVideocard)
	e.GET("/videocard/showf", handler.ShowFullVideocard)
	//e.GET("/videocard/showp", handler.ShowPartVideocard)
	//e.DELETE("/videocard/:videocardid/delete", handler.DeleteVideocard)
	e.POST("/videocard/create", middl.AuthMiddleware(handler.CreateVideocard, "admin"))
	e.PUT("/videocard/:videocardid/update", middl.AuthMiddleware(handler.UpdateVideocard, "admin"))
	e.GET("/videocard/:videocardid/select", middl.AuthMiddleware(handler.SelectVideocard, "admin", "user"))
	//e.GET("/videocard/showf", middl.AuthMiddleware(handler.ShowFullVideocard, "admin", "user"))
	e.GET("/videocard/showp", handler.ShowPartVideocard)
	e.DELETE("/videocard/:videocardid/delete", middl.AuthMiddleware(handler.DeleteVideocard, "admin"))
}
