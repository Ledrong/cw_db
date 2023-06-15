package delivery

import (
	"Kursach_project/apiii/models"
	cpuUsecase "Kursach_project/apiii/src/cpu/usecase"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type Delivery struct {
	CpuUC cpuUsecase.UseCaseI
}

func (delivery *Delivery) CreateCpu(c echo.Context) error {
	var cpu models.Cpu
	e := c.Bind(&cpu)
	if e != nil {
		fmt.Println(e.Error())
		//c.Logger().Error(err) //а мне все это надо ли надо ли 50 на 50
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	e = delivery.CpuUC.CreateCpu(&cpu)
	if e != nil {
		if e == models.ErrConflict {
			//c.Logger().Error(err)
			return c.JSON(http.StatusConflict, cpu)
		} else {
			//c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
		}
	}

	return c.JSON(http.StatusCreated, cpu)
}

func (delivery *Delivery) UpdateCpu(c echo.Context) error {
	var cpu models.Cpu
	e := c.Bind(&cpu)
	if e != nil {
		fmt.Println(e.Error())
		//c.Logger().Error(err) //а мне все это надо ли надо ли 50 на 50
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	id, e := strconv.ParseInt(c.Param("cpuid"), 10, 64)

	if e != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	e = delivery.CpuUC.UpdateCpu(&cpu, id)
	if e != nil {
		if e == models.ErrNotFound { //возвращает ли err not found возвращает
			//c.Logger().Error(e)
			return echo.NewHTTPError(http.StatusNotFound, models.ErrNotFound.Error())
		} else {
			//c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
		}
	}

	return c.JSON(http.StatusOK, cpu)

}

//func (delivery *Delivery) CreateCpu(w http.ResponseWriter, r http.Request) {
//	decoder := json.NewDecoder(r.Body)
//	cpu := new(models.Cpu)
//	err := decoder.Decode(cpu)
//	if err != nil {
//		http.Error(w, "bad request", http.StatusBadRequest)
//		return
//	}
//	err := delivery.CpuUC.CreateCpu(cpu)
//	if err != nil {
//		http.Error(w, "db err", http.StatusInternalServerError)
//		return
//	}
//	encoder := json.NewEncoder(w)
//	err = encoder.Encode(id)
//	if err != nil {
//		http.Error(w, "encode err:", http.StatusInternalServerError)
//		return
//	}
//}

// cpuid с маленькой? да
func (delivery *Delivery) SelectCpu(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("cpuid"), 10, 64)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	cpu, e := delivery.CpuUC.SelectCpu(id)
	if e != nil {
		if e == models.ErrNotFound { //возвращает ли err not found возвращает
			//c.Logger().Error(e)
			return echo.NewHTTPError(http.StatusNotFound, models.ErrNotFound.Error())
		} else {
			//c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
		}
	}

	return c.JSON(http.StatusOK, cpu)
}

func (delivery *Delivery) ShowFullCpu(c echo.Context) error {
	cpu, e := delivery.CpuUC.ShowFullCpu()
	if e != nil {
		//c.Logger().Error(e)
		return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
	}

	return c.JSON(http.StatusOK, cpu)
}

func (delivery *Delivery) ShowCompatibilityRam(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("cpuid"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	ram, e := delivery.CpuUC.ShowCompatibilityRam(id)
	if e != nil {
		//c.Logger().Error(e)
		return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
	}

	return c.JSON(http.StatusOK, ram)
}

func (delivery *Delivery) ShowCompatibilityMotherboard(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("cpuid"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	motherboard, e := delivery.CpuUC.ShowCompatibilityMotherboard(id)
	if e != nil {
		//c.Logger().Error(e)
		return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
	}

	return c.JSON(http.StatusOK, motherboard)
}

func (delivery *Delivery) ShowPartCpu(c echo.Context) error {
	cpu, e := delivery.CpuUC.ShowPartCpu()
	if e != nil {
		//c.Logger().Error(e)
		return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
	}

	return c.JSON(http.StatusOK, cpu)

} //вот тут что делать с двумя переменными и что returнить
func (delivery *Delivery) DeleteCpu(c echo.Context) error {

	id, err := strconv.ParseInt(c.Param("cpuid"), 10, 64)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	e := delivery.CpuUC.DeleteCpuById(id)
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

func NewDelivery(e *echo.Echo, cpuUC cpuUsecase.UseCaseI) {
	handler := &Delivery{
		CpuUC: cpuUC,
	}
	//вот тут что
	e.POST("/cpu/create", handler.CreateCpu)
	e.PUT("/cpu/:cpuid/update", handler.UpdateCpu)
	e.GET("/cpu/:cpuid/select", handler.SelectCpu)
	e.GET("/cpu/showf", handler.ShowFullCpu)
	e.GET("/cpu/showp", handler.ShowPartCpu)
	e.GET("/cpu/:cpuid/showram", handler.ShowCompatibilityRam)
	e.GET("/cpu/:cpuid/showmotherboard", handler.ShowCompatibilityMotherboard)
	e.DELETE("/cpu/:cpuid/delete", handler.DeleteCpu)
	//e.POST("/cpu/create", middl.AuthMiddleware(handler.CreateCpu, "admin"))
	//e.PUT("/cpu/:cpuid/update", middl.AuthMiddleware(handler.UpdateCpu, "admin"))
	//e.GET("/cpu/:cpuid/select", middl.AuthMiddleware(handler.SelectCpu, "admin", "user"))
	//e.GET("/cpu/showf", middl.AuthMiddleware(handler.ShowFullCpu, "admin", "user"))
	//e.GET("/cpu/:cpuid/showram", middl.AuthMiddleware(handler.ShowCompatibilityRam, "admin", "user"))
	//e.GET("/cpu/:cpuid/showmotherboard", middl.AuthMiddleware(handler.ShowCompatibilityMotherboard, "admin", "user"))
	//e.GET("/cpu/showp", handler.ShowPartCpu)
	//e.DELETE("/cpu/:cpuid/delete", middl.AuthMiddleware(handler.DeleteCpu, "admin"))
}
