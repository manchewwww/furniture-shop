import { Alert, Button, Card, Form, Input, Select, Typography, message } from 'antd'
import { useCart } from '../store/CartContext'
import { createOrder, payByCard } from '../api/orders'
import { useState } from 'react'

export default function Checkout() {
  const { items, clear } = useCart()
  const [orderId, setOrderId] = useState<number | null>(null)
  const [paymentMethod, setPaymentMethod] = useState<string>('наложенплатеж')

  const onFinish = async (values: any) => {
    const payload = {
      name: values.name, email: values.email, phone: values.phone, address: values.address,
      payment_method: values.payment_method,
      items: items.map(it => ({ product_id: it.product.id, quantity: it.quantity, options: it.options }))
    }
    try {
      const res = await createOrder(payload)
      setOrderId(res.order_id)
      setPaymentMethod(values.payment_method)
      message.success('Поръчката е създадена')
    } catch {
      message.error('Грешка при създаване на поръчка')
    }
  }

  const onPayCard = async (values: any) => {
    try {
      await payByCard({ order_id: orderId, ...values })
      message.success('Плащането е успешно')
      clear()
    } catch {
      message.error('Плащането е отказано')
    }
  }

  return (
    <div>
      <Typography.Title level={2}>Поръчка</Typography.Title>
      {!orderId && (
        <Card title="Данни за клиент и доставка">
          <Form layout="vertical" onFinish={onFinish} initialValues={{ payment_method: 'наложенплатеж' }}>
            <Form.Item name="name" label="Име" rules={[{ required: true, message: 'Въведете име' }]}><Input /></Form.Item>
            <Form.Item name="email" label="Имейл" rules={[{ required: true, message: 'Въведете имейл' }]}><Input /></Form.Item>
            <Form.Item name="phone" label="Телефон" rules={[{ required: true, message: 'Въведете телефон' }]}><Input /></Form.Item>
            <Form.Item name="address" label="Адрес" rules={[{ required: true, message: 'Въведете адрес' }]}><Input /></Form.Item>
            <Form.Item name="payment_method" label="Метод на плащане" rules={[{ required: true }]}>
              <Select options={[{value:'карта',label:'Карта'},{value:'наложенплатеж',label:'Наложен платеж'},{value:'банковпревод',label:'Банков превод'}]} />
            </Form.Item>
            <Button type="primary" htmlType="submit">Създай поръчка</Button>
          </Form>
        </Card>
      )}
      {orderId && paymentMethod !== 'карта' && (
        <Alert type="success" message={`Поръчка №${orderId} е създадена. Очаква плащане.`} />
      )}
      {orderId && paymentMethod === 'карта' && (
        <Card title="Плащане с карта">
          <Form layout="vertical" onFinish={onPayCard}>
            <Form.Item name="cardholder_name" label="Име на картодържателя" rules={[{ required: true }]}><Input /></Form.Item>
            <Form.Item name="card_number" label="Номер на карта" rules={[{ required: true }]}><Input /></Form.Item>
            <Form.Item name="expiry_month" label="Месец" rules={[{ required: true }]}><Input /></Form.Item>
            <Form.Item name="expiry_year" label="Година" rules={[{ required: true }]}><Input /></Form.Item>
            <Form.Item name="cvv" label="CVV" rules={[{ required: true }]}><Input /></Form.Item>
            <Button type="primary" htmlType="submit">Плати</Button>
          </Form>
        </Card>
      )}
    </div>
  )
}

