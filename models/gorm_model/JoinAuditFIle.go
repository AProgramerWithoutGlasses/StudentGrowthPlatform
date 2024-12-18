package gorm_model

import (
	"gorm.io/gorm"
)

type JoinAuditFile struct {
	gorm.Model
	Username        string        `json:"username"` //照片归属
	FileName        string        `json:"name"`
	FileHash        string        `json:"hash"` //照片hash值
	FilePath        string        `json:"path"` //路径
	JoinAuditDutyID uint          `json:"join_audit_duty_id"`
	JoinAuditDuty   JoinAuditDuty `gorm:"foreignKey:JoinAuditDutyID"`
	//JoinAuditID     uint          `json:"join_audit_id"`
	//JoinAudit       JoinAudit     `gorm:"foreignKey:JoinAuditID"`
	Note string `json:"image_note"` //备注
}

//func (j JoinAuditFile) BeforeDelete(tx *gorm.DB) (err error) {
//	err, getHeader := fileProcess.DelOssFile(j.FilePath)
//	if err != nil {
//		fmt.Println(err)
//		return err
//	}
//	fmt.Println(getHeader)
//	fmt.Println("删除钩子使用成功")
//	return nil
//}
