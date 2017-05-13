package AuthKey_test

import (
	"fmt"
	"testing"
	"time"
)

const rfc2822 = "Mon Jan 02 15:04:05 -0700 2006"

func TestTimeError(*testing.T) {
	fmt.Println("rfc2822 NOW: ", time.Now().Format(rfc2822))

}
