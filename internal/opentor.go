package internal

import (
	"fmt"

	"github.com/tebeka/selenium"
)

// TorSession yapısını tanımlayın
type TorSession struct {
	driver selenium.WebDriver
}

func StartTorSession() (*TorSession, error) {
	const (
		geckoDriverPath = `C:\Users\HP\Downloads\geckodriver-v0.35.0-win64\geckodriver.exe`
		port            = 4444
		firefoxBinary   = `B:\Tor Browser\Browser\firefox.exe`
	)

	caps := selenium.Capabilities{
		"browserName": "firefox",
		"moz:firefoxOptions": map[string]interface{}{
			"args": []string{
				"--proxy-server=socks5://127.0.0.1:9150",
				"-no-remote",
				"-private",
			},
			"binary": firefoxBinary,
		},
	}

	// Geckodriver servisini başlat
	service, err := selenium.NewGeckoDriverService(geckoDriverPath, port)
	if err != nil {
		return nil, fmt.Errorf("Geckodriver başlatma hatası: %v", err)
	}

	// WebDriver'ı başlat
	driver, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d", port))
	if err != nil {
		service.Stop()
		return nil, fmt.Errorf("WebDriver'a bağlanılamadı: %v", err)
	}

	// TorSession oluştur
	torSession := &TorSession{driver: driver}
	return torSession, nil
}

// TorSession'ı kapat
func (t *TorSession) Stop() {
	t.driver.Quit()
}
