package ball

import (
	"fmt"
	"strings"
	"testing"
)

func generateBall() Ball {

	return &ball{
		camp:      1234,
		id:        0,
		bType:     AirPlane,
		special:   99,
		state:     Alive,
		role:      1,
		hp:        100,
		damage:    10,
		speed:     10,
		attackDir: 400,
		location:  location{99.9, 99.9},
	}
}

func compare(std Ball, b Ball) error {
	dball := std.(*ball)
	nball := b.(*ball)

	if dball.camp != nball.camp {
		return fmt.Errorf("Hope %v, get %v.", dball.camp, nball.camp)
	}
	if dball.location != nball.location {
		return fmt.Errorf("Hope %v, get %v.", dball.location, nball.location)
	}
	if dball.attackDir != nball.attackDir {
		return fmt.Errorf("Hope %v, get %v.", dball.attackDir, nball.attackDir)
	}
	if dball.speed != nball.speed {
		return fmt.Errorf("Hope %v, get %v.", dball.speed, nball.speed)
	}
	if dball.damage != nball.damage {
		return fmt.Errorf("Hope %v, get %v.", dball.damage, nball.damage)
	}
	if dball.hp != nball.hp {
		return fmt.Errorf("Hope %v, get %v.", dball.hp, nball.hp)
	}
	if dball.role != nball.role {
		return fmt.Errorf("Hope %v, get %v.", dball.role, nball.role)
	}
	if dball.state != nball.state {
		return fmt.Errorf("Hope %v, get %v.", dball.state, nball.state)
	}
	if dball.special != nball.special {
		return fmt.Errorf("Hope %v, get %v.", dball.special, nball.special)
	}
	if dball.bType != nball.bType {
		return fmt.Errorf("Hope %v, get %v.", dball.bType, nball.bType)
	}
	if dball.id != nball.id {
		return fmt.Errorf("Hope %v, get %v.", dball.id, nball.id)
	}

	return nil
}

func TestNewUserAirplane(t *testing.T) {
	defaultBall := generateBall()
	newBall, _ := NewUserAirplane(1234, 1, 99, 99.9, 99.9)

	if err := compare(defaultBall, newBall); err != nil {
		t.Error(err)
	}

	b, _ := newBall.MarshalBinary()
	t.Logf("MarshalBinary result: % x", b)

	_, err := NewUserAirplane(1234, 255, 99, 99.9, 99.9)
	if err != ErrInvalidRole {
		t.Errorf("Hope get '%v', but get '%v'.", ErrInvalidRole, err)
	}
}

func TestNewBallFromBytes(t *testing.T) {
	defaultBall := generateBall()
	b, _ := defaultBall.MarshalBinary()
	newBall, _ := NewBallFromBytes(b)

	if err := compare(defaultBall, newBall); err != nil {
		t.Error(err)
	}

	defer func() {
		t.Log("Running panic testing.")
		if err := recover(); !strings.Contains(err.(error).Error(), "index out of range") {
			t.Errorf("Hope get panic with error 'index out of range', but get '%v'.", err.(error))
		}
	}()
	NewBallFromBytes(b[30:])
}
