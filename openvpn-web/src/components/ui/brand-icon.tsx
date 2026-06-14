import * as React from "react";

/**
 * Aegis 品牌图标：紫蓝渐变盾牌 + 白色钥匙孔（安全接入）。
 * 自带渐变填充，可直接作为 Logo 使用，无需外层容器。
 */
export function BrandIcon({
  size = 36,
  className,
  style,
  title = "Aegis",
}: {
  size?: number;
  className?: string;
  style?: React.CSSProperties;
  title?: string;
}) {
  return (
    <svg
      width={size}
      height={size}
      viewBox="0 0 40 40"
      fill="none"
      role="img"
      aria-label={title}
      className={className}
      style={style}
    >
      <defs>
        <linearGradient id="aegis-grad" x1="0" y1="0" x2="1" y2="1">
          <stop offset="0%" stopColor="#5e72e4" />
          <stop offset="100%" stopColor="#825ee4" />
        </linearGradient>
      </defs>
      {/* 盾牌主体 */}
      <path
        d="M20 3.2 L33.2 7.8 V18.5 C33.2 27.2 27.6 33.4 20 36.4 C12.4 33.4 6.8 27.2 6.8 18.5 V7.8 Z"
        fill="url(#aegis-grad)"
      />
      {/* 内描边高光 */}
      <path
        d="M20 6 L30.6 9.7 V18.5 C30.6 25.6 26 30.7 20 33.4 C14 30.7 9.4 25.6 9.4 18.5 V9.7 Z"
        fill="none"
        stroke="#ffffff"
        strokeOpacity="0.22"
        strokeWidth="1.1"
      />
      {/* 钥匙孔：环 + 锥形锁孔 */}
      <circle cx="20" cy="16.3" r="3.4" fill="#ffffff" />
      <path d="M18.1 18.4 L21.9 18.4 L23 26 L17 26 Z" fill="#ffffff" />
    </svg>
  );
}

export default BrandIcon;
