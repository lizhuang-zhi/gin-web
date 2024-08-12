package main

/*
行为树主要由以下几种节点组成：
Composite nodes：这些节点有多个孩子节点。他们控制它们孩子节点的执行顺序或条件。
Decorator nodes：这些节点只有一个孩子节点。他们改变或扩展其孩子的行为，例如反转其成功/失败状态，或延迟/重复其执行。
Leaf nodes：这些节点是实际的行动或判断，例如移动到某个地方，跟随某个角色，判断是否有敌人在视野中等。

想象你正在制作一个角色AI，它的主要任务是寻找食物。但是，如果有敌人出现，它应该忘记食物，转而尽可能避开敌人。这个过程可以用一个简单的行为树来表示：
ROOT
|
+- SEQUENCE

	|
	+- Selector
	|  |
	|  +- Is there an enemy nearby?
	|  |
	|  +- Run away
	|
	+- Find food

在这个树中，SEQUENCE节点会依次运行其子节点，直到一个失败为止。首先，它运行Selector节点，它会尝试运行其子节点，
直到找到一个成功的节点。如果“Is there an enemy nearby?”节点返回成功（也就是说有敌人在附近），
那么它就运行“Run away”节点。否则，Selector节点失败，然后SEQUENCE节点运行“Find food”。

这个示例定义了一个Node接口，这个接口有一个方法Tick，用来更新节点的状态并返回结果。
CompositeNode是一个简单的结构，它包含一些Node儿子。Sequence和Selector都是CompositeNode，分别实现了它们的行为。LeafNode是一个具体的行为，它接受一个函数，并在Tick时运行它。
编写实际的任务函数的具体内容取决于你的游戏环境。对于Is there an enemy nearby?节点，该函数可能会检查周围一定半径内是否有敌人。这些函数必须匹配func() NodeState签名。
以上就是一个非常基础的行为树实现。这只是开始，你可以根据需要添加更多的节点类型和行为。
例如，你可能需要实现的某些行为包括：并行执行子节点、只运行一次、运行直到成功、限制运行次数等。你也需要一种方法来构建和存储树结构。
*/
type NodeState int

const (
	SUCCESS NodeState = iota
	FAILURE
	RUNNING
)

type Node interface {
	Tick() NodeState
}

type CompositeNode struct {
	children []Node
}

type Selector struct {
	CompositeNode
}

func (s *Selector) Tick() NodeState {
	for _, child := range s.children {
		if child.Tick() == SUCCESS {
			return SUCCESS
		}
	}
	return FAILURE
}

type Sequence struct {
	CompositeNode
}

func (s *Sequence) Tick() NodeState {
	for _, child := range s.children {
		if child.Tick() != SUCCESS {
			return FAILURE
		}
	}
	return SUCCESS
}

type LeafNode struct {
	task func() NodeState
}

func (l *LeafNode) Tick() NodeState {
	return l.task()
}
