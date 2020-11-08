package sensitive

// Trie 短语组成的Trie树.
type Trie struct {
	Root *Node
}

// Node Trie树上的一个节点.
type Node struct {
	isRootNode bool
	isPathEnd  bool
	Character  rune
	Children   map[rune]*Node
}

// NewTrie 新建一棵Trie
func NewTrie() *Trie {
	return &Trie{
		Root: NewRootNode(0),
	}
}

// Add 添加若干个词
func (tree *Trie) Add(words ...string) {
	for _, word := range words {
		tree.add(word)
	}
}

func (tree *Trie) add(word string) {
	var current = tree.Root
	var runes = []rune(word)
	for position := 0; position < len(runes); position++ {
		r := runes[position]
		if next, ok := current.Children[r]; ok {
			current = next
		} else {
			newNode := NewNode(r)
			current.Children[r] = newNode
			current = newNode
		}
		if position == len(runes)-1 {
			current.isPathEnd = true
		}
	}
}

func (tree *Trie) Del(words ...string) {
	for _, word := range words {
		tree.del(word)
	}
}

func (tree *Trie) del(word string) {
	var current = tree.Root
	var runes = []rune(word)
	for position := 0; position < len(runes); position++ {
		r := runes[position]
		if next, ok := current.Children[r]; !ok {
			return
		} else {
			current = next
		}

		if position == len(runes)-1 {
			current.SoftDel()
		}
	}
}

// Replace 词语替换
func (tree *Trie) Replace(text string, character rune) string {
	var (
		parent  = tree.Root
		current *Node
		runes   = []rune(text)
		length  = len(runes)
		left    = 0
		found   bool
	)

	for position := 0; position < len(runes); position++ {
		current, found = parent.Children[runes[position]]

		if !found || (!current.IsPathEnd() && position == length-1) {
			parent = tree.Root
			position = left
			left++
			continue
		}

		// println(string(current.Character), current.IsPathEnd(), left)
		if current.IsPathEnd() && left <= position {
			for i := left; i <= position; i++ {
				runes[i] = character
			}
		}

		parent = current
	}

	return string(runes)
}

// Filter 直接过滤掉字符串中的敏感词
func (tree *Trie) Filter(text string) string {
	var (
		parent      = tree.Root
		current     *Node
		left        = 0
		found       bool
		runes       = []rune(text)
		length      = len(runes)
		resultRunes = make([]rune, 0, length)
	)

	for position := 0; position < length; position++ {
		current, found = parent.Children[runes[position]]

		if !found || (!current.IsPathEnd() && position == length-1) {
			resultRunes = append(resultRunes, runes[left])
			parent = tree.Root
			position = left
			left++
			continue
		}

		if current.IsPathEnd() {
			left = position + 1
			parent = tree.Root
		} else {
			parent = current
		}

	}

	resultRunes = append(resultRunes, runes[left:]...)
	return string(resultRunes)
}

//HighLight 高亮敏感词
func (tree *Trie) HighLight(text string) string {
	var (
		parent  = tree.Root
		current *Node
		runes   = []rune(text)
		length  = len(runes)
		left    = 0
		found   bool

		runesRightTemp []rune
		runesLeftTemp  []rune
		runesSensitive []rune
		//runesTemp      []rune
		//runesPart = make([][]rune, 0)
		//sensitiveRunes []rune
		leftPositions = make([]int, 0) //用于记录所有敏感词的左侧位置
		//leftIndex     = 0
		rightPostions = make([]int, 0) //用于记录所有敏感词的右侧位置
		//rightIndex    = 0
		//runesRight []rune
		highLightLeft   = []rune(`<span class="sensitive_word" style="background-color: rgb(247, 218, 100);">`)
		hightLightRight = []rune("</span>")
		isNeedToTag     = true
	)

	for position := 0; position < len(runes); position++ {
		current, found = parent.Children[runes[position]]

		if !found || (!current.IsPathEnd() && position == length-1) {
			parent = tree.Root
			position = left
			left++
			continue
		}

		if current.IsPathEnd() && left <= position {

			leftPositions = append(leftPositions, left)
			rightPostions = append(rightPostions, position)
			left = position + 1
		}

		parent = current
	}
	//在敏感词左右侧添加html标签<span style=\"color: rgb(230, 0, 0);\"> xxx </span>
	//方案1： 按照敏感词位置，将runes数组拆分成段，并重新进行组合,貌似有点麻烦
	//你好傻逼呀，真的号傻逼要大姐夫也傻逼呀  2 9 13 16       3 10 14 17
	//runes[:2] runes[2:4] runes[4:9] runes[9:11] runes[11:13] runes[13:15] runes[15:16] runes[16:18] runes{18:]
	//for i := 0; i < len(leftPositions)*2; i++ {
	//	if i == 0 {
	//		runesPart[i] = runes[:leftPositions[leftIndex]]
	//	} else if i == len(leftPositions)*2-1 {
	//		if rightPostions[i] == len(text) {
	//			runesPart[i] = make([]rune, 0)
	//		}
	//		runesPart[i] = runes[rightPostions[rightIndex]:]
	//	} else {
	//		if i%2 == 0 {
	//			runesPart[i] = runes[rightPostions[rightIndex]+1 : leftPositions[leftIndex]]
	//			rightIndex++
	//		} else {
	//			runesPart[i] = runes[leftPositions[leftIndex] : rightPostions[rightIndex]+1]
	//			leftIndex++
	//		}
	//	}
	//
	//}
	//方案2：按照反向顺序添加标签，避免影响之前记录的顺序
	for i := len(leftPositions) - 1; i >= 0; i-- {
		isNeedToTag = false
		//判断是否已添加了标记
		if leftPositions[i] > 76 {
			runesHighLight := runes[leftPositions[i]-76 : leftPositions[i]]
			for i := 0; i < len(runesHighLight)-1; i++ {
				if runesHighLight[i+1] != highLightLeft[i] {
					isNeedToTag = true
					break
				}
			}
			if !isNeedToTag {
				isNeedToTag = false
				continue
			}
		}

		if rightPostions[i] == length-1 {
			runesRightTemp = make([]rune, 0)
		} else {
			temp := runes[rightPostions[i]+1:]
			runesRightTemp = make([]rune, len(temp))
			copy(runesRightTemp, runes[rightPostions[i]+1:])
		}

		//println("右侧为：" + string(runesRightTemp))
		runesSensitive = runes[leftPositions[i] : rightPostions[i]+1]
		//println("高亮前的敏感词为:" + string(runesSensitive))
		if leftPositions[i] == 0 {
			runesLeftTemp = make([]rune, 0)
		} else {
			temp := runes[:leftPositions[i]]
			runesLeftTemp = make([]rune, len(temp))
			copy(runesLeftTemp, temp)
		}
		//println("左侧为：" + string(runesLeftTemp))
		runesSensitive = append(append(highLightLeft, runesSensitive...), hightLightRight...)
		//println("高亮后的敏感词为:" + string(runesSensitive))
		runesLeftTemp = append(runesLeftTemp, runesSensitive...)
		//println("此时的左侧为：" + string(runesLeftTemp))
		//println("此时的右侧为：" + string(runesRightTemp))
		//备注：切片后，更改会影响原数组的值
		runes = append(runesLeftTemp, runesRightTemp...)
		//println("此时的语句为：" + string(runes))
	}
	return string(runes)
}

// Validate 验证字符串是否合法，如不合法则返回false和检测到
// 的第一个敏感词
func (tree *Trie) Validate(text string) (bool, string) {
	const (
		Empty = ""
	)
	var (
		parent  = tree.Root
		current *Node
		runes   = []rune(text)
		length  = len(runes)
		left    = 0
		found   bool
	)

	for position := 0; position < len(runes); position++ {
		current, found = parent.Children[runes[position]]

		if !found || (!current.IsPathEnd() && position == length-1) {
			parent = tree.Root
			position = left
			left++
			continue
		}

		if current.IsPathEnd() && left <= position {
			return false, string(runes[left : position+1])
		}

		parent = current
	}

	return true, Empty
}

// FindIn 判断text中是否含有词库中的词
func (tree *Trie) FindIn(text string) (bool, string) {
	validated, first := tree.Validate(text)
	return !validated, first
}

// FindAll 找有所有包含在词库中的词
func (tree *Trie) FindAll(text string) []string {
	var matches []string
	var (
		parent  = tree.Root
		current *Node
		runes   = []rune(text)
		length  = len(runes)
		left    = 0
		found   bool
	)

	for position := 0; position < length; position++ {
		current, found = parent.Children[runes[position]]

		if !found {
			parent = tree.Root
			position = left
			left++
			continue
		}

		if current.IsPathEnd() && left <= position {
			matches = append(matches, string(runes[left:position+1]))
		}

		if position == length-1 {
			parent = tree.Root
			position = left
			left++
			continue
		}

		parent = current
	}

	var i = 0
	if count := len(matches); count > 0 {
		set := make(map[string]struct{})
		for i < count {
			_, ok := set[matches[i]]
			if !ok {
				set[matches[i]] = struct{}{}
				i++
				continue
			}
			count--
			copy(matches[i:], matches[i+1:])
		}
		return matches[:count]
	}

	return nil
}

// NewNode 新建子节点
func NewNode(character rune) *Node {
	return &Node{
		Character: character,
		Children:  make(map[rune]*Node, 0),
	}
}

// NewRootNode 新建根节点
func NewRootNode(character rune) *Node {
	return &Node{
		isRootNode: true,
		Character:  character,
		Children:   make(map[rune]*Node, 0),
	}
}

// IsLeafNode 判断是否叶子节点
func (node *Node) IsLeafNode() bool {
	return len(node.Children) == 0
}

// IsRootNode 判断是否为根节点
func (node *Node) IsRootNode() bool {
	return node.isRootNode
}

// IsPathEnd 判断是否为某个路径的结束
func (node *Node) IsPathEnd() bool {
	return node.isPathEnd
}

// SoftDel 置软删除状态
func (node *Node) SoftDel() {
	node.isPathEnd = false
}
