package versionmanager

import (
	"context"
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

type Manager struct {
	baseDir     string
	versionsDir string
	binaryName  string
	versionFile string
}

func NewManager() (*Manager, error) {
	base := os.Getenv("NEVA_HOME")
	if strings.TrimSpace(base) == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("determine user home directory: %w", err)
		}
		base = filepath.Join(home, ".neva")
	}

	base = filepath.Clean(base)

	binaryName := "neva"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	return &Manager{
		baseDir:     base,
		versionsDir: filepath.Join(base, "versions"),
		binaryName:  binaryName,
		versionFile: filepath.Join(base, "version"),
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

	if err := os.MkdirAll(m.baseDir, 0o755); err != nil {
		return fmt.Errorf("create version directory: %w", err)
	}

	if err := os.WriteFile(m.versionFile, []byte(normalized+"\n"), 0o644); err != nil {
		return fmt.Errorf("write active version: %w", err)
	}

	return nil
}

func (m *Manager) Use(ctx context.Context, requestedVersion string, currentVersion string) (string, bool, error) {
	normalized, err := Normalize(requestedVersion)
	if err != nil {
		return "", false, err
	}

	currentNormalized, err := Normalize(currentVersion)
	if err != nil {
		return "", false, err
	}

	installed := false
	if normalized != currentNormalized {
		installed, err = m.ensureVersionInstalled(ctx, normalized)
		if err != nil {
			return "", false, err
		}
	}

	if err := m.SetActiveVersion(normalized); err != nil {
		return "", false, err
	}

	return normalized, installed, nil
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

	url := fmt.Sprintf("https://github.com/nevalang/neva/releases/download/%s/%s", version, assetName)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, fmt.Errorf("create download request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
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

	currentNormalized, err := Normalize(currentVersion)
	if err != nil {
		return false, err
	}

	if active == currentNormalized {
		return false, nil
	}

	return true, manager.RunVersion(context.Background(), args, active)
}
