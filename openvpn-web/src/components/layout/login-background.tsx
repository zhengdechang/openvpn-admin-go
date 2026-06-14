"use client";

import React from "react";

/**
 * 登录页风景大图背景。
 * 仿 luci-theme-argon：从必应每日壁纸里随机取最多 4 张组成轮动池，
 * 进页面后每隔若干秒随机切到不同的一张（带预加载，避免空白闪烁）；
 * 网络不通时回退到本地自托管图 public/login-bg.jpg。
 * 壁纸列表由 /login-wallpaper 路由在服务端拉取并缓存。
 */
const LOCAL_FALLBACK = "/login-bg.jpg";
const ROTATE_COUNT = 4; // 轮动图片数量
const ROTATE_INTERVAL = 8000; // 切换间隔(ms)

function shuffle<T>(arr: T[]): T[] {
  const a = [...arr];
  for (let i = a.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1));
    [a[i], a[j]] = [a[j], a[i]];
  }
  return a;
}

export default function LoginBackground() {
  const [src, setSrc] = React.useState(LOCAL_FALLBACK);

  React.useEffect(() => {
    let cancelled = false;
    let timer: ReturnType<typeof setInterval> | undefined;

    fetch("/login-wallpaper")
      .then((r) => r.json())
      .then((data: { urls?: string[] }) => {
        if (cancelled) return;
        const all = data?.urls ?? [];
        if (all.length === 0) return;

        // 随机取最多 4 张作为轮动池
        const pool = shuffle(all).slice(0, ROTATE_COUNT);
        let idx = 0;

        // 预加载成功再切换，避免空白闪烁；失败则保留当前图。
        const show = (url: string) => {
          const img = new window.Image();
          img.onload = () => {
            if (!cancelled) setSrc(url);
          };
          img.src = url;
        };

        show(pool[idx]); // 立即展示第一张

        if (pool.length > 1) {
          timer = setInterval(() => {
            if (cancelled) return;
            // 随机切到与当前不同的一张
            let next = idx;
            while (next === idx) {
              next = Math.floor(Math.random() * pool.length);
            }
            idx = next;
            show(pool[idx]);
          }, ROTATE_INTERVAL);
        }
      })
      .catch(() => {
        /* 保留本地回退图 */
      });

    return () => {
      cancelled = true;
      if (timer) clearInterval(timer);
    };
  }, []);

  return (
    <div aria-hidden="true" className="argon-login-bg">
      <img
        src={src}
        alt=""
        style={{
          position: "absolute",
          inset: 0,
          height: "100%",
          width: "100%",
          objectFit: "cover",
          objectPosition: "center",
          transition: "opacity 0.4s ease",
        }}
      />
      {/* 柔和暗角，提升整体质感与左侧面板可读性 */}
      <div
        style={{
          position: "absolute",
          inset: 0,
          background:
            "linear-gradient(90deg, rgba(23,43,77,0.18) 0%, rgba(23,43,77,0.04) 32%, rgba(23,43,77,0) 60%)",
        }}
      />
    </div>
  );
}
