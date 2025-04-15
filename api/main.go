package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/google/uuid"
)

type Config struct {
	Port string
}

type LabStatus struct {
	Assignment string    `json:"assignment"`
	Completed  bool      `json:"completed"`
	Status     string    `json:"status"`
	UpdatedAt  time.Time `json:"updated_at"`
	Message    string    `json:"message,omitempty"`
}

type LabRequest struct {
	ID string `json:"id"`
}

type Lab struct {
	ID        string    `json:"id"`
	StartedAt time.Time `json:"started_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Status    string    `json:"status"`
	Completed bool      `json:"completed"`
	Message   string    `json:"message,omitempty"`
}

var labInstances = make(map[string]*Lab)

func main() {
	config := Config{
		Port: getEnv("PORT", "14880"),
	}

	app := setupApp()
	setupRoutes(app)

	log.Printf("Starting server on 0.0.0.0:%s", config.Port)
	if err := app.Listen("0.0.0.0:" + config.Port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func setupApp() *fiber.App {
	return fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})
}

func setupRoutes(app *fiber.App) {
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path}\n",
	}))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"time":   time.Now().UTC(),
		})
	})

	app.Post("/lab/start", startLab)
	app.Get("/lab/:id/status", getLabStatusById)
	app.Post("/lab/:id/status", updateLabStatusById)
}

func startLab(c *fiber.Ctx) error {
	id := uuid.New().String()

	lab := &Lab{
		ID:        id,
		StartedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Status:    "provisioning",
	}
	labInstances[id] = lab

	go provisionInfrastructure(id)

	return c.JSON(lab)
}

func provisionInfrastructure(labId string) {
	cmd := exec.Command("terraform", "-chdir=../terraform", "apply",
		"-var", fmt.Sprintf("lab_id=%s", labId),
		"-var", "debug=true",
		"-auto-approve")
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to provision lab %s: %v", labId, err)
	}
}

func shouldDestroyLab(lab *Lab) bool {
	if lab.Completed {
		return true
	}
	return time.Since(lab.StartedAt) > time.Hour
}

func destroyInfrastructure(labId string) error {
	cmd := exec.Command("terraform", "-chdir=../terraform", "destroy",
		"-var", fmt.Sprintf("lab_id=%s", labId),
		"-auto-approve")

	if err := cmd.Run(); err != nil {
		log.Printf("Failed to destroy lab %s: %v", labId, err)
		return err
	}

	delete(labInstances, labId)
	return nil
}

func getLabStatusById(c *fiber.Ctx) error {
	id := c.Params("id")
	lab, exists := labInstances[id]
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Лабораторная работа не найдена",
		})
	}
	return c.JSON(lab)
}

func updateLabStatusById(c *fiber.Ctx) error {
	id := c.Params("id")
	lab, exists := labInstances[id]
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Лабораторная работа не найдена",
		})
	}

	var request struct {
		Completed bool   `json:"completed"`
		Message   string `json:"message"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный формат запроса",
		})
	}

	lab.Completed = request.Completed
	lab.Status = getStatus(request.Completed)
	lab.UpdatedAt = time.Now().UTC()
	lab.Message = request.Message

	if shouldDestroyLab(lab) {
		go destroyInfrastructure(id)
		lab.Status = "destroying"
	}

	return c.JSON(lab)
}

func getStatus(completed bool) string {
	if completed {
		return "completed"
	}
	return "in_progress"
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
