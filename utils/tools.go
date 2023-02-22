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
func IntToInt64(num int)int64{
	return int64(num)
}
func  IntToInt32(num int)int32{
	return int32(num)
}
func Int64ToString(num int64)string{
	return strconv.FormatInt(num, 10)
}
func Int64ToUInt(num int64)uint{
	return StringToUint(Int64ToString(num))
}
func StringToInt(num string)int{
	temp ,_ := strconv.Atoi(num)
	return temp
}

func UintToInt64(num uint)int64{
	ans,_:= strconv.ParseInt(UintToString(num), 10, 64)
	return ans
}
func StringToInt64(num string)int64{
	i, _ := strconv.ParseInt(num, 10, 64)
	return  i
}
