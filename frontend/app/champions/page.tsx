"use client";

import { useEffect, useState } from "react";
import Image from "next/image";
import styles from "./page.module.css";

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

    useEffect(() => {
        // バージョンとチャンピオンを並列取得
        Promise.all([
            fetch("http://localhost:8081/api/version").then((res) => res.json()),
            fetch("http://localhost:8081/api/champions").then((res) => res.json()),
        ]).then(([versionData, championsData]) => {
            setVersion(versionData.version);
            setChampions(championsData);
            setLoading(false);
        });
    }, []);

    const sortedChampions = [...champions].sort((a, b) => {
        if (sortBy === "key") {
            return order === "asc" ? a.key - b.key : b.key - a.key;
        } else {
            const aVal = sortBy === "name" ? a.name : a.id;
            const bVal = sortBy === "name" ? b.name : b.id;
            return order === "asc"
                ? aVal.localeCompare(bVal)
                : bVal.localeCompare(aVal);
        }
    });

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

    return (
        <div style={{ padding: "24px" }}>
            <h1>チャンピオン一覧({version})</h1>
            {loading ? (
                <p>読み込み中...</p>
            ) : (
                <table style={{ borderCollapse: "collapse", width: "100%" }}>
                    <thead>
                        <tr style={{ backgroundColor: "#f0f0f0" }}>
                            <th onClick={() => handleSort("key")} className={styles.th}>
                                Key {getSortIcon("key")}
                            </th>
                            <td className={styles.td}>画像</td>
                            <th onClick={() => handleSort("name")} className={styles.th}>
                                名前（日本語）{getSortIcon("name")}
                            </th>
                            <th onClick={() => handleSort("id")} className={styles.th}>
                                名前（英語）{getSortIcon("id")}
                            </th>
                            <th className={styles.th}>ロール</th>
                        </tr>
                    </thead>
                    <tbody>
                        {sortedChampions.map((champ) => (
                            <tr key={champ.key} style={{ borderBottom: "1px solid #ddd" }}>
                                <td className={styles.th}>{champ.key}</td>
                                <td className={styles.th}>
                                    {version && (
                                        <Image
                                            src={`https://ddragon.leagueoflegends.com/cdn/${version}/img/champion/${champ.image_full}`}
                                            alt={champ.name}
                                            width={48}
                                            height={48}
                                        />
                                    )}
                                </td>
                                <td className={styles.th}>{champ.name}</td>
                                <td className={styles.th}>{champ.id}</td>
                                <td className={styles.th}>{champ.tags?.join(", ")}</td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            )}
        </div>
    );
}