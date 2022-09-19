package main

//go:generate sqlboiler --wipe psql

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ClubCedille/monitoring.serreets.com/api/models"
	"github.com/gin-contrib/cors"
	"github.com/golang-jwt/jwt"
	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"golang.org/x/crypto/bcrypt"
)

type Error struct {
	Message string
	Status  int
	Error   string
}

var (
	router           = gin.Default()
	connectionString = "dbname=" + os.Getenv("DATABASE_NAME") + " host=" + os.Getenv("DATABASE_HOST") + " user=" + os.Getenv("DATABASE_USER") + " password=" + os.Getenv("DATABASE_PASSWORD") + " sslmode=disable"
	db, err          = sql.Open("postgres", connectionString)
)

func main() {
	if err != nil {
		NewInternalServerError("Cannot connect to the database")
		panic(err)
	} else {
		fmt.Println("Connected to the database")
	}

	// CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{os.Getenv("ORIGINS")},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == os.Getenv("ORIGINS")
		},
		MaxAge: 12 * time.Hour,
	}))

	api := router.Group("/api")
	{
		api.POST("/signup", SignUp)
		api.POST("/login", LogIn)
		api.GET("/user", GetUserInformationWithToken)
		api.GET("/logout", LogOut)
	}

	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(NewPageNotFound("Page not found").Status, gin.H{})
	})

	router.Run(":" + os.Getenv("PORT"))
}

func GetUserInformationWithToken(ctx *gin.Context) {
	cookie, err := ctx.Cookie("jwt")
	if err != nil {
		errors := NewInternalServerError("Could not retrieve cookie")
		ctx.JSON(errors.Status, errors)
		return
	}

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(*jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		errors := NewInternalServerError("Error parsing the cookie")
		ctx.JSON(errors.Status, errors)
		return
	}

	claims := token.Claims.(*jwt.StandardClaims)
	issuer, err := strconv.Atoi(claims.Issuer)
	if err != nil {
		errors := NewBadRequestError("User not found")
		ctx.JSON(errors.Status, errors)
		return
	}

	user, restErr := GetUserById(issuer)
	if err != nil {
		ctx.JSON(restErr.Status, restErr)
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func GetUserById(userId int) (*models.User, *Error) {
	res, err := models.Users(qm.Where("id=?", userId)).One(context.Background(), db)
	if err != nil {
		return nil, NewBadRequestError("User not found")
	} else {
		resUser := &models.User{ID: res.ID, Email: res.Email, Password: ""}
		return resUser, nil
	}
}

func LogOut(ctx *gin.Context) {
	ctx.SetCookie("jwt", "", -1, "", "", false, true) // remove token
	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func SignUp(ctx *gin.Context) {
	var user models.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		err := NewBadRequestError("Invalid json body")
		ctx.JSON(err.Status, err)
		return
	}

	result, err := CreateUser(&user)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	} else {
		ctx.JSON(http.StatusOK, result)
	}
}

func CreateUser(user *models.User) (*models.User, *Error) {
	user.Email = strings.TrimSpace(user.Email)
	if user.Email == "" {
		return nil, NewBadRequestError("Invalid name")
	}

	user.Password = strings.TrimSpace(user.Password)
	if user.Password == "" {
		return nil, NewBadRequestError("Invalid password")
	}

	pwdSlice, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return nil, NewBadRequestError("Failed to encrypt password")
	}

	user.Password = string(pwdSlice[:])
	err = user.Insert(context.Background(), db, boil.Infer())
	if err != nil {
		NewBadRequestError("Cannot create the user")
	}

	return user, nil
}

func LogIn(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		err := NewBadRequestError("Invalid json")
		ctx.JSON(err.Status, err)
		return
	}

	currentUser, err := GetUserByEmail(&user)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}

	// JWT
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa((currentUser.ID)),
		ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
	})

	token, errors := claims.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if errors != nil {
		errors := NewInternalServerError("Login failed")
		ctx.JSON(errors.Status, errors)
		return
	}

	ctx.SetCookie("jwt", token, 3600, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, currentUser)
}

func GetUserByEmail(user *models.User) (*models.User, *Error) {
	res, err := models.Users(qm.Where("email=?", user.Email)).One(context.Background(), db)
	if err != nil {
		return nil, NewUnauthorized("Invalid username and password")
	} else {
		if err := bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(user.Password)); err != nil {
			return nil, NewUnauthorized("Invalid username and password")
		}
		resUser := &models.User{ID: res.ID, Email: res.Email, Password: ""}
		return resUser, nil
	}
}

func NewInternalServerError(message string) *Error {
	return &Error{
		Message: message,
		Status:  http.StatusInternalServerError,
		Error:   "internal server error",
	}
}

func NewBadRequestError(message string) *Error {
	return &Error{
		Message: message,
		Status:  http.StatusBadRequest,
		Error:   "bad request error",
	}
}

func NewPageNotFound(message string) *Error {
	return &Error{
		Message: message,
		Status:  http.StatusNotFound,
		Error:   "not found",
	}
}

func NewUnauthorized(message string) *Error {
	return &Error{
		Message: message,
		Status:  http.StatusUnauthorized,
		Error:   "unauthorized",
	}
}
