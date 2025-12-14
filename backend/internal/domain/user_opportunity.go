package domain

import (
	"time"
)

// UserOpportunity 用户-机会关联实体
type UserOpportunity struct {
	ID               int64      `json:"id" db:"id"`
	UserID           int64      `json:"user_id" db:"user_id"`
	OpportunityID    int64      `json:"opportunity_id" db:"opportunity_id"`
	Status           string     `json:"status" db:"status"` // saved/applied/completed/abandoned

	// 双维度评分
	AccessibilityScore float64    `json:"accessibility_score" db:"accessibility_score"` // 匹配度 0-100
	RelevanceScore     float64    `json:"relevance_score" db:"relevance_score"`         // 专业度 0-100
	ScoreDetails       ScoreDetail `json:"score_details" db:"score_details"`             // JSONB评分详情

	// 推送记录
	PushedAt    *time.Time `json:"pushed_at" db:"pushed_at"`
	PushChannel string     `json:"push_channel" db:"push_channel"` // desktop/email/wechat

	// 用户行为
	ViewedAt    *time.Time `json:"viewed_at" db:"viewed_at"`
	SavedAt     *time.Time `json:"saved_at" db:"saved_at"`
	AppliedAt   *time.Time `json:"applied_at" db:"applied_at"`
	CompletedAt *time.Time `json:"completed_at" db:"completed_at"`

	// 用户反馈
	UserFeedback   string `json:"user_feedback" db:"user_feedback"`     // interested/not_interested/irrelevant
	FeedbackReason string `json:"feedback_reason" db:"feedback_reason"`
}

// ScoreDetail 评分详情
type ScoreDetail struct {
	Accessibility AccessibilityDetail `json:"accessibility"`
	Relevance     RelevanceDetail     `json:"relevance"`
}

// AccessibilityDetail 匹配度详情
type AccessibilityDetail struct {
	Eligibility   float64 `json:"eligibility"`    // 硬性门槛 0-1
	SkillsMatch   float64 `json:"skills_match"`   // 技能匹配度 0-1
	TimeCost      float64 `json:"time_cost"`      // 时间冲突成本 0-1
	Total         float64 `json:"total"`          // 总分 0-100
}

// RelevanceDetail 专业度详情
type RelevanceDetail struct {
	MajorMatch        float64 `json:"major_match"`         // 专业匹配度 0-1
	SkillOverlap      float64 `json:"skill_overlap"`       // 技能重叠度 0-1
	CareerAlignment   float64 `json:"career_alignment"`    // 职业发展相关性 0-1
	PeerParticipation float64 `json:"peer_participation"`  // 同专业参与率 0-1
	Total             float64 `json:"total"`               // 总分 0-100
}

// Schedule 日程实体
type Schedule struct {
	ID             int64     `json:"id" db:"id"`
	UserID         int64     `json:"user_id" db:"user_id"`
	Title          string    `json:"title" db:"title"`
	Type           string    `json:"type" db:"type"` // 课程/考试/活动/机会
	StartTime      time.Time `json:"start_time" db:"start_time"`
	EndTime        time.Time `json:"end_time" db:"end_time"`
	Location       string    `json:"location" db:"location"`
	Description    string    `json:"description" db:"description"`
	IsRecurring    bool      `json:"is_recurring" db:"is_recurring"`
	RecurrenceRule string    `json:"recurrence_rule" db:"recurrence_rule"` // WEEKLY_MON_WED
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// CrawlTask 爬虫任务实体
type CrawlTask struct {
	ID             int64      `json:"id" db:"id"`
	TargetURL      string     `json:"target_url" db:"target_url"`
	SiteName       string     `json:"site_name" db:"site_name"`
	SelectorConfig JSONB      `json:"selector_config" db:"selector_config"`
	Frequency      string     `json:"frequency" db:"frequency"` // hourly/daily/weekly
	LastCrawledAt  *time.Time `json:"last_crawled_at" db:"last_crawled_at"`
	NextCrawlAt    *time.Time `json:"next_crawl_at" db:"next_crawl_at"`
	Status         string     `json:"status" db:"status"` // pending/running/success/failed
	ErrorMessage   string     `json:"error_message" db:"error_message"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
}

// CompetitionLevelRule 竞赛级别认定规则
type CompetitionLevelRule struct {
	ID                   int64     `json:"id" db:"id"`
	CompetitionName      string    `json:"competition_name" db:"competition_name"`
	ShortName            string    `json:"short_name" db:"short_name"`
	Level                string    `json:"level" db:"level"` // 国家级A类/B类/省级/国际级
	CertificationSource  string    `json:"certification_source" db:"certification_source"`
	CertificationDocument string   `json:"certification_document" db:"certification_document"`
	Keywords             []string  `json:"keywords" db:"keywords"`
	OrganizerPatterns    []string  `json:"organizer_patterns" db:"organizer_patterns"`
	URLPatterns          []string  `json:"url_patterns" db:"url_patterns"`
	PointsValue          int       `json:"points_value" db:"points_value"`
	DifficultyLevel      int       `json:"difficulty_level" db:"difficulty_level"`
	ParticipationCount   int       `json:"participation_count" db:"participation_count"`
	TargetMajors         []string  `json:"target_majors" db:"target_majors"`
	SkillRequirements    []string  `json:"skill_requirements" db:"skill_requirements"`
	IsActive             bool      `json:"is_active" db:"is_active"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}

// SaveOpportunityRequest 保存机会请求
type SaveOpportunityRequest struct {
	OpportunityID int64 `json:"opportunity_id" binding:"required"`
}

// UpdateOpportunityStatusRequest 更新机会状态请求
type UpdateOpportunityStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=saved applied completed abandoned"`
}

// FeedbackRequest 反馈请求
type FeedbackRequest struct {
	Feedback string `json:"feedback" binding:"required,oneof=interested not_interested irrelevant"`
	Reason   string `json:"reason"`
}
