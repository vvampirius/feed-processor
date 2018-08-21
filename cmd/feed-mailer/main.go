package main

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"time"
	"log"
	"net/mail"
	"flag"
	"os"
	"github.com/jaytaylor/html2text"
	"github.com/vvampirius/feed-processor/email"
	fpFeed "github.com/vvampirius/feed-processor/feed"
	"github.com/vvampirius/feed-processor/filterDSL"
)

func makeMessage(item *gofeed.Item, from *mail.Address, to *mail.Address, prefix string) []byte {
	message := fmt.Sprintf("From: %s\r\n", from.String())
	message += fmt.Sprintf("To: %s\r\n", to.String())
	message += "Subject: "+makeSubject(item, prefix)+"\r\n"
	message += "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	message += "\r\n"
	message += item.Description
	message += "\r\n"
	message += item.Link
	return []byte(message)
}

func updateTitle(item *gofeed.Item) {
	if item.Title != `` { return }
	if s, err := html2text.FromString(item.Description); err==nil {
		if fw := firstWords(s, 10); fw!=`` { item.Title = fw }
	} else { log.Println(err) }
}

func makeSubject(item *gofeed.Item, prefix string) string {
	if prefix != `` { return fmt.Sprintf("%s %s", prefix, item.Title) }
	return item.Title
}

// https://www.dotnetperls.com/first-words-go
func firstWords(value string, count int) string {
	// Loop over all indexes in the string.
	for i := range value {
		// If we encounter a space, reduce the count.
		if value[i] == ' ' {
			count -= 1
			// When no more words required, return a substring.
			if count == 0 {
				return value[0:i] + `...`
			}
		}
	}
	// Return the entire string.
	return value
}

func getAfterTime(s string) time.Time {
	if s == `` { return time.Time{} }
	if s == `now` { return time.Now() }
	if after, err := time.Parse(`2006-01-02 15:04:05 -07`, s); err == nil {
		return after
	} else { log.Println(err) }
	return time.Time{}
}

func getNewAfterTime(after time.Time, items []*gofeed.Item) time.Time {
	for _, item := range items {
		if item.PublishedParsed.After(after) { after = *item.PublishedParsed }
	}
	return after
}



func main() {
	fromFlag := flag.String(`f`, os.Getenv(`FROM`), `From:`)
	toFlag := flag.String(`t`, os.Getenv(`TO`), `To:`)
	afterFlag := flag.String(`a`, os.Getenv(`AFTER`), `Posts after this time (2006-01-02 15:04:05 -07)`)
	onceFlag := flag.Bool(`o`, false, `Once time (no cycle)`)
	intervalFlag := flag.Int(`i`, 30, `Interval`)
	hostPortFlag := flag.String(`m`, `smtp.gmail.com:465`, `SMTP TLS host:port`)
	userFlag := flag.String(`u`, os.Getenv(`SMTP_USERNAME`), `SMTP username`)
	passwordFlag := flag.String(`p`, os.Getenv(`SMTP_PASSWORD`), `SMTP password`)
	prefixFlag := flag.String(`q`, os.Getenv(`PREFIX`), `Prefix in mail subject`)
	fdslFlag := flag.String(`k`, os.Getenv(`FDSL`), `Filter DSL filename`)
	flag.Parse()

	url := flag.Arg(0)
	if url == `` { panic(`no url specified!`) }

	from, err := mail.ParseAddress(*fromFlag)
	if err != nil { log.Panic(err) }

	to, err := mail.ParseAddress(*toFlag)
	if err != nil { log.Panic(err) }

	logger := log.New(os.Stdout, ``, 3)

	fdsl := filterDSL.FilterDSL{FileName: *fdslFlag, FileTimestamp: time.Now()}
	fdsl.Reload()
	go fdsl.CheckUpdate()

	after := getAfterTime(*afterFlag)

	for true {
		messages := make([][]byte, 0)

		logger.Printf("Fetching url %s\n", url)

		fp := gofeed.NewParser()
		feed, err := fp.ParseURL(url)
		if err != nil {
			logger.Println(err)
			time.Sleep(time.Minute * 1)
			continue
		}

		logger.Printf("Got %d items\n", len(feed.Items))

		feed.Items = fpFeed.AfterFilter(feed.Items, after, logger)

		if regexpFilterPipeline, err := fdsl.RegexpFilterPipeline(url); err==nil {
			feed.Items = fpFeed.TitleFilter(feed.Items, regexpFilterPipeline, logger)
		} else { log.Printf("main: fdsl.RegexpFilterPipeline(%s) return error: %s\n", url, err.Error()) }

		for _, b := range feed.Items {
			updateTitle(b)
			logger.Println(b.Title)
			messages = append(messages, makeMessage(b, from, to, *prefixFlag))
		}

		after = getNewAfterTime(after, feed.Items) //TODO: поднять вверх, и встроить в AfterFilter функцию

		if len(messages)>0 {
			email.SendMessages(*hostPortFlag, *userFlag, *passwordFlag, from, to, messages)
			//fmt.Println(*hostPortFlag, *userFlag, *passwordFlag, from, to, messages)
		}

		if *onceFlag { break }

		time.Sleep(time.Duration(*intervalFlag) * time.Minute)
	}


}
