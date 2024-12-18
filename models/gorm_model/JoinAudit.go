package gorm_model

import (
	"gorm.io/gorm"
)

// JoinAudit 组织部审核表

type JoinAudit struct {
	gorm.Model
	ActivityName            string        `json:"activity_name"` //期数
	Username                string        `json:"username" `
	UserClass               string        `json:"user_class"`
	Name                    string        `json:"name"`
	Gender                  string        `json:"gender"`
	MoralCoin               float64       `json:"moral_coin" gorm:"default:0"`                    //道德币
	ComprehensiveScore      float64       `json:"comprehensive_score" gorm:"default:0"`           //综测成绩
	TrainScore              int           `json:"training_score" gorm:"default:0"`                //培训测试成绩
	ClassIsPass             string        `json:"class_is_pass" gorm:"default:null"`              //班级审核
	RulerIsPass             string        `json:"ruler_is_pass" gorm:"default:null"`              //纪权部综测成绩审核结果
	OrganizerMaterialIsPass string        `json:"organizer_material_is_pass" gorm:"default:null"` //组织部材料审核结果
	OrganizerTrainIsPass    string        `json:"organizer_train_is_pass" gorm:"default:null"`    //组织部培训审核结果
	JoinAuditDutyID         uint          `json:"join_audit_duty_id"`
	JoinAuditDuty           JoinAuditDuty `gorm:"foreignKey:JoinAuditDutyID"`
	Note                    string        `json:"join_audit_note"` //备注
}
