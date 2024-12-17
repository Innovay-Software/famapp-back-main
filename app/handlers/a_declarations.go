package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/innovay-software/famapp-main/app/dto"
	"github.com/innovay-software/famapp-main/app/models"
)

// An public api handler for gin
type ApiHandler func(*gin.Context) (dto.ApiResponse, error)

// A file handler display an file, it doens't return any ApiResponses
type ApiFileHandler func(*gin.Context) error

// An private api handler where an authenticated user is required
type ApiUserHandler func(*gin.Context, *models.User) (dto.ApiResponse, error)

// A file handler display an file, it doens't return any ApiResponses
type ApiUserFileHandler func(*gin.Context, *models.User) error

// An private api handler where an authenticated admin is required
type ApiAdminHandler func(*gin.Context, *models.User) (dto.ApiResponse, error)
