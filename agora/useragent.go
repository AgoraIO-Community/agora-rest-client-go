package agora

import (
	"fmt"
	"strings"

	"github.com/AgoraIO-Community/agora-rest-client-go/version"
)

// BuildUserAgent returns the user agent string used by the application.
// e.g. "AgoraRESTClient  Language/go LanguageVersion/1.22.2 Arch/arm64 OS/darwin SDKVersion/0.0.2"
func BuildUserAgent() string {
	return fmt.Sprintf("AgoraRESTClient  Language/go LanguageVersion/%s Arch/%s OS/%s SDKVersion/%s", strings.TrimPrefix(version.GetGoVersion(), "go"), version.GetArch(), version.GetOS(), version.GetSDKVersion())
}
