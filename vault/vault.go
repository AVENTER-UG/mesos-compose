package vault

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/sirupsen/logrus"
)

// Vault struct about the vault connection
type Vault struct {
	Token     string
	URL       string
	Timeout   time.Duration
	Client    *api.Client
	Connected bool
}

// New will create a new Vault object
func New(token, url string, timeout time.Duration) *Vault {
	e := &Vault{
		Token:   token,
		URL:     url,
		Timeout: timeout,
	}

	return e
}

// Connect - with the vault
func (e *Vault) Connect() bool {
	var httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}

	var err error
	e.Client, err = api.NewClient(&api.Config{Address: e.URL, HttpClient: httpClient})
	if err != nil {
		logrus.WithField("func", "vault: Connect").Error("Could not reach Vault: ", err.Error())
		e.Connected = false
		return false
	}
	e.Client.SetToken(e.Token)

	e.Connected = true
	return true
}

// ReadString - get the data of the given secret
func (e *Vault) ReadString(secret string) interface{} {
	secret = strings.ReplaceAll(secret, "vault://", "")
	data, err := e.Client.Logical().Read(secret)
	if err != nil {
		logrus.WithField("func", "vault: ReadString").Error("Could not get secret: ", err.Error())
		return nil
	}

	if data != nil {
		return data.Data["data"]
	}
	return nil
}

// GetKey - extract the specific key of a secret
func (e *Vault) GetKey(secret string) string {
	secret = strings.ReplaceAll(secret, "vault://", "")
	if strings.Contains(secret, ":") {
		param := strings.Split(secret, ":")
		value := e.ReadString(param[0])

		if value != nil {
			return fmt.Sprintf("%v", value.(map[string]interface{})[param[1]])
		}
		return ""
	}
	return ""
}

// IsConnected - give out the connection status
func (e *Vault) IsConnected() bool {
	return e.Connected
}
