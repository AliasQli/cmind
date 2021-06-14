package model

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/groovili/gogtrends"
	"net/http"
	"time"
	"unicode/utf16"
)

var (
	ctxBG = context.Background()
)

func getContext() (context.Context, context.CancelFunc){
	return context.WithTimeout(ctxBG, 5 * time.Second)
}

func getGoogleRelated(word string) ([]string, error) {
	ctx, cancel := getContext()
	defer cancel()
	explore, err := gogtrends.Explore(ctx,
		&gogtrends.ExploreRequest{
			ComparisonItems: []*gogtrends.ComparisonItem{
				{
					Keyword: word,
					Geo:     "CN",
					Time:    "today 12-m",
				},
			},
			Category: 0,
			Property: "",
		}, "ZH")
	if err != nil {
		return nil, err
	}

	var index int // Usually it would be 0
	for i, v := range explore {
		if v.ID == "RELATED_TOPICS" {
			index = i
			break
		}
	}

	relT, err := gogtrends.Related(ctx, explore[index], "ZH")
	if err !=nil {
		return nil, err
	}
	relT = removeDuplicate(relT)
	relL := make([]string, len(relT))
	for i, v := range relT {
		relL[i] = v.Topic.Title
	}
	return relL, nil
}

func removeDuplicate(slice []*gogtrends.RankedKeyword) []*gogtrends.RankedKeyword {
	type Indexed struct {
		Ix int
		Kw *gogtrends.RankedKeyword
	}

	ix := 0
	m := make(map[string]*Indexed)
	for _, v := range slice {
		if _, ok := m[v.Topic.Title]; !ok {
			m[v.Topic.Title] = &Indexed{
				Ix: ix,
				Kw: v,
			}
			ix++
		}
	}
	ret := make([]*gogtrends.RankedKeyword, ix)
	for _, ixd := range m {
		ret[ixd.Ix] = ixd.Kw
	}
	return ret
}

func getAizhanRelated(word string) ([]string, error) {
	encoded := ""
	for _, i := range utf16.Encode([]rune(word)) {
		if i < 256 {
			encoded += "n"
		}
		encoded += fmt.Sprintf("%x", i)
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://ci.aizhan.com/%s/", encoded), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	related := make([]string, 0, 5)
	doc.Find(".title a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		related = append(related, s.Text())
		return i < 4
	})
	return related, nil
}

func GetRelatedWordList(word string) ([]string, error) {
	g, err := getGoogleRelated(word)
	fmt.Println("Google: ", err)
	if err != nil && len(g) > 5 {
		return g, nil
	}

	a, err2 := getAizhanRelated(word)
	fmt.Println("Aizhan: ", err2)
	if err != nil && err2 != nil {
		return nil, err
	}

	return append(g, a...), nil
}
