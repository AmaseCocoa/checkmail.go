package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"strings"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	key, found := os.LookupEnv("CHECKMAIL_KEY")
	if !found {
		key = "NIL"
	}

	if (os.Getenv("DEBUG_MODE") == "1") {
		app.Get("/version", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{
				"core":    "0.2.0",
				"platform": fmt.Sprintf("Go %s", "1.17"),
				"version": "0.2.0",
			})
		})
	}

	app.Get("/healthcheck", func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})

	app.Get("/mode", func(c *fiber.Ctx) error {
		if (os.Getenv("LEGACY_MODE") == "1") {
			return c.SendString("legacy")
		} else {
			return c.SendString("default")
		}
	})

	if (os.Getenv("LEGACY_MODE") == "1") {
		app.Get("/", func(c *fiber.Ctx) error {
			token := c.Get("Authorization", "NIL")
			email := c.Query("email")
			now := time.Now().Format("2006-01-02 15:04:05 -0700")

			if token != key {
				return c.SendStatus(403)
			}

			if email == "" {
				return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"error": "email is required parameter"})
			}

			for _, whitelistedEmail := range WHITELISTED_EMAILS {
				if email == whitelistedEmail {
					return c.JSON(fiber.Map{
						"date":          now,
						"email":         email,
						"validation_type": "whitelist",
						"success":       true,
						"errors":        nil,
						"smtp_debug":    nil,
						"configuration": fiber.Map{
							"validation_type_by_domain": nil,
							"whitelisted_emails":        WHITELISTED_EMAILS,
							"blacklisted_emails":        nil,
							"whitelisted_domains":       nil,
							"blacklisted_domains":       nil,
							"whitelist_validation":      false,
							"blacklisted_mx_ip_addresses": nil,
							"dns":                       nil,
							"email_pattern":             "default gem value",
							"not_rfc_mx_lookup_flow":   false,
							"smtp_error_body_pattern":   "default gem value",
							"smtp_fail_fast":            false,
							"smtp_safe_check":           false,
						},
					})
				}
			}

			if !compiledRegexEmailPattern.MatchString(email) {
				return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"error": "invalid email address"})
			}

			domain := email[strings.LastIndex(email, "@")+1:]
			exists, err := CheckDomainExists(domain)
			if err != nil {
				return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
			}

			if exists {
				return c.JSON(fiber.Map{
					"date":          now,
					"email":         email,
					"validation_type": "is_disposable",
					"success":       false,
					"errors":        nil,
					"smtp_debug":    nil,
					"configuration": fiber.Map{
						"validation_type_by_domain": nil,
						"whitelisted_emails":        WHITELISTED_EMAILS,
						"blacklisted_emails":        nil,
						"whitelisted_domains":       nil,
						"blacklisted_domains":       nil,
						"whitelist_validation":      false,
						"blacklisted_mx_ip_addresses": nil,
						"dns":                       nil,
						"email_pattern":             "default gem value",
						"not_rfc_mx_lookup_flow":   false,
						"smtp_error_body_pattern":   "default gem value",
						"smtp_fail_fast":            false,
						"smtp_safe_check":           false,
					},
				})
			}

			return c.JSON(fiber.Map{
				"date":          now,
				"email":         email,
				"validation_type": "is_disposable",
				"success":       true,
				"errors":        nil,
				"smtp_debug":    nil,
				"configuration": fiber.Map{
					"validation_type_by_domain": nil,
					"whitelisted_emails":        WHITELISTED_EMAILS,
					"blacklisted_emails":        nil,
					"whitelisted_domains":       nil,
					"blacklisted_domains":       nil,
					"whitelist_validation":      false,
					"blacklisted_mx_ip_addresses": nil,
					"dns":                       nil,
					"email_pattern":             "default gem value",
					"not_rfc_mx_lookup_flow":   false,
					"smtp_error_body_pattern":   "default gem value",
					"smtp_fail_fast":            false,
					"smtp_safe_check":           false,
				},
			})
		})
	} else {
		app.Get("/", func(c *fiber.Ctx) error {
			token := c.Get("Authorization", "NIL")
			email := c.Query("email")
			now := time.Now().Format("2006-01-02 15:04:05 -0700")

			if email == "" {
				return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"error": "email is required parameter"})
			}

			if token != key {
				return c.SendStatus(403)
			}

			for _, whitelistedEmail := range WHITELISTED_EMAILS {
				if email == whitelistedEmail {
					return c.JSON(fiber.Map{
						"date":          now,
						"email":         email,
						"validation_type": "whitelist",
						"pass":       true,
					})
				}
			}

			if !compiledRegexEmailPattern.MatchString(email) {
				return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"error": "invalid email address"})
			}

			domain := email[strings.LastIndex(email, "@")+1:]
			exists, err := CheckDomainExists(domain)
			if err != nil {
				return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
			}

			if exists {
				return c.JSON(fiber.Map{
					"date":          now,
					"email":         email,
					"validation_type": "disposable",
					"pass":       false,
				})
			}

			return c.JSON(fiber.Map{
				"date":          now,
				"email":         email,
				"validation_type": "disposable",
				"pass":       true,
			})
		})
	}


	InitializeDB()
	port, found := os.LookupEnv("CHECKMAIL_PORT")
	if !found {
		port = "3000"
	}


	log.Fatal(app.Listen(":" + port))
}
