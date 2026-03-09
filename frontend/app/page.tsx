"use client";

import { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import styles from "./page.module.css";

export default function Home() {
  const [query, setQuery] = useState("");
  const router = useRouter();

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    if (!query.trim()) return;
    const [name, tag] = query.split("#");
    if (name && tag) {
      router.push(`/summoner/${encodeURIComponent(name)}-${encodeURIComponent(tag)}`);
    }
  };

  return (
    <div className={styles.page}>
      <nav className={styles.nav}>
        <Link href="/" className={styles.navLogo}>
          LOL
        </Link>
        <div className={styles.navLinks}>
          <Link href="/champions" className={styles.navLink}>チャンピオン</Link>
        </div>
      </nav>

      <main className={styles.hero}>
        <h1 className={styles.title}>League of Legends </h1>
        <p className={styles.sub}>戦績検索・チャンピオン統計</p>

        <div className={styles.searchWrapper}>
          <form onSubmit={handleSearch}>
            <div className={styles.searchBox}>
              <span className={styles.searchIcon}>⚔</span>
              <input
                className={styles.searchInput}
                type="text"
                placeholder="サモナー名を検索 サモナー名#JP1"
                value={query}
                onChange={(e) => setQuery(e.target.value)}
                autoComplete="off"
              />
              <button className={styles.searchBtn} type="submit">
                検索
              </button>
            </div>
          </form>
        </div>
      </main>
    </div>
  );
}