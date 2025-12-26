package test
import (
	"os"
	"testing"
)
func TestWD(t *testing.T) {
	wd, _ := os.Getwd()
	t.Logf("Working directory: %s", wd)
}
