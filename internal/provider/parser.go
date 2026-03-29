package provider

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"aiguard/internal/config"
	"aiguard/internal/strutil"
)

type RepoInfo struct {
	Provider     string   `json:"provider"`
	Host         string   `json:"host"`
	Owner        string   `json:"owner"`
	Name         string   `json:"name"`
	Number       string   `json:"number"`
	Path         string   `json:"path"`
	RepoURL      string   `json:"repoUrl"`
	RepoSSHURL   string   `json:"repoSshUrl"`
	RepoHTTPSURL string   `json:"repoHttpsUrl"`
	RepoURLs     []string `json:"repoUrls"`
}

var (
	githubPRPattern = regexp.MustCompile(`^/([^/]+)/([^/]+)/pull/(\d+)/?$`)
	gitlabMRPattern = regexp.MustCompile(`^/(.+)/-/merge_requests/(\d+)/?$`)
)

func Parse(raw string, gitCfg config.GitConfig) RepoInfo {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return RepoInfo{}
	}

	u, err := url.Parse(raw)
	if err != nil {
		return RepoInfo{}
	}

	sourceHost := strings.ToLower(strings.TrimSpace(u.Hostname()))
	sourcePort := strings.TrimSpace(u.Port())
	path := strings.TrimSuffix(u.EscapedPath(), "/")
	if path == "" {
		path = strings.TrimSuffix(u.Path, "/")
	}

	if match := githubPRPattern.FindStringSubmatch(path); len(match) == 4 {
		repoPath := fmt.Sprintf("%s/%s", match[1], strings.TrimSuffix(match[2], ".git"))
		return buildRepoInfo("github", sourceHost, sourcePort, repoPath, match[3], gitCfg.GitHub, gitCfg.PreferredProtocol)
	}

	if match := gitlabMRPattern.FindStringSubmatch(path); len(match) == 3 {
		repoPath := strings.TrimSuffix(match[1], ".git")
		return buildRepoInfo("gitlab", sourceHost, sourcePort, repoPath, match[2], gitCfg.GitLab, gitCfg.PreferredProtocol)
	}

	// generic fallback：即使未匹配 MR/PR 模式，也尝试用 host+path 拼接地址方便调试
	trimmedPath := strings.Trim(path, "/")
	if trimmedPath == "" || sourceHost == "" {
		return RepoInfo{
			Provider: "generic",
			Host:     sourceHost,
		}
	}
	repoPath := strings.TrimSuffix(trimmedPath, ".git")
	sshURL := buildSSHRepoURL(sourceHost, sourcePort, "git", repoPath)
	httpsURL := buildHTTPSRepoURL("https", sourceHost, sourcePort, repoPath)
	repoURLs := orderedRepoURLs(gitCfg.PreferredProtocol, sshURL, httpsURL)

	return RepoInfo{
		Provider:     "generic",
		Host:         sourceHost,
		Path:         repoPath,
		RepoURL:      strutil.FirstNonEmpty(repoURLs...),
		RepoSSHURL:   sshURL,
		RepoHTTPSURL: httpsURL,
		RepoURLs:     repoURLs,
	}
}

func buildRepoInfo(provider, fallbackHost, fallbackPort, repoPath, number string, providerCfg config.GitProviderConfig, preferredProtocol string) RepoInfo {
	repoPath = strings.Trim(strings.TrimSpace(repoPath), "/")
	if repoPath == "" {
		return RepoInfo{Provider: provider, Host: fallbackHost, Number: number}
	}

	parts := strings.Split(repoPath, "/")
	owner := ""
	name := ""
	if len(parts) == 1 {
		name = parts[0]
	} else {
		owner = strings.Join(parts[:len(parts)-1], "/")
		name = parts[len(parts)-1]
	}

	sshHost, sshPort, sshUser := effectiveEndpoint(providerCfg.SSH, fallbackHost, fallbackPort)
	httpsHost, httpsPort, _, httpsScheme := effectiveEndpointWithScheme(providerCfg.HTTPS, fallbackHost, fallbackPort)
	sshURL := buildSSHRepoURL(sshHost, sshPort, sshUser, repoPath)
	httpsURL := buildHTTPSRepoURL(httpsScheme, httpsHost, httpsPort, repoPath)
	repoURLs := orderedRepoURLs(preferredProtocol, sshURL, httpsURL)

	return RepoInfo{
		Provider:     provider,
		Host:         fallbackHost,
		Owner:        owner,
		Name:         name,
		Number:       strings.TrimSpace(number),
		Path:         repoPath,
		RepoURL:      strutil.FirstNonEmpty(repoURLs...),
		RepoSSHURL:   sshURL,
		RepoHTTPSURL: httpsURL,
		RepoURLs:     repoURLs,
	}
}

func effectiveEndpoint(cfg config.GitEndpointConfig, fallbackHost, fallbackPort string) (string, string, string) {
	host := strutil.FirstNonEmpty(strings.TrimSpace(cfg.Host), strings.TrimSpace(fallbackHost))
	port := strutil.FirstNonEmpty(strings.TrimSpace(cfg.Port), strings.TrimSpace(fallbackPort))
	user := strutil.FirstNonEmpty(strings.TrimSpace(cfg.User), "git")
	return host, port, user
}

func effectiveEndpointWithScheme(cfg config.GitEndpointConfig, fallbackHost, fallbackPort string) (string, string, string, string) {
	host, port, user := effectiveEndpoint(cfg, fallbackHost, fallbackPort)
	scheme := strings.TrimSpace(cfg.Scheme)
	if scheme == "" {
		scheme = "https"
	}
	return host, port, user, scheme
}

func buildSSHRepoURL(host, port, user, repoPath string) string {
	host = strings.TrimSpace(host)
	repoPath = strings.Trim(strings.TrimSpace(repoPath), "/")
	if host == "" || repoPath == "" {
		return ""
	}
	if user == "" {
		user = "git"
	}
	if strings.TrimSpace(port) == "" {
		return fmt.Sprintf("%s@%s:%s.git", user, host, repoPath)
	}
	return fmt.Sprintf("ssh://%s@%s:%s/%s.git", user, host, port, repoPath)
}

func buildHTTPSRepoURL(scheme, host, port, repoPath string) string {
	host = strings.TrimSpace(host)
	repoPath = strings.Trim(strings.TrimSpace(repoPath), "/")
	if host == "" || repoPath == "" {
		return ""
	}
	if strings.TrimSpace(scheme) == "" {
		scheme = "https"
	}
	if strings.TrimSpace(port) == "" {
		return fmt.Sprintf("%s://%s/%s.git", scheme, host, repoPath)
	}
	return fmt.Sprintf("%s://%s:%s/%s.git", scheme, host, port, repoPath)
}

func orderedRepoURLs(preferredProtocol, sshURL, httpsURL string) []string {
	urls := []string{}
	appendIfPresent := func(value string) {
		value = strings.TrimSpace(value)
		if value == "" {
			return
		}
		for _, existing := range urls {
			if existing == value {
				return
			}
		}
		urls = append(urls, value)
	}

	if strings.EqualFold(strings.TrimSpace(preferredProtocol), "https") {
		appendIfPresent(httpsURL)
		appendIfPresent(sshURL)
	} else {
		appendIfPresent(sshURL)
		appendIfPresent(httpsURL)
	}
	return urls
}

