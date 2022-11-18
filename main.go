package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	database "github.com/snykk/go-echo-crud/database"
	"github.com/snykk/go-echo-crud/dto"
	"github.com/snykk/go-echo-crud/entity"
	"github.com/snykk/go-echo-crud/utils"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func root(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, echo.Map{
		"message": "Welcome to amazing api",
	})
}

func (r *Repository) CreateMovie(ctx echo.Context) error {
	var movie dto.Movie
	if err := ctx.Bind(&movie); err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, dto.BaseResponse{
			Status:  false,
			Message: err.Error(),
		})
	}

	if isValid, err := utils.IsPayloadValid(movie); !isValid {
		return ctx.JSON(http.StatusBadRequest, dto.BaseResponse{
			Status:  false,
			Message: err.Error(),
		})
	}

	entity := movie.ToEntity()
	err := r.DB.Create(&entity).Error
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, dto.BaseResponse{
			Status:  false,
			Message: "could not create movie",
		})
	}

	return ctx.JSON(http.StatusOK, dto.BaseResponse{
		Status:  true,
		Message: "movie data created successfully",
	})
}

func (r *Repository) GetMovies(ctx echo.Context) error {
	var movies []entity.Movie

	err := r.DB.Find(&movies).Error
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, dto.BaseResponse{
			Status:  false,
			Message: "could not get movies data",
		})
	}

	return ctx.JSON(http.StatusOK, dto.BaseResponse{
		Status:  true,
		Message: "movies data fetched successfully",
		Data:    movies,
	})
}

func (r *Repository) GetMovieById(ctx echo.Context) error {
	id := ctx.Param("id")
	var movie entity.Movie
	if id == "" {
		return ctx.JSON(http.StatusBadRequest, dto.BaseResponse{
			Status:  false,
			Message: "id cannot be empty",
		})
	}

	if err := r.DB.Where("id = ?", id).First(&movie).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, dto.BaseResponse{
			Status:  false,
			Message: "could not get movie data with id " + id,
		})
	}

	return ctx.JSON(http.StatusOK, dto.BaseResponse{
		Status:  true,
		Message: fmt.Sprintf("movie with id %s fetched successfully", id),
		Data:    movie,
	})
}

func (r *Repository) UpdateMovieById(ctx echo.Context) error {
	id := ctx.Param("id")
	var movie dto.MovieUpdateRequest
	if err := ctx.Bind(&movie); err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, dto.BaseResponse{
			Status:  false,
			Message: err.Error(),
		})
	}

	entityReq := movie.UpdateToEntity()
	if err := r.DB.Model(&entity.Movie{}).Where("id = ?", id).Updates(&entityReq).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, dto.BaseResponse{
			Status:  false,
			Message: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, dto.BaseResponse{
		Status:  true,
		Message: fmt.Sprintf("movie data with id %s updated successfully", id),
	})
}

func (r *Repository) DeleteMovieById(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return ctx.JSON(http.StatusBadRequest, dto.BaseResponse{
			Status:  false,
			Message: "id cannot be empty",
		})
	}

	idInt, _ := strconv.Atoi(id)
	err := r.DB.Delete(&entity.Movie{}, idInt)
	if err.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, dto.BaseResponse{
			Status:  false,
			Message: "could not delete movie",
		})
	}

	return ctx.JSON(http.StatusOK, dto.BaseResponse{
		Status:  true,
		Message: fmt.Sprintf("movie with id %s delete successfully", id),
	})
}

func (r *Repository) SetupRoutes(app *echo.Echo) {
	app.POST("/movies", r.CreateMovie)
	app.GET("/movies", r.GetMovies)
	app.GET("/movie/:id", r.GetMovieById)
	app.PUT("/movie/:id", r.UpdateMovieById)
	app.DELETE("movie/:id", r.DeleteMovieById)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	configDB := &database.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		DBName:   os.Getenv("DB_NAME"),
	}

	db, err := database.InitializeDatabase(configDB)
	if err != nil {
		log.Fatal("could not load the database")
	}

	r := Repository{
		DB: db,
	}

	app := echo.New()
	app.GET("/", root)
	r.SetupRoutes(app)
	app.Start(fmt.Sprintf(":%s", os.Getenv("APP_PORT")))
}
