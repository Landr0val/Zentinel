"use client";

import { useEffect, useState } from "react";
import { useRouter, usePathname } from "next/navigation";

export default function AuthGuard({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const pathname = usePathname();
  const [authorized, setAuthorized] = useState(false);

  useEffect(() => {
    if (pathname === "/login") {
      return;
    }

    const token = localStorage.getItem("token");

    if (!token) {
      setTimeout(() => setAuthorized(false), 0);
      router.push("/login");
    } else {
      setTimeout(() => setAuthorized(true), 0);
    }
  }, [pathname, router]);

  if (pathname === "/login") {
    return <>{children}</>;
  }

  if (!authorized) {
    return null;
  }

  return <>{children}</>;
}
