package JoinAudit

import (
	"errors"
	"studentGrow/dao/mysql"
)

type Score struct {
	ID    int
	Score float64 `json:"score"`
}
type scoreMsg struct {
	ID       int
	NowScore int `json:"now_score"`
}

func UpdateScore(cr []Score) (resList []scoreMsg, err error) {
	resList = make([]scoreMsg, 0)
	for _, v := range cr {
		err = mysql.UpdateTrainScoreWithID(v.ID, v.Score)
		if err != nil {
			err = errors.New("成绩修改失败")
			return nil, err
		}
		trainScore, err := mysql.GetTrainScoreWithID(v.ID)
		if err != nil {
			err = errors.New("成绩查询失败")
			return nil, err
		}

		resList = append(resList, scoreMsg{
			ID:       v.ID,
			NowScore: trainScore.TrainScore,
		})
	}
	return
}
