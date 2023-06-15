package main

import (
	coolingDelivery "Kursach_project/apiii/src/cooling/delivery"
	coolingRep "Kursach_project/apiii/src/cooling/repository"
	coolingUsecase "Kursach_project/apiii/src/cooling/usecase"
	cpuDelivery "Kursach_project/apiii/src/cpu/delivery"
	cpuRep "Kursach_project/apiii/src/cpu/repository"
	cpuUsecase "Kursach_project/apiii/src/cpu/usecase"
	motherboardDelivery "Kursach_project/apiii/src/motherboard/delivery"
	motherboardRep "Kursach_project/apiii/src/motherboard/repository"
	motherboardUsecase "Kursach_project/apiii/src/motherboard/usecase"
	pcbuildDelivery "Kursach_project/apiii/src/pcbuild/delivery"
	pcbuildRep "Kursach_project/apiii/src/pcbuild/repository"
	pcbuildUsecase "Kursach_project/apiii/src/pcbuild/usecase"
	powerboxDelivery "Kursach_project/apiii/src/powerbox/delivery"
	powerboxRep "Kursach_project/apiii/src/powerbox/repository"
	powerboxUsecase "Kursach_project/apiii/src/powerbox/usecase"
	ramDelivery "Kursach_project/apiii/src/ram/delivery"
	ramRep "Kursach_project/apiii/src/ram/repository"
	ramUsecase "Kursach_project/apiii/src/ram/usecase"
	userDelivery "Kursach_project/apiii/src/user/delivery"
	userRep "Kursach_project/apiii/src/user/repository"
	userUsecase "Kursach_project/apiii/src/user/usecase"
	videocardDelivery "Kursach_project/apiii/src/videocard/delivery"
	videocardRep "Kursach_project/apiii/src/videocard/repository"
	videocardUsecase "Kursach_project/apiii/src/videocard/usecase"
	"fmt"
	_ "github.com/gorilla/mux"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	elog "github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
	_ "log"
)

func main() {
	//db, err := sql.Open("postgres", fmt.Sprintf("host=localhost"))
	//if err != nil {
	//	log.Fatal(err)
	//}
	c := connect()
	if c != nil {
		fmt.Println(c.Error())

		return
	}

	defer db.Close()

	//forumDB := forumRep.New(db)
	//userDB := userRep.New(db)
	//postDB := postRep.New(db)
	//threadDB := threadRep.New(db)
	//serviceDB := serviceRep.New(db)

	//cpuRepo := cpuRep.New(db)
	//coolingRepo := coolingRep.New(db)
	//motherboardRepo := motherboardRep.New(db)
	//powerboxRepo := powerboxRep.New(db)
	//ramRepo := ramRep.New(db)
	//userRep := userRep.New(db)
	//videocardRepo := videocardRep.New(db)
	//
	//m := mux.NewRouter()
	//
	////cpuHandler.NewDelivery(m, driverUcase)
	////teamHandler.NewTeamHandler(m, teamUcase)
	////grandPrixHandler.NewDriverHandler(m, gpUcase)
	//
	//mMiddleware := AccessLogMiddleware(m)
	//
	//fmt.Println("starting server at :8080")
	//http.ListenAndServe(":8080", mMiddleware)

	cpuRepo := cpuRep.New(db)
	coolingRepo := coolingRep.New(db)
	motherboardRepo := motherboardRep.New(db)
	powerboxRepo := powerboxRep.New(db)
	ramRepo := ramRep.New(db)
	userRepo := userRep.New(db)
	videocardRepo := videocardRep.New(db)
	pcbuildRepo := pcbuildRep.New(db)

	cpuUC := cpuUsecase.New(cpuRepo, ramRepo)
	userUC := userUsecase.New(userRepo)
	coolingUC := coolingUsecase.New(coolingRepo)
	motherboardUC := motherboardUsecase.New(motherboardRepo)
	powerboxUC := powerboxUsecase.New(powerboxRepo)
	ramUC := ramUsecase.New(ramRepo, cpuRepo)
	videocardUC := videocardUsecase.New(videocardRepo)
	pcbuildUC := pcbuildUsecase.New(pcbuildRepo, cpuRepo, ramRepo, videocardRepo, powerboxRepo, motherboardRepo, coolingRepo, userRepo)

	//m := mux.NewRouter()

	e := echo.New()

	e.Logger.SetHeader(`time=${time_rfc3339} level=${level} prefix=${prefix} ` +
		`file=${short_file} line=${line} message:`)
	e.Logger.SetLevel(elog.INFO)

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `time=${time_custom} remote_ip=${remote_ip} ` +
			`host=${host} method=${method} uri=${uri} user_agent=${user_agent} ` +
			`status=${status} error="${error}" ` +
			`bytes_in=${bytes_in} bytes_out=${bytes_out}` + "\n",
		CustomTimeFormat: "2006-01-02 15:04:05",
	}))

	e.Use(middleware.Recover())

	cpuDelivery.NewDelivery(e, cpuUC)
	coolingDelivery.NewDelivery(e, coolingUC)
	motherboardDelivery.NewDelivery(e, motherboardUC)
	pcbuildDelivery.NewDelivery(e, pcbuildUC)
	powerboxDelivery.NewDelivery(e, powerboxUC)
	ramDelivery.NewDelivery(e, ramUC)
	userDelivery.NewDelivery(e, userUC)
	videocardDelivery.NewDelivery(e, videocardUC)
	//postDelivery.NewDelivery(e, postUC)
	//threadDelivery.NewDelivery(e, threadUC)
	//serviceDelivery.NewDelivery(e, serviceUC)
	//e.POST("/login", func(c echo.Context) error {
	//	var user models.User
	//	ec := c.Bind(&user)
	//	if ec != nil {
	//		fmt.Println(ec.Error())
	//		//c.Logger().Error(err) //а мне все это надо ли надо ли 50 на 50
	//		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	//	}
	//	// Проверяем, есть ли такой пользователь в базе данных
	//	err := db.QueryRow(`SELECT userid, role FROM users WHERE login=$1 AND password=$2`, user.Login, user.Password).Scan(&user.Userid, &user.Role)
	//	//row := db.QueryRow(query, user.Login, user.Password)
	//
	//	if err != nil {
	//		return c.String(http.StatusUnauthorized, models.ERRUnauthorized.Error())
	//	}
	//
	//	login := user.Login
	//	password := user.Password
	//	role := user.Role
	//	cookies := []string{login, password, role}
	//	cookieString := strings.Join(cookies, " ")
	//
	//	cookie := http.Cookie{
	//		Name:   "oreo",
	//		Value:  cookieString,
	//		MaxAge: 60 * 60 * 24,
	//		Path:   "/",
	//	}
	//	//http.SetCookie(w, &cookie)
	//	echoCookie := new(http.Cookie)
	//	echoCookie.Name = cookie.Name
	//	echoCookie.Value = cookie.Value
	//	echoCookie.Path = cookie.Path
	//	if cookie.MaxAge > 0 {
	//		echoCookie.Expires = time.Now().Add(time.Duration(cookie.MaxAge) * time.Second)
	//	}
	//	c.SetCookie(echoCookie)
	//	return c.String(http.StatusOK, "authorized")
	//})

	s := NewServer(e)
	if err := s.Start(); err != nil {
		e.Logger.Fatal(err)
	}
}

//func adminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
//	return func(c echo.Context) error {
//	// Получаем из контекста логин пользователя
//	login := c.Get("login").(string)
//
//	// Ищем пользователя в базе данных
//	row := db.QueryRow("SELECT role FROM users WHERE login = $1", login)
//	var role string
//	err := row.Scan(&role)
//	if err != nil {
//	return c.JSON(http.StatussrcServerError, err.Error())
//}
//
//	// Проверяем, что роль пользователя является админом или менеджером
//	if role != "admin" && role != "manager" {
//	return c.JSON(http.StatusUnauthorized, "Access denied.")
//}
//
//	// Если пользователь имеет права на выполнение маршрута, то выполняем его
//	return next(c)
//}
//}

//dr, _ := db.Query(`SELECT * FROM videocard`)
//if e != nil {
//	fmt.Println(e.Error())
//}
//var videocards string
//dr.Scan(&videocards)

//db.Exec(`INSERT INTO "videocard" ("vmemory", "videocardid") VALUES ('3400', '1')`)
//http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
//	fmt.Println(cfg.ServerHost + ":" + cfg.ServerPort)
//})
//http.ListenAndServe(":8080", nil)

//	rows, e := db.Query(`select pc_buildid, cpu.name as "Процессор", ram.name as "Оперативная Память", cooling.name as "Охлаждение", powerbox.name as "Блок питания", motherboard.name as "Материнская плата", videocard.name as "Видеокарта", compatibility, pc_build.price as "Стоимость"
//from pc_build
//join cpu on pc_build.cpu_id = cpu.cpuid
//join ram on pc_build.ram_id = ram.ramid
//join cooling on pc_build.cooling_id = cooling.coolingid
//join powerbox on pc_build.powerbox_id = powerbox.powerboxid
//join motherboard on pc_build.motherboard_id = motherboard.motherboardid
//join videocard on pc_build.videocard_id = videocard.videocardid`)
//	if e != nil {
//		//fmt.Println(e.Error())
//		fmt.Println("параша")
//		return
//	}
//	defer rows.Close()
//
//	messiv := make([]*models.Pcbuildinfo, 0)
//	for rows.Next() {
//		var pcbuildinfo models.Pcbuildinfo
//		e = rows.Scan(&pcbuildinfo.Pcbuildid, &pcbuildinfo.Cpuname, &pcbuildinfo.Ramname, &pcbuildinfo.Coolingname, &pcbuildinfo.Powerboxname, &pcbuildinfo.Motherboardname, &pcbuildinfo.Videocardname, &pcbuildinfo.Compatibility, &pcbuildinfo.Price)
//		if e != nil {
//			fmt.Println(e.Error())
//			fmt.Println("ппараша")
//			return
//		}
//		messiv = apiiiend(messiv, &pcbuildinfo)
//		//fmt.Println(pcbuildinfo.Cpuname)
//	}
//
//	fmt.Println(messiv)
//}
