"use client";

import { useEffect, useState } from "react";
import Image from "next/image";
import Nav from "@/components/Nav";
import styles from "./page.module.css";

const RUNE_IMG_BASE = "https://ddragon.leagueoflegends.com/cdn/img/";

type RuneData = {
    id: number;
    key: string;
    name: string;
    icon: string;
    short_desc: string;
    long_desc: string;
};

type RuneSlot = {
    slot: number;
    runes: RuneData[];
};

type RuneStyle = {
    id: number;
    key: string;
    name: string;
    icon: string;
    slots: RuneSlot[];
};

const STAT_RUNES = [
    [
        { id: 5008, name: "アダプティブフォース", icon: "perk-images/StatMods/StatModsAdaptiveForceIcon.png" },
        { id: 5005, name: "攻撃速度", icon: "perk-images/StatMods/StatModsAttackSpeedIcon.png" },
        { id: 5007, name: "スキルヘイスト", icon: "perk-images/StatMods/StatModsCDRScalingIcon.png" },
    ],
    [
        { id: 5008, name: "アダプティブフォース", icon: "perk-images/StatMods/StatModsAdaptiveForceIcon.png" },
        { id: 5010, name: "移動速度", icon: "perk-images/StatMods/StatModsMovementSpeedIcon.png" },
        { id: 5001, name: "体力の伸び", icon: "perk-images/StatMods/StatModsHealthPlusIcon.png" },
    ],
    [
        { id: 5011, name: "体力", icon: "perk-images/StatMods/StatModsHealthScalingIcon.png" },
        { id: 5013, name: "行動妨害&スロウ耐性", icon: "perk-images/StatMods/StatModsTenacityIcon.png" },
        { id: 5002, name: "体力の伸び", icon: "perk-images/StatMods/StatModsHealthPlusIcon.png" },
    ],
];

const STAT_ROW_LABELS = ["攻撃枠", "フレックス枠", "防御枠"];

type BuildState = {
    mainStyle: RuneStyle | null;
    keystone: RuneData | null;
    mainSlots: (RuneData | null)[];
    subStyle: RuneStyle | null;
    subRunes: (RuneData | null)[];
    stats: (number | null)[];
};

function stripTags(html: string): string {
    return html
        .replace(/<br\s*\/?>/gi, "\n")
        .replace(/<[^>]+>/g, "")
        .replace(/&lt;/g, "<").replace(/&gt;/g, ">")
        .trim();
}

export default function RunesPage() {
    const [allStyles, setAllStyles] = useState<RuneStyle[]>([]);
    const [loading, setLoading] = useState(true);
    const [modal, setModal] = useState<RuneData | null>(null);
    const [activeTab, setActiveTab] = useState<"list" | "build">("list");
    const [listStyle, setListStyle] = useState<RuneStyle | null>(null);

    const [build, setBuild] = useState<BuildState>({
        mainStyle: null, keystone: null, mainSlots: [null, null, null],
        subStyle: null, subRunes: [null, null], stats: [null, null, null],
    });

    useEffect(() => {
        fetch("http://localhost:8081/api/runes")
            .then((r) => r.json())
            .then((data: RuneStyle[]) => {
                setAllStyles(data);
                setListStyle(data[0] ?? null);
                setBuild((b) => ({ ...b, mainStyle: data[0] ?? null }));
                setLoading(false);
            });
    }, []);

    const selectMainStyle = (s: RuneStyle) => setBuild((b) => ({
        ...b, mainStyle: s, keystone: null, mainSlots: [null, null, null],
        subStyle: b.subStyle?.id === s.id ? null : b.subStyle,
        subRunes: b.subStyle?.id === s.id ? [null, null] : b.subRunes,
    }));

    const selectKeystone = (rune: RuneData) =>
        setBuild((b) => ({ ...b, keystone: b.keystone?.id === rune.id ? null : rune }));

    const selectMainSlot = (slotIdx: number, rune: RuneData) =>
        setBuild((b) => {
            const next = [...b.mainSlots];
            next[slotIdx - 1] = next[slotIdx - 1]?.id === rune.id ? null : rune;
            return { ...b, mainSlots: next };
        });

    const selectSubStyle = (s: RuneStyle) => {
        if (s.id === build.mainStyle?.id) return;
        setBuild((b) => ({ ...b, subStyle: b.subStyle?.id === s.id ? null : s, subRunes: [null, null] }));
    };

    const selectSubRune = (rune: RuneData, slotNum: number) =>
        setBuild((b) => {
            const next = [...b.subRunes];
            const existIdx = next.findIndex((r) => r?.id === rune.id);
            if (existIdx !== -1) { next[existIdx] = null; return { ...b, subRunes: next }; }
            const sameSlotIdx = next.findIndex((r) => {
                if (!r || !b.subStyle) return false;
                return b.subStyle.slots.find((sl) => sl.runes.some((ru) => ru.id === r.id))?.slot === slotNum;
            });
            if (sameSlotIdx !== -1) { next[sameSlotIdx] = rune; return { ...b, subRunes: next }; }
            const emptyIdx = next.findIndex((r) => r === null);
            if (emptyIdx !== -1) { next[emptyIdx] = rune; return { ...b, subRunes: next }; }
            next[0] = next[1]; next[1] = rune;
            return { ...b, subRunes: next };
        });

    const selectStat = (row: number, col: number) =>
        setBuild((b) => {
            const next = [...b.stats];
            next[row] = next[row] === col ? null : col;
            return { ...b, stats: next };
        });

    const isSubRuneSelected = (rune: RuneData) => build.subRunes.some((r) => r?.id === rune.id);
    const subSlots = build.subStyle?.slots.filter((s) => s.slot > 0) ?? [];

    if (loading) return (
        <div className={styles.page}>
            <Nav />
            <div className={styles.loading}><div className={styles.spinner} /><p>Loading...</p></div>
        </div>
    );

    return (
        <div className={styles.page}>
            <Nav />
            <div className={styles.container}>

                <div className={styles.pageHeader}>
                    <h1 className={styles.title}><span className={styles.titleAccent}>RUNES</span></h1>
                    <div className={styles.tabs}>
                        <button className={`${styles.tab} ${activeTab === "list" ? styles.tabActive : ""}`} onClick={() => setActiveTab("list")}>一覧</button>
                        <button className={`${styles.tab} ${activeTab === "build" ? styles.tabActive : ""}`} onClick={() => setActiveTab("build")}>組み立て</button>
                    </div>
                </div>

                {activeTab === "list" && (
                    <div className={styles.layout}>
                        <div className={styles.styleTabs}>
                            {allStyles.map((s) => (
                                <button key={s.id}
                                    className={`${styles.styleTab} ${listStyle?.id === s.id ? styles.styleTabActive : ""}`}
                                    onClick={() => setListStyle(s)}
                                >
                                    <Image src={`${RUNE_IMG_BASE}${s.icon}`} alt={s.name} width={36} height={36} unoptimized className={styles.styleIcon} />
                                    <span>{s.name}</span>
                                </button>
                            ))}
                        </div>

                        {listStyle && (
                            <div className={styles.slotsArea}>
                                <div className={styles.styleHeader}>
                                    <Image src={`${RUNE_IMG_BASE}${listStyle.icon}`} alt={listStyle.name} width={48} height={48} unoptimized className={styles.styleHeaderIcon} />
                                    <h2 className={styles.styleName}>{listStyle.name}</h2>
                                </div>
                                {(listStyle.slots ?? []).map((slot) => (
                                    <div key={slot.slot} className={styles.slotRow}>
                                        <div className={`${styles.runeList} ${slot.slot === 0 ? styles.keystoneList : ""}`}>
                                            {(slot.runes ?? []).map((rune) => (
                                                <button key={rune.id}
                                                    className={`${styles.runeBtn} ${slot.slot === 0 ? styles.keystoneBtn : ""}`}
                                                    onClick={() => setModal(rune)}
                                                    title={rune.name}
                                                >
                                                    <Image src={`${RUNE_IMG_BASE}${rune.icon}`} alt={rune.name} width={slot.slot === 0 ? 64 : 44} height={slot.slot === 0 ? 64 : 44} unoptimized className={styles.runeIcon} />
                                                    <span className={styles.runeName}>{rune.name}</span>
                                                </button>
                                            ))}
                                        </div>
                                    </div>
                                ))}
                            </div>
                        )}
                    </div>
                )}

                {activeTab === "build" && (
                    <div className={styles.outerLayout}>
                        <div className={styles.panel}>
                            <div className={styles.panelLabel}>メイン</div>
                            <div className={styles.styleRow}>
                                {allStyles.map((s) => (
                                    <button key={s.id}
                                        className={`${styles.styleCircle} ${build.mainStyle?.id === s.id ? styles.styleCircleActive : ""}`}
                                        onClick={() => selectMainStyle(s)} title={s.name}
                                    >
                                        <Image src={`${RUNE_IMG_BASE}${s.icon}`} alt={s.name} width={32} height={32} unoptimized />
                                    </button>
                                ))}
                            </div>
                            {build.mainStyle && <>
                                {build.mainStyle.slots.filter((sl) => sl.slot === 0).map((sl) => (
                                    <div key={sl.slot} className={styles.slotSection}>
                                        <div className={styles.runeRow}>
                                            {(sl.runes ?? []).map((rune) => (
                                                <button key={rune.id}
                                                    className={`${styles.runeCircle} ${styles.keystoneCircle} ${build.keystone?.id === rune.id ? styles.runeSelected : styles.runeUnselected}`}
                                                    onClick={() => selectKeystone(rune)}
                                                    onContextMenu={(e) => { e.preventDefault(); setModal(rune); }}
                                                    title={`${rune.name}（右クリックで説明）`}
                                                >
                                                    <Image src={`${RUNE_IMG_BASE}${rune.icon}`} alt={rune.name} width={56} height={56} unoptimized />
                                                </button>
                                            ))}
                                        </div>
                                    </div>
                                ))}
                                {build.mainStyle.slots.filter((sl) => sl.slot > 0).map((sl) => (
                                    <div key={sl.slot} className={styles.slotSection}>
                                        <div className={styles.runeRow}>
                                            {(sl.runes ?? []).map((rune) => (
                                                <button key={rune.id}
                                                    className={`${styles.runeCircle} ${build.mainSlots[sl.slot - 1]?.id === rune.id ? styles.runeSelected : styles.runeUnselected}`}
                                                    onClick={() => selectMainSlot(sl.slot, rune)}
                                                    onContextMenu={(e) => { e.preventDefault(); setModal(rune); }}
                                                    title={`${rune.name}（右クリックで説明）`}
                                                >
                                                    <Image src={`${RUNE_IMG_BASE}${rune.icon}`} alt={rune.name} width={40} height={40} unoptimized />
                                                </button>
                                            ))}
                                        </div>
                                    </div>
                                ))}
                            </>}
                        </div>

                        <div className={styles.rightCol}>
                            <div className={styles.panel}>
                                <div className={styles.panelLabel}>サブ</div>
                                <div className={styles.styleRow}>
                                    {allStyles.map((s) => (
                                        <button key={s.id}
                                            className={`${styles.styleCircle} ${build.subStyle?.id === s.id ? styles.styleCircleActive : ""} ${s.id === build.mainStyle?.id ? styles.styleCircleDisabled : ""}`}
                                            onClick={() => selectSubStyle(s)}
                                            disabled={s.id === build.mainStyle?.id}
                                            title={s.id === build.mainStyle?.id ? "メインと同じスタイルは選択不可" : s.name}
                                        >
                                            <Image src={`${RUNE_IMG_BASE}${s.icon}`} alt={s.name} width={32} height={32} unoptimized />
                                        </button>
                                    ))}
                                </div>
                                {build.subStyle && subSlots.map((sl) => (
                                    <div key={sl.slot} className={styles.slotSection}>
                                        <div className={styles.runeRow}>
                                            {(sl.runes ?? []).map((rune) => (
                                                <button key={rune.id}
                                                    className={`${styles.runeCircle} ${isSubRuneSelected(rune) ? styles.runeSelected : styles.runeUnselected}`}
                                                    onClick={() => selectSubRune(rune, sl.slot)}
                                                    onContextMenu={(e) => { e.preventDefault(); setModal(rune); }}
                                                    title={`${rune.name}（右クリックで説明）`}
                                                >
                                                    <Image src={`${RUNE_IMG_BASE}${rune.icon}`} alt={rune.name} width={40} height={40} unoptimized />
                                                </button>
                                            ))}
                                        </div>
                                    </div>
                                ))}
                            </div>

                            <div className={styles.panel}>
                                <div className={styles.panelLabel}>ステータス</div>
                                {STAT_RUNES.map((row, rowIdx) => (
                                    <div key={rowIdx} className={styles.statSection}>
                                        <div className={styles.statRowLabel}>{STAT_ROW_LABELS[rowIdx]}</div>
                                        <div className={styles.statRow}>
                                            {row.map((stat, colIdx) => (
                                                <button key={colIdx}
                                                    className={`${styles.statCircle} ${build.stats[rowIdx] === colIdx ? styles.statSelected : styles.statUnselected}`}
                                                    onClick={() => selectStat(rowIdx, colIdx)}
                                                    title={stat.name}
                                                >
                                                    <Image src={`${RUNE_IMG_BASE}${stat.icon}`} alt={stat.name} width={28} height={28} unoptimized />
                                                </button>
                                            ))}
                                        </div>
                                    </div>
                                ))}
                            </div>
                        </div>
                    </div>
                )}
            </div>

            {modal && (
                <div className={styles.modalOverlay} onClick={() => setModal(null)}>
                    <div className={styles.modal} onClick={(e) => e.stopPropagation()}>
                        <button className={styles.modalClose} onClick={() => setModal(null)}>✕</button>
                        <div className={styles.modalHeader}>
                            <Image src={`${RUNE_IMG_BASE}${modal.icon}`} alt={modal.name} width={64} height={64} unoptimized className={styles.modalIcon} />
                            <div>
                                <h3 className={styles.modalName}>{modal.name}</h3>
                                <p className={styles.modalKey}>{modal.key}</p>
                            </div>
                        </div>
                        <div className={styles.modalShortDesc}>{stripTags(modal.short_desc)}</div>
                        <div className={styles.modalDivider} />
                        <div className={styles.modalLongDesc}>{stripTags(modal.long_desc)}</div>
                    </div>
                </div>
            )}
        </div>
    );
}