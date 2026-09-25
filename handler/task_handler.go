package handler

import (
	"net/http"
	"strconv"

	"flow-desk/dto"
	"flow-desk/service"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	taskService service.TaskService
}

func NewTaskHandler(taskService service.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

func (h *TaskHandler) Create(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var req dto.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.taskService.CreateTask(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Task berhasil dibuat",
		"data":    res,
	})
}

func (h *TaskHandler) GetAll(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	res, err := h.taskService.GetTasks(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil daftar task",
		"data":    res,
	})
}

func (h *TaskHandler) GetByID(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID Task tidak valid"})
		return
	}

	res, err := h.taskService.GetTaskByID(uint(taskID), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil detail task",
		"data":    res,
	})
}

func (h *TaskHandler) Update(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID Task tidak valid"})
		return
	}

	var req dto.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.taskService.UpdateTask(uint(taskID), userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task berhasil diperbarui",
		"data":    res,
	})
}

func (h *TaskHandler) Delete(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID Task tidak valid"})
		return
	}

	err = h.taskService.DeleteTask(uint(taskID), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task berhasil dihapus",
	})
}
