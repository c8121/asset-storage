package thumbnails_test

import (
	"os"
	"testing"

	"github.com/c8121/asset-storage/internal/thumbnails"
)

func TestFFMpeg(t *testing.T) {

	in := "/tmp/test.mp4"
	out := "/tmp/test-thumb-2334dg.png"

	stat, err := os.Stat(in)
	if err == nil && !stat.IsDir() {
		err = thumbnails.FFmpegThumb(in, out, 150, 0)
		if err != nil {
			t.Errorf("Thumbnails.FFmpegThumb(%q, %q, 150, 0): %s", in, out, err)
		}
	}

}
