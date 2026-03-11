"use client";

import { useEffect, useState } from "react";
import Image from "next/image";
import Link from "next/link";
import styles from "./page.module.css";
import Nav from "@/components/Nav";

type ItemData = {
    id: number;
    name: string;
    plaintext: string;
    description: string;
    image_full: string;
    gold_base: number;
    gold_total: number;
    gold_sell: number;
    gold_purchasable: boolean;
    into_items: number[];
    from_items: number[];
    tags: string[];
    stats: Record<string, number>;
};

type SortKey = "id" | "name" | "gold_total" | "gold_sell";
type SortOrder = "asc" | "desc";

function stripHtml(html: string): string {
    return html.replace(/<[^>]*>/g, "").trim();
}

const STAT_LABELS: Record<string, string> = {
    FlatHPPoolMod: "HP",
    FlatMPPoolMod: "MP",
    FlatArmorMod: "アーマー",
    FlatSpellBlockMod: "魔法防御",
    FlatPhysicalDamageMod: "物理攻撃力",
    FlatMagicDamageMod: "魔法攻撃力",
    FlatMovementSpeedMod: "移動速度",
    PercentMovementSpeedMod: "移動速度%",
    FlatCritChanceMod: "クリ率",
    FlatAttackSpeedMod: "攻撃速度",
    PercentAttackSpeedMod: "攻撃速度%",
    FlatHPRegenMod: "HP回復",
    FlatMPRegenMod: "MP回復",
    PercentLifeStealMod: "ライフスティール",
    FlatEXPBonus: "経験値ボーナス",
};

function formatStatLabel(key: string): string {
    return STAT_LABELS[key] ?? key;
}

export default function ItemsPage() {
    const [items, setItems] = useState<ItemData[]>([]);
    const [itemMap, setItemMap] = useState<Record<number, ItemData>>({});
    const [version, setVersion] = useState<string>("");
    const [sortBy, setSortBy] = useState<SortKey>("id");
    const [order, setOrder] = useState<SortOrder>("asc");
    const [loading, setLoading] = useState(true);
    const [searchQuery, setSearchQuery] = useState("");

    useEffect(() => {
        Promise.all([
            fetch("http://localhost:8081/api/version").then((res) => res.json()),
            fetch("http://localhost:8081/api/items").then((res) => res.json()),
        ]).then(([versionData, itemsData]) => {
            setVersion(versionData.version);
            const raw: ItemData[] = Array.isArray(itemsData) ? itemsData : [];
            const seen = new Set<number>();
            const safeItems = raw.filter((item) => {
                if (seen.has(item.id)) return false;
                seen.add(item.id);
                return true;
            });
            setItems(safeItems);
            const map: Record<number, ItemData> = {};
            for (const item of safeItems) {
                map[item.id] = item;
            }
            setItemMap(map);
            setLoading(false);
        }).catch((err) => {
            console.error(err);
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

    const filteredAndSorted = [...items]
        .filter((item) => {
            if (!searchQuery) return true;
            const q = searchQuery.toLowerCase();
            return (
                stripHtml(item.name).toLowerCase().includes(q) ||
                String(item.id).includes(q)
            );
        })
        .sort((a, b) => {
            if (sortBy === "id") return order === "asc" ? a.id - b.id : b.id - a.id;
            if (sortBy === "gold_total")
                return order === "asc"
                    ? a.gold_total - b.gold_total
                    : b.gold_total - a.gold_total;
            if (sortBy === "gold_sell")
                return order === "asc"
                    ? a.gold_sell - b.gold_sell
                    : b.gold_sell - a.gold_sell;
            const aName = stripHtml(a.name);
            const bName = stripHtml(b.name);
            return order === "asc"
                ? aName.localeCompare(bName, "ja")
                : bName.localeCompare(aName, "ja");
        });

    const scrollToItem = (id: number) => {
        const el = document.getElementById(`item-${id}`);
        if (el) el.scrollIntoView({ behavior: "smooth", block: "center" });
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
                        <span className={styles.titleAccent}>ITEM</span>
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
                                    <th className={styles.th} onClick={() => handleSort("id")} style={{ minWidth: 70, width: 70 }}>
                                        KEY {getSortIcon("id")}
                                    </th>
                                    <th className={styles.th} onClick={() => handleSort("name")} style={{ minWidth: 180 }}>
                                        アイテム名 {getSortIcon("name")}
                                    </th>
                                    <th className={styles.th} onClick={() => handleSort("gold_total")}>
                                        合計価格 {getSortIcon("gold_total")}
                                    </th>
                                    <th className={styles.th} onClick={() => handleSort("gold_sell")}>
                                        売値 {getSortIcon("gold_sell")}
                                    </th>
                                    <th className={styles.th} style={{ minWidth: 180 }}>材料</th>
                                    <th className={styles.th} style={{ minWidth: 160 }}>能力</th>
                                    <th className={styles.th} style={{ minWidth: 200 }}>メモ</th>
                                    <th className={styles.th} style={{ minWidth: 180 }}>合成先</th>
                                </tr>
                            </thead>
                            <tbody>
                                {filteredAndSorted.map((item, rowIdx) => {
                                    const statEntries = Object.entries(item.stats ?? {}).filter(
                                        ([, v]) => v !== 0
                                    );
                                    return (
                                        <tr key={`row-${item.id}-${rowIdx}`} id={`item-${item.id}`} className={styles.row}>
                                            <td className={styles.td}>
                                                <span className={styles.keyBadge}>{item.id}</span>
                                            </td>

                                            <td className={styles.td}>
                                                <div className={styles.itemCell}>
                                                    {version && item.image_full && (
                                                        <Image
                                                            src={`https://ddragon.leagueoflegends.com/cdn/${version}/img/item/${item.image_full}`}
                                                            alt={stripHtml(item.name)}
                                                            width={40}
                                                            height={40}
                                                            className={styles.itemIcon}
                                                            unoptimized
                                                        />
                                                    )}
                                                    <span className={styles.itemName}>
                                                        {stripHtml(item.name)}
                                                    </span>
                                                </div>
                                            </td>

                                            <td className={styles.td}>
                                                <div className={styles.goldCell}>
                                                    <span className={styles.goldTotal}>
                                                        {item.gold_total > 0 ? `${item.gold_total}G` : "—"}
                                                    </span>
                                                    {item.gold_base > 0 && (
                                                        <span className={styles.goldBase}>
                                                            単体: {item.gold_base}G
                                                        </span>
                                                    )}
                                                </div>
                                            </td>

                                            <td className={styles.td}>
                                                <span className={styles.goldSell}>
                                                    {item.gold_sell > 0 ? `${item.gold_sell}G` : "—"}
                                                </span>
                                            </td>

                                            <td className={styles.td}>
                                                {item.from_items && item.from_items.length > 0 ? (
                                                    <div className={styles.componentList}>
                                                        {[...new Set(item.from_items)].map((fId, idx) => {
                                                            const fItem = itemMap[fId];
                                                            return (
                                                                <button
                                                                    key={`from-${fId}-${idx}`}
                                                                    className={styles.componentBtn}
                                                                    onClick={() => scrollToItem(fId)}
                                                                    title={fItem ? stripHtml(fItem.name) : String(fId)}
                                                                >
                                                                    {version && fItem?.image_full ? (
                                                                        <Image
                                                                            src={`https://ddragon.leagueoflegends.com/cdn/${version}/img/item/${fItem.image_full}`}
                                                                            alt={stripHtml(fItem.name)}
                                                                            width={32}
                                                                            height={32}
                                                                            className={styles.componentIcon}
                                                                            unoptimized
                                                                        />
                                                                    ) : (
                                                                        <span className={styles.unknownIcon}>?</span>
                                                                    )}

                                                                </button>
                                                            );
                                                        })}
                                                    </div>
                                                ) : (
                                                    <span className={styles.empty}>—</span>
                                                )}
                                            </td>

                                            <td className={styles.td}>
                                                {statEntries.length > 0 ? (
                                                    <ul className={styles.statList}>
                                                        {statEntries.map(([key, val]) => (
                                                            <li key={key} className={styles.statItem}>
                                                                <span className={styles.statLabel}>
                                                                    {formatStatLabel(key)}
                                                                </span>
                                                                <span className={styles.statValue}>
                                                                    +{val}
                                                                </span>
                                                            </li>
                                                        ))}
                                                    </ul>
                                                ) : (
                                                    <span className={styles.empty}>—</span>
                                                )}
                                            </td>

                                            <td className={styles.td}>
                                                <span className={styles.memo}>
                                                    {item.plaintext
                                                        ? stripHtml(item.plaintext)
                                                        : stripHtml(item.description).slice(0, 80) || "—"}
                                                </span>
                                            </td>

                                            <td className={styles.td}>
                                                {item.into_items && item.into_items.length > 0 ? (
                                                    <div className={styles.componentList}>
                                                        {[...new Set(item.into_items)].map((iId, idx) => {
                                                            const iItem = itemMap[iId];
                                                            return (
                                                                <button
                                                                    key={`into-${iId}-${idx}`}
                                                                    className={styles.componentBtn}
                                                                    onClick={() => scrollToItem(iId)}
                                                                    title={iItem ? stripHtml(iItem.name) : String(iId)}
                                                                >
                                                                    {version && iItem?.image_full ? (
                                                                        <Image
                                                                            src={`https://ddragon.leagueoflegends.com/cdn/${version}/img/item/${iItem.image_full}`}
                                                                            alt={stripHtml(iItem.name)}
                                                                            width={32}
                                                                            height={32}
                                                                            className={styles.componentIcon}
                                                                            unoptimized
                                                                        />
                                                                    ) : (
                                                                        <span className={styles.unknownIcon}>?</span>
                                                                    )}

                                                                </button>
                                                            );
                                                        })}
                                                    </div>
                                                ) : (
                                                    <span className={styles.empty}>—</span>
                                                )}
                                            </td>
                                        </tr>
                                    );
                                })}
                            </tbody>
                        </table>
                    </div>
                )}
            </div>
        </div>
    );
}