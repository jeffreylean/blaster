package blast

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlast(t *testing.T) {
	err := os.Setenv("PAYLOAD", "file:///payload.txt")
	assert.NoError(t, err)
	Blast("http://127.0.0.1:8081/com.snowplowanalytics.snowplow/tp2", "", 50, 10, 100)
}
