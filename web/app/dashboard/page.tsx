'use client';

/**
 * 项目监控Dashboard
 * 实时显示项目进度和系统状态
 */
import { useEffect, useState } from 'react';
import { Card, Statistic, Row, Col, Table, Tag, Spin } from 'antd';
import { 
  CheckCircleOutlined, 
  ClockCircleOutlined, 
  DatabaseOutlined,
  ApiOutlined 
} from '@ant-design/icons';
import { metricsAPI, SystemMetrics } from '@/lib/api/metrics';
import { opportunitiesAPI, Opportunity } from '@/lib/api/opportunities';

export default function DashboardPage() {
  const [metrics, setMetrics] = useState<SystemMetrics | null>(null);
  const [opportunities, setOpportunities] = useState<Opportunity[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadData();
    // 每30秒刷新一次
    const interval = setInterval(loadData, 30000);
    return () => clearInterval(interval);
  }, []);

  const loadData = async () => {
    try {
      const [metricsData, oppsData] = await Promise.all([
        metricsAPI.getMetrics(),
        opportunitiesAPI.list({ limit: 10 }),
      ]);
      setMetrics(metricsData);
      setOpportunities(oppsData.data);
    } catch (error) {
      console.error('Failed to load dashboard data:', error);
    } finally {
      setLoading(false);
    }
  };

  const columns = [
    {
      title: '标题',
      dataIndex: 'title',
      key: 'title',
      ellipsis: true,
    },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
      render: (type: string) => <Tag>{type}</Tag>,
    },
    {
      title: '级别',
      dataIndex: 'competition_level',
      key: 'competition_level',
      render: (level?: string) => level ? <Tag color="blue">{level}</Tag> : '-',
    },
    {
      title: '浏览量',
      dataIndex: 'view_count',
      key: 'view_count',
    },
    {
      title: '保存量',
      dataIndex: 'save_count',
      key: 'save_count',
    },
  ];

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: '50px' }}>
        <Spin size="large" />
      </div>
    );
  }

  return (
    <div style={{ padding: '24px' }}>
      <h1>项目监控 Dashboard</h1>

      {/* 系统状态卡片 */}
      <Row gutter={16} style={{ marginBottom: '24px' }}>
        <Col span={6}>
          <Card>
            <Statistic
              title="系统状态"
              value={metrics?.status || 'unknown'}
              prefix={<CheckCircleOutlined style={{ color: '#3f8600' }} />}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="运行时间"
              value={metrics?.system?.uptime || '0s'}
              prefix={<ClockCircleOutlined />}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="机会总数"
              value={opportunities.length}
              prefix={<DatabaseOutlined />}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="API状态"
              value="正常"
              prefix={<ApiOutlined style={{ color: '#1890ff' }} />}
            />
          </Card>
        </Col>
      </Row>

      {/* 最近机会列表 */}
      <Card title="最近机会" style={{ marginBottom: '24px' }}>
        <Table
          columns={columns}
          dataSource={opportunities}
          rowKey="id"
          pagination={false}
          size="small"
        />
      </Card>

      {/* 项目进度 */}
      <Card title="开发进度">
        <Row gutter={16}>
          <Col span={12}>
            <h3>已完成功能 ✅</h3>
            <ul>
              <li>数据库连接池和Redis客户端</li>
              <li>用户认证系统（注册/登录/JWT）</li>
              <li>机会管理服务（CRUD + 查询）</li>
              <li>基础爬虫框架</li>
              <li>NLP文本提取服务</li>
              <li>用户画像服务</li>
            </ul>
          </Col>
          <Col span={12}>
            <h3>待开发功能 ⏳</h3>
            <ul>
              <li>双维度评分引擎</li>
              <li>竞赛级别识别服务</li>
              <li>推送服务和定时匹配任务</li>
              <li>API文档生成（Swagger）</li>
            </ul>
          </Col>
        </Row>
      </Card>
    </div>
  );
}


