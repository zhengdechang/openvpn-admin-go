import React, { useEffect, useRef, useState } from "react";
import * as docx from "docx-preview";

interface DocxPreviewProps {
  file?: File;
  url?: string;
  onHtmlGenerated?: (html: string) => void;
}

export const DocxPreview: React.FC<DocxPreviewProps> = ({
  file,
  url,
  onHtmlGenerated,
}) => {
  const containerRef = useRef<HTMLDivElement | null>(null);
  const [docxHtml, setDocxHtml] = useState<string>("");

  useEffect(() => {
    const renderDocx = async () => {
      if (!containerRef.current) return;

      let arrayBuffer: ArrayBuffer | null = null;

      if (file) {
        arrayBuffer = await file.arrayBuffer();
      } else if (url) {
        try {
          const response = await fetch(url);
          if (!response.ok) throw new Error("无法加载文件");
          arrayBuffer = await response.arrayBuffer();
        } catch (error) {
          console.error("加载 DOCX 失败:", error);
          return;
        }
      }

      if (!arrayBuffer || !containerRef.current) return;

      containerRef.current.innerHTML = "";

      await docx.renderAsync(arrayBuffer, containerRef.current, undefined, {
        className: "docx",
        inWrapper: true,
        hideWrapperOnPrint: false,
        ignoreWidth: false,
        ignoreHeight: false,
        ignoreFonts: false,
        breakPages: true,
        ignoreLastRenderedPageBreak: true,
        experimental: false,
        trimXmlDeclaration: true,
        useBase64URL: false,
        renderChanges: false,
        renderHeaders: true,
        renderFooters: true,
        renderFootnotes: true,
        renderEndnotes: true,
        renderComments: false,
        renderAltChunks: true,
        debug: false,
      });

      const docxWrapper = containerRef.current.querySelector(".docx-wrapper");
      if (docxWrapper) {
        const generatedHtml = docxWrapper.innerHTML;
        setDocxHtml(generatedHtml);

        if (onHtmlGenerated) {
          onHtmlGenerated(generatedHtml);
        }
      }
    };

    renderDocx();
  }, [file, url]);

  return (
    <div className="bg-white shadow-md rounded-lg p-6 border border-gray-300">
      <div ref={containerRef} className="hidden w-full" />

      {/* ✅ 让 `docxHtml` 内容居中 */}
      <div className="min-h-[600px] flex justify-center items-center">
        <div
          dangerouslySetInnerHTML={{
            __html: docxHtml,
          }}
        ></div>
      </div>
    </div>
  );
};
