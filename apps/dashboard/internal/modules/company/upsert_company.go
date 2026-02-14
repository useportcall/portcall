package company

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type UpdateCompanyRequest struct {
	Name             *string `json:"name"`
	Alias            *string `json:"alias"`
	FirstName        *string `json:"first_name"`
	LastName         *string `json:"last_name"`
	Email            *string `json:"email"`
	Phone            *string `json:"phone"`
	VATNumber        *string `json:"vat_number"`
	BusinessCategory *string `json:"business_category"`
	RemoveLogo       *bool   `json:"remove_logo"`
}

func isAllowedImageType(contentType string) bool {
	allowed := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}
	return allowed[contentType]
}

func UpsertCompany(c *routerx.Context) {

	// Parse multipart form (limit to e.g., 2 MB total)
	if err := c.Request.ParseMultipartForm(2 << 20); err != nil {
		log.Println("Error parsing multipart form:", err)
		c.BadRequest("Failed to parse form")
		return
	}

	// Get JSON data part
	jsonPart := c.PostForm("data")
	if jsonPart == "" {
		c.BadRequest("Missing data JSON")
		return
	}

	company := models.Company{}
	if err := c.DB().FindFirstForAppID(c.AppID(), &company); err != nil {
		c.NotFound("Company not found")
		return
	}

	var body UpdateCompanyRequest
	if err := json.Unmarshal([]byte(jsonPart), &body); err != nil {
		c.BadRequest("Invalid JSON in data field")
		return
	}

	// Handle logo removal request
	if body.RemoveLogo != nil && *body.RemoveLogo && company.IconLogoURL != "" {
		log.Printf("Removing logo for company")
		// Note: We could delete the old file from S3 here, but for simplicity we just clear the URL
		company.IconLogoURL = ""
	}

	// Handle optional logo upload
	file, header, err := c.Request.FormFile("logo")
	if err == nil { // File provided
		defer file.Close()

		log.Printf("Received logo upload: %s (%d bytes)", header.Filename, header.Size)

		// Validate: small image, allowed types/sizes
		if header.Size > 2<<20 { // e.g., max 2MB
			c.JSON(http.StatusBadRequest, gin.H{"error": "Image too large"})
			return
		}

		log.Printf("Image content type: %s", header.Header.Get("Content-Type"))

		if !isAllowedImageType(header.Header.Get("Content-Type")) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image type"})
			return
		}

		// Generate unique filename
		filename := uuid.New().String()

		imgBytes, err := io.ReadAll(file)
		if err != nil {
			c.ServerError("Failed to read image file", err)
			return
		}

		// Save the image to the icon logos bucket
		if err := c.Store().PutInIconLogoBucket(filename, imgBytes, c); err != nil {
			c.ServerError("Failed to save icon logo image", err)
			return
		}

		// Generate accessible URL using DigitalOcean Spaces CDN
		// The bucket is public-read, so we can use the direct Spaces URL
		s3Endpoint := os.Getenv("S3_ENDPOINT")
		if s3Endpoint != "" {
			// Remove https:// prefix and construct the bucket URL
			// For DigitalOcean Spaces: https://bucket-name.region.digitaloceanspaces.com/key
			company.IconLogoURL = fmt.Sprintf("%s/%s/%s.png", s3Endpoint, "icon-logos", filename)
		} else {
			// Fallback to FILE_API_URL if S3_ENDPOINT is not set
			company.IconLogoURL = fmt.Sprintf("%s/icon-logos/%s.png", os.Getenv("FILE_API_URL"), filename)
		}

		log.Printf("Generated icon logo URL: %s", company.IconLogoURL)
	}

	if body.Name != nil {
		company.Name = *body.Name
	}

	if body.FirstName != nil {
		company.FirstName = *body.FirstName
	}

	if body.LastName != nil {
		company.LastName = *body.LastName
	}

	if body.Email != nil {
		company.Email = *body.Email
	}

	if body.Phone != nil {
		company.Phone = *body.Phone
	}

	if body.VATNumber != nil {
		company.VATNumber = *body.VATNumber
	}

	if body.Alias != nil {
		company.Alias = *body.Alias
	}

	if body.BusinessCategory != nil {
		company.BusinessCategory = *body.BusinessCategory
	}

	if err := c.DB().Save(&company); err != nil {
		c.ServerError("Failed to update company", err)
		return
	}

	c.OK(new(apix.Company).Set(&company))
}
