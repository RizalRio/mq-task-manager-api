package utils

// DefaultResponse adalah kerangka standar untuk semua balasan API
type DefaultResponse struct {
	Status  string      `json:"status"`            // "success" atau "error"
	Message string      `json:"message"`           // Pesan untuk user/frontend
	Data    interface{} `json:"data,omitempty"`    // Payload data (dihilangkan jika kosong)
	Errors  interface{} `json:"errors,omitempty"`  // Detail error (dihilangkan jika kosong)
}

// SuccessResponse digunakan untuk memformat response berhasil (2xx)
func SuccessResponse(message string, data interface{}) DefaultResponse {
	return DefaultResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	}
}

// ErrorResponse digunakan untuk memformat response gagal (4xx, 5xx)
func ErrorResponse(message string, errors interface{}) DefaultResponse {
	return DefaultResponse{
		Status:  "error",
		Message: message,
		Errors:  errors,
	}
}