package pmc

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/oneagent/pmc/ruxit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	t.Run("destination file has fixed 0600 permissions", func(t *testing.T) {
		srcDir := t.TempDir()
		dstDir := t.TempDir()

		srcConf := ruxit.ProcConf{
			Properties: []ruxit.Property{
				{Section: "test", Key: "key", Value: "value"},
			},
		}

		srcPath := filepath.Join(srcDir, "ruxitagentproc.conf")
		require.NoError(t, os.WriteFile(srcPath, []byte(srcConf.ToString()), 0644))

		dstPath := filepath.Join(dstDir, "ruxitagentproc.conf")

		err := Create(testLog, srcPath, dstPath, ruxit.ProcConf{})
		require.NoError(t, err)

		info, err := os.Stat(dstPath)
		require.NoError(t, err)
		assert.Equal(t, os.FileMode(filePerm), info.Mode().Perm())
	})

	t.Run("merges source and override configs", func(t *testing.T) {
		srcDir := t.TempDir()
		dstDir := t.TempDir()

		srcConf := ruxit.ProcConf{
			Properties: []ruxit.Property{
				{Section: "test", Key: "key", Value: "original"},
				{Section: "test", Key: "only-in-src", Value: "src"},
			},
		}

		override := ruxit.ProcConf{
			Properties: []ruxit.Property{
				{Section: "test", Key: "key", Value: "overridden"},
				{Section: "test", Key: "only-in-override", Value: "override"},
			},
		}

		srcPath := filepath.Join(srcDir, "ruxitagentproc.conf")
		require.NoError(t, os.WriteFile(srcPath, []byte(srcConf.ToString()), 0600))

		dstPath := filepath.Join(dstDir, "ruxitagentproc.conf")

		err := Create(testLog, srcPath, dstPath, override)
		require.NoError(t, err)

		content, err := os.ReadFile(dstPath)
		require.NoError(t, err)
		assert.Equal(t, srcConf.Merge(override).ToString(), string(content))
	})

	t.Run("missing source file returns error", func(t *testing.T) {
		dstDir := t.TempDir()

		err := Create(testLog, "/nonexistent/path/ruxitagentproc.conf", filepath.Join(dstDir, "out.conf"), ruxit.ProcConf{})
		require.Error(t, err)
	})
}
