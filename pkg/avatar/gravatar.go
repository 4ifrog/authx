package avatar

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/cybersamx/authx/pkg/utils"
)

const (
	fetchTimeout = 5 * time.Second
)

func FetchGravatarURL(parent context.Context, email string) (string, error) {
	hashEmail := utils.MD5(email)
	gravatarURL, err := url.Parse(fmt.Sprintf("http://gravatar.com/avatar/%s", hashEmail))
	if err != nil {
		return "", fmt.Errorf("malformed")
	}

	// Append d=404 to tell Gravatar to return status code 404 and not a default image if Gravatar not found.
	headURL := gravatarURL
	query := headURL.Query()
	query.Set("d", "404")
	ctx, cancel := context.WithTimeout(parent, fetchTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, headURL.String(), nil)
	if err != nil {
		return "", fmt.Errorf("can't instantiate request to connect to gravatar: %v", err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("can't connect to gravatar: %v", err)
	}
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("problem connecting to gravatar, status code: %d", res.StatusCode)
	}

	return gravatarURL.String(), nil
}
