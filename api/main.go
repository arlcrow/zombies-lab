package main

import (
    "log"
    "time"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
    // Создаем приложение с настройками
    app := fiber.New(fiber.Config{
        ServerHeader: "FiberLabServer",
        AppName:      "Lab API v1.0.0",
    })

    // Добавляем middleware для логирования запросов
    app.Use(logger.New(logger.Config{
        Format: "[${time}] ${status} - ${method} ${path}\n",
        TimeFormat: "2006-01-02 15:04:05",
    }))

    // GET для корня
    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Hello, World!")
    })

    // Структура для хранения статуса задания
    type LabStatus struct {
        Assignment string    `json:"assignment"`
        Completed  bool      `json:"completed"`
        Status     string    `json:"status"`
        UpdatedAt  time.Time `json:"updated_at"`
    }

    // POST для /lab/status
    app.Post("/lab/status", func(c *fiber.Ctx) error {
        type StatusRequest struct {
            Completed bool `json:"completed"`
        }
        
        req := new(StatusRequest)
        if err := c.BodyParser(req); err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "error": "Invalid JSON format",
                "message": err.Error(),
            })
        }

        // Здесь можно добавить логику обновления статуса
        return c.JSON(fiber.Map{
            "status":    "ok",
            "timestamp": time.Now().UTC(),
        })
    })

    // GET для /lab/status
    app.Get("/lab/status", func(c *fiber.Ctx) error {
        // Пример статуса задания
        status := LabStatus{
            Assignment: "zombies-lab",
            Completed:  false,
            Status:     "in_progress",
            UpdatedAt:  time.Now().UTC(),
        }
        
        return c.JSON(status)
    })

    // GET для /lab/info - дополнительный эндпоинт
    app.Get("/lab/info", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "server_time": time.Now().UTC(),
            "version":     "1.0.0",
            "active_labs": []string{"zombies-lab", "network-lab"},
        })
    })

    // Обработка 404
    app.Use(func(c *fiber.Ctx) error {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error":   "Not Found",
            "message": "This route does not exist",
            "path":    c.Path(),
        })
    })

    // Запуск сервера с graceful shutdown
    log.Printf("Starting server on 0.0.0.0:1488")
    err := app.Listen("0.0.0.0:1488")
    if err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}
