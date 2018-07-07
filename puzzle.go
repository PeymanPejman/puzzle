package main

import (
	"container/heap"
	"fmt"
	"log"
	"math/rand"
	"time"
)

////////////// Priority Queue //////////////
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// The lower the sum, the higher the priority
	return pq[i].J < pq[j].J
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*Item)
	*pq = append(*pq, item)
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

////////////// BFS Algorithm //////////////

type Item struct {
	I int64
	J int64
}

// BfsHeuritic uses BFS with a simple heuristic to make the search faster.
func BfsHeuristic(mktcap []int64, sum int64) {

	childOf := make(map[*Item]*Item)
	var priorityQueue PriorityQueue
	priorityQueue = append(priorityQueue, &Item{I: 0, J: sum})
	heap.Init(&priorityQueue)
	n := int64(len(mktcap))
	var found *Item

	for priorityQueue.Len() > 0 {
		curr := heap.Pop(&priorityQueue).(*Item)
		//fmt.Println(curr.J)

		if curr.J == 0 {
			found = curr
			break
		}

		with := Item{I: curr.I + 1, J: curr.J}
		if with.I < n && with.J >= 0 {
			if _, ok := childOf[&with]; !ok {
				priorityQueue = append(priorityQueue, &with)
				childOf[&with] = curr
			}
		}

		without := Item{I: curr.I + 1, J: curr.J - mktcap[curr.I]}
		if without.I < n && without.J >= 0 {
			if _, ok := childOf[&without]; !ok {
				priorityQueue = append(priorityQueue, &without)
				childOf[&without] = curr
			}
		}
	}

	var res []int64
	if found == nil {
		fmt.Println("Subset does not exist")
		return
	}

	iter := found
	for iter.J != sum {
		toAdd := childOf[iter].J - iter.J
		if toAdd > 0 {
			res = append(res, toAdd)
		}
		iter = childOf[iter]
	}

	fmt.Printf("BFS Solution: %v\n", res)
}

////////////// DFS Algorithm //////////////

// SubsetSumDFS recursively fills out the dp "table" and uses memoization
func SubsetSumDFS(dp map[int64]map[int64]bool, mktcap []int64, n int64, sum int64) bool {
	if sum == 0 {
		return true
	}

	if n < 0 || sum < 0 {
		return false
	}
	if dp[n] == nil {
		dp[n] = make(map[int64]bool)
	}
	_, memoized := dp[n][sum]

	if memoized == false {
		dp[n][sum] = SubsetSumDFS(dp, mktcap, n-1, sum) ||
			SubsetSumDFS(dp, mktcap, n-1, sum-mktcap[n])
	}

	return dp[n][sum]
}

// PrintSubset uses the filled out dp table to print out one subset summing to sum
func PrintSubset(mktcap []int64, dp map[int64]map[int64]bool, sum int64) {
	var subset []int64
	n := int64(len(mktcap)) - 1

	for sum > 0 {
		if n == 0 && mktcap[n] == sum {
			subset = append(subset, mktcap[n])
			break
		}

		if dp[n-1][sum] == true {
			n = n - 1
		} else {
			subset = append(subset, mktcap[n])
			sum = sum - mktcap[n]
			n = n - 1

		}
	}

	fmt.Printf("DFS Solution: %v\n", subset)
}

////////////// Test Drivers //////////////

// TestDFS runs the DFS algorithm with the given inputs and returns the duration
func TestDFS(mktcap []int64, sum int64, n int64) time.Duration {
	start := time.Now()
	dp := make(map[int64]map[int64]bool)

	if SubsetSumDFS(dp, mktcap, n-1, sum) {
		PrintSubset(mktcap, dp, sum)
	} else {
		fmt.Printf("A subset summing to %v does not exist\n", sum)
	}

	elapsed := time.Since(start)
	//log.Printf("DFS with memoization took %s", elapsed)
	return elapsed
}

// TestBFSHeuristic runs the BFS algorithm with the given inputs and returns the duration
func TestBFSHeuristic(mktcap []int64, sum int64, n int64) time.Duration {
	start := time.Now()
	BfsHeuristic(mktcap, sum)
	elapsed := time.Since(start)
	//log.Printf("BFS with Heuristic took %s", elapsed)
	return elapsed
}

// Average takes in a collection of durations and returns their average
func Average(a []time.Duration) time.Duration {
	var total time.Duration
	for i := 0; i < len(a); i++ {
		total += a[i]
	}
	return total / time.Duration(len(a))
}

func main() {
	//Solution to the original problem solved using two different algorithms
	sum := int64(100000000000)
	mktcap := []int64{178990553235, 95104612655, 47003797210, 26824713718, 14326106534, 9605184103, 8404690765, 8400270113, 8377919999, 7725964999, 6451400968, 5896028330, 4762072812, 4455321066, 3100349607, 3026104544, 2829847251, 2813479534, 2385674223, 2280233615, 1958243508, 1638915652, 1530697657, 1520698987, 1364603792, 1308870062, 1253339015, 1177503043, 1144710872, 1138275645}
	n := int64(len(mktcap))
	log.Printf("DFS with memoization took %s", TestDFS(mktcap, sum, n))
	log.Printf("BFS with Heuristic took %s", TestBFSHeuristic(mktcap, sum, n))

	//Several trials based on smaller input to show avergae time for each algorithm
	sum = int64(30)
	mktcap = []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18}
	n = int64(len(mktcap))
	var dfstimes, bfstimes []time.Duration

	for trial := 0; trial < 50; trial++ {
		shuffled := make([]int64, len(mktcap))
		perm := rand.Perm(len(mktcap))
		for i, v := range perm {
			shuffled[v] = mktcap[i]
		}
		fmt.Printf("Input Trial %v: %v\n", trial, shuffled)
		dfstimes = append(dfstimes, TestDFS(shuffled, sum, n))
		bfstimes = append(bfstimes, TestBFSHeuristic(shuffled, sum, n))
	}

	fmt.Printf("Average DFS = %v\n", Average(dfstimes))
	fmt.Printf("Average BFS = %v\n", Average(bfstimes))
}
