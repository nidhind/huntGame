package db

import "gopkg.in/mgo.v2/bson"

const PuzzleColl = "puzzles"

// puzzle schema for puzzles collection
type Puzzle struct {
	Level        int    `json:"level",string`
	Image        string `json:"image"`
	Clue         string `json:"clue"`
	SolutionHash string `json:"solutionHash"`
}

//puzzle scheme for insert queyr
type InsertPuzzleQuery struct {
	Level        int    `json:"level"`
	Image        string `json:"image"`
	Clue         string `json:"clue"`
	SolutionHash string `json:"solutionHash"`
}

func GetPuzzleByLevel(l int) (Puzzle, error) {
	s := GetSession()
	defer s.Close()
	c := s.DB(DB).C(PuzzleColl)

	var puzzle Puzzle
	err := c.Find(bson.M{"level": l}).One(&puzzle)
	if err != nil {
		return Puzzle{}, err
	}
	return puzzle, nil
}

func InsertNewPuzzle(p *InsertPuzzleQuery) error {
	s := GetSession()
	defer s.Close()
	c := s.DB(DB).C(PuzzleColl)
	err := c.Insert(p)
	if err != nil {
		return err
	}
	return nil
}
