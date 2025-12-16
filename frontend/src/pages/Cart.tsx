import { Button, Card, List, Typography } from 'antd'
import { Link } from 'react-router-dom'
import { useCart } from '../store/CartContext'

export default function Cart() {
  const { items, remove, clear } = useCart()
  const total = items.reduce((s, it) => s + it.product.base_price * it.quantity, 0)
  return (
    <div>
      <Typography.Title level={2}>Количка</Typography.Title>
      <List
        dataSource={items}
        renderItem={(it) => (
          <List.Item actions={[<Button danger onClick={() => remove(it.product.id)}>Премахни</Button>] }>
            <List.Item.Meta title={it.product.name} description={`Количество: ${it.quantity}`} />
            <div>{(it.product.base_price * it.quantity).toFixed(2)} лв.</div>
          </List.Item>
        )}
      />
      <Card>
        <p>Общо: {total.toFixed(2)} лв.</p>
        <Button type="primary" disabled={items.length===0}><Link to="/checkout">Към поръчка</Link></Button>
        <Button style={{ marginLeft: 8 }} onClick={clear}>Изчисти</Button>
      </Card>
    </div>
  )
}

