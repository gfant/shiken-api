package p2

func Problem(num int) int {
	if num < 2 {
		return 0
	}
	primes := make([]int, num)

	solution := 0
	primes[0] = 1
	primes[1] = 1
	for i := 2; i < num; i++ {
		if primes[i] == 0 {
			solution++
			for j := 2 * i; j < num; j += i {
				primes[j] = 1
			}
		}
	}
	return solution
}
