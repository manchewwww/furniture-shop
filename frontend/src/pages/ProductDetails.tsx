import { Button, Card, Col, InputNumber, Row, Select, Typography, message } from 'antd'
import { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import { fetchProduct, fetchRecommendations } from '../api/catalog'
import { useCart } from '../store/CartContext'

export default function ProductDetails() {
  const { id } = useParams()
  const [product, setProduct] = useState<any>()
  const [rec, setRec] = useState<any[]>([])
  const [qty, setQty] = useState<number>(1)
  const [selected, setSelected] = useState<number[]>([])
  const { add } = useCart()
  useEffect(() => {
    if (id) {
      fetchProduct(Number(id)).then(setProduct)
      fetchRecommendations(Number(id)).then(setRec)
    }
  }, [id])
  if (!product) return null
  return (
    <div>
      <Row gutter={24}>
        <Col md={14} xs={24}>
          <img src={product.image_url} alt={product.name} style={{ maxWidth:'100%' }} />
          <Typography.Title level={3}>{product.name}</Typography.Title>
          <p>{product.long_description}</p>
          <p>Базова цена: {product.base_price.toFixed(2)} лв.</p>
          <p>Базово време за изработка: {product.base_production_time_days} дни</p>
          <div style={{ margin: '12px 0' }}>
            <Typography.Text>Опции:</Typography.Text>
            <Select
              mode="multiple"
              placeholder="Изберете опции"
              style={{ minWidth: 320 }}
              onChange={(vals) => setSelected(vals as number[])}
              options={(product.options || []).map((o: any) => ({ value: o.id, label: `${o.option_name} (${o.option_type})` }))}
            />
          </div>
          <div style={{ display: 'flex', gap: 12, alignItems: 'center' }}>
            <Typography.Text>Количество:</Typography.Text>
            <InputNumber min={1} value={qty} onChange={(v) => setQty(Number(v))} />
            <Button type="primary" onClick={() => {
              add({ product, quantity: qty, options: selected.map(id => ({ id, type: 'extra' })) })
              message.success('Добавено в количката')
            }}>Добави в количката</Button>
          </div>
        </Col>
        <Col md={10} xs={24}>
          <Typography.Title level={4}>Подобни продукти</Typography.Title>
          {rec.map(r => (<Card key={r.id} size="small" style={{ marginBottom: 8 }} title={r.name}></Card>))}
        </Col>
      </Row>
    </div>
  )
}

