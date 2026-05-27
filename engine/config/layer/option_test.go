package layer

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-reacterm-core/engine/config/chunk"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
)

func TestDefaultEntry(t *testing.T) {
	mock := drawable_test.MockUnit{}
	cfg := defaultConfig[winsize.Rows](mock.ToUnit())

	assert.False(t, cfg.config.static)
	assert.Equal(t, 0, cfg.Value)
	assert.True(t, cfg.Status)
	assert.Equal(t, chunk.Dynamic[winsize.Rows]().Sized, cfg.config.chunk.Sized)
}

func TestWithChunkOption(t *testing.T) {
	mock := drawable_test.MockUnit{}
	cfg := defaultConfig[winsize.Rows](mock.ToUnit())

	opt := Fixed[winsize.Rows](10)
	opt(&cfg)

	assert.Equal(t, 10, cfg.Chunk().Adapter(15))
}

func TestWithValueOption(t *testing.T) {
	mock := drawable_test.MockUnit{}
	cfg := defaultConfig[winsize.Rows](mock.ToUnit())

	value := winsize.Rows(5)

	opt := WithValue(value)
	opt(&cfg)

	assert.Equal(t, value, cfg.Value)
}

func TestStaticOption(t *testing.T) {
	mock := drawable_test.MockUnit{}
	cfg := defaultConfig[winsize.Rows](mock.ToUnit())

	opt := Static[winsize.Rows]()
	opt(&cfg)

	assert.True(t, cfg.config.static)
}
