package handler

import (
	"net/http"

	"go-minio/internal/service"
	"go-minio/pkg/response"
)

type StorageHandler struct {
	storageService service.StorageService
}

func NewStorageHandler(service service.StorageService) *StorageHandler {
	return &StorageHandler{storageService: service}
}

func (h *StorageHandler) Upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "Method tidak diizinkan")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "file tidak terdeteksi")
		return
	}
	defer file.Close()

	client := r.FormValue("client")
	folder := r.FormValue("folder")

	fileURL := h.storageService.UploadAsync(
		file,
		header.Filename,
		header.Size,
		header.Header.Get("Content-Type"),
		client,
		folder,
	)

	response.JSON(w, http.StatusAccepted, map[string]interface{}{
		"success": true,
		"message": "Sukses Upload (Queued)",
		"url":     fileURL,
	})
}

func (h *StorageHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, http.StatusMethodNotAllowed, "Method tidak diizinkan")
		return
	}

	client := r.URL.Query().Get("client")
	folder := r.URL.Query().Get("folder")

	files, err := h.storageService.ListFiles(r.Context(), client, folder)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"files":   files,
	})
}

func (h *StorageHandler) View(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "Method tidak diizinkan")
		return
	}

	fileURL := r.FormValue("url")
	if fileURL == "" {
		response.Error(w, http.StatusBadRequest, "URL file kosong")
		return
	}

	data, err := h.storageService.GetFileStat(r.Context(), fileURL)
	if err != nil {
		response.Error(w, http.StatusNotFound, "File tidak ditemukan: "+err.Error())
		return
	}

	response.JSON(w, http.StatusOK, data)
}

func (h *StorageHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "Method tidak diizinkan")
		return
	}

	fileURL := r.FormValue("file")
	if fileURL == "" {
		response.Error(w, http.StatusBadRequest, "File URL kosong")
		return
	}

	if err := h.storageService.DeleteFile(r.Context(), fileURL); err != nil {
		response.Error(w, http.StatusInternalServerError, "Gagal menghapus file: "+err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "File deleted successfully",
	})
}
