-- 001_init_schema.up.sql
-- UniFocus 数据库初始化Schema

-- 启用必要的扩展
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm"; -- 用于模糊搜索

-- ============================================
-- 1. 用户表
-- ============================================
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    school VARCHAR(100),
    major VARCHAR(100),
    grade SMALLINT CHECK (grade >= 1 AND grade <= 4),
    avatar_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);

-- ============================================
-- 2. 用户画像表
-- ============================================
CREATE TABLE user_profiles (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    resume_text TEXT,
    skills JSONB DEFAULT '[]'::jsonb,
    certificates JSONB DEFAULT '[]'::jsonb,
    interests JSONB DEFAULT '[]'::jsonb,
    resume_vector REAL[], -- 向量表示（可用pgvector扩展优化）
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_user_profiles_user_id ON user_profiles(user_id);

-- ============================================
-- 3. 竞赛级别认定规则表
-- ============================================
CREATE TABLE competition_level_rules (
    id BIGSERIAL PRIMARY KEY,
    competition_name VARCHAR(255) NOT NULL,
    short_name VARCHAR(100),
    level VARCHAR(50) NOT NULL, -- 国家级A类/国家级B类/省级/校级/国际级

    -- 认定信息
    certification_source VARCHAR(100),
    certification_document TEXT,

    -- 识别规则
    keywords TEXT[] DEFAULT ARRAY[]::TEXT[],
    organizer_patterns TEXT[] DEFAULT ARRAY[]::TEXT[],
    url_patterns TEXT[] DEFAULT ARRAY[]::TEXT[],

    -- 价值评估
    points_value INT DEFAULT 0,
    difficulty_level INT CHECK (difficulty_level >= 1 AND difficulty_level <= 10),
    participation_count INT DEFAULT 0,

    -- 专业相关性
    target_majors TEXT[] DEFAULT ARRAY[]::TEXT[],
    skill_requirements TEXT[] DEFAULT ARRAY[]::TEXT[],

    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_competition_rules_level ON competition_level_rules(level);
CREATE INDEX idx_competition_rules_keywords ON competition_level_rules USING GIN(keywords);
CREATE INDEX idx_competition_rules_name ON competition_level_rules(competition_name);

-- 插入示例数据
INSERT INTO competition_level_rules
(competition_name, short_name, level, certification_source, keywords, organizer_patterns, target_majors, points_value, difficulty_level)
VALUES
('全国大学生数学建模竞赛', '数学建模', '国家级A类', '教育部认定',
 ARRAY['数学建模', 'Mathematical Contest in Modeling', 'MCM'],
 ARRAY['教育部%', '%高等教育学会%'],
 ARRAY['数学', '统计', '计算机', '工科'],
 10, 8),
('中国国际"互联网+"大学生创新创业大赛', '互联网+', '国家级A类', '教育部主办',
 ARRAY['互联网+', '创新创业'],
 ARRAY['%教育部%'],
 ARRAY['全部专业'],
 15, 9),
('全国大学生电子设计竞赛', '电子设计', '国家级A类', '教育部认定',
 ARRAY['电子设计', 'NUEDC'],
 ARRAY['%教育部%', '%工业和信息化部%'],
 ARRAY['电子', '自动化', '通信', '计算机'],
 10, 8);

-- ============================================
-- 4. 机会表
-- ============================================
CREATE TABLE opportunities (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL, -- 竞赛/实习/项目/奖学金
    description TEXT,
    source_url TEXT NOT NULL,
    source_type VARCHAR(50), -- 官网/公众号/端侧

    -- 竞赛级别认定字段
    competition_level VARCHAR(50),
    certification_type VARCHAR(100),
    organizer VARCHAR(200),
    organizer_type VARCHAR(50),
    award_level VARCHAR(100),
    points_value INT DEFAULT 0,
    is_official BOOLEAN DEFAULT false,

    -- 结构化字段
    start_date DATE,
    deadline DATE,
    event_date DATE,
    location VARCHAR(200),

    -- 要求字段
    requirements JSONB DEFAULT '{}'::jsonb,
    eligibility_rules JSONB DEFAULT '{}'::jsonb,
    target_majors TEXT[] DEFAULT ARRAY[]::TEXT[],

    -- 元数据
    tags TEXT[] DEFAULT ARRAY[]::TEXT[],
    attachments JSONB DEFAULT '[]'::jsonb,
    description_vector REAL[],
    is_active BOOLEAN DEFAULT true,
    view_count INT DEFAULT 0,
    save_count INT DEFAULT 0,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_opportunities_deadline ON opportunities(deadline);
CREATE INDEX idx_opportunities_type ON opportunities(type);
CREATE INDEX idx_opportunities_level ON opportunities(competition_level);
CREATE INDEX idx_opportunities_tags ON opportunities USING GIN(tags);
CREATE INDEX idx_opportunities_majors ON opportunities USING GIN(target_majors);
CREATE INDEX idx_opportunities_is_active ON opportunities(is_active);
CREATE INDEX idx_opportunities_created_at ON opportunities(created_at DESC);

-- ============================================
-- 5. 用户-机会关联表
-- ============================================
CREATE TABLE user_opportunities (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    opportunity_id BIGINT REFERENCES opportunities(id) ON DELETE CASCADE,

    status VARCHAR(20) DEFAULT 'saved', -- saved/applied/completed/abandoned

    -- 双维度评分
    accessibility_score REAL CHECK (accessibility_score >= 0 AND accessibility_score <= 100),
    relevance_score REAL CHECK (relevance_score >= 0 AND relevance_score <= 100),
    score_details JSONB DEFAULT '{}'::jsonb,

    -- 推送记录
    pushed_at TIMESTAMP,
    push_channel VARCHAR(50),

    -- 用户行为
    viewed_at TIMESTAMP,
    saved_at TIMESTAMP,
    applied_at TIMESTAMP,
    completed_at TIMESTAMP,

    -- 用户反馈
    user_feedback VARCHAR(20),
    feedback_reason TEXT,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(user_id, opportunity_id)
);

CREATE INDEX idx_user_opportunities_user ON user_opportunities(user_id);
CREATE INDEX idx_user_opportunities_opportunity ON user_opportunities(opportunity_id);
CREATE INDEX idx_user_opportunities_status ON user_opportunities(status);
CREATE INDEX idx_user_opportunities_scores ON user_opportunities(accessibility_score, relevance_score);

-- ============================================
-- 6. 日程表
-- ============================================
CREATE TABLE schedules (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(200) NOT NULL,
    type VARCHAR(50), -- 课程/考试/活动/机会
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    location VARCHAR(200),
    description TEXT,
    is_recurring BOOLEAN DEFAULT false,
    recurrence_rule VARCHAR(100),

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_schedules_user_time ON schedules(user_id, start_time, end_time);

-- ============================================
-- 7. 爬虫任务表
-- ============================================
CREATE TABLE crawl_tasks (
    id BIGSERIAL PRIMARY KEY,
    target_url TEXT NOT NULL,
    site_name VARCHAR(100),
    selector_config JSONB DEFAULT '{}'::jsonb,
    frequency VARCHAR(50) DEFAULT 'daily', -- hourly/daily/weekly
    last_crawled_at TIMESTAMP,
    next_crawl_at TIMESTAMP,
    status VARCHAR(20) DEFAULT 'pending', -- pending/running/success/failed
    error_message TEXT,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_crawl_tasks_next_crawl ON crawl_tasks(next_crawl_at);
CREATE INDEX idx_crawl_tasks_status ON crawl_tasks(status);

-- ============================================
-- 8. NLP任务队列表
-- ============================================
CREATE TABLE nlp_tasks (
    id BIGSERIAL PRIMARY KEY,
    task_type VARCHAR(50) NOT NULL, -- extract/classify/vectorize
    input_data JSONB NOT NULL,
    output_data JSONB,
    status VARCHAR(20) DEFAULT 'pending', -- pending/processing/completed/failed
    retry_count INT DEFAULT 0,
    error_message TEXT,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP
);

CREATE INDEX idx_nlp_tasks_status ON nlp_tasks(status);
CREATE INDEX idx_nlp_tasks_created_at ON nlp_tasks(created_at);

-- ============================================
-- 触发器：自动更新 updated_at
-- ============================================
-- 功能说明:
-- 此触发器函数在每次UPDATE操作时自动更新 updated_at 字段为当前时间戳
-- 触发时机: BEFORE UPDATE（在更新操作执行之前）
-- 作用范围: 以下5个表都使用此触发器
--   - users
--   - user_profiles
--   - opportunities
--   - user_opportunities
--   - competition_level_rules
-- 性能影响:
--   - 对单行更新影响极小（微秒级）
--   - 批量更新时，每行都会触发，可能略微影响性能
--   - 如果需要进行大量批量更新，可以考虑临时禁用触发器
--   示例: ALTER TABLE users DISABLE TRIGGER update_users_updated_at;
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_profiles_updated_at BEFORE UPDATE ON user_profiles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_opportunities_updated_at BEFORE UPDATE ON opportunities
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_opportunities_updated_at BEFORE UPDATE ON user_opportunities
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_competition_rules_updated_at BEFORE UPDATE ON competition_level_rules
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================
-- 完成
-- ============================================
COMMENT ON DATABASE postgres IS 'UniFocus - 高校机会助手数据库';
