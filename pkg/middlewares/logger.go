package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var CustomTags = map[string]logger.LogFunc{
	"real_ip": func(output logger.Buffer, c *fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
		ip := c.Get("x-real-ip")
		if ip == "" {
			ip = c.IP()
		}
		return output.WriteString(ip)
	},
}

var Logger = logger.New(logger.Config{
	Format:     "${time} [${real_ip}] ${status} - ${method} ${path}\n",
	TimeFormat: "2006/01/02 15:04:05",
	CustomTags: CustomTags,
})
