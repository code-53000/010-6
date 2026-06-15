import request from '../utils/request'

export const getPricingList = () => request.get('/pricing')
export const updatePricing = (data) => request.put('/pricing', data)
export const calculateFee = (params) => request.get('/pricing/calculate', { params })

export const createOrder = (data) => request.post('/queue', data)
export const getOrderList = (params) => request.get('/queue', { params })
export const getOrder = (id) => request.get(`/queue/${id}`)
export const callNext = () => request.post('/queue/call-next')
export const callSpecific = (id) => request.post(`/queue/${id}/call`)
export const completeOrder = (id, data) => request.post(`/queue/${id}/complete`, data)
export const cancelOrder = (id) => request.post(`/queue/${id}/cancel`)
export const pickupOil = (id, data) => request.post(`/queue/${id}/pickup`, data)

export const createStorage = (data) => request.post('/storage', data)
export const getStorageList = (params) => request.get('/storage', { params })
export const getStorage = (id) => request.get(`/storage/${id}`)
export const pickupStorage = (id) => request.post(`/storage/${id}/pickup`)
