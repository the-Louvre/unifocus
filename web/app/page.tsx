'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { Button, Card, Typography, Space } from 'antd';
import { DashboardOutlined, LoginOutlined } from '@ant-design/icons';

const { Title, Paragraph } = Typography;

export default function HomePage() {
  const router = useRouter();

  return (
    <div style={{ 
      minHeight: '100vh', 
      display: 'flex', 
      alignItems: 'center', 
      justifyContent: 'center',
      background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)'
    }}>
      <Card style={{ width: 600, textAlign: 'center' }}>
        <Title level={1}>UniFocus</Title>
        <Title level={3}>高校机会助手</Title>
        <Paragraph>
          智能匹配竞赛、奖学金等高校机会，让您不错过任何机会
        </Paragraph>
        <Space size="large">
          <Button 
            type="primary" 
            size="large" 
            icon={<DashboardOutlined />}
            onClick={() => router.push('/dashboard')}
          >
            监控面板
          </Button>
          <Button 
            size="large" 
            icon={<LoginOutlined />}
            onClick={() => router.push('/login')}
          >
            登录
          </Button>
        </Space>
      </Card>
    </div>
  );
}


