package utils

import "strconv"

func StringToUint(num string)uint{
	intNum, _ := strconv.Atoi(num)
	return uint(intNum)
}

func UintToString(num uint)string{
	b := strconv.Itoa(int(num))
	return string(b)
}

func IntToString(num int)string{
	return strconv.Itoa(num)
}

func Int64ToString(num int64)string{
	return strconv.FormatInt(num, 10)
}

func StringToInt(num string)int{
	temp ,_ := strconv.Atoi(num)
	return temp
}


