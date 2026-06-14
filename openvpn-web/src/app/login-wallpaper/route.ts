import { NextResponse } from "next/server";

// 拉取必应每日壁纸（最近 8 天），服务端缓存 12 小时，供登录页随机选用。
// 仿 luci-theme-argon 的「每次进登录页背景都不同」效果。
const BING_API =
  "https://www.bing.com/HPImageArchive.aspx?format=js&idx=0&n=8&mkt=zh-CN";

export const revalidate = 43200; // 12h

export async function GET() {
  try {
    const res = await fetch(BING_API, { next: { revalidate } });
    if (!res.ok) throw new Error(`bing responded ${res.status}`);
    const data = (await res.json()) as { images?: { url?: string }[] };
    const urls = (data.images ?? [])
      .map((img) => img.url)
      .filter((u): u is string => Boolean(u))
      .map((u) => (u.startsWith("http") ? u : `https://www.bing.com${u}`));
    return NextResponse.json({ urls });
  } catch {
    // 网络不通时返回空数组，前端回退到本地自托管图。
    return NextResponse.json({ urls: [] });
  }
}
