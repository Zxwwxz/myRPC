package main

import "fmt"

//树结构
type TrieTree struct {
    root *TrieNode
}

//树节点结构
type TrieNode struct {
    //值
    value      int
    //26个字母的下一个节点
    dictionary [26]*TrieNode
    //到当前为止是否为一个单词
    over       bool
}

func main() {
    arrList := []string{"how", "hi", "her", "hello", "so", "see"}
    tree := createTree(arrList)
    flag := tree.findWord("h")
    fmt.Println(flag)
    flag = tree.findWord("hell")
    fmt.Println(flag)
    flag = tree.findWord("hello")
    fmt.Println(flag)
    flag = tree.findWord("helloe")
    fmt.Println(flag)
}

func createTree(arrList []string) *TrieTree {
    myTree := &TrieTree{}
    for _, value := range arrList {
        myTree.addWord(value)
    }
    return myTree
}

func (t *TrieTree) addWord(word string) {
    if t == nil {
        return
    }
    //根结点没有保存值
    if t.root == nil {
        t.root = &TrieNode{}
    }
    nowNode := t.root
    //得到a的ascii值
    a := int('a')
    var char int
    //遍历每个字母
    for i := 0; i < len(word); i++ {
        char = int(word[i])
        if nowNode.dictionary[char-a] != nil {
            nowNode = nowNode.dictionary[char-a]
        } else {
            newNode := &TrieNode{}
            nowNode.dictionary[char-a] = newNode
            nowNode = newNode
        }
        if i == len(word)-1 {
            nowNode.over = true
        }
    }
}

func (t *TrieTree) findWord(word string) bool {
    if t == nil {
        return false
    }
    nowNode := t.root
    if nowNode == nil {
        return false
    }
    a := int('a')
    var char int
    for i := 0; i < len(word); i++ {
        char = int(word[i])
        if nowNode.dictionary[char-a] == nil {
            return false
        } else {
            nowNode = nowNode.dictionary[char-a]
        }
        if i == len(word)-1 && nowNode.over {
            return true
        }
    }
    return false
}
