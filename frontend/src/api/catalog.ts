import { api } from './client'

export const fetchDepartments = () => api.get('/departments').then(r => r.data)
export const fetchCategories = (deptId: number) => api.get(`/departments/${deptId}/categories`).then(r => r.data)
export const fetchProductsByCategory = (catId: number) => api.get(`/categories/${catId}/products`).then(r => r.data)
export const searchProducts = (q: string) => api.get('/products/search', { params: { query: q } }).then(r => r.data)
export const fetchProduct = (id: number) => api.get(`/products/${id}`).then(r => r.data)
export const fetchRecommendations = (id: number) => api.get(`/products/${id}/recommendations`).then(r => r.data)

