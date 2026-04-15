/*
 * @Description: 用户状态管理
 * @Author: Devin
 * @Date: 2025-03-19 10:27:01
 */
import { create } from "zustand";
import { persist } from "zustand/middleware";
import { User } from "@/types/types";
import Cookies from "js-cookie";

// Zustand Store 的类型
interface UserState {
  user: User | null;
  isLogin: boolean;
  accessToken: string;
  updateAccessToken: (accessToken: string) => void;
  updateUser: (user: User) => void;
  updateIsLogin: (isLogin: boolean) => void;
  clearLoginInfo: () => void;
}

export const USER_KEY = "user-info";

export const useUserStore = create<UserState>()(
  persist(
    (set) => ({
      user: null,
      isLogin: false,
      accessToken: "",

      updateAccessToken: (accessToken) => set({ accessToken }),

      updateUser: (user) => set({ user }),

      updateIsLogin: (isLogin) => set({ isLogin }),

      clearLoginInfo: () => {
        set({
          isLogin: false,
          accessToken: "",
          user: null,
        });
        Cookies.remove("token");
      },
    }),
    {
      name: USER_KEY,
      version: 1,
    }
  )
);
