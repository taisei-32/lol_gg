"use client";

import { useEffect, useState } from "react";
import Image from "next/image";
import Link from "next/link";
import styles from "./page.module.css";
import Nav from "@/components/Nav";

type Champion = {
    key: number;
    id: string;
    name: string;
    image_full: string;
    tags: string[];
};

type SortKey = "key" | "name" | "id";
type SortOrder = "asc" | "desc";

export default function ChampionsPage() {
    const [champions, setChampions] = useState<Champion[]>([]);
    const [version, setVersion] = useState<string>("");
    const [sortBy, setSortBy] = useState<SortKey>("key");
    const [order, setOrder] = useState<SortOrder>("asc");
    const [loading, setLoading] = useState(true);
    const [searchQuery, setSearchQuery] = useState("");

    useEffect(() => {
        Promise.all([
            fetch("http://localhost:8081/api/version").then((res) => res.json()),
            fetch("http://localhost:8081/api/champions").then((res) => res.json()),
        ]).then(([versionData, championsData]) => {
            setVersion(versionData.version);
            setChampions(championsData);
            setLoading(false);
        });
    }, []);

    const handleSort = (key: SortKey) => {
        if (sortBy === key) {
            setOrder(order === "asc" ? "desc" : "asc");
        } else {
            setSortBy(key);
            setOrder("asc");
        }
    };

    const getSortIcon = (key: SortKey) => {
        if (sortBy !== key) return "↕";
        return order === "asc" ? "↑" : "↓";
    };

    const filteredAndSorted = [...champions]
        .filter((champ) => {
            if (!searchQuery) return true;
            const q = searchQuery.toLowerCase();
            return (
                champ.name.toLowerCase().includes(q) ||
                champ.id.toLowerCase().includes(q) ||
                String(champ.key).includes(q)
            );
        })
        .sort((a, b) => {
            if (sortBy === "key") return order === "asc" ? a.key - b.key : b.key - a.key;
            const aVal = sortBy === "name" ? a.name : a.id;
            const bVal = sortBy === "name" ? b.name : b.id;
            return order === "asc" ? aVal.localeCompare(bVal, "ja") : bVal.localeCompare(aVal, "ja");
        });

    const TAG_LABELS: Record<string, string> = {
        Fighter: "ファイター",
        Tank: "タンク",
        Mage: "メイジ",
        Assassin: "アサシン",
        Support: "サポート",
        Marksman: "マークスマン",
    };

    const TAG_CLASS: Record<string, string> = {
        Fighter: styles.tagFighter,
        Tank: styles.tagTank,
        Mage: styles.tagMage,
        Assassin: styles.tagAssassin,
        Support: styles.tagSupport,
        Marksman: styles.tagMarksman,
    };

    return (
        <div className={styles.page}>
            <Nav
                searchQuery={searchQuery}
                onSearchChange={setSearchQuery}
                searchPlaceholder="チャンピオン名・IDで検索..."
            />

            <div className={styles.container}>
                <div className={styles.header}>
                    <h1 className={styles.title}>
                        <span className={styles.titleAccent}>CHAMPION</span>
                        {version && <span className={styles.version}>v{version}</span>}
                    </h1>
                </div>

                {loading ? (
                    <div className={styles.loading}>
                        <div className={styles.spinner} />
                        <p>Loading...</p>
                    </div>
                ) : (
                    <div className={styles.tableWrapper}>
                        <table className={styles.table}>
                            <thead>
                                <tr className={styles.headerRow}>
                                    <th className={styles.th} onClick={() => handleSort("key")} style={{ minWidth: 70, width: 70 }}>
                                        KEY {getSortIcon("key")}
                                    </th>
                                    <th className={styles.th} style={{ width: 56 }}>画像</th>
                                    <th className={styles.th} onClick={() => handleSort("name")} style={{ minWidth: 140 }}>
                                        名前 {getSortIcon("name")}
                                    </th>
                                    <th className={styles.th} onClick={() => handleSort("id")} style={{ minWidth: 140 }}>
                                        英語名 {getSortIcon("id")}
                                    </th>
                                    <th className={styles.th} style={{ minWidth: 160 }}>ロール</th>
                                </tr>
                            </thead>
                            <tbody>
                                {filteredAndSorted.map((champ) => (
                                    <tr key={champ.key} className={styles.row}>
                                        <td className={styles.td}>
                                            <span className={styles.keyBadge}>{champ.key}</span>
                                        </td>
                                        <td className={styles.td}>
                                            {version && (
                                                <Image
                                                    src={`https://ddragon.leagueoflegends.com/cdn/${version}/img/champion/${champ.image_full}`}
                                                    alt={champ.name}
                                                    width={40}
                                                    height={40}
                                                    className={styles.champIcon}
                                                    unoptimized
                                                />
                                            )}
                                        </td>
                                        <td className={styles.td}>
                                            <span className={styles.champName}>{champ.name}</span>
                                        </td>
                                        <td className={styles.td}>
                                            <span className={styles.champId}>{champ.id}</span>
                                        </td>
                                        <td className={styles.td}>
                                            <div className={styles.tagList}>
                                                {champ.tags?.map((tag) => (
                                                    <span key={tag} className={`${styles.tag} ${TAG_CLASS[tag] ?? styles.tagMage}`}>
                                                        {TAG_LABELS[tag] ?? tag}
                                                    </span>
                                                ))}
                                            </div>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                )}
            </div>
        </div>
    );
}