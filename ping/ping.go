package ping

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/laof/lite-speed-test/web"
)

func Test(url string) ([]int, error) {
	link := flag.String("link", url, "link to test")
	mode := flag.String("mode", "pingonly", "speed test mode")
	flag.Parse()
	// link := "vmess://aHR0cHM6Ly9naXRodWIuY29tL3h4ZjA5OC9MaXRlU3BlZWRUZXN0"
	if len(*link) < 1 {
		log.Fatal("link required")
	}
	opts := web.ProfileTestOptions{
		GroupName:     "Default",
		SpeedTestMode: *mode,        //  pingonly speedonly all
		PingMethod:    "googleping", // googleping
		SortMethod:    "rspeed",     // speed rspeed ping rping
		Concurrency:   2,
		TestMode:      2, // 2: ALLTEST 3: RETEST
		Subscription:  *link,
		Language:      "en", // en cn
		FontSize:      24,
		Theme:         "rainbow",
		Timeout:       10 * time.Second,
		OutputMode:    0, // 0: base64 1:file path 2: no pic 3: json 4: txt
	}
	ctx := context.Background()

	arr := make([]int, 0)

	nodes, err := web.TestContext(ctx, opts, &web.EmptyMessageWriter{})
	if err != nil {
		return arr, err
	}

	for _, node := range nodes {
		// tested node info here
		if node.IsOk {
			fmt.Println("id:", node.Id, node.Remarks, "ping:", node.Ping)
			o, _ := strconv.Atoi(node.Ping)
			arr = append(arr, o)
		}
	}
	return arr, nil
}
