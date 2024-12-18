package JoinAudit

import "studentGrow/dao/mysql"

type RecList struct {
	True  []int `form:"true"`
	False []int `form:"false"`
	Null  []int `json:"null"`
}
type ResListWithIsPass struct {
	ID        int
	NowStatus string
}

func IsPassWithJSON(cr RecList, passType string) (resList []ResListWithIsPass) {
	resList = make([]ResListWithIsPass, 0)
	if len(cr.True) != 0 {
		for _, id := range cr.True {
			var resMsg ResListWithIsPass
			resMsg.ID = id
			updatedJoinAudit := mysql.IsPass(id, passType, "true")
			switch passType {
			case "class_is_pass":
				resMsg.NowStatus = updatedJoinAudit.ClassIsPass
			case "ruler_is_pass":
				resMsg.NowStatus = updatedJoinAudit.RulerIsPass
			case "organizer_material_is_pass":
				resMsg.NowStatus = updatedJoinAudit.OrganizerMaterialIsPass
			case "organizer_train_is_pass":
				resMsg.NowStatus = updatedJoinAudit.OrganizerTrainIsPass
			}
			resList = append(resList, resMsg)
		}
	}
	if len(cr.False) != 0 {
		for _, id := range cr.False {
			var resMsg ResListWithIsPass
			resMsg.ID = id
			updatedJoinAudit := mysql.IsPass(id, passType, "false")
			switch passType {
			case "class_is_pass":
				resMsg.NowStatus = updatedJoinAudit.ClassIsPass
			case "ruler_is_pass":
				resMsg.NowStatus = updatedJoinAudit.RulerIsPass
			case "organizer_material_is_pass":
				resMsg.NowStatus = updatedJoinAudit.OrganizerMaterialIsPass
			case "organizer_train_is_pass":
				resMsg.NowStatus = updatedJoinAudit.OrganizerTrainIsPass
			}
			resList = append(resList, resMsg)
		}
	}
	if len(cr.Null) != 0 {
		for _, id := range cr.Null {
			var resMsg ResListWithIsPass
			resMsg.ID = id
			updatedJoinAudit := mysql.IsPass(id, passType, "null")
			switch passType {
			case "class_is_pass":
				resMsg.NowStatus = updatedJoinAudit.ClassIsPass
			case "ruler_is_pass":
				resMsg.NowStatus = updatedJoinAudit.RulerIsPass
			case "organizer_material_is_pass":
				resMsg.NowStatus = updatedJoinAudit.OrganizerMaterialIsPass
			case "organizer_train_is_pass":
				resMsg.NowStatus = updatedJoinAudit.OrganizerTrainIsPass
			}
			resList = append(resList, resMsg)
		}
	}
	return
}
