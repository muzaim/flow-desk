package handler

import (
	"net/http"
	"strconv"

	"flow-desk/dto"
	"flow-desk/service"
	"flow-desk/utils"

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
		appErr := utils.NewBadRequestError("INVALID_INPUT", "Format payload request tidak valid", err)
		utils.RespondWithError(c, appErr)
		return
	}

	res, err := h.taskService.CreateTask(userID, req)
	if err != nil {
		utils.RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Task berhasil dibuat",
		"data":    res,
	})
}

func (h *TaskHandler) GetAll(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	var query dto.TaskQueryParam
	if err := c.ShouldBindQuery(&query); err != nil {
		appErr := utils.NewBadRequestError("INVALID_QUERY_PARAM", "Parameter query tidak valid", err)
		utils.RespondWithError(c, appErr)
		return
	}

	res, err := h.taskService.GetTasks(userID, query)
	if err != nil {
		utils.RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil daftar task",
		"data":    res.Data,
		"meta":    res.Meta,
	})
}

func (h *TaskHandler) GetByID(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		appErr := utils.NewBadRequestError("INVALID_TASK_ID", "ID Task tidak valid", err)
		utils.RespondWithError(c, appErr)
		return
	}

	res, err := h.taskService.GetTaskByID(uint(taskID), userID)
	if err != nil {
		utils.RespondWithError(c, err)
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
		appErr := utils.NewBadRequestError("INVALID_TASK_ID", "ID Task tidak valid", err)
		utils.RespondWithError(c, appErr)
		return
	}

	var req dto.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := utils.NewBadRequestError("INVALID_INPUT", "Format payload request tidak valid", err)
		utils.RespondWithError(c, appErr)
		return
	}

	res, err := h.taskService.UpdateTask(uint(taskID), userID, req)
	if err != nil {
		utils.RespondWithError(c, err)
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
		appErr := utils.NewBadRequestError("INVALID_TASK_ID", "ID Task tidak valid", err)
		utils.RespondWithError(c, appErr)
		return
	}

	err = h.taskService.DeleteTask(uint(taskID), userID)
	if err != nil {
		utils.RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task berhasil dihapus",
	})
}
