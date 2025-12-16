import { Button, Card, Form, Input, Modal, Table, Tabs, Tag, message } from 'antd'
import { useEffect, useState } from 'react'
import { api } from '../api/client'

export default function AdminDashboard() {
  const [depts, setDepts] = useState<any[]>([])
  const [orders, setOrders] = useState<any[]>([])
  const [open, setOpen] = useState(false)
  const [form] = Form.useForm()
  const load = async () => {
    try {
      const [d, o] = await Promise.all([
        api.get('/admin/departments'),
        api.get('/admin/orders')
      ])
      setDepts(d.data); setOrders(o.data)
    } catch { message.error('Нужен е админ вход') }
  }
  useEffect(() => { load() }, [])

  const createDept = async () => {
    const v = await form.validateFields()
    await api.post('/admin/departments', v)
    setOpen(false); form.resetFields(); load()
  }

  const setStatus = async (id: number, status: string) => {
    await api.patch(`/admin/orders/${id}/status`, { status })
    load()
  }

  return (
    <div>
      <Tabs items={[
        { key: 'depts', label: 'Отдели', children: (
          <Card title="Отдели" extra={<Button onClick={() => setOpen(true)}>Нов отдел</Button>}>
            <Table rowKey="id" dataSource={depts} columns={[{ title:'Име', dataIndex:'name' },{ title:'Описание', dataIndex:'description' }]} />
            <Modal title="Нов отдел" open={open} onOk={createDept} onCancel={() => setOpen(false)}>
              <Form layout="vertical" form={form}>
                <Form.Item name="name" label="Име" rules={[{ required:true }]}><Input /></Form.Item>
                <Form.Item name="description" label="Описание"><Input /></Form.Item>
              </Form>
            </Modal>
          </Card>
        ) },
        { key: 'orders', label: 'Поръчки', children: (
          <Card title="Всички поръчки">
            <Table rowKey="id" dataSource={orders} columns={[
              { title:'№', dataIndex:'id' },
              { title:'Статус', dataIndex:'status', render:(s, r) => <Tag>{s}</Tag> },
              { title:'Плащане', dataIndex:'payment_status' },
              { title:'Общо', dataIndex:'total_price' },
              { title:'Действия', render: (_:any, r:any) => (
                <>
                  {['нова','потвърдена','впроизводство','изпратена','доставена','отказана'].map(st => (
                    <Button key={st} size="small" style={{ marginRight: 4 }} onClick={() => setStatus(r.id, st)}>{st}</Button>
                  ))}
                </>
              )}
            ]} />
          </Card>
        ) },
      ]} />
    </div>
  )
}

