"use client";

import React, { createContext, useContext, useEffect, useState } from "react";
import { User, LoginCredentials, RegisterCredentials } from "../types/types";
import { userAPI } from "../services/api";
import { useRouter } from "next/navigation";
import { useUserStore } from "@/store";
import { showToast } from "@/lib/toast-utils";

interface AuthContextType {
  user: User | null;
  loading: boolean;
  error: string | null;
  login: (credentials: LoginCredentials) => Promise<User | undefined>;
  register: (credentials: RegisterCredentials) => Promise<boolean>;
  logout: () => void;
  updateUserInfo: (userData: Partial<User>) => Promise<boolean>;
  setUser: (user: User | null) => void;
  isLogin: boolean;
  refreshToken: () => Promise<boolean>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: React.ReactNode }) => {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const router = useRouter();

  const {
    user: storeUser,
    updateUser,
    isLogin,
    accessToken,
    updateIsLogin,
    clearLoginInfo,
  } = useUserStore();

  const fetchUser = async () => {
    setLoading(true);
    try {
      const response = await userAPI.getMe();
      if (response.success && response.data) {
        setUser(response.data);
        updateUser(response.data);
        return response.data;
      }
    } catch (error) {
      setError("Failed to fetch user information");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (isLogin && storeUser) {
      setUser(storeUser);
    }
  }, [isLogin, storeUser]);

  const refreshToken = async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await userAPI.refreshToken();
      if (response.success) {
        updateIsLogin(true);
        await fetchUser();
        return true;
      } else {
        setError(response.error || "Please login first");
        return false;
      }
    } catch (error) {
      setError("Please login first");
      return false;
    } finally {
      setLoading(false);
    }
  };

  const login = async (credentials: LoginCredentials) => {
    setLoading(true);
    setError(null);
    try {
      const response = await userAPI.login(credentials);
      if (response.success) {
        updateIsLogin(true);
        const user = await fetchUser();
        return user;
      } else {
        showToast.error(
          response.error || "Login failed. Please check your credentials"
        );
        setError(
          response.error || "Login failed. Please check your credentials"
        );
        return;
      }
    } catch (error) {
      setError("Login failed. Please check your credentials");
      return;
    } finally {
      setLoading(false);
    }
  };

  const register = async (
    credentials: RegisterCredentials
  ): Promise<boolean> => {
    setLoading(true);
    setError(null);
    try {
      const response = await userAPI.register(credentials);
      if (!response.success) {
        setError(response.error || "Registration failed");
        showToast.error(response.error || "Registration failed");
        return false;
      }
      if (response.data && response.message) {
        showToast.success(
          response.message ||
            "Registration successful. Please verify your email"
        );
        return true;
      }
      return false;
    } catch (error) {
      setError("Registration failed");
      return false;
    } finally {
      setLoading(false);
    }
  };

  const logout = async () => {
    await userAPI.logout();
    setUser(null);
    clearLoginInfo();
    router.push("/");
  };

  const updateUserInfo = async (userData: Partial<User>) => {
    setLoading(true);
    setError(null);
    try {
      const response = await userAPI.updateMe(userData);
      if (response.success && response.data) {
        setUser(response.data);
        return true;
      } else {
        setError(response.error || "Failed to update user information");
        router.push("/auth/login");
        return false;
      }
    } catch (error) {
      setError("Failed to update user information");
      router.push("/auth/login");
      return false;
    } finally {
      setLoading(false);
    }
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        loading,
        error,
        login,
        register,
        logout,
        updateUserInfo,
        setUser,
        isLogin,
        refreshToken,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};
