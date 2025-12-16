import { api } from './client'

export const createOrder = (payload: any) => api.post('/orders', payload).then(r => r.data)
export const payByCard = (payload: any) => api.post('/payments/card', payload).then(r => r.data)
export const myOrders = () => api.get('/user/orders').then(r => r.data)
export const myOrder = (id: number) => api.get(`/user/orders/${id}`).then(r => r.data)

