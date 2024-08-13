/*

ROOT
|
+- Selector
|
+- Sequence
|  |
|  +- See player?
|  |
|  +- Too close?
|  |
|  +- Attack player
|
+- Wait

*/

package main

import (
	"fmt"
	"math"
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
	root Node
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

func main() {
	player := Player{Position{5, 5}}
	enemy := Enemy{Position{0, 0}, nil}

	/*
		叶子节点: 这些是树中的叶子节点，代表AI的实际行为。在这个例子中，
		有4个不同的LeafNode，分别代表敌人AI的四种行为：是否看到玩家、是否靠近玩家、攻击玩家和等待
	*/
	seePlayer_Node := LeafNode{func() NodeState { return enemy.SeePlayer(&player) }}
	closeToPlayer_Node := LeafNode{func() NodeState { return enemy.CloseToPlayer(&player) }}
	attackPlayer_Node := LeafNode{func() NodeState { return enemy.AttackPlayer(&player) }}
	wait_Node := LeafNode{func() NodeState { return enemy.Wait() }}

	// Sequence：这种节点也会依次尝试它的所有子节点，但只有当所有子节点都成功时，它才会成功
	attack_Sequence := Sequence{CompositeNode{[]Node{&seePlayer_Node, &closeToPlayer_Node, &attackPlayer_Node}}}
	// Selector：这种节点会依次尝试它的所有子节点，只要有一个成功，它就会成功
	root_Selector := Selector{CompositeNode{[]Node{&attack_Sequence, &wait_Node}}}

	enemy.root = &root_Selector

	for {
		enemy.root.Tick()
	}
}
