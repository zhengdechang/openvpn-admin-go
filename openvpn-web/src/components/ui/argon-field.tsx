"use client";

import React, { useEffect, useRef, useState } from "react";

type ArgonFieldProps = React.InputHTMLAttributes<HTMLInputElement> & {
  label?: string;
  error?: string;
  icon?: React.ReactNode;
};

/**
 * 浅底下划线输入框，风格对齐登录页（argon）。供后台白卡/弹窗表单复用。
 * 透传所有原生 input 属性，支持 react-hook-form register（forwardRef）与受控 value/onChange。
 */
export const ArgonField = React.forwardRef<HTMLInputElement, ArgonFieldProps>(
  ({ label, error, icon, className, ...rest }, ref) => (
    <div className="argon-field">
      {label && <label className="argon-field-label">{label}</label>}
      <div className="argon-field-control">
        {icon && <span className="argon-field-icon">{icon}</span>}
        <input
          ref={ref}
          className={[
            "argon-field-input",
            icon ? "has-icon" : "",
            className || "",
          ]
            .filter(Boolean)
            .join(" ")}
          {...rest}
        />
        <span className="argon-field-underline" />
      </div>
      {error && <div className="argon-field-error">{error}</div>}
    </div>
  )
);
ArgonField.displayName = "ArgonField";

type ArgonSelectOption = { value: string; label: React.ReactNode; disabled?: boolean };

type ArgonSelectProps = {
  label?: string;
  error?: string;
  value?: string;
  onChange?: (e: { target: { value: string } }) => void;
  disabled?: boolean;
  className?: string;
  placeholder?: string;
  /** light=后台白卡（默认）；auth=登录/注册页磨砂面板（深色文字 + 左侧图标） */
  variant?: "light" | "auth";
  icon?: React.ReactNode;
  children?: React.ReactNode;
};

/**
 * 浅底下划线下拉框，风格对齐登录页（argon）。
 * 自定义 in-DOM 下拉（绝对定位于相对定位的容器内），默认在触发器正下方，
 * 下方空间不足时翻转到正上方。这样可避免 body{zoom:0.8} 破坏原生 select / MUI 弹层的定位。
 * API 兼容：接收 <option> 子元素，onChange 回调形如 { target: { value } }。
 */
export function ArgonSelect({
  label,
  error,
  value,
  onChange,
  disabled,
  className,
  placeholder,
  variant = "light",
  icon,
  children,
}: ArgonSelectProps) {
  const [open, setOpen] = useState(false);
  const [openUp, setOpenUp] = useState(false);
  const controlRef = useRef<HTMLDivElement>(null);
  const triggerRef = useRef<HTMLButtonElement>(null);

  const options: ArgonSelectOption[] = [];
  React.Children.forEach(children, (child) => {
    if (
      React.isValidElement(child) &&
      (child.type === "option" || (child.props as any)?.value !== undefined)
    ) {
      const props = child.props as {
        value?: string | number;
        children?: React.ReactNode;
        disabled?: boolean;
      };
      options.push({
        value: String(props.value ?? ""),
        label: props.children,
        disabled: props.disabled,
      });
    }
  });

  const selected = options.find((o) => o.value === String(value ?? ""));
  const selectedLabel = selected?.label;
  const isPlaceholder = !selected || selected.value === "";

  const computePlacement = () => {
    const trigger = triggerRef.current;
    if (!trigger) return;
    const rect = trigger.getBoundingClientRect();
    // 以最近的裁剪祖先（弹窗/滚动容器）为边界，而不是视口：
    // 触发器贴近弹窗底部时，视口下方虽空旷，菜单却会被弹窗裁掉，
    // 必须按弹窗边界判断下方空间是否够，不够就翻到上方。
    let clipBottom = window.innerHeight;
    let clipTop = 0;
    let node: HTMLElement | null = trigger.parentElement;
    while (node) {
      const oy = window.getComputedStyle(node).overflowY;
      if (oy === "auto" || oy === "scroll" || oy === "hidden") {
        const r = node.getBoundingClientRect();
        clipBottom = Math.min(clipBottom, r.bottom);
        clipTop = Math.max(clipTop, r.top);
      }
      node = node.parentElement;
    }
    const spaceBelow = clipBottom - rect.bottom;
    const spaceAbove = rect.top - clipTop;
    const menuH = Math.min(options.length * 36 + 8, 256);
    setOpenUp(spaceBelow < menuH && spaceAbove > spaceBelow);
  };

  const toggleOpen = () => {
    if (disabled) return;
    if (!open) computePlacement();
    setOpen((prev) => !prev);
  };

  const handleSelect = (opt: ArgonSelectOption) => {
    if (opt.disabled) return;
    onChange?.({ target: { value: opt.value } });
    setOpen(false);
    triggerRef.current?.focus();
  };

  useEffect(() => {
    if (!open) return;
    const onPointerDown = (e: MouseEvent) => {
      if (controlRef.current && !controlRef.current.contains(e.target as Node)) {
        setOpen(false);
      }
    };
    const onKeyDown = (e: KeyboardEvent) => {
      if (e.key === "Escape") setOpen(false);
    };
    document.addEventListener("mousedown", onPointerDown);
    document.addEventListener("keydown", onKeyDown);
    return () => {
      document.removeEventListener("mousedown", onPointerDown);
      document.removeEventListener("keydown", onKeyDown);
    };
  }, [open]);

  const isAuth = variant === "auth";

  const inner = (
    <>
      {icon && (
        <span className={isAuth ? "argon-input-icon" : "argon-field-icon"}>{icon}</span>
      )}
        <button
          type="button"
          ref={triggerRef}
          disabled={disabled}
          onClick={toggleOpen}
          aria-haspopup="listbox"
          aria-expanded={open}
          className={[
            isAuth ? "argon-input" : "argon-field-input",
            "argon-select-trigger",
            !isAuth && icon ? "has-icon" : "",
            isPlaceholder ? "is-placeholder" : "",
            className || "",
          ]
            .filter(Boolean)
            .join(" ")}
        >
          <span className="argon-select-value">
            {isPlaceholder ? placeholder || selectedLabel ||" " : selectedLabel}
          </span>
        </button>
        <span className={["argon-field-caret", open ? "is-open" : ""].filter(Boolean).join(" ")}>
          <svg
            width="16"
            height="16"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
          >
            <path d="m6 9 6 6 6-6" />
          </svg>
        </span>
        <span
          className={[
            isAuth ? "argon-input-underline" : "argon-field-underline",
            open ? "is-open" : "",
          ]
            .filter(Boolean)
            .join(" ")}
        />
        {open && (
          <ul
            className={["argon-select-menu", openUp ? "is-up" : "is-down"].join(" ")}
            role="listbox"
          >
            {options.map((opt, i) => (
              <li
                key={`${opt.value}-${i}`}
                role="option"
                aria-selected={opt.value === String(value ?? "")}
                onClick={() => handleSelect(opt)}
                className={[
                  "argon-select-option",
                  opt.value === String(value ?? "") ? "is-selected" : "",
                  opt.disabled ? "is-disabled" : "",
                ]
                  .filter(Boolean)
                  .join(" ")}
              >
                {opt.label}
              </li>
            ))}
          </ul>
        )}
    </>
  );

  if (isAuth) {
    return (
      <div className="argon-input-group" ref={controlRef}>
        {inner}
        {error && (
          <div
            className="argon-login-error"
            style={{ paddingTop: 4, paddingBottom: 0, textAlign: "left" }}
          >
            {error}
          </div>
        )}
      </div>
    );
  }

  return (
    <div className="argon-field">
      {label && <label className="argon-field-label">{label}</label>}
      <div className="argon-field-control" ref={controlRef}>
        {inner}
      </div>
      {error && <div className="argon-field-error">{error}</div>}
    </div>
  );
}
ArgonSelect.displayName = "ArgonSelect";
