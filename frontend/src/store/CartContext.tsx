import React, { createContext, useContext, useEffect, useMemo, useState } from 'react'

export type CartItem = {
  product: any
  quantity: number
  options: { id: number; type: string }[]
}

type CartCtxType = {
  items: CartItem[]
  add: (item: CartItem) => void
  remove: (productId: number) => void
  clear: () => void
}

const CartCtx = createContext<CartCtxType | undefined>(undefined)

export const CartProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [items, setItems] = useState<CartItem[]>(() => {
    const saved = localStorage.getItem('cart')
    return saved ? JSON.parse(saved) : []
  })
  useEffect(() => { localStorage.setItem('cart', JSON.stringify(items)) }, [items])
  const value = useMemo(() => ({
    items,
    add: (item: CartItem) => setItems(prev => {
      const existing = prev.find(p => p.product.id === item.product.id)
      if (existing) {
        return prev.map(p => p.product.id === item.product.id ? { ...p, quantity: p.quantity + item.quantity } : p)
      }
      return [...prev, item]
    }),
    remove: (id: number) => setItems(prev => prev.filter(p => p.product.id !== id)),
    clear: () => setItems([])
  }), [items])
  return <CartCtx.Provider value={value}>{children}</CartCtx.Provider>
}

export const useCart = () => {
  const ctx = useContext(CartCtx)
  if (!ctx) throw new Error('CartProvider is missing')
  return ctx
}

