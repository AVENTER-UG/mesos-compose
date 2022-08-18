package vault

import (
	"net/http"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/sirupsen/logrus"
)

// Vault struct about the vault connection
type Vault struct {
	Token   string
	URL     string
	Timeout time.Duration
	Client  *api.Client
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
		return false
	}
	e.Client.SetToken(e.Token)

	return true
}

// ReadString - get the data of the given secret
func (e *Vault) ReadString(secret string) map[string]interface{} {
	data, err := e.Client.Logical().Read(secret)
	if err != nil {
		logrus.WithField("func", "vault: ReadString").Error("Could not get secret: ", err.Error())
		return nil
	}

	return data.Data
}
