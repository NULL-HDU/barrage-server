package ball

// roleConf defines base data of the ball.
type roleConf struct {
	id        role
	hp        hp
	damage    damage
	speed     speed
	attackDir attackDir
}

var roleConfTable = map[role]*roleConf{}

var confList = []*roleConf{
	&roleConf{ // 1
		id:        1,
		hp:        100,
		damage:    10,
		speed:     10,
		attackDir: 400,
	},
}

// init read role conf from confList then put them into roleConfTable.
func init() {
	for _, v := range confList {
		roleConfTable[v.id] = v
	}
}
