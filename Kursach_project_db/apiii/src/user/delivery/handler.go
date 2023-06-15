package delivery

import (
	"Kursach_project/apiii/middl"
	"Kursach_project/apiii/models"
	userUsecase "Kursach_project/apiii/src/user/usecase"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"time"
)

type Delivery struct {
	UserUC userUsecase.UseCaseI
}

func (delivery *Delivery) CreateUser(c echo.Context) error {
	var user models.User
	e := c.Bind(&user)
	if e != nil {
		//c.Logger().Error(err) //а мне все это надо ли надо ли 50 на 50
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	e = delivery.UserUC.CreateUser(&user)
	if e != nil {
		if e == models.ErrConflict {
			//c.Logger().Error(err)
			return c.JSON(http.StatusConflict, user)
		} else {
			//c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
		}
	}

	return c.JSON(http.StatusCreated, user)
}

func (delivery *Delivery) Logout(c echo.Context) error {

	cookie := http.Cookie{
		Name:    "oreo",
		Value:   "",
		Expires: time.Unix(1, 0),
		MaxAge:  -1,
		Path:    "/",
	}
	//http.SetCookie(w, &cookie)
	echoCookie := new(http.Cookie)
	echoCookie.Name = cookie.Name
	echoCookie.Value = cookie.Value
	echoCookie.Path = cookie.Path
	echoCookie.Expires = time.Now().Add(time.Duration(cookie.MaxAge) * time.Second)
	c.SetCookie(echoCookie)

	return c.JSON(http.StatusOK, "logout")
}

func (delivery *Delivery) Login(c echo.Context) error {
	var user models.User
	e := c.Bind(&user)
	if e != nil {
		//c.Logger().Error(err) //а мне все это надо ли надо ли 50 на 50
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	err := delivery.UserUC.Login(&user)
	if err != nil {
		if err == models.ErrNotFound {
			return echo.NewHTTPError(http.StatusUnauthorized, models.ErrNotFound.Error()) //изменить на ерранафтерайзд
		} else {
			return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
		}
	}

	//login := user.Login
	//password := user.Password
	role := user.Role
	id := user.Userid
	cookieString := fmt.Sprintf("%s %d", role, id)
	//cookies := []string{login, password, role}
	//cookieString := strings.Join(cookies, " ")

	cookie := http.Cookie{
		Name:   "oreo",
		Value:  cookieString,
		MaxAge: 60 * 60 * 24,
		Path:   "/",
	}
	//http.SetCookie(w, &cookie)
	echoCookie := new(http.Cookie)
	echoCookie.Name = cookie.Name
	echoCookie.Value = cookie.Value
	echoCookie.Path = cookie.Path
	if cookie.MaxAge > 0 {
		echoCookie.Expires = time.Now().Add(time.Duration(cookie.MaxAge) * time.Second)
	}
	c.SetCookie(echoCookie)

	return c.JSON(http.StatusOK, "authorized")
}

func (delivery *Delivery) UpdateUser(c echo.Context) error {
	var user models.User
	e := c.Bind(&user)
	if e != nil {
		fmt.Println(e.Error())
		//c.Logger().Error(err) //а мне все это надо ли надо ли 50 на 50
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	//id, err := middl.GetIdFromCookie(c)
	//if err != nil {
	//	return echo.NewHTTPError(http.StatusUnauthorized, models.ErrBadRequest.Error()) //errunauthorized
	//}

	id, e := strconv.ParseInt(c.Param("userid"), 10, 64)

	if e != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	e = delivery.UserUC.UpdateUser(&user, id)
	if e != nil {
		if e == models.ErrNotFound { //возвращает ли err not found возвращает
			//c.Logger().Error(e)
			return echo.NewHTTPError(http.StatusNotFound, models.ErrNotFound.Error())
		} else {
			//c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
		}
	}

	return c.JSON(http.StatusOK, user)

}

// userid с маленькой? да
func (delivery *Delivery) SelectUser(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("userid"), 10, 64)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}
	user, e := delivery.UserUC.SelectUser(id)
	if e != nil {
		if e == models.ErrNotFound { //возвращает ли err not found возвращает
			//c.Logger().Error(e)
			return echo.NewHTTPError(http.StatusNotFound, models.ErrNotFound.Error())
		} else {
			//c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
		}
	}

	return c.JSON(http.StatusOK, user)
}

func (delivery *Delivery) GetroleUser(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("userid"), 10, 64)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}
	userrole, e := delivery.UserUC.SelectUser(id)
	if e != nil {
		if e == models.ErrNotFound { //возвращает ли err not found возвращает
			//c.Logger().Error(e)
			return echo.NewHTTPError(http.StatusNotFound, models.ErrNotFound.Error())
		} else {
			//c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
		}
	}

	return c.JSON(http.StatusOK, userrole.Role) //c.json мб???
}

func (delivery *Delivery) ShowFullUser(c echo.Context) error {
	user, e := delivery.UserUC.ShowFullUser()
	if e != nil {
		//c.Logger().Error(e)
		return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
	}

	return c.JSON(http.StatusOK, user)
}

// вот тут что делать с двумя переменными и что returнить
func (delivery *Delivery) DeleteUser(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("userid"), 10, 64)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}
	e := delivery.UserUC.DeleteUserById(id)
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

func NewDelivery(e *echo.Echo, userUC userUsecase.UseCaseI) {
	handler := &Delivery{
		UserUC: userUC,
	}
	//вот тут что

	//e.POST("/user/create", handler.CreateUser)
	//e.POST("/login", handler.Login)
	//e.GET("/logout", handler.Logout)
	//e.PUT("/user/:userid/update", handler.UpdateUser)
	//e.GET("/user/:userid/select", handler.SelectUser)
	//e.GET("/user/showf", handler.ShowFullUser)
	//e.GET("/user/:userid/role", handler.GetroleUser)
	//e.DELETE("/user/:userid/delete", handler.DeleteUser)
	e.POST("/user/create", middl.AuthMiddleware(handler.CreateUser, "admin"))
	e.POST("/login", handler.Login)
	e.GET("/logout", handler.Logout)
	e.PUT("/user/:userid/update", middl.AuthMiddleware(handler.UpdateUser, "admin"))
	e.GET("/user/:userid/select", middl.AuthMiddleware(handler.SelectUser, "admin"))
	e.GET("/user/showf", middl.AuthMiddleware(handler.ShowFullUser, "admin"))
	e.GET("/user/:userid/role", middl.AuthMiddleware(handler.GetroleUser, "admin"))
	e.DELETE("/user/:userid/delete", middl.AuthMiddleware(handler.DeleteUser, "admin"))
}
