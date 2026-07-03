'use client';

import { useEffect, useState } from 'react';
import { fetchDevices, type Device } from '@/lib/api';

const POLL_MS = 5000;

function formatTime(iso: string): string {
  if (!iso) return '-';
  const d = new Date(iso);
  return Number.isNaN(d.getTime()) ? iso : d.toLocaleString();
}

export default function DeviceTable() {
  const [devices, setDevices] = useState<Device[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);
  const [updatedAt, setUpdatedAt] = useState<string>('');

  useEffect(() => {
    let active = true;

    async function load() {
      try {
        const data = await fetchDevices();
        if (!active) return;
        setDevices(data);
        setError(null);
        setUpdatedAt(new Date().toLocaleTimeString());
      } catch (e) {
        if (!active) return;
        setError(e instanceof Error ? e.message : 'unknown error');
      } finally {
        if (active) setLoading(false);
      }
    }

    load();
    const id = setInterval(load, POLL_MS);
    return () => {
      active = false;
      clearInterval(id);
    };
  }, []);

  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between text-sm text-slate-500">
        <span>
          총 <span className="font-semibold text-slate-700">{devices.length}</span> 대
        </span>
        <span>{updatedAt && `자동 갱신 · ${updatedAt}`}</span>
      </div>

      {error && (
        <div className="rounded-md border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
          ACS API 호출 실패: {error} — Go 서버(:7547)가 실행 중인지 확인하세요.
        </div>
      )}

      <div className="overflow-x-auto rounded-lg border border-slate-200 bg-white shadow-sm">
        <table className="min-w-full divide-y divide-slate-200 text-sm">
          <thead className="bg-slate-100 text-left text-xs font-semibold uppercase tracking-wide text-slate-500">
            <tr>
              <th className="px-4 py-3">시리얼번호</th>
              <th className="px-4 py-3">제조사</th>
              <th className="px-4 py-3">Product Class</th>
              <th className="px-4 py-3">IP</th>
              <th className="px-4 py-3">SW 버전</th>
              <th className="px-4 py-3">마지막 이벤트</th>
              <th className="px-4 py-3 text-right">Inform 횟수</th>
              <th className="px-4 py-3">마지막 Inform</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-100">
            {devices.map((d) => (
              <tr key={d.serialNumber} className="hover:bg-slate-50">
                <td className="px-4 py-3 font-mono font-medium text-slate-800">{d.serialNumber}</td>
                <td className="px-4 py-3">{d.manufacturer || '-'}</td>
                <td className="px-4 py-3">{d.productClass || '-'}</td>
                <td className="px-4 py-3 font-mono">{d.ip || '-'}</td>
                <td className="px-4 py-3">{d.softwareVersion || '-'}</td>
                <td className="px-4 py-3">
                  <span className="rounded bg-slate-100 px-2 py-0.5 text-xs text-slate-600">
                    {d.lastEvent || '-'}
                  </span>
                </td>
                <td className="px-4 py-3 text-right tabular-nums">{d.informCount}</td>
                <td className="px-4 py-3 text-slate-500">{formatTime(d.lastInformAt)}</td>
              </tr>
            ))}
            {devices.length === 0 && !loading && (
              <tr>
                <td colSpan={8} className="px-4 py-10 text-center text-slate-400">
                  아직 Inform 을 보낸 장치가 없습니다.
                </td>
              </tr>
            )}
            {loading && (
              <tr>
                <td colSpan={8} className="px-4 py-10 text-center text-slate-400">
                  불러오는 중…
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
