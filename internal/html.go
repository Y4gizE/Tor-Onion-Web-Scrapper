package internal

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/tebeka/selenium"
)

// FetchAndSaveHTML, belirtilen URL'den HTML içeriğini alır ve kaydeder.
func FetchAndSaveHTML(url string, torSession *TorSession) {

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

	// URL'yi aç ve HTML içeriğini al
	fmt.Println("Siteye gidiliyor:", url)
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

	// HTML içeriğini al
	html, err := torSession.driver.PageSource()
	if err != nil {
		// HTML kaynağı alınamadıysa
		fmt.Println("Sayfa kaynağı alınamadı:", err)
		return
	}

	// HTML içeriğini dosyaya kaydet
	err = ioutil.WriteFile("output.html", []byte(html), 0644)
	if err != nil {
		// Dosyaya yazma işlemi başarısız olursa
		fmt.Println("HTML dosyaya yazılamadı:", err)
		return
	}

	// İşlem başarılı olursa
	fmt.Println("HTML içeriği output.html dosyasına başarıyla kaydedildi.")
}
