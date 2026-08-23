import { useState } from "react";
import type { Model } from "../core/types";
import { modelProfileAxes } from "../domain/model-profile";
import { tx } from "../i18n/runtime";

export function ModelProfileRadar({ model, size = 260 }: { model: Model; size?: number }) {
  const axes = modelProfileAxes(model);
  const [hover, setHover] = useState<number | null>(null);
  const center = size / 2;
  const radius = size * 0.34;
  const labelRadius = size * 0.46;
  const count = axes.length;
  const angle = (index: number) => -Math.PI / 2 + (index * 2 * Math.PI) / count;
  const point = (index: number, ratio: number) => {
    const a = angle(index);
    return [center + Math.cos(a) * radius * ratio, center + Math.sin(a) * radius * ratio] as const;
  };
  const labelPoint = (index: number) => {
    const a = angle(index);
    return [center + Math.cos(a) * labelRadius, center + Math.sin(a) * labelRadius] as const;
  };
  const polyPoints = (ratio: number) =>
    axes.map((_, index) => point(index, ratio).join(",")).join(" ");
  const dataPoints = axes.map((axis, index) => point(index, axis.score / 100).join(",")).join(" ");
  const hovered = hover !== null ? axes[hover] : null;

  return (
    <div className="profile-radar" style={{ width: size + 28 }}>
      <svg viewBox={`0 0 ${size} ${size}`} width={size} height={size} role="img" aria-label={`${model.name} ${tx("模型画像")}`}>
        <defs>
          <linearGradient id="radarFill" x1="0" y1="0" x2="1" y2="1">
            <stop offset="0" stopColor="#3B82F6" stopOpacity="0.32" />
            <stop offset="1" stopColor="#1E40AF" stopOpacity="0.16" />
          </linearGradient>
        </defs>
        {/* grid rings */}
        {[0.25, 0.5, 0.75, 1].map((ratio) => (
          <polygon key={ratio} points={polyPoints(ratio)} className="profile-radar-ring" />
        ))}
        {/* spokes */}
        {axes.map((_, index) => {
          const [x, y] = point(index, 1);
          return <line key={index} x1={center} y1={center} x2={x} y2={y} className="profile-radar-spoke" />;
        })}
        {/* data polygon */}
        <polygon points={dataPoints} fill="url(#radarFill)" stroke="#1E40AF" strokeWidth="2" strokeLinejoin="round" className="profile-radar-data" />
        {/* vertices */}
        {axes.map((axis, index) => {
          const [x, y] = point(index, axis.score / 100);
          const active = hover === index;
          return (
            <circle
              key={axis.key}
              cx={x}
              cy={y}
              r={active ? 5.5 : 3.5}
              className={active ? "profile-radar-vertex active" : "profile-radar-vertex"}
              onMouseEnter={() => setHover(index)}
              onMouseLeave={() => setHover(null)}
            />
          );
        })}
        {/* axis labels */}
        {axes.map((axis, index) => {
          const [x, y] = labelPoint(index);
          const anchor = Math.abs(Math.cos(angle(index))) < 0.35 ? "middle" : Math.cos(angle(index)) > 0 ? "start" : "end";
          return (
            <text key={axis.key} x={x} y={y} textAnchor={anchor} dominantBaseline="middle" className="profile-radar-label">
              {tx(axis.label)}
            </text>
          );
        })}
        {/* score labels on vertices */}
        {axes.map((axis, index) => {
          const [x, y] = point(index, axis.score / 100);
          const inner = point(index, axis.score / 100 - (axis.score / 100 >= 0.5 ? 0.13 : -0.13));
          return (
            <text key={axis.key} x={inner[0]} y={inner[1]} textAnchor="middle" dominantBaseline="middle" className="profile-radar-score">
              {axis.score}
            </text>
          );
        })}
      </svg>
      {hovered ? (
        <div className="profile-radar-tip" style={{ left: 24, top: 8 }}>
          <strong>{tx(hovered.label)}</strong>
          <em>{hovered.score} / 100</em>
          <span>{hovered.evidence}</span>
        </div>
      ) : null}
    </div>
  );
}
