package main

import "math"

type Position struct {
	X, Y float64
}

type Player struct {
	Position
}

type Enemy struct {
	Position
}

// 感知: 计算距离玩家的距离
func (e *Enemy) DistanceToPlayer(p *Player) float64 {
	dx := p.X - e.X
	dy := p.Y - e.Y
	return math.Sqrt(dx*dx + dy*dy)
}

// 决策/行动: 向玩家移动
func (e *Enemy) MoveTowardsPlayer(p *Player) {
	if e.X < p.X {
		e.X += 1
	}
	if e.X > p.X {
		e.X -= 1
	}
	if e.Y < p.Y {
		e.Y += 1
	}
	if e.Y > p.Y {
		e.Y -= 1
	}
}

// 决策/行动: 攻击玩家
func (e *Enemy) AttackPlayer(p *Player) {
	// 正常情况下这里会有攻击玩家的代码
	// 在此示例中我们只会打印一行文本
	println("The enemy attacks the player!")
}

func main() {
	player := Player{Position{5, 5}}
	enemy := Enemy{Position{0, 0}}

	for {
		distance := enemy.DistanceToPlayer(&player)

		if distance <= 3 {
			enemy.AttackPlayer(&player)
			break
		} else {
			enemy.MoveTowardsPlayer(&player)
			println("The enemy moves towards the player.")
			println("The enemy position:", enemy.X, enemy.Y)
		}
	}
}
