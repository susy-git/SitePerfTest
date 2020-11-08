package main

import (
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"
)

func main() {
	url := os.Args[1]                            //"https://www.baidu.com"
	totalNum, _ := strconv.Atoi(os.Args[2])      //100
	concurrentNum, _ := strconv.Atoi(os.Args[3]) //10

	var totalRespTimes []int64
	for i := 0; i < totalNum; i = i + concurrentNum {
		responseTimes := connectSite(url, concurrentNum)
		totalRespTimes = append(totalRespTimes, responseTimes...)
	}
	sort.Slice(totalRespTimes, func(i, j int) bool { return totalRespTimes[i] < totalRespTimes[j] })

	fmt.Println("Performance test result for " + url)
	fmt.Println("Total tests: " + strconv.Itoa(totalNum) + ". Concurrent tests: " + strconv.Itoa(concurrentNum))
	p95 := totalRespTimes[totalNum*95/100]
	fmt.Println(" -- %95 response time: " + strconv.FormatInt(p95, 10) + "ms")

	sum := int64(0)
	for i := 0; i < totalNum; i++ {
		sum += totalRespTimes[i]
	}
	avg := sum / int64(totalNum)
	fmt.Println(" -- average response time: " + strconv.FormatInt(avg, 10) + "ms")
}

func connectSite(url string, concurrentNum int) []int64 {
	c := make(chan int64)
	var wg sync.WaitGroup

	for i := 0; i < concurrentNum; i++ {
		wg.Add(1)
		go checkSite(url, c, &wg)
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	var respTimes []int64
	for responseTimeInMs := range c {
		//fmt.Println(responseTimeInMs)
		respTimes = append(respTimes, responseTimeInMs)
	}

	return respTimes
}

func checkSite(url string, c chan int64, wg *sync.WaitGroup) {
	defer (*wg).Done()
	start := time.Now()
	_, err := http.Get(url)

	if err != nil {
		c <- 0
	} else {
		c <- time.Now().Sub(start).Milliseconds()
	}
}
