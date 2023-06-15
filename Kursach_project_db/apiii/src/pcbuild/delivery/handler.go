package delivery

import (
	"Kursach_project/apiii/middl"
	"Kursach_project/apiii/models"
	pcbuildUsecase "Kursach_project/apiii/src/pcbuild/usecase"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type Delivery struct {
	PcbuildUC pcbuildUsecase.UseCaseI
}

func (delivery *Delivery) CreatePcbuild(c echo.Context) error {
	var pcbuild models.Pcbuild
	e := c.Bind(&pcbuild)
	if e != nil {
		//c.Logger().Error(err) //а мне все это надо ли надо ли 50 на 50
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	userid, err := middl.GetIdFromCookie(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, models.ErrBadRequest.Error()) //errunauthorized
	}

	e = delivery.PcbuildUC.CreatePcbuild(&pcbuild, userid)
	if e != nil {
		if e == models.ErrConflict {
			//c.Logger().Error(err)
			return c.JSON(http.StatusConflict, pcbuild)
		} else {
			//c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
		}
	}

	return c.JSON(http.StatusCreated, pcbuild)
}

func (delivery *Delivery) UpdatePcbuild(c echo.Context) error {
	var pcbuild models.Pcbuild
	e := c.Bind(&pcbuild)
	if e != nil {
		fmt.Println(e.Error())
		//c.Logger().Error(err) //а мне все это надо ли надо ли 50 на 50
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	pcbuildid, e := strconv.ParseInt(c.Param("pc_buildid"), 10, 64)

	if e != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}
	userid, err := middl.GetIdFromCookie(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, models.ERRUnauthorized.Error()) //errunauthorized
	}
	e = delivery.PcbuildUC.UpdatePcbuild(&pcbuild, userid, pcbuildid)
	if e != nil {
		if e == models.ErrNotFound { //возвращает ли err not found возвращает
			//c.Logger().Error(e)
			return echo.NewHTTPError(http.StatusNotFound, models.ErrNotFound.Error())
		} else {
			//c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
		}
	}

	return c.JSON(http.StatusOK, pcbuild)

}

// pcbuildid с маленькой? да
func (delivery *Delivery) SelectPcbuild(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("pc_buildid"), 10, 64)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}
	pcbuild, e := delivery.PcbuildUC.SelectPcbuildById(id)
	if e != nil {
		if e == models.ErrNotFound { //возвращает ли err not found возвращает
			//c.Logger().Error(e)
			return echo.NewHTTPError(http.StatusNotFound, models.ErrNotFound.Error())
		} else {
			//c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
		}
	}

	return c.JSON(http.StatusOK, pcbuild)
}

func (delivery *Delivery) ShowFullPcbuild(c echo.Context) error {
	pcbuild, e := delivery.PcbuildUC.ShowFullPcbuild()
	if e != nil {
		//c.Logger().Error(e)
		return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
	}

	return c.JSON(http.StatusOK, pcbuild)
}

func (delivery *Delivery) ShowMyPcbuild(c echo.Context) error {
	userid, err := middl.GetIdFromCookie(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, models.ErrBadRequest.Error()) //errunauthorized
	}

	pcbuild, e := delivery.PcbuildUC.ShowMyPcbuild(userid)
	if e != nil {
		//c.Logger().Error(e)
		return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
	}

	return c.JSON(http.StatusOK, pcbuild)
}

func (delivery *Delivery) ShowPcbuildById(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("pc_buildid"), 10, 64)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}
	pcbuild, e := delivery.PcbuildUC.ShowPcbuildById(id)
	if e != nil {
		//c.Logger().Error(e)
		return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
	}

	return c.JSON(http.StatusOK, pcbuild)

}

// вот тут что делать с двумя переменными и что returнить
func (delivery *Delivery) DeletePcbuild(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("pc_buildid"), 10, 64)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}
	userid, err := middl.GetIdFromCookie(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, models.ERRUnauthorized.Error()) //errunauthorized
	}
	e := delivery.PcbuildUC.DeletePcbuildById(userid, id)
	if e != nil {
		if e == models.ErrNotFound { //проверить вроде возвращается
			//c.Logger().Error(e)
			return echo.NewHTTPError(http.StatusNotFound, models.ErrNotFound.Error())
		} else {
			//`c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
		}
	}

	return c.NoContent(http.StatusOK) //мб заменить на http.StatusNoContent
}

func NewDelivery(e *echo.Echo, pcbuildUC pcbuildUsecase.UseCaseI) {
	handler := &Delivery{
		PcbuildUC: pcbuildUC,
	}
	//вот тут что
	e.POST("/pcbuild/create", handler.CreatePcbuild)
	e.PUT("/pcbuild/:pc_buildid/update", handler.UpdatePcbuild)
	e.GET("/pcbuild/:pc_buildid/select", handler.UpdatePcbuild)
	e.GET("/pcbuild/showmy", handler.ShowMyPcbuild)
	e.GET("/pcbuild/showf", handler.ShowFullPcbuild)
	e.GET("/pcbuild/:pc_buildid/showid", handler.ShowPcbuildById)
	e.DELETE("/pcbuild/:pc_buildid/delete", handler.DeletePcbuild)
	//e.POST("/pcbuild/create", middl.AuthMiddleware(handler.CreatePcbuild, "admin", "user"))
	//e.PUT("/pcbuild/:pc_buildid/update", middl.AuthMiddleware(handler.UpdatePcbuild, "admin", "user"))
	//e.GET("/pcbuild/:pc_buildid/select", middl.AuthMiddleware(handler.SelectPcbuild, "admin", "user"))
	//e.GET("/pcbuild/showmy", middl.AuthMiddleware(handler.ShowMyPcbuild, "admin", "user"))
	//e.GET("/pcbuild/showf", middl.AuthMiddleware(handler.ShowFullPcbuild, "admin", "user"))
	//e.GET("/pcbuild/:pc_buildid/showid", middl.AuthMiddleware(handler.ShowPcbuildById, "admin", "user"))
	//e.DELETE("/pcbuild/:pc_buildid/delete", middl.AuthMiddleware(handler.DeletePcbuild, "admin", "user"))
}
