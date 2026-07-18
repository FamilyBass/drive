import { apiClient } from './api'

// Service d'administration
class AdminService {
  async validateUser(userId) {
    return apiClient.validateUser(userId)
  }
}

export const adminService = new AdminService()

