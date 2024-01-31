package webdriver

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/gan-of-culture/get-sauce/request"
)

type Session struct {
	Value struct {
		SessionID string `json:"sessionId"`
	} `json:"value"`
}

type SessionCookies struct {
	Value []struct {
		Name     string `json:"name"`
		Value    string `json:"value"`
		Path     string `json:"path"`
		Domain   string `json:"domain"`
		Secure   bool   `json:"secure"`
		HTTPOnly bool   `json:"httpOnly"`
		Expiry   int    `json:"expiry,omitempty"`
		SameSite string `json:"sameSite"`
	} `json:"value"`
}

type SessionStringValue struct {
	Value string `json:"value"`
}

type WebDriver struct {
	sessionID string
	cmd       *exec.Cmd
}

func New() (*WebDriver, error) {
	cmd := exec.Command("geckodriver")
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	var res *http.Response
	var err error
	for i := 0; i < 10; i++ {
		res, err = request.Request(http.MethodPost, "http://localhost:4444/session", map[string]string{"Content-Type": "application/json"}, strings.NewReader(`{"capabilities":{"alwaysMatch":{"acceptInsecureCerts":true,"moz:firefoxOptions":{"args":["-headless"]}}}}`))
		if err == nil {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		if err != io.ErrUnexpectedEOF {
			return nil, err
		}
	}

	session := Session{}
	err = json.Unmarshal(body, &session)
	if err != nil {
		return nil, err
	}

	return &WebDriver{
		sessionID: session.Value.SessionID,
		cmd:       cmd,
	}, nil
}

func (webDriver *WebDriver) navigateTo(URL string) error {
	_, err := webDriver.command(http.MethodPost, "url", map[string]string{"Content-Type": "application/json"}, strings.NewReader(fmt.Sprintf(`{"url": "%s"}`, URL)))
	if err != nil {
		return err
	}

	return nil
}

func (webDriver *WebDriver) Close() error {

	_, err := webDriver.command(http.MethodDelete, "", nil, nil)
	if err != nil {
		return err
	}

	webDriver.sessionID = ""

	return webDriver.cmd.Process.Kill()
}

func (webDriver *WebDriver) source() (string, error) {
	body, err := webDriver.command(http.MethodGet, "source", nil, nil)
	if err != nil {
		return "", err
	}

	sessionSource := SessionStringValue{}
	err = json.Unmarshal(body, &sessionSource)
	if err != nil {
		return "", err
	}
	return sessionSource.Value, nil
}

func (webDriver *WebDriver) title() (string, error) {
	body, err := webDriver.command(http.MethodGet, "title", nil, nil)
	if err != nil {
		return "", err
	}

	sessionSource := SessionStringValue{}
	err = json.Unmarshal(body, &sessionSource)
	if err != nil {
		return "", err
	}
	return sessionSource.Value, nil
}

func (webDriver *WebDriver) getCookies() ([]*http.Cookie, error) {
	body, err := webDriver.command(http.MethodGet, "cookie", nil, nil)
	if err != nil {
		return nil, err
	}

	sessionCookies := SessionCookies{}
	err = json.Unmarshal(body, &sessionCookies)
	if err != nil {
		return nil, err
	}

	cookies := []*http.Cookie{}
	for _, c := range sessionCookies.Value {
		cookies = append(cookies, &http.Cookie{
			Name:     c.Name,
			Value:    c.Value,
			Path:     c.Path,
			Domain:   c.Domain,
			Expires:  time.Unix(int64(c.Expiry), 0),
			Secure:   c.Secure,
			HttpOnly: c.HTTPOnly,
			SameSite: http.SameSiteNoneMode,
		})
	}

	return cookies, nil
}

func (webDriver *WebDriver) command(method string, command string, headers map[string]string, body io.Reader) ([]byte, error) {
	if webDriver.sessionID == "" {
		return nil, fmt.Errorf("webdriver session has been closed")
	}

	res, err := request.Request(method, fmt.Sprintf("http://localhost:4444/session/%s/%s", webDriver.sessionID, command), headers, body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return io.ReadAll(res.Body)
}

// SolveChallenge from Cloudflare or DDoS-Guard
func (webDriver *WebDriver) SolveChallenge(URL string) ([]*http.Cookie, error) {

	webDriver.navigateTo(URL)
	select {
	case <-time.After(10 * time.Second):
		break
	default:
		title, err := webDriver.title()
		if err != nil {
			return nil, err
		}
		if title != "DDoS-Guard" && title != "Just a moment..." {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	time.Sleep(5 * time.Second)
	return webDriver.getCookies()
}

// Get HTTP response body as string
func (webDriver *WebDriver) Get(URL string) (string, error) {

	webDriver.navigateTo(URL)
	select {
	case <-time.After(10 * time.Second):
		break
	default:
		title, err := webDriver.title()
		if err != nil {
			return "", err
		}
		if title != "DDoS-Guard" && title != "Just a moment..." {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	return webDriver.source()
}
