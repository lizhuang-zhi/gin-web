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

type State int

const (
	Wait State = iota
	Chase
	Attack
)

type Enemy struct {
	Position
	State // 和1.硬编码.go相比, 多了一个State字段
}

func (e *Enemy) DistanceToPlayer(p *Player) float64 {
	dx := p.X - e.X
	dy := p.Y - e.Y
	return math.Sqrt(dx*dx + dy*dy)
}

func (e *Enemy) MoveTowardsPlayer(p *Player) {
	if e.X < p.X {
		e.X++
	}
	if e.X > p.X {
		e.X--
	}
	if e.Y < p.Y {
		e.Y++
	}
	if e.Y > p.Y {
		e.Y--
	}
}

func (e *Enemy) AttackPlayer(p *Player) {
	fmt.Println("The enemy attacks the player!")
}

func main() {
	player := Player{Position{5, 5}}
	enemy := Enemy{Position{0, 0}, Wait}

	for {
		distance := enemy.DistanceToPlayer(&player)

		switch enemy.State {
		case Wait: // 等待状态
			if distance < 6 {
				enemy.State = Chase
			}
		case Chase: // 追逐状态
			if distance <= 3 {
				enemy.State = Attack
			} else if distance > 6 {
				enemy.State = Wait
			} else {
				enemy.MoveTowardsPlayer(&player)
				fmt.Println("The enemy moves towards the player.")
			}
		case Attack: // 攻击状态
			if distance > 3 {
				enemy.State = Chase
			} else {
				enemy.AttackPlayer(&player)
			}
		}
	}
}

/*
	在这个例子中，每次循环时，AI都会根据其当前状态进行相应的行为，并依据特定的规则转换状态。例如，当敌人AI在Chase状态下，并且玩家在3个单位以内时，AI会转换到Attack状态。
	这就是有限状态机的基础概念。用有限状态机管理AI行为是清晰且有条理的，特别是当你有更多的状态和更复杂的规则时。然而，有限状态机也有其局限性，
	比如难以处理过于复杂的状态转换规则，或者处理多个同时发生的状态。在这些情况下，可能需要其他的AI设计模式，如行为树或者实用系统。
*/
