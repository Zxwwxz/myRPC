package main

import (
    "math"
)

//树结构
type Tree struct {
    //指向根节点
    root   *Node
    //树长度
    length int
}

//每个节点结构
type Node struct {
    //节点值
    value  int
    //左子树
    left   *Node
    //右子树
    right  *Node
    //父节点
    parent *Node
}

func main() {
    arrList := []int{14, 2, 5, 7, 23, 35, 12, 17, 31}
    CreatTree(arrList)
}

func CreatTree(arrList []int)(*Tree) {
    tree := &Tree{}
    for i := 0; i < len(arrList); i++ {
        tree = insertNode(tree, arrList[i])
    }
    return tree
}

func AddTreeNode(tree *Tree, insertValue int) {
    insertNode(tree, insertValue)
}

//二叉搜索树
func insertNode(tree *Tree, insertValue int) *Tree {
    var currentNode *Node
    var tmp *Node
    //没有节点
    if tree.length == 0 {
        tmp = new(Node)
        tmp.value = insertValue
        tree.root = tmp
        tree.length = 1
        return tree
    } else {
        currentNode = tree.root
    }
    for {
        //右子树
        if currentNode.value < insertValue {
            if currentNode.right == nil {
                tmp = new(Node)
                tmp.value = insertValue
                currentNode.right = tmp
                tmp.parent = currentNode
                break
            } else {
                currentNode = currentNode.right
                continue
            }
            //左子树
        } else {
            if currentNode.left == nil {
                tmp = new(Node)
                tmp.value = insertValue
                currentNode.left = tmp
                tmp.parent = currentNode
                break
            } else {
                currentNode = currentNode.left
                continue
            }
        }
    }
    tree.length = tree.length + 1
    return tree
}

func DelTreeNode(tree *Tree, delValue int) {
    if tree == nil {
        return
    }
}

func TreeHeight(tree *Tree)int {
    if tree == nil {
        return 0
    }
    return heightMax(tree.root, 0)
}

func heightMax(node *Node, h int) int {
    var hL = h
    var hR = h
    if node.left == nil && node.right == nil {
        return h
    }
    if node.left != nil {
        h++
        hL = heightMax(node.left, h)
    }
    if node.right != nil {
        h++
        hR = heightMax(node.right, h)
    }
    return int(math.Max(float64(hL), float64(hR)))
}

func LDR(tree *Tree) {

}