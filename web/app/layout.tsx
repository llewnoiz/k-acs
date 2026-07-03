import type { Metadata } from 'next';
import './globals.css';

export const metadata: Metadata = {
  title: 'TR-069 장치 목록',
  description: 'Inform 을 보낸 CPE 장치 목록'
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="ko">
      <body>{children}</body>
    </html>
  );
}
