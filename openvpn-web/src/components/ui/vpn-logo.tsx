import React from "react";

interface VPNLogoProps {
  size?: number;
  className?: string;
}

/**
 * VPN Admin Logo — shield + lock + network nodes + green check badge
 * Inline SVG so it scales and inherits context without extra HTTP requests.
 */
export default function VPNLogo({ size = 40, className }: VPNLogoProps) {
  const id = React.useId().replace(/:/g, "");
  return (
    <svg
      viewBox="0 0 80 80"
      xmlns="http://www.w3.org/2000/svg"
      fill="none"
      width={size}
      height={size}
      className={className}
      aria-label="VPN Admin Logo"
      role="img"
    >
      <defs>
        <linearGradient id={`sg-${id}`} x1="0%" y1="0%" x2="100%" y2="100%">
          <stop offset="0%" stopColor="#012a4a" />
          <stop offset="60%" stopColor="#0369a1" />
          <stop offset="100%" stopColor="#0ea5e9" />
        </linearGradient>
        <linearGradient id={`gg-${id}`} x1="0%" y1="0%" x2="100%" y2="100%">
          <stop offset="0%" stopColor="#16a34a" />
          <stop offset="100%" stopColor="#22c55e" />
        </linearGradient>
        <filter id={`sh-${id}`} x="-10%" y="-10%" width="120%" height="130%">
          <feDropShadow dx="0" dy="3" stdDeviation="3" floodColor="#012a4a" floodOpacity="0.3" />
        </filter>
      </defs>

      {/* Shield */}
      <path
        d="M40 4 L10 16 V42 C10 58 23 71 40 76 C57 71 70 58 70 42 V16 Z"
        fill={`url(#sg-${id})`}
        filter={`url(#sh-${id})`}
      />

      {/* Inner shimmer */}
      <path
        d="M40 11 L16 21.5 V42 C16 55 27 66 40 70 C53 66 64 55 64 42 V21.5 Z"
        fill="white"
        opacity="0.07"
      />

      {/* Network line */}
      <line x1="24" y1="34" x2="56" y2="34" stroke="white" strokeWidth="1.5" opacity="0.3" />

      {/* Lock body */}
      <rect x="27" y="41" width="26" height="20" rx="4" fill="white" />

      {/* Lock shackle */}
      <path
        d="M31.5 41 V33 C31.5 24.5 48.5 24.5 48.5 33 V41"
        stroke="white"
        strokeWidth="3.5"
        strokeLinecap="round"
        strokeLinejoin="round"
      />

      {/* Keyhole */}
      <circle cx="40" cy="51.5" r="3.5" fill="#0369a1" />
      <rect x="38.5" y="53.5" width="3" height="4.5" rx="1" fill="#0369a1" />

      {/* Left node */}
      <circle cx="21" cy="34" r="3.5" fill="#0ea5e9" opacity="0.9" />
      <circle cx="21" cy="34" r="1.5" fill="white" />

      {/* Right node */}
      <circle cx="59" cy="34" r="3.5" fill="#0ea5e9" opacity="0.9" />
      <circle cx="59" cy="34" r="1.5" fill="white" />

      {/* Dashed connectors */}
      <line x1="21" y1="37.5" x2="30" y2="42" stroke="white" strokeWidth="1" opacity="0.3" strokeDasharray="2,2" />
      <line x1="59" y1="37.5" x2="50" y2="42" stroke="white" strokeWidth="1" opacity="0.3" strokeDasharray="2,2" />

      {/* Green badge */}
      <circle cx="58" cy="62" r="11" fill={`url(#gg-${id})`} stroke="white" strokeWidth="2" />
      <path d="M53 62 L57.5 66.5 L64 58" stroke="white" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  );
}
