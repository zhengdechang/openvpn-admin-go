import axios from "axios";

/**
 * 从未知错误中提取可读的错误消息
 * @param error 捕获的未知错误
 * @param fallback 默认错误消息
 */
export function extractErrorMessage(error: unknown, fallback: string): string {
  if (axios.isAxiosError(error)) {
    return (
      error.response?.data?.error ||
      error.response?.data?.message ||
      error.message ||
      fallback
    );
  }
  if (error instanceof Error) {
    return error.message;
  }
  return fallback;
}
