import { Card, Col, Row, Typography } from 'antd'
import { useEffect, useState } from 'react'
import { fetchDepartments } from '../api/catalog'

export default function Home() {
  const [depts, setDepts] = useState<any[]>([])
  useEffect(() => { fetchDepartments().then(setDepts).catch(() => setDepts([])) }, [])
  return (
    <div>
      <Typography.Title level={2}>Добре дошли в магазина за мебели</Typography.Title>
      <Typography.Paragraph>Разгледайте отделите и открийте подходящите мебели.</Typography.Paragraph>
      <Row gutter={[16,16]}>
        {depts.map(d => (
          <Col key={d.id} xs={24} sm={12} md={8}>
            <Card title={d.name}>
              <p>{d.description}</p>
            </Card>
          </Col>
        ))}
      </Row>
    </div>
  )
}

