package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"Dedenruslan19/med-project/cmd/echo-server/controller"
	"Dedenruslan19/med-project/cmd/echo-server/middleware"
	"Dedenruslan19/med-project/repository/billing"
	"Dedenruslan19/med-project/repository/diagnose"
	"Dedenruslan19/med-project/repository/exercise"
	"Dedenruslan19/med-project/repository/gemini"
	"Dedenruslan19/med-project/repository/logs"
	"Dedenruslan19/med-project/repository/rapidAPI/bmi"
	"Dedenruslan19/med-project/repository/user"
	"Dedenruslan19/med-project/repository/workout"
	billingService "Dedenruslan19/med-project/service/billings"
	diagnosisService "Dedenruslan19/med-project/service/diagnoses"
	exerciseService "Dedenruslan19/med-project/service/exercises"
	logService "Dedenruslan19/med-project/service/logs"
	userService "Dedenruslan19/med-project/service/users"
	workoutService "Dedenruslan19/med-project/service/workouts"

	"github.com/labstack/echo/v4"
	mdw "github.com/labstack/echo/v4/middleware"
	cfg "github.com/pobyzaarif/go-config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	loggerOption = slog.HandlerOptions{AddSource: true}
	logger       = slog.New(slog.NewJSONHandler(os.Stdout, &loggerOption))
)

type Config struct {
	AppPort     string `env:"APP_PORT"`
	DBDriver    string `env:"DB_DRIVER"`
	DBHost      string `env:"DB_HOST"`
	DBPort      string `env:"DB_PORT"`
	DBUser      string `env:"DB_USER"`
	DBPass      string `env:"DB_PASS"`
	DBName      string `env:"DB_NAME"`
	DatabaseURL string `env:"DATABASE_URL"`

	RapidAPIBMI string `env:"BMI_API_KEY"`
	GEMINI      string `env:"GEMINI_API_KEY"`
}

func main() {
	config := Config{}
	if err := cfg.LoadConfig(&config); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	logger.Info("Config loaded")

	var dsn string
	if config.DatabaseURL != "" {
		dsn = config.DatabaseURL
		logger.Info("Using DATABASE_URL from Railway environment")
	} else {
		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			config.DBHost, config.DBPort, config.DBUser, config.DBPass, config.DBName,
		)
		logger.Info("Using local database connection")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	logger.Info("Database client connected!")

	// RapidAPI BMI
	bmiRepo := bmi.NewRapidAPIRepository(logger, config.RapidAPIBMI)
	geminiRepo := gemini.NewGeminiRepository(logger, config.GEMINI)

	// Repository & Service
	userRepo := user.NewUserRepo(db, logger)
	userSvc := userService.NewService(logger, userRepo, bmiRepo)
	userController := controller.NewUserController(userSvc, logger)

	workoutRepo := workout.NewWorkoutRepo(db, logger)
	workoutSvc := workoutService.NewService(logger, workoutRepo, geminiRepo)
	workoutController := controller.NewWorkoutController(workoutSvc, logger)

	exerciseRepo := exercise.NewExerciseRepo(db, logger)
	exerciseSvc := exerciseService.NewService(logger, exerciseRepo, workoutSvc, geminiRepo)
	exerciseController := controller.NewExerciseController(exerciseSvc, logger)

	logRepo := logs.NewLogsRepo(db, logger)
	logSvc := logService.NewService(logger, logRepo)
	logController := controller.NewLogController(logSvc, logger)

	billingRepo := billing.NewBillingRepo(db, logger)
	billingSvc := billingService.NewService(logger, billingRepo)
	billingController := controller.NewBillingController(billingSvc, logger)

	diagnosisRepo := diagnose.NewDiagnosisRepo(db, logger)
	diagnosisSvc := diagnosisService.NewService(logger, diagnosisRepo)
	diagnosisController := controller.NewDiagnosisController(diagnosisSvc, billingSvc, logger)

	// Setup Echo
	e := echo.New()
	e.HideBanner = true

	// Middleware
	e.Use(mdw.CORS())
	e.Use(mdw.LoggerWithConfig(mdw.LoggerConfig{
		Skipper:          mdw.DefaultSkipper,
		CustomTimeFormat: "2006-01-02 15:04:05.00000",
		Format: `{"time":"${time_rfc3339_nano}","level":"INFO","id":"${id}","remote_ip":"${remote_ip}",` +
			`"host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` +
			`"status":${status},"error":"${error}","latency":${latency},"latency_human":"${latency_human}"` +
			`,"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n",
	}))
	e.Pre(mdw.RemoveTrailingSlash())
	e.Pre(mdw.Recover())

	// Routes
	// users
	userGroup := e.Group("/users")
	userGroup.POST("/register", userController.Register)
	userGroup.POST("/login", userController.Login)

	userMiddleware := userGroup.Group("", middleware.JWTMiddleware)
	userMiddleware.GET("", userController.GetMe)

	// workouts
	workoutGroup := e.Group("/workouts", middleware.JWTMiddleware)
	workoutGroup.POST("", workoutController.CreateWorkout, middleware.ValidateContentType)
	workoutGroup.GET("", workoutController.GetAllWorkouts)
	workoutGroup.GET("/:id", workoutController.GetWorkoutByID)
	workoutGroup.PUT("/:id", workoutController.UpdateWorkout, middleware.ValidateContentType)
	workoutGroup.DELETE("/:id", workoutController.DeleteWorkout)

	// exercises
	exerciseGroup := e.Group("/exercises", middleware.JWTMiddleware)
	exerciseGroup.POST("", exerciseController.CreateExercise, middleware.ValidateContentType)
	exerciseGroup.DELETE("/:id", exerciseController.DeleteExercise)

	// logs
	logGroup := e.Group("/logs", middleware.JWTMiddleware)
	logGroup.POST("", logController.CreateLog)
	logGroup.GET("", logController.GetAllLogs)

	// billings
	billingGroup := e.Group("/billings", middleware.JWTMiddleware)
	billingGroup.POST("", billingController.CreateBilling, middleware.ValidateContentType)
	billingGroup.GET("/:id", billingController.GetBillingByID)
	billingGroup.GET("/appointment/:appointment_id", billingController.GetBillingByAppointmentID)
	billingGroup.PUT("/:id/payment-status", billingController.UpdatePaymentStatus, middleware.ValidateContentType)

	// diagnoses
	diagnosisGroup := e.Group("/diagnoses", middleware.JWTMiddleware)
	diagnosisGroup.POST("", diagnosisController.CreateDiagnosis, middleware.ValidateContentType)
	diagnosisGroup.GET("/:id", diagnosisController.GetDiagnosisByID)
	diagnosisGroup.GET("/appointment/:appointment_id", diagnosisController.GetDiagnosisByAppointmentID)
	diagnosisGroup.PUT("/:id", diagnosisController.UpdateDiagnosis, middleware.ValidateContentType)

	// Detect port from Railway
	port := os.Getenv("PORT")
	if port == "" {
		port = config.AppPort
		if port == "" {
			port = "8080"
		}
	}

	// Start server
	go func() {
		if err := e.Start(":" + port); err != http.ErrServerClosed {
			log.Fatalf("Failed to start HTTP server on port %s: %v", port, err)
		}
	}()

	logger.Info("API service running on port " + port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// a timeout of 10 seconds to shutdown the server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatal("Failed to shutting down echo server", "err", err)
	} else {
		logger.Info("Successfully shutting down echo server")
	}
}
