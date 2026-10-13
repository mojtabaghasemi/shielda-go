# Shielda Captcha for Go (کتابخانه رسمی Go برای شیلدا)

[![Go Reference](https://pkg.go.dev/badge/github.com/shielda-ir/shielda-go.svg)](https://pkg.go.dev/github.com/shielda-ir/shielda-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

کتابخانه رسمی زبان **Go (Golang)** برای اعتبارسنجی کپچای امنیتی **شیلدا (Shielda)**. سبک، فوق‌العاده سریع و بدون هیچ‌گونه وابستگی خارجی (تنها با استفاده از کتابخانه استاندارد Go).

---

## 📦 نصب

```bash
go get github.com/shielda-ir/shielda-go
```

---

## 🚀 نحوه استفاده

### ۱. نمونه کد پایه

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shielda-ir/shielda-go"
)

func main() {
	client := shielda.New("YOUR_SITE_KEY", "YOUR_SECRET_KEY")

	ctx := context.Background()
	token := "TOKEN_RECEIVED_FROM_FRONTEND"

	// بررسی سریع صحت کپچا (true/false)
	valid, err := client.Verify(ctx, token)
	if err != nil {
		log.Fatalf("خطا در ارتباط با سرور: %v", err)
	}

	if valid {
		fmt.Println("کپچا با موفقیت تایید شد!")
	} else {
		fmt.Println("کپچا نامعتبر است.")
	}
}
```

### ۲. نمونه در فریم‌ورک Gin

```go
package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/shielda-ir/shielda-go"
)

func main() {
	r := gin.Default()
	client := shielda.New("YOUR_SITE_KEY", "YOUR_SECRET_KEY")

	r.POST("/api/login", func(c *gin.Context) {
		token := c.PostForm("darban-captcha-response")

		valid, _ := client.Verify(c.Request.Context(), token, c.ClientIP())
		if !valid {
			c.JSON(http.StatusBadRequest, gin.H{"error": "کپچا تایید نشد."})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "ورود موفق"})
	})

	r.Run(":8080")
}
```

---

## 📄 لایسنس
[MIT License](LICENSE) - [https://shielda.ir](https://shielda.ir)
