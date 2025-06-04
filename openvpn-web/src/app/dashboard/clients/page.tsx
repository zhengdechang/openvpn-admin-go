"use client";

import React, { useEffect } from "react";
import { useRouter } from "next/navigation";

export default function ClientsPage() {
  const router = useRouter();
  useEffect(() => {
    router.replace("/dashboard/users");
  }, [router]);
  return null;
}