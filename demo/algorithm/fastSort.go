package main
//快排

import "fmt"

func main()  {
    arr := []int{1, 9, 10, 30, 2, 5, 45, 8, 63, 234, 12}
    fmt.Println(QuickSort2(arr))
}

//解析：第一个数据作为比较，将当前数据列表区分大中小三类
//将大和小分别递归，返回时已经排好序
//进行拼凑
func QuickSort1(arr []int) []int {
    if len(arr) <= 1 {
        return arr
    }
    splitdata := arr[0]
    low := make([]int, 0, 0)
    hight := make([]int, 0, 0)
    mid := make([]int, 0, 0)
    mid = append(mid, splitdata)
    for i := 1; i < len(arr); i++ {
        if arr[i] < splitdata {
            low = append(low, arr[i])
        } else if arr[i] > splitdata {
            hight = append(hight, arr[i])
        } else {
            mid = append(mid, arr[i])
        }
    }
    low, hight = QuickSort1(low), QuickSort1(hight)
    myarr := append(append(low, mid...), hight...)
    return myarr
}

func QuickSort2(arr []int) []int {
    return _quickSort(arr, 0, len(arr)-1)
}

func _quickSort(arr []int, left, right int) []int {
    if left < right {
        //用index将数组列表区分大小两类
        partitionIndex := partition(arr, left, right)
        //小于index的
        _quickSort(arr, left, partitionIndex-1)
        //大于index的
        _quickSort(arr, partitionIndex+1, right)
    }
    return arr
}

//要区分哪些是比第一个值小的
func partition(arr []int, left, right int) int {
    //最左边当做对比值
    pivot := left
    //0-index-1的值都是比对比值小的
    index := pivot + 1
    for i := index; i <= right; i++ {
        // i=index时，准备比较
        // i>index时，当前位置肯定是比对比值大的
        //发现有比对比值小的，与index位置调换，index增加1
        if arr[i] < arr[pivot] {
            swap(arr, i, index)
            index += 1
        }
    }
    //index-1是最后一个小于对比值，将对比值放到合适位置
    swap(arr, pivot, index-1)
    return index - 1
}

func swap(arr []int, i, j int) {
    arr[i], arr[j] = arr[j], arr[i]
}
