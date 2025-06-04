/*
 * @Description:
 * @Author: Devin
 * @Date: 2025-03-17 10:05:46
 */
// User role enum
export enum UserRole {
  ADMIN = "admin",
  MANAGER = "manager",
  USER = "user",
  SUPERADMIN = "superadmin",
}

// User type definition
export interface User {
  id: string;
  name: string;
  email: string;
  role: UserRole;
  createdAt: string;
  updatedAt: string;
  // Department ID for user, if applicable
  departmentId?: string;
  avatar?: string;
  bio?: string;
}

// API response type
export interface ApiResponse<T> {
  success: boolean;
  data?: T;
  message?: string;
  error?: string;
}

// File info type
export interface FileInfo {
  id: string;
  name: string;
  extension: string;
  path: string;
  size: number;
  storage: string;
  storage_key: string;
  created_at: Date;
  updated_at: Date;
}

// Authentication types
export interface LoginCredentials {
  email: string;
  password: string;
}

export interface RegisterCredentials {
  name: string;
  email: string;
  password: string;
  confirmPassword: string;
}

// Dashboard types
export interface DashboardStats {
  totalUsers: number;
  activeUsers: number;
  totalFiles: number;
  storageUsed: number;
}

// Settings types
export interface UserSettings {
  theme: 'light' | 'dark' | 'system';
  notifications: boolean;
  language: string;
}

export interface EnterpriseInfo {
  name?: string; // 企业名称
  enterprise_location?: string; // 企业注册地
  postal_address?: string; // 邮政地址
  credit_code?: string; // 统一社会信用代码
  zip_code?: string; // 邮编
  legal_representative?: string; // 法定代表人
  legal_rep_phone?: string; // 法定代表人电话
  legal_rep_mobile?: string; // 法定代表人手机
  holding_shareholder?: string; // 控股股东
  actual_controller?: string; // 实际控制人
  actual_controller_nation?: string; // 实际控制人国籍
  city?: string; // 所在设区市
  county?: string; // 所在县
  contact_person?: string; // 联系人
  contact_mobile?: string; // 联系人手机
  fax?: string; // 传真
  email?: string; // E-mail
  registration_date?: string; // 注册时间
  registered_capital?: string; // 注册资本
  enterprise_scale?: string; // 企业规模
  industry?: string; // 所属行业
  sub_industry?: string; // 细分领域
  enterprise_type?: string; // 企业类型
  listing_status?: string; // 上市情况
  listing_progress?: string; // 上市进度
  planned_listing_place?: string; // 拟上市地
  employee_count?: number; // 员工数
  rd_personnel_count?: number; // 研发人员数
}

// 标签类型定义
export interface Industries {
  id: string;
  name: string;
  description?: string;
  createdAt: string;
  updatedAt: string;
}

export interface PolicyLabel {
  id: string;
  name: string;
  industries?: Industries[];
}

// 政策类型定义
export interface Policy {
  id: string;
  title: string;
  content: string;
  labels: string;
  policy_labels: PolicyLabel[];
  industries: Industries[];
  type: string;
  issuingUnit: string;
  publishedAt: string;
  url: string;
  createdAt: string;
  updatedAt: string;
  file_id: string;
  description?: string;
  file?: FileInfo;
}

export interface AddPolicyRequest {
  title: string;
  type: string;
  issuingUnit: string;
  createTime: string;
  url: string;
  file: File;
  labels: PolicyLabel[];
  description: string;
}

export interface EditPolicyRequest {
  title: string;
  type: string;
  issuingUnit: string;
  createTime: string;
  url: string;
  file: File;
  labels: PolicyLabel[];
  description: string;
}

export const validEnterpriseSizes = ["微型", "小型", "中型", "大型"] as const;
export const validEnterpriseType = ["国有", "民营", "其他"] as const;

export interface OpenVPNClient {
  id: string;
  name: string;
  email: string;
  status: 'active' | 'inactive' | 'pending';
  createdAt: string;
  lastConnected?: string;
  ipAddress?: string;
  notes?: string;
  // Department ID associated with this client
  departmentId?: string;
}

// 部门类型
export interface Department {
  id: string;
  name: string;
  // 部门负责人ID
  headId?: string;
  // 部门负责人信息
  head?: {
    id: string;
    name: string;
  };
  // 上级部门ID
  parentId?: string;
  createdAt: string;
  updatedAt: string;
}

// 管理用户类型
export interface AdminUser {
  id: string;
  name: string;
  email: string;
  role: UserRole;
  departmentId?: string;
}

// 服务器状态
export interface ServerStatus {
  name: string;
  status: string;
  uptime: string;
  connected: number;
  total: number;
  lastUpdated: string;
}
