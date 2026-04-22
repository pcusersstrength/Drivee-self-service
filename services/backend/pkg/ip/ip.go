package ip
import (
	"net"
	"net/http"
	"strings"
)

func GetRealIP(r *http.Request) string {
	// 1. Самые надёжные заголовки (в порядке приоритета)
	headers := []string{
		"CF-Connecting-IP", // Cloudflare (самый надёжный)
		"X-Real-IP",        // часто ставит nginx / caddy
		"X-Forwarded-For",  // самый распространённый
		"X-Original-Forwarded-For",
		"True-Client-IP", // иногда у Akamai, Cloudfront и т.д.
	}

	for _, header := range headers {
		if ip := r.Header.Get(header); ip != "" {
			// X-Forwarded-For может содержать несколько IP через запятую
			// Берём самый левый (первый) — это обычно и есть клиент
			if header == "X-Forwarded-For" || header == "X-Original-Forwarded-For" {
				parts := strings.Split(ip, ",")
				if len(parts) > 0 {
					candidate := strings.TrimSpace(parts[0])
					if isValidIP(candidate) {
						return candidate
					}
				}
			} else if isValidIP(ip) {
				return strings.TrimSpace(ip)
			}
		}
	}

	// 2. Если заголовков нет — берём то, что видит Go напрямую
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil && isValidIP(ip) {
		return ip
	}

	return r.RemoteAddr // fallback
}

// Простая проверка, что это действительно IP (а не мусор)
func isValidIP(ip string) bool {
	return net.ParseIP(strings.TrimSpace(ip)) != nil
}
