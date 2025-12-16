import { Table, Tag, Typography, message } from 'antd'
import { useEffect, useState } from 'react'
import { myOrders } from '../api/orders'

export default function Orders() {
  const [orders, setOrders] = useState<any[]>([])
  useEffect(() => { myOrders().then(setOrders).catch(() => message.error('Моля, влезте в профила си')) }, [])
  return (
    <div>
      <Typography.Title level={2}>Моите поръчки</Typography.Title>
      <Table rowKey="id" dataSource={orders} columns={[
        { title: '№', dataIndex: 'id' },
        { title: 'Статус', dataIndex: 'status', render: (s) => <Tag>{s}</Tag> },
        { title: 'Плащане', dataIndex: 'payment_status' },
        { title: 'Общо', dataIndex: 'total_price' },
        { title: 'Изработка (дни)', dataIndex: 'estimated_production_time_days' },
      ]} />
    </div>
  )
}

