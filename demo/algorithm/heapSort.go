package main

import "fmt"

//堆排序
func main() {
    arr := []int{1, 9, 10, 30, 2, 5, 45, 8, 63, 234, 12}
    fmt.Println(HeapSort(arr))
}

func HeapSortMax(arr []int, lastIndex int) []int {
    if lastIndex <= 0 {
        return arr
    }
    //找到最后一个节点的父节点
    depth := (lastIndex-1)/2
    //遍历每个父节点
    for i := depth; i >= 0; i-- {
        topmax := i
        leftchild := 2*i + 1
        rightchild := 2*i + 2
        //找到左右孩子中最大的，而且大于父节点的，与父节点交换
        if leftchild <= lastIndex && arr[leftchild] > arr[topmax] {
            topmax = leftchild
        }
        if rightchild <= lastIndex && arr[rightchild] > arr[topmax] { //防止越过界限
            topmax = rightchild
        }
        if topmax != i {
            arr[i], arr[topmax] = arr[topmax], arr[i]
        }
    }
    return arr
}

func HeapSort(arr []int) []int {
    length := len(arr)
    //从最后一位开始，找到最大的，推选到顶部，然后换到最后一位
    //再从倒数第二位开始，循环往复
    for i := 0; i < length; i++ {
        //未排序的最后一个
        lastIndex := length - 1 - i
        //找到最大的推到顶部
        HeapSortMax(arr, lastIndex)
        //最大的和最后一个交换
        arr[0], arr[lastIndex] = arr[lastIndex], arr[0]
    }
    return arr
}
