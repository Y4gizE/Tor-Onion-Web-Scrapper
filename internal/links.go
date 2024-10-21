package internal

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/tebeka/selenium"
)

// FetchAndSaveLinks, belirtilen URL'den bağlantıları alır ve links.txt dosyasına kaydeder.
func FetchAndSaveLinks(url string, torSession *TorSession) {

	// Tor ağına bağlantı kontrolü
	var torConnected bool
	for i := 0; i < 30; i++ {
		time.Sleep(2 * time.Second)
		currentURL, err := torSession.driver.CurrentURL()
		if err != nil {
			fmt.Println("Geçerli URL alınırken hata oluştu:", err)
			continue
		}
		// Eğer Tor ağına bağlanmışsak veya .onion uzantılı bir siteye yönlenmişsek
		if currentURL == "about:tor" || strings.Contains(currentURL, ".onion") {
			torConnected = true
			break
		}

	}
	if !torConnected {
		// Eğer Tor ağına bağlanılamadıysa
		fmt.Println("Tor ağına bağlanılamadı.")
		return
	}

	// URL'yi aç ve bağlantıları al
	fmt.Println("Bağlantılar alınıyor:", url)
	if err := torSession.driver.Get(url); err != nil {
		// URL yüklenirken hata oluşursa
		fmt.Println("URL yüklenemedi:", err)
		return
	}

	// Sayfanın tamamen yüklenmesini bekle
	err := torSession.driver.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		currentURL, err := wd.CurrentURL()
		if err != nil {
			return false, nil
		}

		// URL yüklendiyse ve sayfa durumu 'complete' ise
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
		// Sayfa yüklenirken hata oluşursa
		fmt.Println("Sayfanın yüklenmesi sırasında hata oluştu:", err)
		return
	}

	// Sayfadaki bağlantıları al
	links, err := torSession.driver.FindElements(selenium.ByTagName, "a")
	if err != nil {
		// Bağlantılar bulunamazsa
		fmt.Println("Bağlantılar bulunamadı:", err)
		return
	}

	// links.txt dosyasını aç veya oluştur
	file, err := os.Create("links.txt")
	if err != nil {
		// Dosya oluşturulamazsa
		fmt.Println("links.txt dosyası oluşturulamadı:", err)
		return
	}
	defer file.Close()

	// Bağlantıları dosyaya yaz
	for _, link := range links {
		href, err := link.GetAttribute("href")
		if err != nil {
			// href özelliği alınamazsa
			fmt.Println("href özelliği alınamadı:", err)
			continue
		}
		if strings.HasPrefix(href, "http") {
			// Eğer bağlantı 'http' ile başlıyorsa
			fmt.Println("Bulunan bağlantı:", href)
			// Bağlantıyı dosyaya yaz
			_, err := file.WriteString(href + "\n")
			if err != nil {
				// Dosyaya yazılamazsa
				fmt.Println("Bağlantı dosyaya yazılamadı:", err)
				return
			}
		}
	}

	// İşlem başarılı olursa
	fmt.Println("Bağlantılar links.txt dosyasına kaydedildi.")
}
