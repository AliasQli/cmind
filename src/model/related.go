package model

import (
	"context"
	"github.com/groovili/gogtrends"
	"time"
)

var (
	ctxBG = context.Background()
)

func getContext() (context.Context, context.CancelFunc){
	return context.WithTimeout(ctxBG, 5 * time.Second)
}

func getSingleRelatedTopics(word string) ([]*gogtrends.RankedKeyword, error) {
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

	return gogtrends.Related(ctx, explore[index], "ZH")
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

func GetSingleRelatedWordList(word string) ([]string, error) {
	relT, err := getSingleRelatedTopics(word)
	if err != nil {
		return nil, err
	}

	relT = removeDuplicate(relT)

	relL := make([]string, len(relT))
	for i, v := range relT {
		relL[i] = v.Topic.Title
	}
	return relL, nil
}

func GetSingleRelatedCategorized(word string) (map[string][]string, error) {
	relT, err := getSingleRelatedTopics(word)
	if err != nil {
		return nil, err
	}

	relT = removeDuplicate(relT)

	relM := make(map[string][]string)
	for _, v := range relT {
		if relM[v.Topic.Type] == nil {
			relM[v.Topic.Type] = []string {v.Topic.Title}
		} else {
			relM[v.Topic.Type] = append(relM[v.Topic.Type], v.Topic.Title)
		}
	}
	return relM, nil
}