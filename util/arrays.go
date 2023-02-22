package util

import "unsafe"

func ReverseArray(arr []byte) []byte {
	revArr := make([]byte, len(arr))
	j := 0
	for i := len(arr) - 1; i >= 0; i-- {
		revArr[j] = arr[i]
		j++
	}
	return revArr
}

func IntToByteArray(num int) []byte {
	size := int(unsafe.Sizeof(num))
	arr := make([]byte, size)
	for i := 0; i < size; i++ {
		byt := *(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&num)) + uintptr(i)))
		arr[i] = byt
	}
	return arr
}

func FillArrayWithValue[V any](arr []V, value V) {
	for i := range arr {
		arr[i] = value
	}
}
