package main

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	name, _ := os.Hostname()
	log.Println("start stress test", name)

	users := 20
	if os.Getenv("USERS") != "" {
		if i, err := strconv.Atoi(os.Getenv("USERS")); err == nil && i > 0 {
			users = i
		}
	}

	if u := os.Getenv("LINK"); u != "" {
		link, err := url.Parse(u)
		if err != nil {
			log.Println(err)
			return
		}

		s := strings.Split(link.Path, "/")
		roomID := s[len(s)-1]
		log.Println("roomID", roomID)

		_, err = os.Stat(roomID)

		if os.IsNotExist(err) {
			errDir := os.MkdirAll(roomID, 0755)
			if errDir != nil {
				log.Println(err)
				return
			}
		}

		for i := 0; i < users; i++ {
			go func() {
				name := fmt.Sprintf("%v-%v", name, i)
				l := launcher.New().
					//Set("--disable-gpu").
					//Set("disable-web-security").
					Set("start-fullscreen").
					Set("--use-fake-device-for-media-stream").
					Set("--use-fake-ui-for-media-stream").
					Headless(true)

				url := l.MustLaunch()

				rd := rod.New().ControlURL(url).MustConnect()
				//defer rd.Close()

				page := rd.MustPage(u).MustWaitLoad()
				page.MustElement(`.input-group input[placeholder="Enter your name!"]`).MustWaitVisible().
					MustInput(name)
				page.MustElement(`button.join-form`).MustWaitVisible().MustClick()
				page.MustElement(`button[class="lg--Q7ufB buttonWrapper--x8uow button--qv0Xy btn--29prju"`).MustWaitVisible().MustClick()
				page.MustElement(`button[aria-label="Start sharing"`).MustWaitVisible().MustClick()
				log.Println(name)
				time.Sleep(5 * time.Second)
			}()
			time.Sleep(5 * time.Second)
		}
	}

	var idle = 100
	if os.Getenv("IDLE") != "" {
		if id, err := strconv.Atoi(os.Getenv("IDLE")); err == nil && id > 0 {
			idle = id
		}
	}
	log.Println("waiting idle ", idle, "minutes")
	time.Sleep(time.Duration(idle) * time.Minute)
	log.Println("DONE")
}
