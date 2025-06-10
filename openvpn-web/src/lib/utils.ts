import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";
import mammoth from "mammoth";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export const getPolicyFileUrl = (fileId: string): string => {
  const baseUrl = process.env.NEXT_PUBLIC_API_URL
  return `${baseUrl}/api/policy/download/${fileId}`;
};

/**
 * 解析 DOCX 文件并提取纯文本
 * @param fileOrUrl {File | string} - 可以是本地 File 对象或远程 URL
 * @returns {Promise<string>} - 返回提取的纯文本
 */
export async function convertDocxToText(
  fileOrUrl: File | string
): Promise<string> {
  let arrayBuffer: ArrayBuffer;

  try {
    if (fileOrUrl instanceof File) {
      // 📂 处理本地上传的 DOCX 文件
      arrayBuffer = await fileOrUrl.arrayBuffer();
    } else if (typeof fileOrUrl === "string") {
      // 🌍 处理远程 DOCX 文件（URL）
      const response = await fetch(fileOrUrl);
      if (!response.ok) throw new Error("无法加载 DOCX 文件");
      arrayBuffer = await response.arrayBuffer();
    } else {
      throw new Error("无效的输入类型");
    }

    // 📜 使用 mammoth 解析 DOCX 并提取纯文本
    const { value: text } = await mammoth.extractRawText({ arrayBuffer });
    return text.trim(); // 去掉首尾空格
  } catch (error) {
    console.error("DOCX 转换失败:", error);
    return "";
  }
}
