package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"fitgenie/pkg/logger"
	"fitgenie/pkg/storage"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UploadHandler struct {
	s3Client *storage.S3Client
	log      *logger.Logger
}

func NewUploadHandler(s3Client *storage.S3Client, log *logger.Logger) *UploadHandler {
	return &UploadHandler{
		s3Client: s3Client,
		log:      log,
	}
}

// UploadImage sube una imagen a S3 y devuelve la URL
func (h *UploadHandler) UploadImage(c *gin.Context) {
	// Obtener userID del contexto (seteado por middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Obtener archivo del form
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No image file provided"})
		return
	}
	defer file.Close()

	// Validar tamaño (max 5MB)
	if header.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File too large (max 5MB)"})
		return
	}

	// Validar tipo
	ext := strings.ToLower(filepath.Ext(header.Filename))
	validExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
	if !validExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type (jpg, png, webp only)"})
		return
	}

	// Generar nombre único
	fileName := fmt.Sprintf("users/%s/%s%s", userID.(uuid.UUID).String(), uuid.New().String(), ext)
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg"
	}

	// Leer archivo en memoria
	data := make([]byte, header.Size)
	if _, err := file.Read(data); err != nil {
		h.log.Error("failed to read file", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process image"})
		return
	}

	// Subir a S3
	ctx := c.Request.Context()
	if err := h.s3Client.Upload(ctx, fileName, data, contentType); err != nil {
		h.log.Error("failed to upload to S3", "error", err, "file", fileName)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image"})
		return
	}

	// Generar URL presignada (válida por 1 hora)
	url, err := h.s3Client.GetPresignedURL(ctx, fileName, time.Hour)
	if err != nil {
		h.log.Error("failed to generate presigned URL", "error", err)
		// Continuamos igual, devolvemos el path
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"file_path": fileName,
		"url":       url,
		"size":      header.Size,
	})
}

// GetImageURL genera URL temporal para una imagen
func (h *UploadHandler) GetImageURL(c *gin.Context) {
	filePath := c.Param("path")
	if filePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file path provided"})
		return
	}

	ctx := c.Request.Context()
	url, err := h.s3Client.GetPresignedURL(ctx, filePath, time.Hour)
	if err != nil {
		h.log.Error("failed to generate presigned URL", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url":       url,
		"expires_in": 3600, // segundos
	})
}
