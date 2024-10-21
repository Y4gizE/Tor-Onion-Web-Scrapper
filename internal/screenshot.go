package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tebeka/selenium"
)

// TakeScreenshotViaTor, belirtilen URL'in ekran görüntüsünü verilen Tor oturumu ile alır.
func TakeScreenshotViaTor(url string, torSession *TorSession) {
	screenshotDir := "screenshots" // Ekran görüntüleri için klasör

	// Tor ağına bağlantı kontrolü
	var torConnected bool
	for i := 0; i < 30; i++ {
		time.Sleep(2 * time.Second)
		currentURL, err := torSession.driver.CurrentURL()
		if err != nil {
			fmt.Println("Geçerli URL alınırken hata oluştu:", err)
			continue
		}
		if currentURL == "about:tor" || strings.Contains(currentURL, ".onion") {
			torConnected = true
			break
		}

	}
	if !torConnected {
		fmt.Println("Tor ağına bağlanılamadı.")
		return
	}

	// URL'ye git
	fmt.Println("Tor ağına bağlanıldı, siteye gidiliyor:", url)
	if err := torSession.driver.Get(url); err != nil {
		fmt.Println("URL yüklenemedi:", err)
		return
	}

	// Sayfanın tamamen yüklenmesini bekle
	err := torSession.driver.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		currentURL, err := wd.CurrentURL()
		if err != nil {
			return false, nil
		}

		if currentURL == url {
			state, err := wd.ExecuteScript("return document.readyState;", nil)
			if err != nil {
				return false, nil
			}
			return state == "complete", nil
		}
		return false, nil
	}, 120*time.Second)

	if err != nil {
		fmt.Println("Sayfa yüklenirken hata oluştu:", err)
		return
	}

	// Sayfa yüksekliğini al
	pageHeightResult, err := torSession.driver.ExecuteScript("return document.body.scrollHeight;", nil)
	if err != nil {
		fmt.Println("Sayfa yüksekliği alınamadı:", err)
		return
	}
	pageHeight := int64(pageHeightResult.(float64))

	// Kaydırma ve ekran görüntüsü değişkenlerini başlat
	scrollPos := int64(0)
	screenshotCount := 0
	maxScreenshots := 10

	// Ekran görüntüsü klasörünü oluştur, yoksa oluştur
	if err := os.MkdirAll(screenshotDir, os.ModePerm); err != nil {
		fmt.Println("Ekran görüntüsü klasörü oluşturulamadı:", err)
		return
	}

	// Sayfayı kaydırarak ekran görüntüsü alma işlemi
	for scrollPos < pageHeight && screenshotCount < maxScreenshots {
		// Ekran görüntüsü al
		screenshot, err := torSession.driver.Screenshot()
		if err != nil {
			fmt.Println("Ekran görüntüsü alınamadı:", err)
			return
		}

		// Ekran görüntüsünü PNG olarak kaydet
		filePath := filepath.Join(screenshotDir, fmt.Sprintf("screenshot_part_%d.png", screenshotCount))
		file, err := os.Create(filePath)
		if err != nil {
			fmt.Println("Ekran görüntüsü dosyası oluşturulamadı:", err)
			return
		}
		defer file.Close()
		file.Write(screenshot)

		fmt.Printf("Ekran görüntüsü bölümü %d %s olarak kaydedildi\n", screenshotCount, filePath)

		// Sayfayı aşağı kaydır
		_, err = torSession.driver.ExecuteScript("window.scrollBy(0, window.innerHeight);", nil)
		if err != nil {
			fmt.Println("Sayfa kaydırılamadı:", err)
			return
		}

		// Kaydırma pozisyonunu güncelle
		newScrollPosResult, err := torSession.driver.ExecuteScript("return window.pageYOffset;", nil)
		if err != nil {
			fmt.Println("Yeni kaydırma pozisyonu alınamadı:", err)
			return
		}
		scrollPos = int64(newScrollPosResult.(float64))

		// Kaydırmalar arasında bekleme süresi
		time.Sleep(500 * time.Millisecond) // Kaydırmalar arasında 0.5 saniye bekle

		screenshotCount++
	}

	fmt.Println("Ekran görüntüleri başarıyla alındı.")
}
