/*
 * @Description:
 * @Author: Devin
 * @Date: 2025-07-02 10:39:45
 */
"use client";

import React, { useState, useEffect, useCallback, useRef } from "react";
import { openvpnAPI } from "@/services/api";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";

interface SimpleLogViewerProps {
  className?: string;
  height?: number;
}

const SimpleLogViewer: React.FC<SimpleLogViewerProps> = ({
  className = "",
  height = 400,
}) => {
  const { t } = useTranslation();
  const [logs, setLogs] = useState<string>("");
  const [loading, setLoading] = useState(true);
  const [loadingMore, setLoadingMore] = useState(false);
  const [totalLines, setTotalLines] = useState(0);
  const [currentOffset, setCurrentOffset] = useState(0);
  const [hasMore, setHasMore] = useState(true);
  const scrollRef = useRef<HTMLDivElement>(null);

  const PAGE_SIZE = 100; // 每次加载100行

  // 初始加载前100行日志
  useEffect(() => {
    const fetchInitialLogs = async () => {
      setLoading(true);
      try {
        // 从第1行开始加载100行
        const response = await openvpnAPI.getClientLogs(0, PAGE_SIZE);
        setLogs(response.logs || "");
        setTotalLines(response.totalLines);
        setCurrentOffset(PAGE_SIZE);
        setHasMore(response.hasMore);
      } catch (error) {
        toast.error(t("dashboard.logs.fetchClientLogsError"));
      } finally {
        setLoading(false);
      }
    };

    fetchInitialLogs();
  }, [t]);

  // 加载更多日志
  const loadMoreLogs = useCallback(async () => {
    if (loadingMore || !hasMore) return;

    setLoadingMore(true);
    try {
      // 加载下一页日志
      const response = await openvpnAPI.getClientLogs(currentOffset, PAGE_SIZE);

      if (response.logs) {
        // 将新日志添加到现有日志后面
        setLogs((prevLogs) => prevLogs + "\n" + response.logs);
        setCurrentOffset(currentOffset + PAGE_SIZE);
        setHasMore(response.hasMore);
      } else {
        setHasMore(false);
      }
    } catch (error) {
      toast.error(t("dashboard.logs.fetchClientLogsError"));
    } finally {
      setLoadingMore(false);
    }
  }, [currentOffset, hasMore, loadingMore, t]);

  // 处理滚动事件
  const handleScroll = useCallback(
    (e: React.UIEvent<HTMLDivElement>) => {
      const { scrollTop, scrollHeight, clientHeight } = e.currentTarget;

      // 当滚动到底部附近时加载更多
      if (
        scrollTop + clientHeight >= scrollHeight - 100 &&
        hasMore &&
        !loadingMore
      ) {
        loadMoreLogs();
      }
    },
    [hasMore, loadingMore, loadMoreLogs]
  );

  if (loading) {
    return (
      <div
        className={`flex items-center justify-center ${className}`}
        style={{ height }}
      >
        <p className="text-gray-600 dark:text-gray-400">
          {t("common.loading")}
        </p>
      </div>
    );
  }

  return (
    <div
      className={`bg-gray-100 dark:bg-gray-800 rounded-md p-4 ${className}`}
      style={{ height }}
    >
      <div
        ref={scrollRef}
        className="h-full overflow-auto"
        onScroll={handleScroll}
      >
        <pre className="whitespace-pre-wrap break-all">
          {logs || t("dashboard.logs.noClientLogs")}
        </pre>
        {loadingMore && hasMore && (
          <div className="text-center py-2">
            <span className="text-xs text-gray-500 dark:text-gray-400">
              {t("common.loading")}...
            </span>
          </div>
        )}
        {!hasMore && logs && (
          <div className="text-center py-2 border-t border-gray-300 dark:border-gray-600 mt-2">
            <span className="text-xs text-gray-500 dark:text-gray-400">
              已显示全部日志 (共 {totalLines} 行)
            </span>
          </div>
        )}
      </div>
    </div>
  );
};

export default SimpleLogViewer;
