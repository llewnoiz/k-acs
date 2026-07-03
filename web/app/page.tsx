import DeviceTable from '@/components/DeviceTable';

export default function Home() {
  return (
    <main className="mx-auto max-w-6xl px-6 py-10">
      <header className="mb-6">
        <h1 className="text-2xl font-bold text-slate-800">TR-069 장치 목록</h1>
        <p className="mt-1 text-sm text-slate-500">
          Inform 메시지를 보낸 CPE 장치들이 실시간으로 표시됩니다.
        </p>
      </header>
      <DeviceTable />
    </main>
  );
}
