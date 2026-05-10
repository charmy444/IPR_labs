import type { Metadata, Viewport } from 'next';
import 'bootstrap/dist/css/bootstrap.min.css';
import 'bootstrap-icons/font/bootstrap-icons.css';
import './globals.css';

export const metadata: Metadata = {
  title: 'Telegram Support Bot',
  description: 'Система технической поддержки на базе Telegram-бота',
  keywords: ['telegram', 'support', 'bot', 'техподдержка', 'чат'],
  authors: [{ name: 'Telegram Support Bot Team' }],
  creator: 'Telegram Support Bot',
  publisher: 'Telegram Support Bot',
  formatDetection: {
    email: false,
    address: false,
    telephone: false,
  },
  openGraph: {
    type: 'website',
    locale: 'ru_RU',
    url: 'https://telegram-support-bot.com',
    title: 'Telegram Support Bot',
    description: 'Система технической поддержки на базе Telegram-бота',
    siteName: 'Telegram Support Bot',
  },
  twitter: {
    card: 'summary_large_image',
    title: 'Telegram Support Bot',
    description: 'Система технической поддержки на базе Telegram-бота',
  },
  icons: {
    icon: '/favicon.ico',
    shortcut: '/favicon-16x16.png',
    apple: '/apple-touch-icon.png',
  },
  manifest: '/site.webmanifest',
};

export const viewport: Viewport = {
  width: 'device-width',
  initialScale: 1,
  maximumScale: 5,
  userScalable: true,
  themeColor: '#667eea',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="ru">
      <head>
        <link rel="icon" href="/favicon.ico" />
      </head>
      <body>{children}</body>
    </html>
  );
}
