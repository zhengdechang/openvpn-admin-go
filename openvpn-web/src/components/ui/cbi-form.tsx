import * as React from "react";
import { cn } from "@/lib/utils";

/**
 * LuCI Argon 风格的 CBI 表单原语。
 * - CbiSection：带标题的白色面板（对应 .cbi-section / .panel-title）
 * - CbiValue：左标签（右对齐）+ 右控件的表单行（对应 .cbi-value / .cbi-value-title）
 */

interface CbiSectionProps {
  title?: React.ReactNode;
  actions?: React.ReactNode;
  className?: string;
  bodyClassName?: string;
  children: React.ReactNode;
}

export function CbiSection({
  title,
  actions,
  className,
  bodyClassName,
  children,
}: CbiSectionProps) {
  return (
    <section className={cn("cbi-panel", className)}>
      {(title || actions) && (
        <div className="cbi-panel-title">
          <span>{title}</span>
          {actions ? <span className="flex items-center gap-2">{actions}</span> : null}
        </div>
      )}
      <div className={cn("cbi-panel-body", bodyClassName)}>{children}</div>
    </section>
  );
}

interface CbiValueProps {
  title?: React.ReactNode;
  htmlFor?: string;
  className?: string;
  children: React.ReactNode;
}

export function CbiValue({ title, htmlFor, className, children }: CbiValueProps) {
  return (
    <div className={cn("cbi-value", className)}>
      <label className="cbi-value-title" htmlFor={htmlFor}>
        {title}
      </label>
      <div className="cbi-value-field">{children}</div>
    </div>
  );
}
