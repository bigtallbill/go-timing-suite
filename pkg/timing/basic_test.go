package timing

import (
	"bufio"
	"bytes"
	"log"
	"strings"
	"testing"
	"time"
)

func TestLogDuration(t *testing.T) {
	t.Run("write to logger", func(t *testing.T) {
		var (
			buffer = &bytes.Buffer{}
			writer = bufio.NewWriter(buffer)
			reader = bufio.NewReader(buffer)
			logger = log.New(writer, "fruits: ", log.Lmsgprefix)
		)

		func() {
			defer LogDuration("banana", time.Now(), logger)
			<-time.After(time.Millisecond * 500)
		}()

		_ = writer.Flush()

		logText, err := reader.ReadString('\n')
		if err != nil {
			t.Errorf("didnt expect err: %s", err)
		}

		if !strings.Contains(logText, "banana") {
			t.Errorf("couldnt find log text: %s", logText)
		}
	})
}
