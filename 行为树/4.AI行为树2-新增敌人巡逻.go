/*

新增敌人的巡逻行为

*/

package main

import (
	"fmt"
	"math"
	"time"
)

type Position struct {
	X, Y float64
}

type Player struct {
	Position
}

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

type Enemy struct {
	Position
	path  []Position // 用于存储敌人的路径
	index int        // 用于存储敌人当前的路径索引
	root  Node
}

func (e *Enemy) DistanceToPlayer(p *Player) float64 {
	dx := p.X - e.X
	dy := p.Y - e.Y
	return math.Sqrt(dx*dx + dy*dy)
}

func (e *Enemy) SeePlayer(p *Player) NodeState {
	distance := e.DistanceToPlayer(p)
	if distance < 6 {
		return SUCCESS
	}
	return FAILURE
}

func (e *Enemy) CloseToPlayer(p *Player) NodeState {
	distance := e.DistanceToPlayer(p)
	if distance < 3 {
		return SUCCESS
	}
	return FAILURE
}

func (e *Enemy) AttackPlayer(p *Player) NodeState {
	fmt.Println("The enemy attacks the player!")
	return SUCCESS
}

func (e *Enemy) Wait() NodeState {
	fmt.Println("The enemy waits.")
	return SUCCESS
}

// Patrol: 敌人的巡逻行为，它会在路径上移动
func (e *Enemy) Patrol() NodeState {
	if len(e.path) == 0 {
		return FAILURE // No path defined
	}

	goal := e.path[e.index]
	dx := goal.X - e.X
	dy := goal.Y - e.Y
	distance := math.Sqrt(dx*dx + dy*dy)

	// 如果敌人到达目标点，则移动到下一个目标点
	if distance < 1 {
		e.index = (e.index + 1) % len(e.path)
	} else { // 否则，向目标点移动
		e.X += dx / distance // Normalize the direction and move
		e.Y += dy / distance // Normalize the direction and move
		fmt.Printf("The enemy patrols to %v, %v.\n", e.X, e.Y)
	}

	return SUCCESS
}

func main() {
	player := Player{Position{5, 5}}
	enemy := Enemy{Position{0, 0}, []Position{{0, 0}, {3, 0}, {3, 3}, {0, 3}}, 0, nil}

	/*
		叶子节点: 这些是树中的叶子节点，代表AI的实际行为。在这个例子中，
		有4个不同的LeafNode，分别代表敌人AI的四种行为：是否看到玩家、是否靠近玩家、攻击玩家和等待
	*/
	seePlayer_Node := LeafNode{func() NodeState { return enemy.SeePlayer(&player) }}
	closeToPlayer_Node := LeafNode{func() NodeState { return enemy.CloseToPlayer(&player) }}
	attackPlayer_Node := LeafNode{func() NodeState { return enemy.AttackPlayer(&player) }}
	// wait_Node := LeafNode{func() NodeState { return enemy.Wait() }}
	patrol_Node := LeafNode{func() NodeState { return enemy.Patrol() }} // 新增巡逻行为

	// Sequence：这种节点也会依次尝试它的所有子节点，但只有当所有子节点都成功时，它才会成功
	attack_Sequence := Sequence{CompositeNode{[]Node{&seePlayer_Node, &closeToPlayer_Node, &attackPlayer_Node}}}
	// Selector：这种节点会依次尝试它的所有子节点，只要有一个成功，它就会成功
	root_Selector := Selector{CompositeNode{[]Node{&attack_Sequence, &patrol_Node}}} // 改等待为巡逻行为

	enemy.root = &root_Selector

	for {
		time.Sleep(1 * time.Second)
		enemy.root.Tick()
	}
}
