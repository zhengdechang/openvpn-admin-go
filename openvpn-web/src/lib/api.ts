import axios from "axios";
import {
  ApiResponse,
  User,
  Industries,
  Policy,
  PolicyLabel,
  AddPolicyRequest,
  LoginCredentials,
  RegisterCredentials,
  OpenVPNClient,
  ServerStatus,
} from "./types";
import Cookies from "js-cookie";
import { useUserStore } from "@/store";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:3000";

// 创建axios实例
const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

// 请求拦截器添加token
api.interceptors.request.use((config) => {
  const token = Cookies.get("token");
  if (token) {
    config.headers["Authorization"] = `Bearer ${token}`;
  }
  return config;
});

// 响应拦截器：处理 401 错误
api.interceptors.response.use(
  (response) => response, // 正常返回
  (error) => {
    if (error.response?.status === 401) {
      console.warn("Unauthorized or session expired, clearing login info...");
      useUserStore.getState().clearLoginInfo();
    }
    return Promise.reject(error); // 继续抛出错误，供业务代码处理
  }
);

// 用户API
export const userAPI = {
  // 获取当前用户信息
  getMe: async (): Promise<ApiResponse<User>> => {
    try {
      const response = await api.get("/api/user/me");
      return response.data;
    } catch (error) {
      return {
        success: false,
        error: "Failed to fetch user information",
      };
    }
  },

  // 用户注册
  register: async (credentials: RegisterCredentials): Promise<ApiResponse<Object>> => {
    try {
      const response = await api.post("/api/user/register", credentials);
      return response.data;
    } catch (error: any) {
      return {
        success: false,
        error: error.response ? error.response.data.error : "Registration failed",
      };
    }
  },

  verifyEmail: async (token: string): Promise<ApiResponse<Object>> => {
    try {
      const response = await api.get(`/api/user/verify-email/${token}`);
      return response.data;
    } catch (error: any) {
      return {
        success: false,
        error: error.response ? error.response.data.error : "Email verification failed",
      };
    }
  },

  forgotPassword: async (email: string): Promise<ApiResponse<Object>> => {
    try {
      const response = await api.post("/api/user/forgot-password", { email });
      return response.data;
    } catch (error: any) {
      return {
        success: false,
        error: error.response ? error.response.data.error : "Password reset failed",
      };
    }
  },

  resetPassword: async (
    resetToken: string,
    password: string,
    confirmPassword: string
  ): Promise<ApiResponse<Object>> => {
    try {
      const response = await api.patch(
        `/api/user/reset-password/${resetToken}`,
        { password, confirmPassword }
      );
      return response.data;
    } catch (error: any) {
      return {
        success: false,
        error: error.response ? error.response.data.error : "Password reset failed",
      };
    }
  },

  // 用户登录
  login: async (credentials: LoginCredentials): Promise<ApiResponse<{ user: User; token: string }>> => {
    try {
      const response = await api.post("/api/user/login", credentials);
      if (response.data.success && response.data.data.token) {
        Cookies.set("token", response.data.data.token, { expires: 7 });
      }
      return response.data;
    } catch (error: any) {
      return {
        success: false,
        error: error.response ? error.response.data.error : "Login failed",
      };
    }
  },

  // Refresh token
  refreshToken: async (): Promise<ApiResponse<{ token: string }>> => {
    try {
      const response = await api.get("/api/user/refresh");
      if (response.data.success && response.data.data.token) {
        Cookies.set("token", response.data.data.token, { expires: 7 });
      }
      return response.data;
    } catch (error) {
      return {
        success: false,
        error: "Token refresh failed",
      };
    }
  },

  getEnterpriseById: async (id: string): Promise<ApiResponse<User>> => {
    try {
      const response = await api.get(`/api/user/info/${id}`);
      return response.data;
    } catch (error: any) {
      return {
        success: false,
        error: error.response.data.error || "获取用户失败",
      };
    }
  },

  // 更新用户信息
  updateMe: async (userData: Partial<User>): Promise<ApiResponse<User>> => {
    try {
      const response = await api.patch("/api/user/me", userData);
      return response.data;
    } catch (error) {
      return {
        success: false,
        error: "Failed to update user information",
      };
    }
  },

  // 获取用户角色
  getRoles: async (): Promise<ApiResponse<string[]>> => {
    try {
      const response = await api.get("/api/user/roles");
      return response.data;
    } catch (error) {
      return {
        success: false,
        error: "Failed to fetch user roles",
      };
    }
  },

  // 退出登录
  logout: async () => {
    try {
      const response = await api.post("/api/user/logout");
      if (response.data.success) {
        Cookies.remove("token");
      }
      return response.data;
    } catch (error) {
      return {
        success: false,
        error: "Logout failed",
      };
    }
  },
};

// OpenVPN 客户端管理 API
export const openvpnAPI = {
  // 获取客户端列表
  getClientList: async (): Promise<OpenVPNClient[]> => {
    const response = await api.get<OpenVPNClient[]>('/api/client/list');
    return response.data;
  },
  // 添加客户端
  addClient: async (username: string): Promise<any> => {
    const response = await api.post('/api/client/add', { username });
    return response.data;
  },
  // 更新客户端
  updateClient: async (username: string): Promise<any> => {
    const response = await api.put('/api/client/update', { username });
    return response.data;
  },
  // 删除客户端
  deleteClient: async (username: string): Promise<any> => {
    const response = await api.delete(`/api/client/delete/${username}`);
    return response.data;
  },
  // 获取客户端配置
  getClientConfig: async (username: string): Promise<{ config: string }> => {
    const response = await api.get<{ config: string }>(`/api/client/config/${username}`);
    return response.data;
  },
  // 吊销客户端证书
  revokeClient: async (username: string): Promise<any> => {
    const response = await api.post('/api/client/revoke', { username });
    return response.data;
  },
  // 续期客户端证书
  renewClient: async (username: string): Promise<any> => {
    const response = await api.post('/api/client/renew', { username });
    return response.data;
  },
  // 获取服务器状态
  getServerStatus: async (): Promise<ServerStatus> => {
    const response = await api.get<ServerStatus>('/api/server/status');
    return response.data;
  },
  // 获取服务器日志
  getServerLogs: async (): Promise<string> => {
    const response = await api.get<{ logs: string }>('/api/logs/server');
    return response.data.logs;
  },
  // 获取指定客户端日志
  getClientLogs: async (username: string): Promise<string[]> => {
    const response = await api.get<{ logs: string[] }>(`/api/logs/client/${username}`);
    return response.data.logs;
  },
};

// 服务器管理 API
export const serverAPI = {
  getStatus: async (): Promise<ServerStatus> => {
    const response = await api.get<ServerStatus>('/api/server/status');
    return response.data;
  },
  start: async (): Promise<any> => {
    const response = await api.post('/api/server/start');
    return response.data;
  },
  stop: async (): Promise<any> => {
    const response = await api.post('/api/server/stop');
    return response.data;
  },
  restart: async (): Promise<any> => {
    const response = await api.post('/api/server/restart');
    return response.data;
  },
  getConfigTemplate: async (): Promise<{ template: string }> => {
    const response = await api.get<{ template: string }>('/api/server/config/template');
    return response.data;
  },
  updateConfig: async (config: string): Promise<any> => {
    const response = await api.put('/api/server/config', { config });
    return response.data;
  },
  updatePort: async (port: number): Promise<any> => {
    const response = await api.put('/api/server/port', { port });
    return response.data;
  },
};

// 部门管理 API
export const departmentAPI = {
  list: async (): Promise<Department[]> => {
    const response = await api.get<Department[]>('/api/departments');
    return response.data;
  },
  create: async (name: string): Promise<any> => {
    const response = await api.post('/api/departments', { name });
    return response.data;
  },
  update: async (id: string, name: string): Promise<any> => {
    const response = await api.put(`/api/departments/${id}`, { name });
    return response.data;
  },
  delete: async (id: string): Promise<any> => {
    const response = await api.delete(`/api/departments/${id}`);
    return response.data;
  },
};

// 用户管理 API
export const userManagementAPI = {
  list: async (): Promise<AdminUser[]> => {
    const response = await api.get<AdminUser[]>('/api/users');
    return response.data;
  },
  create: async (user: Partial<AdminUser> & { password: string }): Promise<any> => {
    const response = await api.post('/api/users', user);
    return response.data;
  },
  update: async (id: string, updates: Partial<AdminUser> & { password?: string }): Promise<any> => {
    const response = await api.put(`/api/users/${id}`, updates);
    return response.data;
  },
  delete: async (id: string): Promise<any> => {
    const response = await api.delete(`/api/users/${id}`);
    return response.data;
  },
};
};

// 标签API
export const industryAPI = {
  // 添加标签
  addIndustry: async (
    tagData: Partial<Industries>
  ): Promise<ApiResponse<Industries>> => {
    try {
      const response = await api.post("/api/industry/add", tagData);
      return response.data;
    } catch (error) {
      return {
        success: false,
        error: "添加标签失败",
      };
    }
  },

  addPolicy: async (
    data: Partial<AddPolicyRequest>
  ): Promise<ApiResponse<Policy>> => {
    try {
      const formData = new FormData();

      // 附加文件（如果有）
      if (data.file) {
        formData.append("file", data.file);
      }

      // 附加其他字段
      if (data.title) formData.append("title", data.title);
      if (data.type) formData.append("type", data.type);
      if (data.issuingUnit) formData.append("issuing_unit", data.issuingUnit);
      if (data.createTime) formData.append("create_time", data.createTime);
      if (data.url) formData.append("url", data.url);
      if (data.description) formData.append("description", data.description);

      // 处理 labels（遍历并逐个添加 name）
      if (data.labels && Array.isArray(data.labels)) {
        data.labels.forEach((label) => {
          formData.append("labels", label.name); // 只存 name
        });
      }

      // 发送 FormData 请求
      const response = await api.post("/api/policy/add", formData, {
        headers: {
          "Content-Type": "multipart/form-data",
        },
      });

      return response.data;
    } catch (error) {
      console.error("添加政策失败:", error);
      return {
        success: false,
        error: "添加政策失败，请稍后重试",
      };
    }
  },

  updatePolicy: async (
    data: Partial<AddPolicyRequest>,
    id: string
  ): Promise<ApiResponse<Policy>> => {
    try {
      const formData = new FormData();

      // 附加文件（如果有）
      if (data.file) {
        formData.append("file", data.file);
      }

      // 附加其他字段
      if (data.title) formData.append("title", data.title);
      if (data.type) formData.append("type", data.type);
      if (data.issuingUnit) formData.append("issuing_unit", data.issuingUnit);
      if (data.createTime) formData.append("create_time", data.createTime);
      if (data.url) formData.append("url", data.url);
      if (data.description) formData.append("description", data.description);

      // 处理 labels（遍历并逐个添加 name）
      if (data.labels && Array.isArray(data.labels)) {
        data.labels.forEach((label) => {
          formData.append("labels", label.name); // 只存 name
        });
      }

      // 发送 FormData 请求
      const response = await api.post(`/api/policy/update/${id}`, formData, {
        headers: {
          "Content-Type": "multipart/form-data",
        },
      });

      return response.data;
    } catch (error) {
      console.error("添加政策失败:", error);
      return {
        success: false,
        error: "添加政策失败，请稍后重试",
      };
    }
  },

  deletePolicy: async (id: string): Promise<ApiResponse<any>> => {
    try {
      // 发送 FormData 请求
      const response = await api.get(`/api/policy/delete/${id}`);

      return response.data;
    } catch (error) {
      console.error("删除政策失败:", error);
      return {
        success: false,
        error: "删除政策失败，请稍后重试",
      };
    }
  },

  // 获取所有标签
  getIndustry: async (): Promise<ApiResponse<Industries[]>> => {
    try {
      const response = await api.get("/api/industry/all");
      return response.data;
    } catch (error) {
      return {
        success: false,
        error: "获取标签失败",
      };
    }
  },
  getAllPolicies: async (): Promise<ApiResponse<Policy[]>> => {
    try {
      const response = await api.get("/api/policy/all");
      let policies = response.data.data.map((policy: any) => {
        return {
          ...policy,
          publishedAt: policy.create_time,
          createdAt: policy.create_time,
          issuingUnit: policy.issuing_unit,
        };
      });

      return {
        ...response.data,
        data: policies,
      };
    } catch (error) {
      return {
        success: false,
        error: "获取政策失败1",
      };
    }
  },
  getAllPolicyLabels: async (): Promise<ApiResponse<PolicyLabel[]>> => {
    try {
      const response = await api.get("/api/policy_label/all");
      return response.data;
    } catch (error) {
      return {
        success: false,
        error: "获取政策标签失败",
      };
    }
  },
  getUserPolicies: async (): Promise<ApiResponse<Policy[]>> => {
    try {
      const response = await api.get("/api/policy/get_user_policies");
      let policies = response.data.data.map((policy: any) => {
        return {
          ...policy,
          publishedAt: policy.create_time,
          createdAt: policy.create_time,
          issuingUnit: policy.issuing_unit,
        };
      });

      return {
        ...response.data,
        data: policies,
      };
    } catch (error) {
      return {
        success: false,
        error: "获取政策失败",
      };
    }
  },
  getGovernmentUserPolicies: async (): Promise<ApiResponse<Policy[]>> => {
    try {
      const response = await api.get("/api/policy/government/all");
      let policies = response.data.data.map((policy: any) => {
        return {
          ...policy,
          publishedAt: policy.create_time,
          createdAt: policy.create_time,
          issuingUnit: policy.issuing_unit,
        };
      });

      return {
        ...response.data,
        data: policies,
      };
    } catch (error) {
      return {
        success: false,
        error: "获取政策失败",
      };
    }
  },
  getPolicyForId: async (id: string): Promise<ApiResponse<Policy>> => {
    try {
      const response = await api.get(`/api/policy/${id}`);
      let policy = {
        ...response.data.data,
        publishedAt: response.data.data.create_time,
        issuingUnit: response.data.data.issuing_unit,
        createdAt: response.data.data.create_time,
      };
      return {
        ...response.data,
        data: policy,
      };
    } catch (error) {
      return {
        success: false,
        error: "获取政策失败",
      };
    }
  },
  getPolicyMatchEnterprise: async (
    id: string
  ): Promise<ApiResponse<User[]>> => {
    try {
      const response = await api.get(`/api/policy/enterprise/match/${id}`);
      return response.data;
    } catch (error) {
      return {
        success: false,
        error: "获取企业失败",
      };
    }
  },
};
