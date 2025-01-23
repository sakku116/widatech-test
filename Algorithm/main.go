package main

import (
	"fmt"
	"strconv"
)

func Solution(l, t int) [][]int {
	results := make([][]int, 0)
	digits := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	states := make(map[string]int) // {id:found}

	stack := make([]int, l)
	for i := 0; i < l; i++ {
		stack[i] = 1
	}

	for {
		fmt.Println("current states", states)
		invalidStatesCount := 0
		for _, v := range states {
			if v >= l {
				invalidStatesCount++
			}
		}
		if invalidStatesCount > 0 && invalidStatesCount >= len(states) {
			break
		}
		for stackIdx := range stack {
			// apply digit
			for _, digit := range digits {
				// check if digit is used
				isUsed := false
				for _, v := range stack {
					if v == digit {
						isUsed = true
					}
				}
				if !isUsed {
					fmt.Println("new digit", digit)
					stack[stackIdx] = digit
					fmt.Println("stack", stack)

					// check if combination reached
					total := 0
					for _, v := range stack {
						total += v
					}
					fmt.Println("total", total)
					// fmt.Println("t", t)

					if total != t { // combination not reached
						nextIndex := stackIdx + 1
						if nextIndex >= len(stack) {
							nextIndex = 0
						}
						fmt.Println("combination not reached, move to next position to index " + strconv.Itoa(nextIndex))
						break
					}

					// check if combination already exists
					fmt.Println("current results", results)
					isExists := false
					var foundDuplicateStack []int
					for _, tmpStack := range results {
						numMap := make(map[int]interface{})
						for _, v := range tmpStack {
							numMap[v] = nil
						}

						fmt.Println("numMap", numMap)

						matchCount := 0
						for _, v := range stack {
							_, ok := numMap[v]
							if ok {
								matchCount++
							}
						}

						fmt.Println("tmpStack, stack", tmpStack, stack)

						if matchCount == l {
							isExists = true
							foundDuplicateStack = make([]int, len(tmpStack))
							copy(foundDuplicateStack, tmpStack)
							break
						}
					}

					stackId := ""
					for _, v := range foundDuplicateStack {
						stackId += strconv.Itoa(v)
					}

					if stackId == "" {
						for _, v := range stack {
							stackId += strconv.Itoa(v)
						}
					}

					if isExists { // combination not exists
						fmt.Println(stack, "already exists")
						states[stackId] += 1
						continue
					} else {
						states[stackId] = 1
					}

					copyStack := make([]int, l)
					copy(copyStack, stack)
					results = append(results, copyStack)
					fmt.Println("new combination", stack)
				}
			}
		}
	}

	return results
}

func main() {
	result := Solution(3, 9)
	fmt.Println(result)
}
