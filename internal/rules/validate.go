package rules

import "unicode"

func LuhnValid(num string) bool {
	digits := onlyDigits(num)
	if len(digits) < 13 || len(digits) > 19 {
		return false
	}
	sum := 0
	alt := false
	for i := len(digits) - 1; i >= 0; i-- {
		n := int(digits[i] - '0')
		if alt {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
		alt = !alt
	}
	return sum%10 == 0
}

func SnilsValid(sn string) bool {
	d := onlyDigits(sn)
	if len(d) != 11 {
		return false
	}
	num := d[:9]
	ctrl := (int(d[9]-'0') * 10) + int(d[10]-'0')

	sum := 0
	for i := 0; i < 9; i++ {
		sum += int(num[i]-'0') * (9 - i)
	}
	var check int
	if sum < 100 {
		check = sum
	} else if sum == 100 || sum == 101 {
		check = 0
	} else {
		check = sum % 101
		if check == 100 {
			check = 0
		}
	}
	return check == ctrl
}

func InnValid(inn string) bool {
	d := onlyDigits(inn)
	if len(d) == 10 {
		k := []int{2, 4, 10, 3, 5, 9, 4, 6, 8}
		s := 0
		for i := 0; i < 9; i++ {
			s += int(d[i]-'0') * k[i]
		}
		c := (s % 11) % 10
		return c == int(d[9]-'0')
	}
	if len(d) == 12 {
		k1 := []int{7, 2, 4, 10, 3, 5, 9, 4, 6, 8}
		k2 := []int{3, 7, 2, 4, 10, 3, 5, 9, 4, 6, 8}
		s1 := 0
		for i := 0; i < 10; i++ {
			s1 += int(d[i]-'0') * k1[i]
		}
		c1 := (s1 % 11) % 10
		if c1 != int(d[10]-'0') {
			return false
		}
		s2 := 0
		for i := 0; i < 11; i++ {
			s2 += int(d[i]-'0') * k2[i]
		}
		c2 := (s2 % 11) % 10
		return c2 == int(d[11]-'0')
	}
	return false
}

func onlyDigits(s string) string {
	out := make([]rune, 0, len(s))
	for _, r := range s {
		if unicode.IsDigit(r) {
			out = append(out, r)
		}
	}
	return string(out)
}