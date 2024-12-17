package main


func repeatLimitedString(s string, repeatLimit int) string {
	freqList := make([]int, 26)
	for _, c := range s {
		freqList[c - 'a']  += 1
	}

	cur, prev := 25, 24
	ans := []rune{}
	for true {
		for cur >= 0 && freqList[cur] == 0 {
			cur -= 1
		}
		prev = min(prev, cur-1)
		for prev >= 0 && freqList[prev] == 0 {
			prev -= 1
		}

		if cur == -1 {
			break
		}
		if freqList[cur] > repeatLimit {
			for i:=0; i<repeatLimit; i++ {
				ans = append(ans, rune('a'+cur))
			}
			freqList[cur] -= repeatLimit
			if prev == -1 {
				break
			}
			ans = append(ans, rune('a'+prev))
			freqList[prev] -= 1
		} else {
			for i:=0; i<freqList[cur]; i++ {
				ans = append(ans, rune('a'+cur))
			}
			freqList[cur] = 0
		}
	}	
	return string(ans)
}