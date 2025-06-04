import { toast } from "sonner";

// 用于跟踪最近显示的消息
const recentMessages: Map<string, number> = new Map();
const MESSAGE_EXPIRY_TIME = 3000; // 3秒内不重复显示相同消息
const DEFAULT_DURATION = 3000; // 默认显示时间

/**
 * 自定义toast函数，避免短时间内显示重复消息
 */
export const showToast = {
  success: (message: string, options?: any) => {
    if (shouldShowMessage(message)) {
      toast.success(message, {
        duration: DEFAULT_DURATION,
        ...options,
      });
    }
  },
  error: (message: string, options?: any) => {
    if (shouldShowMessage(message)) {
      toast.error(message, {
        duration: DEFAULT_DURATION,
        ...options,
      });
    }
  },
  info: (message: string, options?: any) => {
    if (shouldShowMessage(message)) {
      toast.info(message, {
        duration: DEFAULT_DURATION,
        ...options,
      });
    }
  },
  warning: (message: string, options?: any) => {
    if (shouldShowMessage(message)) {
      toast.warning(message, {
        duration: DEFAULT_DURATION,
        ...options,
      });
    }
  },
};

/**
 * 检查是否应该显示消息
 * @param message 消息内容
 * @returns 是否应该显示
 */
function shouldShowMessage(message: string): boolean {
  const now = Date.now();
  const lastShown = recentMessages.get(message);

  // 如果消息最近没有显示过，或者已经过了过期时间，则显示
  if (!lastShown || now - lastShown > MESSAGE_EXPIRY_TIME) {
    recentMessages.set(message, now);

    // 清理过期的消息记录
    setTimeout(() => {
      recentMessages.delete(message);
    }, MESSAGE_EXPIRY_TIME);

    return true;
  }

  return false;
}
