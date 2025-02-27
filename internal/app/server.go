package app

import (
	"PersonalBlog/internal/config"
	"PersonalBlog/internal/handlers"
	usecases "PersonalBlog/internal/usecases/article"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	httpServer *http.Server
	logger     *zap.Logger
}

// NewServer Создает и возвращает новый HTTP сервер
func NewServer(cfg *config.Config, articleUseCase usecases.ArticleUseCase, logger *zap.Logger) *Server {
	router := gin.New()

	router.Use(ginZapLogger(logger))
	router.Use(ginRecoveryWithLogging(logger))

	// настройка маршрутов
	handlers.SetupRoutes(router, articleUseCase)

	// создание HTTP сервера
	address := fmt.Sprintf("%s: %s", cfg.Listen.BindIP, cfg.Listen.Port)
	return &Server{
		httpServer: &http.Server{
			Addr:    address,
			Handler: router,
		},
		logger: logger,
	}
}

// Run запускает сервер c graceful shutdown
func (s *Server) Run() {
	s.logger.Info("Starting server", zap.String("address", s.httpServer.Addr))

	// канал для graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("Server error", zap.Error(err))
		}
	}()

	// ожидание сигнала для graceful shutdown
	<-done
	s.logger.Info("Shutting down server")

	// завершение работы с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Error("Server shutdown error", zap.Error(err))
	}

	s.logger.Info("Application stopped")
}

// возвращает middleware для Gin, который использует Zap для логирования
func ginZapLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		logger.Info("Request handled",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.String("client_ip", c.ClientIP()),
			zap.Duration("duration", time.Since(start)),
			zap.Int("response_size", c.Writer.Size()),
		)
	}
}

// возвращает middleware для обработки panic с логированием
func ginRecoveryWithLogging(logger *zap.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, err interface{}) {
		logger.Error("Recovered from panic", zap.Any("error", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	})
}
