package app

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2/middleware/logger"
)

var logMiddleware = logger.New(logger.Config{
	Format:     "${time} | ${locals:requestid} | ${status} | ${reqHeader:X-Real-IP} | ${latency} | ${method} | ${path}\n",
	Output:     os.Stderr,
	TimeFormat: time.RFC3339,
})
