package controller

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"go-electroshop/internal/payload/request"
	"go-electroshop/internal/payload/response"
	"go-electroshop/internal/service"
	"go-electroshop/internal/utility"
	"io"
	"net/http"
	"regexp"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ChatRequest struct {
	Message string `json:"message"`
}

type ElectroAssistant struct {
	DB           *gorm.DB
	ProductSvc   *service.ProductService
	DashboardSvc *service.DashboardService
}

func NewElectroAssistant(db *gorm.DB) *ElectroAssistant {
	return &ElectroAssistant{
		DB:           db,
		ProductSvc:   service.NewProductService(db),
		DashboardSvc: service.NewDashboardService(db),
	}
}

func cleanResponse(text string) string {
	// Hapus tag XML tetapi pertahankan kontennya
	re := regexp.MustCompile(`<think>(.*?)</think>`)
	text = re.ReplaceAllString(text, "$1")

	re = regexp.MustCompile(`<[^>]*>`)
	text = re.ReplaceAllString(text, "")

	// Pisahkan kata-kata yang menempel
	var result []rune
	var lastRune rune
	for i, r := range text {
		if i > 0 && lastRune != ' ' {
			// Jika bertemu huruf kapital dan sebelumnya huruf kecil, tambahkan spasi
			if unicode.IsUpper(r) && unicode.IsLower(lastRune) {
				result = append(result, ' ')
			}
			// Jika bertemu huruf dan sebelumnya angka atau sebaliknya, tambahkan spasi
			if (unicode.IsLetter(r) && unicode.IsNumber(lastRune)) ||
				(unicode.IsNumber(r) && unicode.IsLetter(lastRune)) {
				result = append(result, ' ')
			}
		}
		result = append(result, r)
		lastRune = r
	}
	text = string(result)

	// Handle markdown formatting
	text = regexp.MustCompile(`\*\*(.*?)\*\*`).ReplaceAllString(text, "**$1**")
	text = regexp.MustCompile(`\*(.*?)\*`).ReplaceAllString(text, "*$1*")
	text = regexp.MustCompile(`\_(.*?)\_`).ReplaceAllString(text, "_$1_")

	// Handle lists
	text = regexp.MustCompile(`(?m)^(\d+)\.\s`).ReplaceAllString(text, "$1. ")
	text = regexp.MustCompile(`(?m)^[-*]\s`).ReplaceAllString(text, "• ")

	// Fix multiple spaces
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")

	// Fix newlines
	text = strings.ReplaceAll(text, `\n`, "\n")
	text = regexp.MustCompile(`\n{3,}`).ReplaceAllString(text, "\n\n")

	return strings.TrimSpace(text)
}

func prosesChunk(chunk string) string {
	var processed []string
	current := ""

	for _, char := range chunk {
		if current == "" {
			current += string(char)
			continue
		}

		lastChar := rune(current[len(current)-1])
		if unicode.IsLower(lastChar) && unicode.IsUpper(char) {
			processed = append(processed, current)
			current = string(char)
		} else {
			current += string(char)
		}
	}

	if current != "" {
		processed = append(processed, current)
	}

	return strings.Join(processed, current)
}

// fetchAppData mengambil data yang diperlukan untuk AI dari berbagai layanan
func (ea *ElectroAssistant) fetchAppData() (map[string]interface{}, error) {
	appData := make(map[string]interface{})

	// Ambil data produk (maksimal 10 produk untuk menjaga performa)
	// Buat filter dengan limit 10 produk
	productFilter := request.ProductFilter{
		Page:  1,
		Limit: 10,
	}

	// Panggil GetProducts dengan filter yang sesuai
	productListResp, err := ea.ProductSvc.GetProducts(productFilter)
	if err != nil {
		return nil, err
	}

	// Simpan produk dari response
	appData["products"] = productListResp.Products
	appData["pagination"] = productListResp.Pagination

	// Ambil kategori produk unik
	categories, err := ea.ProductSvc.GetProductCategories()
	if err != nil {
		return nil, err
	}
	appData["categories"] = categories

	// Ambil statistik produk
	productStats := map[string]interface{}{
		"total_products":   productListResp.Pagination.TotalItems,
		"categories_count": len(categories),
	}
	appData["product_stats"] = productStats

	// Opsi: Tambahkan statistik lain yang relevan (jumlah produk per kategori, rata-rata harga, dll)
	categoryCounts := make(map[string]int)
	for _, p := range productListResp.Products {
		categoryCounts[p.Category]++
	}
	appData["category_counts"] = categoryCounts

	return appData, nil
}

// buildSystemPrompt membuat prompt sistem yang berisi informasi tentang aplikasi
func (ea *ElectroAssistant) buildSystemPrompt(appData map[string]interface{}) string {
	return `Anda adalah Asisten Elektro Shop, asisten virtual untuk aplikasi toko elektronik.

	PERAN ANDA:
	- Anda adalah asisten untuk toko elektronik yang menjual berbagai produk elektronik.
	- Anda membantu pelanggan dan admin dengan informasi tentang produk, kategori, dan harga.
	- Anda dapat memberikan rekomendasi produk berdasarkan kebutuhan pelanggan.
	- Anda dapat menganalisis data penjualan dan memberikan insight untuk meningkatkan penjualan.

	KEMAMPUAN ANDA:
	- Memberikan informasi lengkap tentang produk di toko (nama, kategori, harga, fitur).
	- Membandingkan produk dalam kategori yang sama.
	- Merekomendasikan produk berdasarkan budget, kebutuhan, atau preferensi.
	- Menjelaskan fitur dan spesifikasi produk elektronik.
	- Membantu dengan pertanyaan tentang proses pembelian, ketersediaan, atau kategori produk.
	- Memberikan insight tentang tren penjualan dan performa produk (untuk admin).

	BATASAN PENTING:
	- Anda HANYA menjawab pertanyaan seputar toko elektronik, produk, dan layanan terkait.
	- Anda MENOLAK semua pertanyaan yang tidak terkait dengan toko elektronik, misalnya politik, berita, kesehatan, finansial pribadi, atau topik lain di luar konteks toko elektronik.
	- Jika ditanya tentang topik di luar konteks, katakan dengan sopan bahwa Anda khusus membantu seputar toko elektronik dan produknya.
	- HINDARI memberikan jawaban spekulatif tentang produk atau harga yang tidak ada dalam database.
	- JANGAN membuat klaim tentang produk yang tidak didukung data yang tersedia.

	INFORMASI TOKO:
	- Nama: Elektro Shop
	- Jenis: Toko elektronik/aksesoris
	- Fitur aplikasi: Dashboard admin, manajemen produk, login

	CARA MERESPONS:
	1. Selalu verifikasi apakah pertanyaan terkait dengan toko elektronik.
	2. Jika pertanyaan di luar konteks, tolak dengan sopan dan arahkan kembali ke toko.
	3. Gunakan data produk terbaru untuk memberikan jawaban akurat.
	4. Berikan respons yang sopan, informatif, dan fokus pada kebutuhan pengguna.
	5. Jika memungkinkan, berikan beberapa opsi atau rekomendasi berdasarkan konteks.`
}

// formatAppDataForPrompt memformat data aplikasi untuk dimasukkan ke dalam prompt
func (ea *ElectroAssistant) formatAppDataForPrompt(appData map[string]interface{}) string {
	var sb strings.Builder

	// Format informasi produk
	products := appData["products"].([]response.ProductResponse)
	categories := appData["categories"].([]string)
	stats := appData["product_stats"].(map[string]interface{})
	categoryCounts := appData["category_counts"].(map[string]int)
	pagination := appData["pagination"].(response.Pagination)

	// Format statistik umum
	sb.WriteString(fmt.Sprintf("Jumlah produk: %d\n", stats["total_products"]))
	sb.WriteString(fmt.Sprintf("Jumlah kategori: %d\n", stats["categories_count"]))

	// Format informasi harga
	var minPrice, maxPrice float64
	if len(products) > 0 {
		minPrice = products[0].Price
		maxPrice = minPrice

		for _, p := range products {
			if p.Price < minPrice {
				minPrice = p.Price
			}
			if p.Price > maxPrice {
				maxPrice = p.Price
			}
		}
	}

	sb.WriteString(fmt.Sprintf("Rentang harga: Rp%s - Rp%s\n\n",
		utility.FormatNumber(int64(minPrice)),
		utility.FormatNumber(int64(maxPrice))))

	// Format informasi kategori
	sb.WriteString("Kategori produk dan jumlahnya:\n")
	for _, cat := range categories {
		count := categoryCounts[cat]
		sb.WriteString(fmt.Sprintf("- %s: %d produk\n", cat, count))
	}
	sb.WriteString("\n")

	// Format contoh produk (batasi hanya beberapa produk)
	maxSampleProducts := 5
	if len(products) > maxSampleProducts {
		products = products[:maxSampleProducts]
	}

	sb.WriteString("Contoh produk tersedia:\n")
	for _, p := range products {
		sb.WriteString("- ")
		sb.WriteString(p.Name)
		sb.WriteString(" (Kategori: ")
		sb.WriteString(p.Category)
		sb.WriteString(", Harga: Rp")
		sb.WriteString(utility.FormatNumber(int64(p.Price)))
		sb.WriteString(")\n")
	}

	// Tambahkan petunjuk jika hanya menampilkan sebagian produk
	if len(products) < int(pagination.TotalItems) {
		sb.WriteString(fmt.Sprintf("\n(Data di atas hanya menampilkan %d dari %d total produk)",
			len(products), pagination.TotalItems))
	}

	return sb.String()
}

func (ea *ElectroAssistant) StreamChatHandler(ctx *gin.Context) {
	var req ChatRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utility.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// Set headers for SSE
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("X-Accel-Buffering", "no")

	// Fetch app data
	appData, err := ea.fetchAppData()
	if err != nil {
		ctx.SSEvent("error", gin.H{"error": "Gagal mengambil data aplikasi: " + err.Error()})
		return
	}

	// Bangun prompt lengkap dengan data aplikasi
	systemPrompt := ea.buildSystemPrompt(appData)
	appDataFormatted := ea.formatAppDataForPrompt(appData)

	// Gabungkan prompt sistem dengan pesan pengguna
	fullPrompt := systemPrompt + "\n\nDATA SAAT INI:\n" + appDataFormatted + "\n\nPertanyaan pengguna: " + req.Message

	// Prepare Ollama request body
	ollamaBody := map[string]interface{}{
		"model":  "deepseek-r1:7b",
		"prompt": fullPrompt,
		"stream": true,
		"options": map[string]interface{}{
			"temperature": 0.7,
			"top_p":       0.9,
		},
	}

	ollamaBodyBytes, err := json.Marshal(ollamaBody)
	if err != nil {
		ctx.SSEvent("error", gin.H{"error": err.Error()})
		return
	}

	// Create Ollama request
	ollamaReq, err := http.NewRequest("POST", "http://localhost:11434/api/generate",
		bytes.NewBuffer(ollamaBodyBytes))
	if err != nil {
		ctx.SSEvent("error", gin.H{"error": err.Error()})
		return
	}
	ollamaReq.Header.Set("Content-Type", "application/json")

	// Send request to Ollama
	ollamaResp, err := http.DefaultClient.Do(ollamaReq)
	if err != nil {
		ctx.SSEvent("error", gin.H{"error": err.Error()})
		return
	}
	defer ollamaResp.Body.Close()

	// Create reader for response body
	reader := bufio.NewReader(ollamaResp.Body)
	var responseBuffer strings.Builder

	// Stream response
	for {
		select {
		case <-ctx.Request.Context().Done():
			return
		default:
			line, err := reader.ReadBytes('\n')
			if err == io.EOF {
				if responseBuffer.Len() > 0 {
					processedText := cleanResponse(responseBuffer.String())
					if processedText != "" {
						ctx.SSEvent("message", processedText)
						ctx.Writer.Flush()
					}
				}
				return
			}
			if err != nil {
				ctx.SSEvent("error", gin.H{"error": err.Error()})
				return
			}

			// Parse response
			var response struct {
				Response string `json:"response"`
				Done     bool   `json:"done"`
			}
			if err := json.Unmarshal(line, &response); err != nil {
				continue
			}

			if response.Response != "" {
				processed := prosesChunk(response.Response)
				responseBuffer.WriteString(processed)

				if strings.ContainsAny(response.Response, ".!?\n") {
					processedText := cleanResponse(responseBuffer.String())
					if processedText != "" {
						ctx.SSEvent("message", processedText)
						ctx.Writer.Flush()
						responseBuffer.Reset()
					}
				}
			}

			// Check if generation is complete
			if response.Done {
				if responseBuffer.Len() > 0 {
					processedText := cleanResponse(responseBuffer.String())
					if processedText != "" {
						ctx.SSEvent("message", processedText)
						ctx.Writer.Flush()
					}
				}
				return
			}
		}
	}
}
