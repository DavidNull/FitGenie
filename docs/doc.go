// Package docs FitGenie API Documentation
//
// This is the API documentation for FitGenie - an AI-powered outfit recommendation
// application with Flutter mobile frontend and Go backend.
//
// Terms of Service: http://swagger.io/terms/
//
// Contact: support@fitgenie.local
//
// License: MIT https://opensource.org/licenses/MIT
//
// BasePath: /api/v1
//
// swagger:meta
package docs

import (
	"github.com/swaggo/swag"
)

// SwaggerInfo holds the swagger document info
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{"http"},
	Title:            "FitGenie API",
	Description:      "AI-powered outfit recommendation API with Flutter mobile app",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "API Support",
            "url": "http://www.fitgenie.local/support",
            "email": "support@fitgenie.local"
        },
        "license": {
            "name": "MIT",
            "url": "https://opensource.org/licenses/MIT"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {},
    "definitions": {
        "models.User": {
            "type": "object",
            "properties": {
                "id": {"type": "string"},
                "email": {"type": "string"},
                "name": {"type": "string"},
                "created_at": {"type": "string", "format": "date-time"}
            }
        },
        "models.ClothingItem": {
            "type": "object",
            "properties": {
                "id": {"type": "string"},
                "name": {"type": "string"},
                "category": {"type": "string"},
                "primary_color": {"type": "string"},
                "style": {"type": "string"},
                "image_url": {"type": "string"}
            }
        },
        "models.Outfit": {
            "type": "object",
            "properties": {
                "id": {"type": "string"},
                "name": {"type": "string"},
                "description": {"type": "string"},
                "clothing_items": {
                    "type": "array",
                    "items": {"type": "object"}
                }
            }
        }
    }
}`

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
