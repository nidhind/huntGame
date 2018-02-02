package db

import "gopkg.in/mgo.v2/bson"

const StatisticsColl = "statistics"

// User schema for statistics collection
type Statistics struct {
	Id               string    `bson:"id"`
	Count                int    `bson:"count"`
}

type InsertStatisticsQuery struct {
	Id               string    `bson:"id"`
	Count       	 int    `bson:"count"`
}

func GetCount(id string) (int,error) {
	s := GetSession()
	defer s.Close()
	c := s.DB(DB).C(StatisticsColl)

	var statistics Statistics
	err := c.Find(bson.M{"id": id}).One(&statistics)
	if err != nil {
		return 0, err
	}
	return statistics.Count, nil

} 

func UpdateCount(id string) error {
	s := GetSession()
	defer s.Close()
	c := s.DB(DB).C(StatisticsColl)
	q := bson.M{"id": id}
	u := bson.M{"$inc": bson.M{"count": 1}}
	err := c.Update(&q, &u)
	if err != nil {
		return err
	}
	return nil
}

func InsertStatistics(stat *InsertStatisticsQuery) error {
	s := GetSession()
	defer s.Close()
	c := s.DB(DB).C(StatisticsColl)
	err := c.Insert(stat)
	if err != nil {
		return err
	}
	return nil
}