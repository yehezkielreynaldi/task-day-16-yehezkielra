package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"task-day-16/connection"
	"task-day-16/middleware"
	"text/template"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type Project struct {
	Id           int
	ProjectName  string
	StartDate    time.Time
	EndDate      time.Time
	Duration     string
	Description  string
	Technologies []string
	Image        string
	Tech1        bool
	Tech2        bool
	Tech3        bool
	Tech4        bool
	FormatStart  string
	FormatEnd    string
	Author       string
	UserID       int
}

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}

type SessionData struct {
	IsLogin bool
	Name    string
}

var userData = SessionData{}

// var dataProject = []Project{
// {
// 	ProjectName: "Project 1",
// 	StartDate:   "2023-05-01",
// 	EndDate:     "2023-06-01",
// 	Duration:    "1 Bulan",
// 	Description: "Ini Project 1",
// 	Tech1:       true,
// 	Tech2:       true,
// 	Tech3:       true,
// 	Tech4:       true,
// },
// {
// 	ProjectName: "Project 2",
// 	// StartDate:   "2023-05-02",
// 	// EndDate:     "2023-06-02",
// 	Duration:    "1 Bulan",
// 	Description: "Ini Project 2",
// 	Tech1:       true,
// 	Tech2:       true,
// 	Tech3:       true,
// 	Tech4:       true,
// },
// {
// 	ProjectName: "Project 3",
// 	StartDate:   "2023-05-03",
// 	EndDate:     "2023-06-03",
// 	Duration:    "1 Bulan",
// 	Description: "Ini Project 3",
// 	Tech1:       true,
// 	Tech2:       true,
// 	Tech3:       true,
// 	Tech4:       true,
// },
// {
// 	ProjectName: "Project 4",
// 	StartDate:   "2023-05-04",
// 	EndDate:     "2023-06-04",
// 	Duration:    "1 Bulan",
// 	Description: "Ini Project 4",
// 	Tech1:       false,
// 	Tech2:       false,
// 	Tech3:       true,
// 	Tech4:       true,
// },
// {
// 	ProjectName: "Project 5",
// 	StartDate:   "2023-05-05",
// 	EndDate:     "2023-06-05",
// 	Duration:    "1 Bulan",
// 	Description: "Ini Project 5",
// 	Tech1:       true,
// 	Tech2:       false,
// 	Tech3:       true,
// 	Tech4:       false,
// },
// {
// 	ProjectName: "Project 6",
// 	StartDate:   "2023-05-06",
// 	EndDate:     "2023-06-06",
// 	Duration:    "1 Bulan",
// 	Description: "Ini Project 6",
// 	Tech1:       true,
// 	Tech2:       false,
// 	Tech3:       true,
// 	Tech4:       true,
// },
// }

func main() {
	connection.DatabaseConnect()

	e := echo.New()

	// e = echo package
	// GET/POST = run the method
	// "/" = endpoint/routing (ex. localhost:5000'/' | ex. dumbways.id'/lms')
	// helloWorld = function that will run if the routes are opened

	// Serve a static files from "public" directory
	e.Static("/public", "public")
	e.Static("/uploads", "uploads")

	// To use sessions using echo
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("session"))))

	// Routing

	// GET
	e.GET("/", home)
	e.GET("/contact", contact)
	e.GET("/my-project", myProject)
	e.GET("/project-detail/:id", projectDetail)
	e.GET("/testimonials", testimonials)
	e.GET("/update-project/:id", updateMyProject)

	// Register
	e.GET("/form-register", formRegister)
	e.POST("/register", register)

	// login
	e.GET("/form-login", formLogin)
	e.POST("/login", login)

	// Logout
	e.POST("/logout", logout)

	// POST
	e.POST("/add-project", middleware.UploadFile(addProject))
	e.POST("/project-delete/:id", deleteProject)
	e.POST("/update-project/:id", middleware.UploadFile(updateProject))

	e.Logger.Fatal(e.Start("localhost:5000"))
}

func home(c echo.Context) error {
	sess, _ := session.Get("session", c)
	var result []Project

	if sess.Values["isLogin"] != true {
		userData.IsLogin = false
		data, _ := connection.Conn.Query(context.Background(), "SELECT tb_proyek.id, tb_proyek.name, start_date, end_date, duration, description, tech1, tech2, tech3, tech4, image, tb_users.name AS author FROM tb_proyek JOIN tb_users ON tb_proyek.authorid = tb_users.id ORDER BY tb_proyek.id DESC")
		fmt.Println(data)

		for data.Next() {
			var each = Project{}

			err := data.Scan(&each.Id, &each.ProjectName, &each.StartDate, &each.EndDate, &each.Duration, &each.Description, &each.Tech1, &each.Tech2, &each.Tech3, &each.Tech4, &each.Image, &each.Author)
			if err != nil {
				fmt.Println(err.Error())
				return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
			}
			fmt.Println(each)

			each.FormatStart = each.StartDate.Format("2 January 2006")
			each.FormatEnd = each.EndDate.Format("2 January 2006")

			result = append(result, each)
		}

	} else {
		userData.IsLogin = sess.Values["isLogin"].(bool)
		userData.Name = sess.Values["name"].(string)
		userId := sess.Values["id"]
		data, _ := connection.Conn.Query(context.Background(), "SELECT tb_proyek.id, tb_proyek.name, start_date, end_date, duration, description, tech1, tech2, tech3, tech4, image, tb_users.name AS author FROM tb_proyek JOIN tb_users ON tb_proyek.authorid = tb_users.id WHERE tb_proyek.authorid=$1 ORDER BY tb_proyek.id DESC", userId)

		fmt.Println(data)
		for data.Next() {
			var each = Project{}

			err := data.Scan(&each.Id, &each.ProjectName, &each.StartDate, &each.EndDate, &each.Duration, &each.Description, &each.Tech1, &each.Tech2, &each.Tech3, &each.Tech4, &each.Image, &each.Author)
			if err != nil {
				fmt.Println(err.Error())
				return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
			}
			fmt.Println(each)

			each.FormatStart = each.StartDate.Format("2 January 2006")
			each.FormatEnd = each.EndDate.Format("2 January 2006")

			result = append(result, each)
		}

	}

	projects := map[string]interface{}{
		"Projects":     result,
		"FlashStatus":  sess.Values["status"],
		"FlashMessage": sess.Values["message"],
		"DataSession":  userData,
	}

	delete(sess.Values, "message")
	delete(sess.Values, "status")
	sess.Save(c.Request(), c.Response())

	var tmpl, err = template.ParseFiles("views/index.html")

	if err != nil { // null
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), projects)
}

func contact(c echo.Context) error {

	sess, _ := session.Get("session", c)

	if sess.Values["isLogin"] != true {
		userData.IsLogin = false
	} else {
		userData.IsLogin = sess.Values["isLogin"].(bool)
		userData.Name = sess.Values["name"].(string)
	}

	projects := map[string]interface{}{
		"DataSession": userData,
	}

	sess.Save(c.Request(), c.Response())

	var tmpl, err = template.ParseFiles("views/contact.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"messsage": err.Error()})
	}

	return tmpl.Execute(c.Response(), projects)
}

func myProject(c echo.Context) error {
	sess, _ := session.Get("session", c)

	if sess.Values["isLogin"] != true {
		userData.IsLogin = false
	} else {
		userData.IsLogin = sess.Values["isLogin"].(bool)
		userData.Name = sess.Values["name"].(string)
	}

	projects := map[string]interface{}{
		"DataSession": userData,
	}

	sess.Save(c.Request(), c.Response())

	var tmpl, err = template.ParseFiles("views/my-project.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), projects)
}

func projectDetail(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	var ProjectDetail = Project{}

	err := connection.Conn.QueryRow(context.Background(), "SELECT tb_proyek.id, tb_proyek.name, start_date, end_date, duration, description, tech1, tech2, tech3, tech4, image, tb_users.name AS author FROM tb_proyek JOIN tb_users ON tb_proyek.authorid = tb_users.id WHERE tb_proyek.id=$1", id).Scan(
		&ProjectDetail.Id, &ProjectDetail.ProjectName, &ProjectDetail.StartDate, &ProjectDetail.EndDate, &ProjectDetail.Duration, &ProjectDetail.Description, &ProjectDetail.Tech1, &ProjectDetail.Tech2, &ProjectDetail.Tech3, &ProjectDetail.Tech4, &ProjectDetail.Image, &ProjectDetail.Author)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	ProjectDetail.FormatStart = ProjectDetail.StartDate.Format("2 January 2006")
	ProjectDetail.FormatEnd = ProjectDetail.EndDate.Format("2 January 2006")

	sess, _ := session.Get("session", c)

	if sess.Values["isLogin"] != true {
		userData.IsLogin = false
	} else {
		userData.IsLogin = sess.Values["isLogin"].(bool)
		userData.Name = sess.Values["name"].(string)
	}

	data := map[string]interface{}{
		"Project":     ProjectDetail,
		"DataSession": userData,
	}

	var tmpl, errTemplate = template.ParseFiles("views/project-detail.html")

	if errTemplate != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), data)
}

func testimonials(c echo.Context) error {
	sess, _ := session.Get("session", c)

	if sess.Values["isLogin"] != true {
		userData.IsLogin = false
	} else {
		userData.IsLogin = sess.Values["isLogin"].(bool)
		userData.Name = sess.Values["name"].(string)
	}

	projects := map[string]interface{}{
		"DataSession": userData,
	}

	sess.Save(c.Request(), c.Response())

	var tmpl, err = template.ParseFiles("views/testimonials.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), projects)
}

func updateMyProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	var ProjectDetail = Project{}

	err := connection.Conn.QueryRow(context.Background(), "SELECT id, name, start_date, end_date, duration, description, tech1, tech2, tech3, tech4, image FROM tb_proyek WHERE id=$1", id).Scan(&ProjectDetail.Id, &ProjectDetail.ProjectName, &ProjectDetail.StartDate, &ProjectDetail.EndDate, &ProjectDetail.Duration, &ProjectDetail.Description, &ProjectDetail.Tech1, &ProjectDetail.Tech2, &ProjectDetail.Tech3, &ProjectDetail.Tech4, &ProjectDetail.Image)

	data := map[string]interface{}{
		"Project": ProjectDetail,
	}

	var tmpl, errTmplt = template.ParseFiles("views/update-project.html")

	if errTmplt != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), data)
}

func addProject(c echo.Context) error {
	projectName := c.FormValue("projectName")
	startDate := c.FormValue("startDate")
	endDate := c.FormValue("endDate")
	duration := hitungDurasi(startDate, endDate)
	description := c.FormValue("descriptionProject")
	tech1 := (c.FormValue("tech1") == "tech1")
	tech2 := (c.FormValue("tech2") == "tech2")
	tech3 := (c.FormValue("tech3") == "tech3")
	tech4 := (c.FormValue("tech4") == "tech4")
	// image := c.FormValue("input-image")

	image := c.Get("dataFile").(string)

	sess, _ := session.Get("session", c)

	author := sess.Values["id"].(int)

	_, err := connection.Conn.Exec(context.Background(), "INSERT INTO tb_proyek (name, start_date, end_date, duration, description, tech1, tech2, tech3, tech4, image, authorid) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)", projectName, startDate, endDate, duration, description, tech1, tech2, tech3, tech4, image, author)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	// fmt.Println(dataProject)

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func deleteProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	fmt.Println("Index : ", id)

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_proyek WHERE id=$1", id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func updateProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	fmt.Println("Index :", id)

	projectName := c.FormValue("projectName")
	startDate := c.FormValue("startDate")
	endDate := c.FormValue("endDate")
	duration := hitungDurasi(startDate, endDate)
	description := c.FormValue("descriptionProject")
	tech1 := (c.FormValue("tech1") == "tech1")
	tech2 := (c.FormValue("tech2") == "tech2")
	tech3 := (c.FormValue("tech3") == "tech3")
	tech4 := (c.FormValue("tech4") == "tech4")

	image := c.Get("dataFile").(string)

	sess, _ := session.Get("session", c)

	author := sess.Values["id"].(int)

	_, err := connection.Conn.Exec(context.Background(), "UPDATE tb_proyek SET name=$1, start_date=$2, end_date=$3, duration=$4, description=$5, tech1=$6, tech2=$7, tech3=$8, tech4=$9, authorid=$10, image=$11 WHERE id=$12", projectName, startDate, endDate, duration, description, tech1, tech2, tech3, tech4, author, image, id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func hitungDurasi(startDate, endDate string) string {
	startTime, _ := time.Parse("2006-01-02", startDate)
	endTime, _ := time.Parse("2006-01-02", endDate)

	durationTime := int(endTime.Sub(startTime).Hours())
	durationDays := durationTime / 24
	durationWeeks := durationDays / 7
	durationMonths := durationWeeks / 4
	durationYears := durationMonths / 12

	var duration string

	if durationYears > 1 {
		duration = strconv.Itoa(durationYears) + " Tahun"
	} else if durationYears > 0 {
		duration = strconv.Itoa(durationYears) + " Tahun"
	} else {
		if durationMonths > 1 {
			duration = strconv.Itoa(durationMonths) + " Bulan"
		} else if durationMonths > 0 {
			duration = strconv.Itoa(durationMonths) + " Bulan"
		} else {
			if durationWeeks > 1 {
				duration = strconv.Itoa(durationWeeks) + " Minggu"
			} else if durationWeeks > 0 {
				duration = strconv.Itoa(durationWeeks) + " Minggu"
			} else {
				if durationDays > 1 {
					duration = strconv.Itoa(durationDays) + " Hari"
				} else {
					duration = strconv.Itoa(durationDays) + " Hari"
				}
			}
		}
	}

	return duration
}

func formRegister(c echo.Context) error {
	var tmpl, err = template.ParseFiles("views/form-register.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), nil)
}

func register(c echo.Context) error {
	// to make sure request body is form data format, not JSON, XML, etc.
	err := c.Request().ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	name := c.FormValue("inputName")
	email := c.FormValue("inputEmail")
	password := c.FormValue("inputPassword")

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_users(name, email, password) VALUES ($1, $2, $3)", name, email, passwordHash)

	if err != nil {
		redirectWithMessage(c, "Register failed, please try again.", false, "/form-register")
	}

	return redirectWithMessage(c, "Register success!", true, "/form-login")
}

func formLogin(c echo.Context) error {
	sess, _ := session.Get("session", c)

	flash := map[string]interface{}{
		"FlashStatus":  sess.Values["status"],
		"FlashMessage": sess.Values["message"],
	}

	delete(sess.Values, "message")
	delete(sess.Values, "status")
	sess.Save(c.Request(), c.Response())

	var tmpl, err = template.ParseFiles("views/form-login.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), flash)
}

func login(c echo.Context) error {
	err := c.Request().ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	email := c.FormValue("inputEmail")
	password := c.FormValue("inputPassword")

	user := User{}
	err = connection.Conn.QueryRow(context.Background(), "SELECT * FROM tb_users WHERE email=$1", email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return redirectWithMessage(c, "Email Incorrect!", false, "/form-login")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return redirectWithMessage(c, "Password Incorrect!", false, "/form-login")
	}

	sess, _ := session.Get("session", c)
	sess.Options.MaxAge = 10800 // 3 JAM
	sess.Values["message"] = "Login success!"
	sess.Values["status"] = true
	sess.Values["name"] = user.Name
	sess.Values["email"] = user.Email
	sess.Values["id"] = user.ID
	sess.Values["isLogin"] = true
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func logout(c echo.Context) error {
	sess, _ := session.Get("session", c)
	sess.Options.MaxAge = -1
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func redirectWithMessage(c echo.Context, message string, status bool, path string) error {
	sess, _ := session.Get("session", c)
	sess.Values["message"] = message
	sess.Values["status"] = status
	sess.Save(c.Request(), c.Response())
	return c.Redirect(http.StatusMovedPermanently, path)
}
