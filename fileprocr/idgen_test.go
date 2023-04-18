package fileprocr_test

import (
	"testing"

	"github.com/bilginyuksel/file-processor/fileprocr"
	"github.com/stretchr/testify/assert"
)

func TestIDGeneratorGenerate(t *testing.T) {
	idgen := fileprocr.NewIDGenerator()
	assert.NotEmpty(t, idgen.Generate())
}
