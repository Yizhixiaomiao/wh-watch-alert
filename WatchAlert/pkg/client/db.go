package client

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"watchAlert/internal/global"
	"watchAlert/internal/models"
)

type DBConfig struct {
	Host    string
	Port    string
	User    string
	Pass    string
	DBName  string
	Timeout string
}

func NewDBClient(config DBConfig) *gorm.DB {
	// 初始化本地 test.db 数据库文件
	//db, err := gorm.Open(sqlite.Open("data/sql.db"), &gorm.Config{})

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s",
		config.User,
		config.Pass,
		config.Host,
		config.Port,
		config.DBName,
		config.Timeout)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
	})

	if err != nil {
		logc.Errorf(context.Background(), "failed to connect database: %s", err.Error())
		return nil
	}

	// 设置数据库连接字符集
	sqlDB, err := db.DB()
	if err != nil {
		logc.Errorf(context.Background(), "failed to get database instance: %s", err.Error())
		return nil
	}

	// 设置连接参数
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	// 执行设置字符集的SQL
	db.Exec("SET NAMES utf8mb4")
	db.Exec("SET CHARACTER SET utf8mb4")
	db.Exec("SET character_set_connection=utf8mb4")

	// 修改表字符集为 utf8mb4
	db.Exec("ALTER TABLE ticket CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci")
	db.Exec("ALTER TABLE ticket MODIFY COLUMN title VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci")
	db.Exec("ALTER TABLE ticket MODIFY COLUMN description TEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci")

	// 检查 Product 结构是否变化，变化则进行迁移
	err = db.AutoMigrate(
		&models.DutySchedule{},
		&models.DutyManagement{},
		&models.AlertNotice{},
		&models.AlertDataSource{},
		&models.AlertRule{},
		&models.AlertCurEvent{},
		&models.AlertHisEvent{},
		&models.AlertSilences{},
		&models.Member{},
		&models.UserRole{},
		&models.UserPermissions{},
		&models.NoticeTemplateExample{},
		&models.RuleGroups{},
		&models.RuleTemplateGroup{},
		&models.RuleTemplate{},
		&models.Tenant{},
		&models.Dashboard{},
		&models.AuditLog{},
		&models.Settings{},
		&models.TenantLinkedUsers{},
		&models.DashboardFolders{},
		&models.AlertSubscribe{},
		&models.NoticeRecord{},
		&models.ProbingRule{},
		&models.FaultCenter{},
		&models.AiContentRecord{},
		&models.ProbingHistory{},
		&models.Comment{},
		&models.AlertTicketRule{},
		&models.AlertTicketRuleHistory{},
		&models.Ticket{},
		&models.TicketWorkLog{},
		&models.TicketComment{},
		&models.TicketAttachment{},
		&models.TicketTemplate{},
		&models.TicketSLAPolicy{},
		&models.TicketStep{},
		&models.TicketReview{},
		&models.TicketReviewer{},
		&models.WorkHoursStandard{},
		&models.Knowledge{},
		&models.KnowledgeLike{},
		&models.KnowledgeCategory{},
		&models.AssignmentRule{},
		&models.WechatRepairRequest{},
	)
	if err != nil {
		logc.Error(context.Background(), err.Error())
		return nil
	}

	if global.Config.Server.Mode == "debug" {
		db.Debug()
	} else {
		db.Logger = logger.Default.LogMode(logger.Silent)
	}

	return db
}
