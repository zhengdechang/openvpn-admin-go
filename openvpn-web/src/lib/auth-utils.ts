import { cookies } from 'next/headers';

/**
 * 从cookies中获取认证令牌
 * @returns {Promise<string|null>} 认证令牌或null
 */
export async function getAuthToken(): Promise<string | null> {
  const cookieStore = cookies();
  const token = cookieStore.get('auth_token');
  return token ? token.value : null;
}

/**
 * 从客户端获取认证令牌
 * @returns {string|null} 认证令牌或null
 */
export function getClientAuthToken(): string | null {
  if (typeof window === 'undefined') return null;
  return document.cookie
    .split('; ')
    .find(row => row.startsWith('auth_token='))
    ?.split('=')[1] || null;
} 