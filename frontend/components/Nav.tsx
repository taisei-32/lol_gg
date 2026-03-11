"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import styles from "./Nav.module.css";

type NavProps = {
    searchQuery?: string;
    onSearchChange?: (value: string) => void;
    searchPlaceholder?: string;
};

export default function Nav({ searchQuery, onSearchChange, searchPlaceholder }: NavProps) {
    const pathname = usePathname();

    return (
        <nav className={styles.nav}>
            <Link href="/" className={styles.navLogo}>LOL</Link>
            <div className={styles.navLinks}>
                <Link
                    href="/champions-id"
                    className={`${styles.navLink} ${pathname.startsWith("/champions") ? styles.navLinkActive : ""}`}
                >
                    チャンピオン
                </Link>
                <Link
                    href="/items-id"
                    className={`${styles.navLink} ${pathname.startsWith("/items") ? styles.navLinkActive : ""}`}
                >
                    アイテム
                </Link>
                <Link
                    href="/runes"
                    className={`${styles.navLink} ${pathname.startsWith("/items") ? styles.navLinkActive : ""}`}
                >
                    ルーン
                </Link>
            </div>
            {onSearchChange && (
                <input
                    className={styles.search}
                    type="text"
                    placeholder={searchPlaceholder ?? "検索..."}
                    value={searchQuery ?? ""}
                    onChange={(e) => onSearchChange(e.target.value)}
                />
            )}
        </nav>
    );
}