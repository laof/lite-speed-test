package ping

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/laof/lite-speed-test/web"
)

type TestResponse struct {
	SuccessIndex []int
	ErrorIndex   []int
	Ok           bool
}

func Test(url string) (TestResponse, error) {
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

	res := TestResponse{SuccessIndex: []int{}, ErrorIndex: []int{}, Ok: false}

	nodes, err := web.TestContext(ctx, opts, &web.EmptyMessageWriter{})
	if err != nil {
		return res, err
	}

	for _, node := range nodes {
		// tested node info here
		if node.IsOk {
			fmt.Println("SUCCESS id:", node.Id, node.Remarks, "ping:", node.Ping)
			res.SuccessIndex = append(res.SuccessIndex, node.Id)
		} else {
			fmt.Println("ERROR id:", node.Id)
			res.ErrorIndex = append(res.ErrorIndex, node.Id)
		}
	}

	if len(res.SuccessIndex) > 0 && len(res.ErrorIndex) == 0 {
		res.Ok = true
	}

	return res, nil
}
