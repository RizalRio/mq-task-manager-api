package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// FormatValidationError menerjemahkan error bawaan Gin (go-playground/validator)
// menjadi array of string yang bersih dan terstruktur.
func FormatValidationError(err error) []string {
	var errors []string

	// Melakukan type assertion untuk memastikan error berasal dari kegagalan validasi form/JSON
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			// Membuat pesan error kustom berdasarkan tag validasi yang dilanggar
			switch e.Tag() {
			case "required":
				errors = append(errors, fmt.Sprintf("Kolom '%s' wajib diisi", e.Field()))
			case "email":
				errors = append(errors, fmt.Sprintf("Kolom '%s' harus berupa alamat email yang valid", e.Field()))
			case "min":
				errors = append(errors, fmt.Sprintf("Kolom '%s' minimal terdiri dari %s karakter", e.Field(), e.Param()))
			case "max":
				errors = append(errors, fmt.Sprintf("Kolom '%s' maksimal terdiri dari %s karakter", e.Field(), e.Param()))
			case "oneof":
				errors = append(errors, fmt.Sprintf("Kolom '%s' harus diisi dengan salah satu nilai berikut: %s", e.Field(), e.Param()))
			default:
				// Fallback untuk tag validasi lain yang belum didefinisikan secara spesifik
				errors = append(errors, fmt.Sprintf("Kolom '%s' tidak memenuhi kriteria validasi", e.Field()))
			}
		}
		return errors
	}

	// Jika error bukan dari validasi field (misalnya format JSON cacat/malformed)
	return append(errors, "Format payload tidak valid: "+err.Error())
}