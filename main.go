package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"
)

// ProducerInfo ...
type ProducerInfo struct {
	id      int64
	score   float64
	vote    float64
	chance  int64 // the total blocks the producer produced
	roundId int64
	sumto   float64
	per     float64
	max     float64
}

var oneBlockBonus = 3.41222746 // iwallet -s 18.209.137.246:30002 table bonus.iost blockContrib

func initProducer(votes []int) map[int64]*ProducerInfo {
	producers := make(map[int64]*ProducerInfo)

	sum := float64(0)
	for i, vote := range votes {
		sum += float64(vote)
		// m: [1, 1, 1, ..., 1, 2, 3, ..., 15, 16, 17]
		m := 18 + i - len(votes)
		if m <= 0 {
			m = 1
		}
		producers[int64(i)] = &ProducerInfo{
			id:     int64(i),
			score:  float64(0),
			vote:   float64(vote),
			chance: 0,
			sumto:  sum,
			per:    float64(vote) / sum,
			max:    1.0 / float64(m),
		}
	}
	return producers
}

func iterateOneRound(producers map[int64]*ProducerInfo, producerIds map[int64]bool) {
	//        fmt.Println("producerIds", producerIds)
	// calc score
	totalScore := float64(0)
	var pending []*ProducerInfo
	for id := 0; id < len(producers); id++ {
		producer := producers[int64(id)]
		producer.roundId++
		//fmt.Println("node id", id, "add score", producer.score, "to", producer.score + producer.vote, "vote", producer.vote)
		producer.score += producer.vote
		totalScore += producer.score
		if ok := producerIds[int64(id)]; !ok {
			pending = append(pending, producer)
		}
	}
	// pre sort
	sort.SliceStable(pending, func(i, j int) bool {
		return pending[i].score > pending[j].score
	})
	//fmt.Print("sorted\t\n")
	//        for id := range pending {
	//              fmt.Printf("%v\t", pending[id].score)
	//      }
	//fmt.Print("split\n")

	// limit replace num
	//pending = pending[:17-len(producerIds)+4]
	for id := range producerIds {
		pending = append(pending, producers[id])
		//fmt.Printf("%v\t", producers[id].score)
	}
	//fmt.Printf("end\n")
	// sort
	sort.SliceStable(pending, func(i, j int) bool {
		return pending[i].score > pending[j].score
	})
	// new producerIds
	producerIds = make(map[int64]bool)
	for j := 0; j < 17; j++ {
		//              fmt.Println("score", pending[j].score)
		// producer blocks here
		producerIds[pending[j].id] = true
		//producers[pending[j].id].chance += (int64)((float64)(10 * 60 * 2 * oneBlockBonus)/ 17.0 ) //1000
		//producers[pending[j].id].chance += (int64)((float64)(10 * 60 * 2 * oneBlockBonus)/ 17.0 ) //1000
		producers[pending[j].id].chance += 1 //oneBlockBonus
	}
	// after the round, subtract score
	totalProScore := float64(0)
	for k := range producerIds {
		totalProScore += producers[k].score
	}
	// proAvg := totalProScore / 17
	avg := totalScore / float64(len(producers))
	//fmt.Println("avg", avg)
	// avg := stableTotalScore / 50
	for k := range producerIds {
		// producers[k].score = math.Max(0, producers[k].score-minScore)
		producers[k].score = math.Max(0, producers[k].score-avg/10)
		//              fmt.Println("minus to", producers[k].score)
		// producers[k].score = producers[k].score - avg/10
		// producers[k].score = math.Max(0, producers[k].score-proAvg/100)
		// producers[k].score = producers[k].score/2 + avg/5
	}

}

func main() {
	rand.Seed(time.Now().UnixNano())
	votes := []int{}
	for _, v := range allVotes {
		votes = append(votes, (int)(v))
	}
	sort.Ints(votes)
	//fmt.Println("votes", votes)
	producers := initProducer(votes)
	producerIds := make(map[int64]bool)

	loop := 365 * 24 * 6 // 52560
	//loop = 50000
	for i := 0; i < loop; i++ {
		iterateOneRound(producers, producerIds)
	}

	// after the iterations, print the statistic of producers here
	var prosList []*ProducerInfo
	for _, p := range producers {
		prosList = append(prosList, p)
	}
	sort.Slice(prosList, func(i, j int) bool {
		return prosList[i].id < prosList[j].id
	})
	for _, p := range prosList {
		// chance should be proportional to vote, exception the several producers with most votes since we try to discourage centralization of votes
		prop := float64(p.chance * 10 * 60 * 2) * oneBlockBonus / 17.0 / float64(p.vote) // 年化
		//reverseIdx := len(prosList) - idx
		//fmt.Printf("%v\t%2f\t\t%v\t%5f\t%5f\t%5f\t%5f\n", reverseIdx, p.vote, p.chance, prop, p.per, p.max, p.score)
		//fmt.Printf("%v\t%2f\t%v\t%5f\n", reverseIdx, p.vote, p.chance, prop) //, p.per, p.max, p.score)
		//fmt.Printf("%d,%d,%.2f%%\n", int(p.vote), p.chance, prop*100.0) //, p.per, p.max, p.score)
		fmt.Printf("%d,%.2f%%\n", int(p.vote), prop*100.0) //, p.per, p.max, p.score)
	}
}
