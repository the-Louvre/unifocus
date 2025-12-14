package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Opportunity 机会实体
type Opportunity struct {
	ID          int64  `json:"id" db:"id"`
	Title       string `json:"title" db:"title"`
	Type        string `json:"type" db:"type"` // 竞赛/实习/项目/奖学金
	Description string `json:"description" db:"description"`
	SourceURL   string `json:"source_url" db:"source_url"`
	SourceType  string `json:"source_type" db:"source_type"` // 官网/公众号/端侧

	// 竞赛级别认定字段
	CompetitionLevel  string `json:"competition_level" db:"competition_level"`   // 国家级A类/B类/省级/校级
	CertificationType string `json:"certification_type" db:"certification_type"` // 教育部认定/省教育厅认定
	Organizer         string `json:"organizer" db:"organizer"`                   // 主办方
	OrganizerType     string `json:"organizer_type" db:"organizer_type"`         // 政府/高校/企业/协会
	AwardLevel        string `json:"award_level" db:"award_level"`               // 奖项级别
	PointsValue       int    `json:"points_value" db:"points_value"`             // 学分/加分值
	IsOfficial        bool   `json:"is_official" db:"is_official"`               // 是否官方认定

	// 结构化字段
	StartDate *time.Time `json:"start_date" db:"start_date"`
	Deadline  *time.Time `json:"deadline" db:"deadline"`
	EventDate *time.Time `json:"event_date" db:"event_date"`
	Location  string     `json:"location" db:"location"`

	// 要求字段
	Requirements     Requirements `json:"requirements" db:"requirements"`           // JSONB
	EligibilityRules JSONB        `json:"eligibility_rules" db:"eligibility_rules"` // JSONB
	TargetMajors     []string     `json:"target_majors" db:"target_majors"`         // 数组

	// 元数据
	Tags              []string     `json:"tags" db:"tags"`               // 数组
	Attachments       []Attachment `json:"attachments" db:"attachments"` // JSONB
	DescriptionVector []float32    `json:"-" db:"description_vector"`    // 向量
	IsActive          bool         `json:"is_active" db:"is_active"`
	ViewCount         int          `json:"view_count" db:"view_count"`
	SaveCount         int          `json:"save_count" db:"save_count"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Requirements 要求信息
type Requirements struct {
	Grade        []int    `json:"grade,omitempty"`        // 年级要求
	Major        []string `json:"major,omitempty"`        // 专业要求
	Skills       []string `json:"skills,omitempty"`       // 技能要求
	Certificates []string `json:"certificates,omitempty"` // 证书要求
}

// Attachment 附件信息
type Attachment struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Type string `json:"type"` // pdf/image/doc
}

// JSONB 自定义类型用于处理PostgreSQL的JSONB字段
// 实现 database/sql 的 driver.Valuer 和 sql.Scanner 接口
// 使得Go的map类型可以直接映射到PostgreSQL的JSONB列
// 这样可以在Go代码中使用 map[string]interface{} 类型，而数据库自动处理JSON序列化/反序列化
type JSONB map[string]interface{}

// Value 实现 driver.Valuer 接口
// 将Go的map序列化为JSON字节流，供PostgreSQL存储
// 返回值: JSON字节流或nil（当map为空时）
// 错误场景: JSON序列化失败（通常不会发生，除非包含不可序列化的类型）
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan 实现 sql.Scanner 接口
// 从PostgreSQL读取JSONB数据并反序列化为Go的map
// 参数: value可以是[]byte（JSON字节流）或nil
// 错误场景: 类型断言失败或JSON反序列化失败
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal JSONB value: expected []byte")
	}

	return json.Unmarshal(bytes, j)
}

// Value 实现 Requirements 的 driver.Valuer 接口
// 将Requirements结构体序列化为JSON字节流，供PostgreSQL JSONB列存储
func (r Requirements) Value() (driver.Value, error) {
	return json.Marshal(r)
}

// Scan 实现 Requirements 的 sql.Scanner 接口
// 从PostgreSQL JSONB列读取数据并反序列化为Requirements结构体
func (r *Requirements) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal Requirements value: expected []byte")
	}

	return json.Unmarshal(bytes, r)
}

// CreateOpportunityRequest 创建机会请求
type CreateOpportunityRequest struct {
	Title        string       `json:"title" binding:"required"`
	Type         string       `json:"type" binding:"required"`
	Description  string       `json:"description" binding:"required"`
	SourceURL    string       `json:"source_url" binding:"required,url"`
	Organizer    string       `json:"organizer"`
	StartDate    *time.Time   `json:"start_date"`
	Deadline     *time.Time   `json:"deadline"`
	EventDate    *time.Time   `json:"event_date"`
	Location     string       `json:"location"`
	Requirements Requirements `json:"requirements"`
	TargetMajors []string     `json:"target_majors"`
	Tags         []string     `json:"tags"`
}

// OpportunityFilter 机会筛选条件
type OpportunityFilter struct {
	Type             string     `form:"type"`
	CompetitionLevel string     `form:"competition_level"`
	Major            string     `form:"major"`
	DeadlineAfter    *time.Time `form:"deadline_after"`
	DeadlineBefore   *time.Time `form:"deadline_before"`
	IsActive         *bool      `form:"is_active"`
	Tags             []string   `form:"tags"`
	Limit            int        `form:"limit"`
	Offset           int        `form:"offset"`
}
