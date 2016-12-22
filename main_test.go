package main

import (
	"testing"
)

func TestSetChecker(t *testing.T) {
	t.Log("Testing Set checker")
	ft := Field{}
	ft.SetChecker(0)
	if ft.checkerColor == 1 {
		t.Errorf("Expected white cell. Found black")
	}
}

func TestBeat(t *testing.T) {
	t.Log("Testing Beat")
	b := Board{}
	b.CreateField()
	b.bord[1][2].SetChecker(0)
	b.bord[2][3].SetChecker(1)
	b.Beat(1, 2, 3, 4)
	if b.bord[2][3].hasChecker == true {
		t.Errorf("Expected empty cell. Beat error")
	}
	if b.bord[3][4].hasChecker == false {
		t.Error("Expected checker. Beat error")
	}
}

func TestChain(t *testing.T) {
	t.Log("Testing chain of moves")
	b := Board{}
	b.CreateField()
	b.bord[1][2].SetChecker(0)
	b.bord[2][3].SetChecker(1)
	b.bord[4][5].SetChecker(1)
	if b.CheckBeat(1, 2) == false {
		t.Errorf("Error. Beat Expected")
	}
	b.Beat(1, 2, 3, 4)
	if b.CheckBeat(3, 4) == false {
		t.Errorf("Error. Beat Expeted")
	}
	b.Beat(3, 4, 5, 6)
}
