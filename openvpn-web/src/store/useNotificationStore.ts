"use client";

import { create } from "zustand";
import { Notification } from "@/types";
import { notificationAPI } from "@/services/api";

interface NotificationState {
  notifications: Notification[];
  unreadCount: number;
  isOpen: boolean;
  isLoading: boolean;
  error: string | null;

  // Actions
  fetchNotifications: () => Promise<void>;
  fetchUnreadCount: () => Promise<void>;
  markRead: (id: string) => Promise<void>;
  markAllRead: () => Promise<void>;
  setOpen: (open: boolean) => void;
}

export const useNotificationStore = create<NotificationState>((set, get) => ({
  notifications: [],
  unreadCount: 0,
  isOpen: false,
  isLoading: false,
  error: null,

  fetchNotifications: async () => {
    set({ isLoading: true, error: null });
    try {
      const notifications = await notificationAPI.list();
      set({ notifications, isLoading: false });
    } catch {
      set({ isLoading: false, error: "Failed to load notifications" });
    }
  },

  fetchUnreadCount: async () => {
    try {
      const count = await notificationAPI.getUnreadCount();
      set({ unreadCount: count });
    } catch {
      // Stale count shown — no crash, no error toast
    }
  },

  markRead: async (id: string) => {
    try {
      await notificationAPI.markRead(id);
      set((state) => ({
        notifications: state.notifications.map((n) =>
          n.id === id ? { ...n, isRead: true } : n
        ),
        unreadCount: Math.max(0, state.unreadCount - 1),
      }));
    } catch {
      // Best-effort — UI stays optimistic if needed
    }
  },

  markAllRead: async () => {
    try {
      await notificationAPI.markAllRead();
      set((state) => ({
        notifications: state.notifications.map((n) => ({ ...n, isRead: true })),
        unreadCount: 0,
      }));
    } catch {
      // Partial failure: re-fetch to reconcile
      get().fetchUnreadCount();
    }
  },

  setOpen: (open: boolean) => {
    set({ isOpen: open });
    if (open) {
      // Fetch fresh list every time the panel opens
      get().fetchNotifications();
    }
  },
}));
