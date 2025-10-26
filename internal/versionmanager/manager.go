package versionmanager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	defaultAPIBaseURL      = "https://api.github.com"
	defaultDownloadBaseURL = "https://github.com/nevalang/neva/releases/download"
)

type managerConfig struct {
	baseDir         string
	httpClient      *http.Client
	apiBaseURL      string
	downloadBaseURL string
}

type Manager struct {
	baseDir     string
	versionsDir string
	binaryName  string
	versionFile string

	httpClient      *http.Client
	apiBaseURL      string
	downloadBaseURL string
}

func NewManager() (*Manager, error) {
	return newManager(managerConfig{})
}

func newManager(cfg managerConfig) (*Manager, error) {
	base := strings.TrimSpace(cfg.baseDir)
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("determine user home directory: %w", err)
		}
		base = filepath.Join(home, "neva")
	}

	base = filepath.Clean(base)

	binaryName := "neva"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	client := cfg.httpClient
	if client == nil {
		client = http.DefaultClient
	}

	apiBase := cfg.apiBaseURL
	if apiBase == "" {
		apiBase = defaultAPIBaseURL
	}

	downloadBase := cfg.downloadBaseURL
	if downloadBase == "" {
		downloadBase = defaultDownloadBaseURL
	}

	return &Manager{
		baseDir:         base,
		versionsDir:     filepath.Join(base, "versions"),
		binaryName:      binaryName,
		versionFile:     filepath.Join(base, "active-version"),
		httpClient:      client,
		apiBaseURL:      strings.TrimRight(apiBase, "/"),
		downloadBaseURL: strings.TrimRight(downloadBase, "/"),
	}, nil
}

func Normalize(version string) (string, error) {
	trimmed := strings.TrimSpace(version)
	if trimmed == "" {
		return "", errors.New("version must not be empty")
	}

	if strings.HasPrefix(trimmed, "v") || strings.HasPrefix(trimmed, "V") {
		trimmed = trimmed[1:]
	}

	if trimmed == "" {
		return "", errors.New("version must not be empty")
	}

	return "v" + trimmed, nil
}

func (m *Manager) ActiveVersion() (string, error) {
	data, err := os.ReadFile(m.versionFile)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return "", nil
		}
		return "", fmt.Errorf("read active version: %w", err)
	}

	normalized, err := Normalize(string(data))
	if err != nil {
		return "", err
	}

	return normalized, nil
}

func (m *Manager) SetActiveVersion(version string) error {
	normalized, err := Normalize(version)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(m.versionFile), 0o755); err != nil {
		return fmt.Errorf("prepare version file directory: %w", err)
	}

	// Persist the selected version so the bundled CLI can delegate future invocations.
	if err := os.WriteFile(m.versionFile, []byte(normalized+"\n"), 0o644); err != nil {
		return fmt.Errorf("write active version: %w", err)
	}

	return nil
}

func (m *Manager) Use(ctx context.Context, requestedVersion string, currentVersion string) (string, bool, error) {
	normalizedRequested, err := m.normalizeRequestedVersion(ctx, requestedVersion)
	if err != nil {
		return "", false, err
	}

	normalizedCurrent, err := Normalize(currentVersion)
	if err != nil {
		return "", false, err
	}

	wasInstalledJustNow := false
	if normalizedRequested != normalizedCurrent {
		wasInstalledJustNow, err = m.ensureVersionInstalled(ctx, normalizedRequested)
		if err != nil {
			return "", false, err
		}
	}

	if err := m.SetActiveVersion(normalizedRequested); err != nil {
		return "", false, err
	}

	return normalizedRequested, wasInstalledJustNow, nil
}

func (m *Manager) normalizeRequestedVersion(ctx context.Context, version string) (string, error) {
	trimmed := strings.TrimSpace(version)
	if trimmed == "" {
		return "", errors.New("version must not be empty")
	}

	if strings.EqualFold(trimmed, "latest") {
		tag, err := m.fetchLatestReleaseTag(ctx)
		if err != nil {
			return "", err
		}
		return Normalize(tag)
	}

	return Normalize(trimmed)
}

func (m *Manager) fetchLatestReleaseTag(ctx context.Context) (string, error) {
	url := m.apiBaseURL + "/repos/nevalang/neva/releases/latest"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("create latest release request: %w", err)
	}

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("fetch latest release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("fetch latest release: unexpected status %s", resp.Status)
	}

	var payload struct {
		TagName string `json:"tag_name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", fmt.Errorf("decode latest release response: %w", err)
	}

	if strings.TrimSpace(payload.TagName) == "" {
		return "", errors.New("latest release response missing tag_name")
	}

	return payload.TagName, nil
}

func (m *Manager) RunVersion(ctx context.Context, args []string, version string) error {
	normalized, err := Normalize(version)
	if err != nil {
		return err
	}

	binPath := filepath.Join(m.versionsDir, normalized, m.binaryName)
	if _, err := os.Stat(binPath); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("neva version %s is not installed. run 'neva use %s' first", normalized, normalized)
		}
		return fmt.Errorf("check installed version: %w", err)
	}

	cmd := exec.CommandContext(ctx, binPath, args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func (m *Manager) ensureVersionInstalled(ctx context.Context, version string) (bool, error) {
	binPath := filepath.Join(m.versionsDir, version, m.binaryName)
	if _, err := os.Stat(binPath); err == nil {
		return false, nil
	} else if !errors.Is(err, fs.ErrNotExist) {
		return false, fmt.Errorf("check existing version: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(binPath), 0o755); err != nil {
		return false, fmt.Errorf("create version directory: %w", err)
	}

	assetName, err := releaseAssetName()
	if err != nil {
		return false, err
	}

	downloadURL := fmt.Sprintf("%s/%s/%s", m.downloadBaseURL, version, assetName)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
	if err != nil {
		return false, fmt.Errorf("create download request: %w", err)
	}

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("download neva %s: %w", version, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("download neva %s: unexpected status %s", version, resp.Status)
	}

	tmpFile, err := os.CreateTemp(filepath.Dir(binPath), "neva-download-*")
	if err != nil {
		return false, fmt.Errorf("create temporary file: %w", err)
	}
	defer func() {
		_ = tmpFile.Close()
		_ = os.Remove(tmpFile.Name())
	}()

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		return false, fmt.Errorf("write downloaded binary: %w", err)
	}

	if err := tmpFile.Chmod(0o755); err != nil {
		return false, fmt.Errorf("set permissions on downloaded binary: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		return false, fmt.Errorf("close temporary file: %w", err)
	}

	if err := os.Rename(tmpFile.Name(), binPath); err != nil {
		return false, fmt.Errorf("move downloaded binary: %w", err)
	}

	return true, nil
}

func releaseAssetName() (string, error) {
	arch := runtime.GOARCH
	switch arch {
	case "amd64":
	case "arm64":
	case "loong64":
	default:
		return "", fmt.Errorf("unsupported architecture: %s", arch)
	}

	osName := runtime.GOOS
	switch osName {
	case "darwin", "linux":
	case "windows":
	default:
		return "", fmt.Errorf("unsupported operating system: %s", osName)
	}

	asset := fmt.Sprintf("neva-%s-%s", osName, arch)
	if osName == "windows" {
		asset += ".exe"
	}

	return asset, nil
}

func MaybeDelegate(args []string, currentVersion string) (bool, error) {
	if len(args) > 1 {
		if args[1] == "use" {
			return false, nil
		}
	}

	manager, err := NewManager()
	if err != nil {
		return false, err
	}

	active, err := manager.ActiveVersion()
	if err != nil {
		return false, err
	}

	if active == "" {
		return false, nil
	}

	normalizedCurrent, err := Normalize(currentVersion)
	if err != nil {
		return false, err
	}

	if active == normalizedCurrent {
		return false, nil
	}

	return true, manager.RunVersion(context.Background(), args, active)
}
